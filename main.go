package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"prox/envs"
	"prox/pgsql"

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
	log.Info().Msgf("Init...")
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

	ctx := context.Background()
	pgx := pgsql.Connect(&ctx)
	defer pgx.Close()

	if args.DBInit != nil {
		goose.SetLogger(&gooseLogger{l: &log.Logger})

		if err := goose.RunContext(ctx, *args.DBInit, pgx.DB, "./pgsql/goose", args.Param...); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

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

	app.Use("/notify", func(c *fiber.Ctx) error {
		if c.Get("x-user-liff") == "" {
			return fiberThrowError(c, fiber.StatusUnauthorized, errors.New("Unauthorized"))
		}
		return c.Next()
	})

	app.Put("/notify/:serviceName/:roomName", handlerPutNotify(ctx, pgx))

	app.Get("/health", func(c *fiber.Ctx) error {
		if strings.Contains(string(c.Request().Header.ContentType()), "json") {
			c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
			return c.SendString(`{"ok":"☕"}`)
		}
		return c.SendString(`☕`)
	})

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(map[string]bool{"ok": false})
	})
	go app.Listen(":11535")
	log.Info().Msgf("Starting listen :11535")

	if args.DaemonService {
		RunService(envs.AppName, envs.IsDev, prox)
	}
}
