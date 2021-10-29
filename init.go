package daas

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/tmilewski/goenv"
)

func Initialize(appName string) (appTitle string, appVersion string, appIsProduction bool) {
	appIsProduction = os.Getenv(ENV) == "production"
	if !appIsProduction {
		goenv.Load()
	}
	content, err := ioutil.ReadFile("VERSION")
	if err != nil {
		content, _ = ioutil.ReadFile("../VERSION")
	}

	appVersion = strings.TrimSpace(string(content))
	appTitle = fmt.Sprintf("%s@%s", appName, appVersion)

	Infof("%s (%s) initilize...", appTitle, os.Getenv(ENV))

	hostname, err := os.Hostname()
	if err != nil {
		Fatalf("Hostname: %s", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		ServerName:  hostname,
		Environment: os.Getenv(SENTRY_ENV),
		Release:     appTitle,
		Dsn:         os.Getenv(SENTRY_DSN),
	})

	if err != nil {
		Fatalf("sentry.Init: %s", err)
	}
	return
}

type DataStore struct {
	AppID      string
	CourierMap map[string]*Courier
	BuAPIMap   map[string]bool
	TokenSet   SubSet
}

func (d *DataStore) GetCourierId(auth string) string {
	if ix := d.TokenSet.Find(auth); ix < len(d.TokenSet) {
		return d.CourierMap[d.TokenSet[ix]].ID
	}
	return "1"
}
func (d *DataStore) GetCourierConfig(courierId string) *Courier {
	for _, v := range d.CourierMap {
		if v.ID == courierId {
			return v
		}
	}
	return &Courier{}
}

func (d *DataStore) BUCodeVerifyAPI(buCode string) error {
	if !d.BuAPIMap[strings.ToUpper(buCode)] {
		return fmt.Errorf("bussiness Unit '%s' is't supported", strings.ToUpper(buCode))
	}
	return nil
}

func PrepareDataStore(pgx *PGClient, ctx *context.Context, appTitle string, appName string) *DataStore {
	data := DataStore{}
	pgx.Connect(ctx, appTitle)
	stx, err := pgx.Begin()
	if TraceIsError(err, nil) {
		Fatal("Begin::", err)
	}

	app, err := stx.QueryOne("SELECT id cf_stock_app_id FROM sync.cf_stock_app WHERE s_name = $1", appName)
	if TraceIsError(err, nil) {
		Fatal("cf_stock_app::", err)
	}
	data.AppID = app["cf_stock_app_id"]

	// err = data.GetCourierToken(stx)
	// if TraceIsError(err, nil) {
	// 	Fatal("getCourierToken::", err)
	// }
	err = stx.Commit()
	if TraceIsError(err, nil) {
		Fatalf("Commit: %s", err)
	}
	return &data
}

func (d *DataStore) SetConfigBUAllow(stx *PGTx) error {
	cfBU := make(map[string]bool)
	raw, err := stx.Query(`
		SELECT
			coalesce(o_global->'apis', 'false') apis,
			b.id m_bu_id, b.s_code bu_code
		FROM sync.cf_stock_bu sb
		INNER JOIN main.m_bu b ON b.id = sb.m_bu_id;
	`)
	if err != nil {
		return err
	}

	defer raw.Close()
	for raw.Next() {
		cf, err := stx.FetchRow(raw)
		if err != nil {
			return err
		}

		cfBU[cf["bu_code"]] = cf.ToBoolean("apis")
	}
	d.BuAPIMap = cfBU
	return nil
}

type Courier struct {
	ID     string
	Paging struct {
		Limit  int `json:"limit"`
		Record int `json:"record"`
	}
}

func (d *DataStore) SetConfigCourierToken(stx *PGTx) error {
	cfID := make(map[string]*Courier)
	cfToken := []string{}

	raw, err := stx.Query(`
		SELECT
			ct.cf_courier_app_id, ca.s_name, ca.o_attr, s_client || ':' || s_secret s_auth
		FROM sync.cf_courier_token ct
		INNER JOIN sync.cf_courier_app ca ON ca.id = ct.cf_courier_app_id
		WHERE e_status = 'ACTIVE'
	`)
	if err != nil {
		return err
	}

	defer raw.Close()
	for raw.Next() {
		cf, err := stx.FetchRow(raw)
		if err != nil {
			return err
		}
		cfToken = append(cfToken, cf["s_auth"])
		cfID[cf["s_auth"]] = &Courier{ID: cf["cf_courier_app_id"]}
		json.Unmarshal(cf.ToByte("o_attr"), &cfID[cf["s_auth"]].Paging)
	}
	sort.Strings(cfToken)
	d.CourierMap = cfID
	d.TokenSet = cfToken

	return nil
}
