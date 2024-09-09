package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"prox/envs"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	DB_POSTGRES = "DB_POSTGRES"
	PGHOST      = "PG_HOST"
	PGPORT      = "PG_PORT"
	PGUSER      = "PG_USER"
	PGPASSWORD  = "PG_PASS"
	PGDATABASE  = "PG_DBNAME"
	PGCACHE     = "PG_DBCACHE"
	PGSSLMODE   = "PG_SSLMODE"
	PGLIFETIME  = "PG_LIFETIME"
	PGMAXIDLE   = "PG_MAXIDLE"
	PGMAXCONN   = "PG_MAXCONN"
)

var ErrNoRows error = sql.ErrNoRows

func getSSLMode() string {
	sslmode := os.Getenv(PGSSLMODE)
	if strings.Contains(os.Getenv(PGSSLMODE), "") {
		sslmode = "disable"
	}
	return sslmode
}

func getDSN() string {
	if os.Getenv(DB_POSTGRES) == "" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s application_name='%s'",
			os.Getenv(PGHOST), os.Getenv(PGPORT), os.Getenv(PGUSER), os.Getenv(PGPASSWORD), os.Getenv(PGDATABASE), getSSLMode(), envs.AppName)
	}

	if !strings.Contains(os.Getenv(DB_POSTGRES), "application_name") {
		return fmt.Sprintf("%s application_name='%s'", os.Getenv(DB_POSTGRES), envs.AppName)
	}

	return os.Getenv(DB_POSTGRES)
}

type Client struct {
	DB  *sql.DB
	ctx *context.Context
}

// type PGRow map[string]string
// type PGRecord []PGRow

// type PGTx struct {
// 	Closed bool
// 	// tx     *sql.Tx
// 	// ctx    *context.Context
// }

func Connect(c *context.Context) *Client {
	var err error
	pg := Client{ctx: c}

	pg.DB, err = sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatal().Msgf("Postgres:: Open %v", err)
	}

	if os.Getenv(PGLIFETIME) != "" {
		lifeTimeSecond, err := strconv.ParseInt(os.Getenv(PGLIFETIME), 0, 64)
		if err != nil {
			log.Error().Msgf("ENV::PGLIFETIME ParseInt %v", err)
		}
		pg.DB.SetConnMaxLifetime(time.Second * time.Duration(lifeTimeSecond))
	}

	if os.Getenv(PGMAXIDLE) != "" {
		maxIdle, err := strconv.ParseInt(os.Getenv(PGMAXIDLE), 0, 32)
		if err != nil {
			log.Error().Msgf("ENV::PGMAXIDLE ParseInt %v", err)
		}

		pg.DB.SetMaxIdleConns(int(maxIdle))
	}

	if os.Getenv(PGMAXCONN) != "" {
		maxConn, err := strconv.ParseInt(os.Getenv(PGMAXCONN), 0, 32)
		if err != nil {
			log.Error().Msgf("ENV::PGMAXCONN ParseInt %v", err)
		}
		pg.DB.SetMaxOpenConns(int(maxConn))
	}

	err = pg.DB.PingContext(*pg.ctx)
	if err != nil {
		log.Fatal().Msgf("Postgres:: PingContext %v", err)
	}

	log.Debug().Msgf("pgsql::connected")
	return &pg
}

func (pg *Client) Close() error {
	log.Debug().Msgf("pgsql::closed")
	return pg.DB.Close()
}

// func (pg *Client) QueryOne(query string, args ...any) (PGRow, error) {
// 	rows, err := pg.DB.QueryContext(*pg.ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("Client.QueryOne::%s", err.Error())
// 	}
// 	if !rows.Next() {
// 		return nil, sql.ErrNoRows
// 	}
// 	defer rows.Close()
// 	return fetchRow(rows)

// }

// func (pg *Client) Query(query string, args ...any) (*sql.Rows, error) {
// 	rows, err := pg.DB.QueryContext(*pg.ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("Client.Query::%s", err.Error())
// 	}

// 	return rows, err
// }
// func (pg PGRow) ToByte(name string) []byte {
// 	return []byte(pg[name])
// }

// func (pg PGRow) ToBoolean(name string) bool {
// 	if pg[name] == "" {
// 		return false
// 	}

// 	data, err := strconv.ParseBool(pg[name])
// 	if err != nil {
// 		log.Error().Msgf("PGRow.ToBoolean('%s'): %s", name, err)
// 	}
// 	return data
// }
// func (pg PGRow) ToInt64(name string) int64 {
// 	data, err := strconv.ParseInt(pg[name], 0, 64)
// 	if err != nil {
// 		log.Error().Msgf("PGRow.ToInt64('%s', 0, 64): %s", name, err)
// 	}
// 	return data
// }
// func (pg PGRow) ToFloat64(name string) float64 {
// 	data, err := strconv.ParseFloat(pg[name], 64)
// 	if err != nil {
// 		log.Error().Msgf("PGRow.ToFloat64('%s', 64): %s", name, err)
// 	}
// 	return data
// }

