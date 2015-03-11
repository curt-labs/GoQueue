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

	goapiHandler := &handlers.ConsumerHandler{
		Category:     "GoAPI",
		TrackingCode: "UA-59297117-1",
	}
	v2MockHandler := &handlers.ConsumerHandler{
		Category:     "v2Mock",
		TrackingCode: "UA-59297117-1",
	}

	goapi.AddConcurrentHandlers(goapiHandler, ConsumerConcurrency)
	v2mock.AddConcurrentHandlers(v2MockHandler, ConsumerConcurrency)

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

	if running == 0 {
		wg.Done()
		return
	}
	wg.Wait()

}
