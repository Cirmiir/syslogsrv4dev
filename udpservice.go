package main

import (
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	store "github.com/cirmiir/syslogsrv4dev/storages"
)

func display(conn *net.UDPConn) (string, error) {

	buf := make([]byte, 0)
	tmp := make([]byte, 1024)
	for {
		n, err := conn.Read(tmp)
		if n < len(tmp) {
			buf = append(buf, tmp[:n]...)
			break
		}
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		buf = append(buf, tmp[:n]...)

	}

	return string(buf), nil
}

func getParams(regEx, url string) (paramsMap map[string]interface{}) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]interface{})
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func runUdpListener(wg *sync.WaitGroup, storage *store.LogStorage) {
	defer wg.Done()
	protocol := "udp"

	udpAddr, err := net.ResolveUDPAddr(protocol, ":"+strconv.Itoa(cfg.UdpPort))
	if err != nil {
		log.Println("Wrong Address", err)
		return
	}

	udpConn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		log.Println(err)
	}

	for {
		if val, err := display(udpConn); err == nil {
			dict := getParams(cfg.UdpRegExpTemplate, val)
			(*storage).SaveEvent(store.LogEvent{
				Host:     coalesce(dict, ParseAsString, "unknown", "host").(string),
				AppName:  coalesce(dict, ParseAsString, "unknown", "app").(string),
				Message:  coalesce(dict, ParseAsString, "", "message").(string),
				Facility: coalesce(dict, ParseAsInt, 0, "facilityid").(int),
				Severity: coalesce(dict, ParseAsInt, 0, "levelId").(int),
				Date:     coalesce(dict, ParseAsTime, time.Unix(0, 0), "date").(time.Time),
			})
		}
	}

}
