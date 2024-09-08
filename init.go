package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"prox/envs"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func printVersion() error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("ReadBuildInfo can't read")
	}
	fmt.Printf("goProx %s (commit/%s, %s, %s/%s)", os.Getenv("VERSION"), bi.Main.Version, bi.GoVersion, runtime.GOOS, runtime.GOARCH)
	return nil
}

func initVersion() error {
	if body, err := os.ReadFile("VERSION"); err != nil {
		return err
	} else {
		os.Setenv("VERSION", strings.TrimSpace(string(body)))
	}
	return nil
}

func initLogging() error {
	logLevel := os.Getenv("LOG_LEVEL")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if logLevel == "info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if logLevel == "trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if !envs.IsDev && args.DaemonService {
		var err error
		execFilename, err := os.Executable()
		if err != nil {
			log.Fatal().Err(err)
		}
		baseFilename := strings.ReplaceAll(filepath.Base(execFilename), filepath.Ext(execFilename), "")
		dirname := filepath.Dir(execFilename)
		logPath := path.Join(dirname, fmt.Sprintf("%s.%s", baseFilename, logExt))

		logFile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err)
		}

		log.Logger = log.Output(logFile)
		log.Info().Msgf("goProx starting...")
	} else {
		timeFormat := time.DateTime
		if envs.IsDev {
			timeFormat = time.TimeOnly
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})
		log.Info().Msgf("goProx is Development mode.")
	}

	log.Info().Msgf("os: %s arch: %s", runtime.GOOS, runtime.GOARCH)

	return nil
}
