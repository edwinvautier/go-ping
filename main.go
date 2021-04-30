package main

import (
	"time"

	"github.com/edwinvautier/test-ping/services"
	log "github.com/sirupsen/logrus"
)

func main() {
	ipBase := "192.168.0."
	log.Info("Starting pings for network ", ipBase+"0")
	services.PingNetwork(ipBase)
	log.Info("pinging almost done..")
	time.Sleep(7 * time.Second)
	arpList := services.GetARPTable()
	services.FindAllConstructor(&arpList)
	services.WriteList(arpList)
	log.WithField("file", "results.txt").Info("wrote results to file")
}
