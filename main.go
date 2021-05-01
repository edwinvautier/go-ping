package main

import (
	"github.com/edwinvautier/go-ping/scanner"
	log "github.com/sirupsen/logrus"
)

func main() {
	ipBase := "192.168.0."
	network := scanner.NetworkClient(ipBase)

	log.Info(network.FindIPs())
	log.Info(network.FindDevices())
}
