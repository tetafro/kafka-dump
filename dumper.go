package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Dumper is a main app's entity. It run read-filter-save loop.
type Dumper struct {
	consumer  Consumer
	filter    Filter
	storage   Storage
	logPeriod time.Duration
}

// NewDumper creates new dumper.
func NewDumper(c Consumer, f Filter, s Storage, p time.Duration) *Dumper {
	return &Dumper{consumer: c, filter: f, storage: s, logPeriod: p}
}

// Run starts main read-filter-save loop and logs current state.
func (d *Dumper) Run(ctx context.Context) error {
	defer d.consumer.Close()
	defer d.storage.Close()

	var total, saved int
	var firstMsg, lastMsg time.Time
	lastLog := time.Now()

	logStats := func() {
		log.Infof(
			"Read messages from %s to %s (total %d, saved %d)",
			firstMsg.Local().Format("2006-01-02 15:04:05"),
			lastMsg.Local().Format("2006-01-02 15:04:05"),
			total, saved,
		)
		firstMsg = time.Time{}
		lastMsg = time.Time{}
		lastLog = time.Now()
	}
	defer logStats()

	for {
		if ctx.Err() != nil {
			return nil
		}

		if time.Since(lastLog) > d.logPeriod {
			logStats()
		}

		msg, err := d.consumer.Read(ctx)
		if err == context.Canceled {
			return nil
		}
		if err != nil {
			log.Errorf("Failed to read message: %v", err)
			continue
		}
		total++

		if msg.time.Before(firstMsg) || firstMsg.IsZero() {
			firstMsg = msg.time
		}
		if msg.time.After(lastMsg) || lastMsg.IsZero() {
			lastMsg = msg.time
		}

		if !d.filter.Check(msg) {
			continue
		}

		if err := d.storage.Save(msg); err != nil {
			return fmt.Errorf("save message: %v", err)
		}
		saved++
	}
}
