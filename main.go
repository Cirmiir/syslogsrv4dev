package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "net/http/pprof"

	store "github.com/cirmiir/syslogsrv4dev/storages"
	_ "github.com/mattn/go-sqlite3"
)

type configuration struct {
	SyslogPort        int
	WebPort           int
	UdpPort           int
	StorageAddress    string
	UdpRegExpTemplate string
}

type JSONTime time.Time

type getter func(obj interface{}) (interface{}, error)

var defaultSettings configuration = configuration{
	StorageAddress:    "./test.db",
	SyslogPort:        7570,
	WebPort:           5050,
	UdpPort:           9090,
	UdpRegExpTemplate: `(?s)^(?P<date>.*?)\|(?P<levelId>.*?)\|(?P<facilityid>.*?)\|(?P<message>.*?)\|(?P<host>.*?)\|(?P<app>.*?)$`,
}

var storage store.LogStorage
var cfg configuration

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}

func loadSettings(fileName string) (configuration, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return defaultSettings, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := defaultSettings
	err = decoder.Decode(&conf)

	return conf, err
}

func main() {
	var err error
	cfg, err = loadSettings("conf.json")

	if err != nil {
		log.Println("error:", err)
	}

	var wg sync.WaitGroup
	storage = store.NewSqlLite3Storage(cfg.StorageAddress)
	if cfg.SyslogPort > 0 {
		wg.Add(1)
		go runSysloglistener(&wg, &storage)
	}

	if cfg.UdpPort > 0 {
		wg.Add(1)
		go runUdpListener(&wg, &storage)
	}

	wg.Add(1)
	go runWeblistener(&wg)

	wg.Wait()

}
