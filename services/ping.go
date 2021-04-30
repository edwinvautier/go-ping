package services

import (
	"fmt"
	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func PingNetwork(ipBase string) {
	for adressNumber := 1; adressNumber < 255; adressNumber++ {
		fullIP := ipBase + fmt.Sprint(adressNumber)

		go PingIP(fullIP)
		time.Sleep(60 * time.Millisecond)
	}
}

// PingIP
func PingIP(ip string) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Interval = 3 * time.Second
	pinger.Timeout = 3 * time.Second

	if err := pinger.Run(); err != nil {
	} else {
		if pinger.Statistics().PacketLoss < 2 {
			if err := writeIPtoFile(ip); err != nil {
				log.Error(err)
			}
		}
	}
}

func writeIPtoFile(ip string) error {
	f, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(ip + "\n"); err != nil {
		return err
	}

	return nil
}
