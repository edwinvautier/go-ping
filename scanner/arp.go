package scanner

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

func (n *Network) FindDevices() []Device {
	records := GetARPRecords()
	recordsNumber := int64(len(records))

	constructors := make([]string, 0)
	constructorChan := make(chan string, recordsNumber)
	n.FindAllConstructors(constructorChan, records)
	close(constructorChan)
	for constructor := range constructorChan {
		constructors = append(constructors, constructor)
	}

	devices := make([]Device, 0)
	for i, record := range records {
		device := Device{
			Constructor: constructors[i],
			IP:          record.IP,
			Mac:         record.Mac,
		}
		devices = append(devices, device)
	}

	return devices
}

func GetARPRecords() []*ARPRecord {
	output, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(output)
	lines := strings.Split(outputString, "\n")
	arpList := make([]*ARPRecord, 0)
	for _, line := range lines {
		if strings.Contains(line, "incomplet") || len(line) < 4 {
			continue
		}
		elements := strings.Split(line, " ")

		ip := elements[1]
		mac := elements[3]
		arpList = append(arpList, &ARPRecord{IP: ip, Mac: mac})
	}

	return arpList
}

func GetConstructor(mac string) string {
	url := "https://api.macvendors.com/"
	resp, err := http.Get(url + mac)
	if err != nil || resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			dur := rand.Intn(5)
			time.Sleep(time.Duration(dur) * time.Second)
			return GetConstructor(mac)
		} else {
			return "not-found"
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "not-found"
	}

	company := string(body)
	return company
}

func (n *Network) FindAllConstructors(constructorChan chan string, records []*ARPRecord) {
	recordsNumber := int64(len(records))
	bar := progressbar.Default(recordsNumber)
	bar.Describe("mac lookup API")
	bar.Add(1)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for _, record := range records {
		n.Lock.Acquire(context.TODO(), 1)
		wg.Add(1)

		go func(record ARPRecord) {
			defer n.Lock.Release(1)
			defer wg.Done()
			defer bar.Add(1)
			constructorChan <- GetConstructor(record.Mac)
		}(*record)
	}
}

type MACLookupBody struct {
	Company string
}
