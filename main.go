package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	log.SetLevel(log.InfoLevel)
	log.Info("Starting...")

	conf, err := ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	level, err := log.ParseLevel(conf.Logs.Level)
	if err != nil {
		log.Errorf("Invalid log level '%s', using default: info", conf.Logs.Level)
	}
	log.SetLevel(level)

	// Init kafka consumer
	c, err := NewKafkaConsumer(conf.Kafka)
	if err != nil {
		log.Fatalf("Failed to init consumer: %v", err)
	}

	// Init messages filter
	f := NewFieldFilter(conf.Filter)

	// Init storage
	var s Storage
	if conf.File != "" {
		log.Infof("Saving messages to %s", conf.File)
		s, err = NewFileSystemStorage(conf.File)
		if err != nil {
			log.Fatalf("Failed to init file storage: %v", err)
		}
	} else {
		log.Infof("Saving messages to %s", conf.Mongo.Addr)
		s, err = NewMongoStorage(conf.Mongo)
		if err != nil {
			log.Fatalf("Failed to init mongodb storage: %v", err)
		}
	}

	// Init pipeline
	dmp := NewDumper(c, f, s, conf.Logs.Period)

	// Listen for SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go waitForStop(cancel)

	// Run pipeline
	if err := dmp.Run(ctx); err != nil {
		log.Fatalf("Dump process failed: %v", err)
	}
	log.Info("Shutdown")
}

func waitForStop(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info("Got termination signal")
	cancel()
}