// func (pg PGRow) ToTime(name string) time.Time {
// 	data, err := time.Parse(time.RFC3339Nano, pg[name])
// 	if err != nil {
// 		log.Error().Msgf("PGRow.ToTime('%s'): %s", name, err)
// 	}
// 	return data
// }

// const (
// 	LevelDefault sql.IsolationLevel = iota
// 	LevelReadUncommitted
// 	LevelReadCommitted
// 	LevelWriteCommitted
// 	LevelRepeatableRead
// 	LevelSnapshot
// 	LevelSerializable
// 	LevelLinearizable
// )

// func (pg *Client) Begin(level ...sql.IsolationLevel) (*PGTx, error) {
// 	// defer EstimatedPrint(time.Now(), fmt.Sprintf("Begin: %+v", pg.ctx))
// 	err := pg.DB.PingContext(*pg.ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(level) > 0 {
// 		stx, err := pg.DB.BeginTx(*pg.ctx, &sql.TxOptions{Isolation: level[0]})
// 		pgx := PGTx{tx: stx, ctx: pg.ctx}
// 		return &pgx, err
// 	} else {
// 		stx, err := pg.DB.BeginTx(*pg.ctx, &sql.TxOptions{})
// 		pgx := PGTx{tx: stx, ctx: pg.ctx}
// 		return &pgx, err
// 	}

// }

// func (stx *PGTx) Commit() error {
// 	stx.Closed = true
// 	return stx.tx.Commit()
// }

// func (stx *PGTx) Rollback() error {
// 	stx.Closed = true
// 	return stx.tx.Rollback()
// }

// func (stx *PGTx) QueryOne(query string, args ...any) (PGRow, error) {
// 	rows, err := sctxQuery(stx.tx, stx.ctx, false, query, args...)

// 	if err != nil {
// 		return nil, fmt.Errorf("QueryOne::%s", err.Error())
// 	}
// 	if !rows.Next() {
// 		return nil, sql.ErrNoRows
// 	}
// 	defer rows.Close()
// 	return fetchRow(rows)
// }

// func (stx *PGTx) QueryOnePrint(query string, args ...any) (PGRow, error) {
// 	rows, err := sctxQuery(stx.tx, stx.ctx, true, query, args...)

// 	if err != nil {
// 		return nil, fmt.Errorf("QueryOne::%s", err.Error())
// 	}
// 	if !rows.Next() {
// 		return nil, sql.ErrNoRows
// 	}
// 	defer rows.Close()
// 	return fetchRow(rows)
// }

// func (stx *PGTx) Query(query string, args ...any) (*sql.Rows, error) {
// 	return sctxQuery(stx.tx, stx.ctx, false, query, args...)
// }

// func (stx *PGTx) QueryPrint(query string, args ...any) (*sql.Rows, error) {
// 	return sctxQuery(stx.tx, stx.ctx, true, query, args...)
// }

// func (stx *PGTx) Execute(query string, args ...any) error {
// 	return sctxExecute(stx.tx, stx.ctx, false, query, args...)
// }

// func (stx *PGTx) ExecutePrint(query string, args ...any) error {
// 	return sctxExecute(stx.tx, stx.ctx, true, query, args...)
// }

// func (stx *PGTx) FetchRow(rows *sql.Rows) (PGRow, error) {
// 	return fetchRow(rows)
// }

// func (pg *Client) FetchRow(rows *sql.Rows) (PGRow, error) {
// 	return fetchRow(rows)
// }

// func (stx *PGTx) FetchAll(rows *sql.Rows) (PGRecord, error) {
// 	result := []PGRow{}
// 	for rows.Next() {
// 		data, err := stx.FetchRow(rows)
// 		if err != nil {
// 			return PGRecord{}, nil
// 		}

// 		result = append(result, data)
// 	}
// 	return result, nil
// }
// func (stx *PGTx) FetchOneColumn(rows *sql.Rows, columnName string) (SubSet, error) {
// 	result := SubSet{}
// 	for rows.Next() {
// 		data, err := stx.FetchRow(rows)
// 		if err != nil {
// 			return SubSet{}, nil
// 		}

// 		result = append(result, data[columnName])
// 	}
// 	return result, nil
// }

// func (row PGRecord) Find(columnName string, compareValue string) bool {
// 	for i := 0; i < len(row); i++ {
// 		if row[i][columnName] == compareValue {
// 			return true
// 		}
// 	}
// 	return false
// }

// func fetchRow(rows *sql.Rows) (PGRow, error) {
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return nil, fmt.Errorf("FetchRow::Columns::%v", err)
// 	}

