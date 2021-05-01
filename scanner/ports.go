package scanner

import (
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"net"
	"strings"
	"sync"
	"time"
)

func ScanPort(ip string, port int) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(500 * time.Millisecond)
			return ScanPort(ip, port)
		}

		return false
	}

	conn.Close()
	return true
}

func (n *Network) ScanDevice(ip string, l int, portChan chan Port) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	bar := progressbar.Default(int64(l))
	bar.Describe(ip + ": ports analyzed")
	bar.Add(1)
	for port := 1; port <= l; port++ {
		n.Lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer n.Lock.Release(1)
			defer wg.Done()
			defer bar.Add(1)
			portChan <- Port{
				Number: port,
				Open:   ScanPort(ip, port),
			}
		}(port)
	}
}

type Port struct {
	Number int
	Open   bool
}

func (n *Network) ScanAllDevicesPorts(devices *[]*Device) {
	for _, device := range *devices {
		portChan := make(chan Port, 9999)
		n.ScanDevice(device.IP, 9999, portChan)
		close(portChan)
		ports := make([]Port, 0)
		for port := range portChan {
			if port.Open {
				ports = append(ports, port)
			}
		}
		device.OpenPorts = ports
	}
}
