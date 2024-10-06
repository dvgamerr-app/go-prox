package main

import (
	probing "github.com/prometheus-community/pro-bing"
	"github.com/rs/zerolog/log"
)

type ProxService struct{}

var pingStats []*probing.Pinger

func (p *ProxService) Init() error {
	log.Info().Msgf("Init...")
	pingStats = make([]*probing.Pinger, 1)

	pinger, err := probing.NewPinger("103.206.205.129")
	pinger.SetPrivileged(true)
	if err != nil {
		log.Fatal().Msgf("NewPinger::%v", err)
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		log.Debug().Msgf("%d bytes from %s: icmp_seq=%d time=%v",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	go func() {
		if err := pinger.Run(); err != nil {
			log.Fatal().Msgf("Run::%v", err)
		}
	}()

	pingStats = append(pingStats, pinger)
	log.Info().Msgf("Inited")
	return nil
}

func (p *ProxService) Tick() error {
	return nil
}

func (p *ProxService) Shutdown() error {
	log.Info().Msgf("Shutdownting...")
	if app != nil {
		if err := app.Shutdown(); err != nil {
			log.Error().Msgf("%v", err)
		}
	}

	for _, pinger := range pingStats {
		if pinger != nil {
			pinger.Stop()
		}
	}

	log.Info().Msgf("Shutdown")
	return nil
}