// 	resultMap := make(PGRow)
// 	values := make([]any, len(columns))
// 	pointers := make([]any, len(columns))
// 	for i := range values {
// 		pointers[i] = &values[i]
// 	}
// 	err = rows.Scan(pointers...)
// 	if err == sql.ErrNoRows {
// 		return resultMap, fmt.Errorf("FetchRow::ErrNoRows: %v", err)
// 	} else if err != nil {
// 		return nil, fmt.Errorf("FetchRow::Scan: %v", err)
// 	}

// 	for i, val := range values {
// 		if reflect.TypeOf(val) == nil {
// 			resultMap[columns[i]] = ""
// 			continue
// 		}
// 		switch reflect.TypeOf(val).String() {
// 		case "int64":
// 			resultMap[columns[i]] = fmt.Sprint(val.(int64))
// 		case "float64":
// 			resultMap[columns[i]] = fmt.Sprint(val.(float64))
// 		case "string":
// 			resultMap[columns[i]] = val.(string)
// 		case "[]uint8":
// 			resultMap[columns[i]] = string(val.([]uint8))
// 		case "bool":
// 			resultMap[columns[i]] = fmt.Sprintf("%t", val.(bool))
// 		case "time.Time":
// 			resultMap[columns[i]] = val.(time.Time).Format(time.RFC3339Nano)
// 		default:
// 			log.Error().Msgf("Reflect TypeOf: %s ", reflect.TypeOf(val).String())
// 			resultMap[columns[i]] = ""
// 		}
// 	}
// 	return resultMap, nil
// }

// func sctxQuery(pgstx *sql.Tx, pgctx *context.Context, envDebug bool, query string, args ...any) (*sql.Rows, error) {
// 	// elapsed := time.Now()

// 	if envDebug {
// 		// defer sqlQuery(elapsed, query, args...)
// 	}
// 	// defer EstimatedPrint(elapsed, "pg::query")

// 	return pgstx.QueryContext(*pgctx, query, args...)
// }

// func sctxExecute(pgstx *sql.Tx, pgctx *context.Context, envDebug bool, query string, args ...any) error {
// 	// elapsed := time.Now()
// 	if envDebug {
// 		// defer sqlQuery(elapsed, query, args...)
// 	}

// 	// defer EstimatedPrint(elapsed, "pg::execute")

// 	_, err := pgstx.ExecContext(*pgctx, query, args...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func leadingSpace(line string) int {
// 	count := 0
// 	for _, v := range line {
// 		if v == ' ' || v == '\t' {
// 			count++
// 		} else {
// 			break
// 		}
// 	}
// 	return count
// }

// type SubSet []string

// func (s *SubSet) ToParam() string {
// 	return fmt.Sprintf("{%s}", strings.Join(*s, ","))
// }
// func (s *SubSet) Find(val string) int {
// 	for ix, v := range *s {
// 		if v == val {
// 			return ix
// 		}
// 	}
// 	return len(*s)
// }

// func estimated(start time.Time) int {
// 	duration, _ := elapsedDuration(start)
// 	return int(float64(duration.Microseconds()) / 1000)
// }

// func EstimatedPrint(start time.Time, name string, ctx ...*fiber.Ctx) {
// 	if os.Getenv(DEBUG) == "false" && os.Getenv(ENV) == "production" {
// 		return
// 	}
// 	_, elapsed := elapsedDuration(start)

// 	pc, _, _, _ := runtime.Caller(1)
// 	funcObj := runtime.FuncForPC(pc)
// 	if name == "" {
// 		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
// 		name = runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
// 	}
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)
// 	// Debugf("%s # %s estimated. | alloc: %vMiB (%vMiB), sys: %vMiB, gc: %vMiB", name, elapsed, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)

// 	if len(ctx) != 0 && ctx[0] != nil {
// 		ctx[0].Append("Server-Timing", fmt.Sprintf("app;dur=%v", elapsed))
// 	}
// 	log.Debug().Msgf("%s # %s estimated.", name, elapsed)
// }

// func elapsedDuration(start time.Time) (time.Duration, string) {
// 	duration := time.Since(start)

// 	elapsed := ""
// 	if duration.Nanoseconds() < 1000 {
// 		elapsed = fmt.Sprintf("%dns", duration.Nanoseconds())
// 	} else if duration.Microseconds() < 1000 {
// 		elapsed = fmt.Sprintf("%0.3fμs", round(float64(duration.Nanoseconds())/1000, 2))
// 	} else if duration.Milliseconds() < 1000 {
// 		elapsed = fmt.Sprintf("%0.3fms", round(float64(duration.Microseconds())/1000, 2))
// 	} else if duration.Seconds() < 60 {
// 		elapsed = fmt.Sprintf("%0.3fms", round(float64(duration.Microseconds())/1000, 2))
// 	} else {
// 		elapsed = fmt.Sprintf("%0.3fm", round(float64(duration.Seconds()/60), 2))
// 	}
// 	return duration, elapsed
// }

// // round math round decimal
// func round(n float64, m float64) float64 {
// 	return math.Round(n*math.Pow(10, m)) / math.Pow(10, m)
// }
