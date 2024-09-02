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

	if !envs.IsDev || args.DaemonService {
		log.Info().Msgf("goProx starting...")
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

		log.Output(logFile)
	} else {
		timeFormat := time.DateTime
		if envs.IsDev {
			timeFormat = time.TimeOnly
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: timeFormat})
		log.Info().Msgf("goProx is Development mode.")
	}

	log.Info().Msgf("OS: %s Arch: %s", runtime.GOOS, runtime.GOARCH)

	return nil
}

func main() {
	defer logFile.Close()
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

// type GoProx struct{}

// func (m *GoProx) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

// 	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
// 	tick := time.Tick(5 * time.Second)

// 	status <- svc.Status{State: svc.StartPending}

// 	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

// loop:
// 	for {
// 		select {
// 		case <-tick:
// 			log.Print("Tick Handled...!")
// 		case c := <-r:
// 			switch c.Cmd {
// 			case svc.Interrogate:
// 				status <- c.CurrentStatus
// 			case svc.Stop, svc.Shutdown:
// 				log.Print("Shutting service...!")
// 				break loop
// 			case svc.Pause:
// 				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
// 			case svc.Continue:
// 				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
// 			default:
// 				log.Printf("Unexpected service control request #%d", c)
// 			}
// 		}
// 	}

// 	status <- svc.Status{State: svc.StopPending}
// 	return false, 1
// }

// func runService(name string, isDebug bool) {
// 	if isDebug {
// 		err := debug.Run(name, &GoProx{})
// 		if err != nil {
// 			log.Fatalln("Error running service in debug mode.")
// 		}
// 	} else {
// 		err := svc.Run(name, &GoProx{})
// 		if err != nil {
// 			log.Fatalln("Error running service in Service Control mode.")
// 		}
// 	}
// }
