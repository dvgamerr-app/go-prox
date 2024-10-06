package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type ServiceWindows struct {
	handle ServiceHandled
}

type ServiceHandled interface {
	Init() error
	Tick() error
	Shutdown() error
}

func (w *ServiceWindows) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	tick := time.Tick(5 * time.Second)

	status <- svc.Status{State: svc.StartPending}

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	if err := w.handle.Init(); err != nil {
		log.Fatal().Msgf("Error initializing service. %v", err)
		return true, 1
	}
loop:
	for {
		select {
		case <-tick:
			if err := w.handle.Tick(); err != nil {
				log.Error().Msgf("%v", err)
				return true, 1
			}

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				if err := w.handle.Shutdown(); err != nil {
					log.Error().Msgf("%v", err)
					return true, 0
				}
				break loop
			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Info().Msgf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

func RunService(name string, isDebug bool, handle ServiceHandled) {
	if isDebug {
		err := debug.Run(name, &ServiceWindows{handle: handle})
		if err != nil {
			log.Fatal().Msgf("Error running service in debug mode. %v", err)
		}
	} else {
		err := svc.Run(name, &ServiceWindows{handle: handle})
		if err != nil {
			log.Fatal().Msg("Error running service in Service Control mode.")
		}
	}
}
