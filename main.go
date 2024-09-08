package main

import (
	"context"
	"os"
	"strings"

	"prox/envs"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/fiber/v2"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

type ProxService struct{}

var (
	// configExt string = "yaml"
	logExt  string = "log"
	logFile *os.File
	prox    *ProxService
	app     *fiber.App
)

var args struct {
	DBInit        *string  `arg:"--db"`
	DaemonService bool     `arg:"--daemon,-d" default:"false"`
	Version       bool     `arg:"--version,-v" default:"false"`
	Param         []string `arg:"positional"`
	// ListBrowser   bool   `arg:"--list,-l"`
	// Token         string `arg:"--token" help:"set token one job"`
	// Generate      string `arg:"--gen" help:"want value 'browser:profile'"`
}

func init() {
	arg.MustParse(&args)
	if err := initVersion(); err != nil {
		log.Error().Err(err)
	}
	if args.Version {
		if err := printVersion(); err != nil {
			log.Error().Err(err)
		}
		os.Exit(0)
	}

	if err := envs.Load(); err != nil {
		log.Error().Err(err)
	}

	if err := initLogging(); err != nil {
		log.Error().Err(err)
	}

	goose.SetTableName("db_version")

	// IsDev := os.Getenv("LOG_LEVEL") != ""
}

func (p *ProxService) Init() error {
	log.Info().Msgf("Init...")
	app = fiber.New(fiber.Config{
		ServerHeader:          "go-prox/" + os.Getenv("VERSION"),
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Error().Msgf("%v", err)
			return nil
		},
		DisableDefaultDate:       true,
		DisableHeaderNormalizing: true,
		AppName:                  "Prox",
	})

	// GET /api/register
	app.Get("/health", func(c *fiber.Ctx) error {
		if strings.Contains(string(c.Request().Header.ContentType()), "json") {
			c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
			return c.SendString(`{"ok":"☕"}`)
		}
		return c.SendString(`☕`)
	})

	log.Info().Msgf("Starting listen :3000")

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(501).SendString(`🫗`)
	})
	go app.Listen(":3000")

	return nil
}

func (p *ProxService) Tick() error {
	return nil
}

func (p *ProxService) Shutdown() error {
	log.Info().Msgf("Shutdown()")
	if app != nil {
		if err := app.Shutdown(); err != nil {
			log.Error().Msgf("%v", err)
		}
	}
	return nil
}

func main() {
	defer logFile.Close()

	if args.DBInit != nil {
		ctx := context.Background()

		if err := goose.RunContext(ctx, *args.DBInit, nil, "./database", args.Param...); err != nil {
			log.Error().Msgf("%v", err)
		}

		log.Debug().Msg("connection closed.")
		os.Exit(0)
	}

	if args.DaemonService {
		RunService("GoProx", envs.IsDev, prox)
	}
}

// func main() {
// 	f, err := os.OpenFile("G:/.dvgamerr/go-prox/debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatalln(fmt.Errorf("error opening file: %v", err))
// 	}
// 	defer f.Close()

// 	log.SetOutput(f)
// 	runService("GoProx", false)
// }
