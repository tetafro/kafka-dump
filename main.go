package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const configFile = "./config.yaml"

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	log.SetLevel(log.InfoLevel)
	log.Info("Starting...")

	conf, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	level, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Errorf("Invalid log level '%s', using default: info", conf.LogLevel)
	}
	log.SetLevel(level)

	ctx, cancel := context.WithCancel(context.Background())
	go waitForStop(cancel)

	c, err := newConsumer(conf.Brokers, conf.Topic, conf.GroupID, conf.Offset)
	if err != nil {
		log.Fatalf("Failed to init consumer: %v", err)
	}
	st, err := newStorage(conf.File, conf.LogPeriod, conf.Filter, c.messages)
	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Debug("Start storage")
		if err := st.run(); err != nil {
			log.Fatalf("Storage failure: %v", err)
		}
		wg.Done()
	}()

	log.Debug("Start consumer")
	if err := c.run(ctx); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
	log.Debug("Consumer is done")

	wg.Wait()
	log.Info("Shutdown")
}

func waitForStop(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info("Got termination signal")
	cancel()
}
