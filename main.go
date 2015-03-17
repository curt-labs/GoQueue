package main

import (
	"github.com/bitly/go-nsq"
	"github.com/curt-labs/GoQueue/handlers"

	"log"
	"sync"
)

var (
	NSQDHosts = []string{
		"130.211.131.82:4150",
		"130.211.183.49:4150",
		"130.211.185.177:4150",
		"23.236.48.160:4150",
	}

	ConsumerConcurrency = 100
)

func main() {

	wg := &sync.WaitGroup{}

	config := nsq.NewConfig()
	goapi, err := nsq.NewConsumer("goapi_analytics", "ch", config)
	if err != nil {
		log.Fatal(err.Error())
	}
	v2mock, err := nsq.NewConsumer("v2mock_analytics", "ch", config)
	if err != nil {
		log.Fatal(err.Error())
	}
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

	running := 0
	err = goapi.ConnectToNSQDs(NSQDHosts)
	if err == nil {
		running = running + 1
		wg.Add(1)
	}
	err = v2mock.ConnectToNSQDs(NSQDHosts)
	if err == nil {
		running = running + 1
		wg.Add(1)
	}
	err = admin.ConnectToNSQDs(NSQDHosts)
	if err == nil {
		running = running + 1
		wg.Add(1)
	}

	if running == 0 {
		wg.Done()
		return
	}
	wg.Wait()

}
