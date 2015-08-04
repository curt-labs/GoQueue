package main

import (
	"github.com/bitly/go-nsq"
	"github.com/curt-labs/GoQueue/handlers"

	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	NSQDHosts = []string{
		"146.148.64.5:4150",
		"104.197.73.14:4150",
		"104.197.78.177:4150",
		"104.154.51.41:4150",
		"162.222.182.178:4150",
		// "127.0.0.1:4150",
	}

	ConsumerConcurrency = 100
)

func main() {
	log.Print("running")
	config := nsq.NewConfig()
	goapi, err := nsq.NewConsumer("goapi_analytics", "ch", config)
	if err != nil {
		log.Fatal(err.Error())
	}
	v2mock, err := nsq.NewConsumer("v2mock_analytics", "ch", config)
	if err != nil {
		log.Fatal(err.Error())
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	admin, err := nsq.NewConsumer("admin_change", "ch", config)

	goapiHandler := &handlers.AnalyticsHandler{
		Category:     "GoAPI",
		TrackingCode: "UA-59297117-1",
	}
	v2MockHandler := &handlers.AnalyticsHandler{
		Category:     "v2Mock",
		TrackingCode: "UA-59297117-1",
	}
	adminHandler := &handlers.AdminHandler{}

	goapi.AddConcurrentHandlers(goapiHandler, ConsumerConcurrency)
	v2mock.AddConcurrentHandlers(v2MockHandler, ConsumerConcurrency)
	admin.AddConcurrentHandlers(adminHandler, ConsumerConcurrency)

	err = goapi.ConnectToNSQDs(NSQDHosts)
	if err != nil {
		log.Fatal(err)
	}
	err = v2mock.ConnectToNSQDs(NSQDHosts)
	if err != nil {
		log.Fatal(err)
	}
	err = admin.ConnectToNSQDs(NSQDHosts)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-admin.StopChan:
			return
		case <-v2mock.StopChan:
			return
		case <-goapi.StopChan:
			return
		case <-termChan:
			admin.Stop()
			v2mock.Stop()
			goapi.Stop()
		}
	}
}
