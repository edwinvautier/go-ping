package scanner

import (
	"golang.org/x/sync/semaphore"
	"os/exec"
	"strconv"
	"strings"
)

type NetworkInterface interface {
	FindDevices() []Device
	FindIPs() []string
	PingIP(string)
	Init(string)
}

type ARPRecord struct {
	IP  string
	Mac string
}

type Network struct {
	IP      string
	Mask    string
	Devices []Device
	Lock    *semaphore.Weighted
}

type Device struct {
	IP          string
	Mac         string
	Constructor string
	OpenPorts   []uint
}

func ulimit() int64 {
	out, err := exec.Command("sh", "-c", "ulimit -n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func NetworkClient(ip string) *Network {
	network := Network{}
	network.Init(ip)

	return &network
}
