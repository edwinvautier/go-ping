package scanner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

func (n *Network) Init(ip string) {
	n.IP = ip
	n.Lock = semaphore.NewWeighted(ulimit())
}

func (n *Network) FindIPs() []string {
	addresses := make([]string, 0)
	ipChan := make(chan string, 256)
	n.pingAll(ipChan)
	close(ipChan)
	for ip := range ipChan {
		addresses = append(addresses, ip)
	}

	return addresses
}

func (n *Network) pingAll(ipChan chan string) {
	bar := progressbar.Default(254)
	bar.Describe("pinging IP adresses")
	bar.Add(1)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for adressNumber := 1; adressNumber < 255; adressNumber++ {
		n.Lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		fullIP := n.IP + fmt.Sprint(adressNumber)

		go func(fullIP string) {
			defer n.Lock.Release(1)
			defer wg.Done()
			defer bar.Add(1)
			if n.PingIP(fullIP) {
				ipChan <- fullIP
			}
		}(fullIP)
	}
}

// PingIP pings the given address, returns true if device responded, false otherwise
func (n *Network) PingIP(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Error("pinging ip: ", ip, " : ", err)
		return false
	}
	pinger.Count = 3
	pinger.Interval = 3 * time.Second
	pinger.Timeout = 3 * time.Second

	if err := pinger.Run(); err != nil {
		return false
	}

	return pinger.Statistics().PacketLoss < 2
}
