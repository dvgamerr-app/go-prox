package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"prox/envs"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	// configExt string = "yaml"
	logExt  string = "log"
	logFile *os.File
)

var args struct {
	DaemonService bool `arg:"--daemon,-d" default:"false"`
	// ListBrowser   bool   `arg:"--list,-l"`
	// Token         string `arg:"--token" help:"set token one job"`
	// Generate      string `arg:"--gen" help:"want value 'browser:profile'"`
}

func init() {
	arg.MustParse(&args)

	if err := envs.Load(); err != nil {
		log.Error().Err(err)
	}

	if err := initLogging(); err != nil {
		log.Error().Err(err)
	}
	// IsDev := os.Getenv("LOG_LEVEL") != ""
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

type ProxService struct{}

func (p *ProxService) Init() error {
	log.Info().Msgf("Init()")
	return nil
}

func (p *ProxService) Tick() error {
	log.Info().Msgf("Tick()")
	return nil
}

func (p *ProxService) Shutdown() error {
	log.Info().Msgf("Shutdown()")
	return nil
}

func main() {
	defer logFile.Close()

	var prox *ProxService

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
