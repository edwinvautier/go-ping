package main

import (
	"fmt"

	"github.com/edwinvautier/go-ping/scanner"
	log "github.com/sirupsen/logrus"
)

func main() {
	ipBase := "192.168.0."
	network := scanner.NetworkClient(ipBase)

	network.FindIPs()
	devices := network.FindDevices()
	network.ScanAllDevicesPorts(&devices)

	for i, device := range devices {
		log.Info("------------------")
		log.Info("DEVICE ", i+1)
		log.Info("IP: ", device.IP)
		log.Info("MAC: ", device.Mac)
		log.Info("CONSTRUCTOR: ", device.Constructor)
		log.Info("OPEN PORT: ")
		portsString := ""
		for _, port := range device.OpenPorts {
			portsString += fmt.Sprintf("%d", port.Number) + " - "
		}
		log.Info(portsString)
		log.Info("------------------")
	}
}
