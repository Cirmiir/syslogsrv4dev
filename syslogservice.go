package main

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	store "github.com/cirmiir/syslogsrv4dev/storages"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func ParseAsInt(obj interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, errors.New("wrong type")
	}
}

func ParseAsTime(obj interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case time.Time:
		return v, nil
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return time.Unix(0, 0), errors.New("wrong type")
	}
}

func ParseAsString(obj interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New("wrong type")
	}
}

func coalesce(log map[string]interface{}, fun getter, def interface{}, fields ...string) interface{} {
	for _, field := range fields {
		val, ok := log[field]
		if ok {
			result, err := fun(val)
			if err == nil {
				return result
			}
		}
	}

	return def
}

func runSysloglistener(wg *sync.WaitGroup, storage *store.LogStorage) {
	defer wg.Done()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	err := server.ListenUDP("0.0.0.0:" + strconv.Itoa(cfg.SyslogPort))

	if err != nil {
		log.Println(err.Error())
	}
	err = server.ListenTCP("0.0.0.0:" + strconv.Itoa(cfg.SyslogPort))

	if err != nil {
		log.Println(err.Error())
	}

	err = server.Boot()
	if err != nil {
		log.Println(err.Error())
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			(*storage).SaveEvent(store.LogEvent{
				Host:     coalesce(logParts, ParseAsString, "unknown", "hostname").(string),
				AppName:  coalesce(logParts, ParseAsString, "unknown", "app_name", "tag").(string),
				Message:  coalesce(logParts, ParseAsString, "", "message", "content").(string),
				Facility: coalesce(logParts, ParseAsInt, 0, "facility").(int),
				Severity: coalesce(logParts, ParseAsInt, 0, "severity").(int),
				Date:     coalesce(logParts, ParseAsTime, time.Unix(0, 0), "timestamp").(time.Time),
			})
		}
	}(channel)

	server.Wait()
}
