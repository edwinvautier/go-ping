package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/edwinvautier/go-ping/scanner"
)

func main() {
	ipBase := "192.168.0."
	network := scanner.NetworkClient(ipBase)

	log.Info(network.FindIPs())
}
