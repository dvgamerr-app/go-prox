package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"prox/envs"
	"prox/pgsql"

	"github.com/alexflint/go-arg"
	"github.com/gofiber/contrib/otelfiber"
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

	goose.SetTableName("db_prox")

	// IsDev := os.Getenv("LOG_LEVEL") != ""
}

func (p *ProxService) Init() error {
	log.Info().Msgf("Init Fiber...")
	app = fiber.New(fiber.Config{
		ServerHeader:             fmt.Sprintf("%s/%s", envs.AppName, envs.Version),
		AppName:                  envs.AppName,
		DisableKeepalive:         true,
		DisableStartupMessage:    true,
		DisableDefaultDate:       true,
		DisableHeaderNormalizing: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Error().Msgf("%v", err)
			return nil
		},
	})

	app.Use(otelfiber.Middleware())

	app.Get("/health", func(c *fiber.Ctx) error {
		if strings.Contains(string(c.Request().Header.ContentType()), "json") {
			c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
			return c.SendString(`{"ok":"☕"}`)
		}
		return c.SendString(`☕`)
	})

	app.Use("/syscall", func(c *fiber.Ctx) error {
		c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
		return c.Next()
	})

	app.Get("/syscall/monitor", func(c *fiber.Ctx) error {
		if PostMessage(HWND_BROADCAST, WM_SYSCOMMAND, SC_MONITORPOWER, MONITOR_OFF) {
			return c.JSON(map[string]bool{"error": false})
		}
		return c.JSON(map[string]bool{"error": true})
	})

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(501).SendString(`🫗`)
	})
	go app.Listen(":3000")
	log.Info().Msgf("Starting listen :3000")

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
		goose.SetLogger(&gooseLogger{l: &log.Logger})

		ctx := context.Background()
		pgx := pgsql.Connect(&ctx)
		defer pgx.Close()

		if err := goose.RunContext(ctx, *args.DBInit, pgx.DB, "./pgsql/goose", args.Param...); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	if args.DaemonService {
		RunService(envs.AppName, envs.IsDev, prox)
	}
}
