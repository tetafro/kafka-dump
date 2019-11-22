package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type storage struct {
	storage   io.WriteCloser
	logPeriod time.Duration
	filter    map[string]interface{}
	input     chan message
	counter   counter
}

type counter struct {
	total int
	saved int
}

func newStorage(
	file string,
	p time.Duration,
	filter map[string]interface{},
	in chan message,
) (*storage, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %v", file, err)
	}
	st := &storage{
		storage:   f,
		logPeriod: p,
		filter:    filter,
		input:     in,
	}
	return st, nil
}

func (s *storage) run() error {
	var lastMsg time.Time
	var lastLog time.Time
	for msg := range s.input {
		s.counter.total++

		if msg.time.After(lastMsg) {
			lastMsg = msg.time
		}
		if time.Since(lastLog) > s.logPeriod {
			if !lastLog.IsZero() { // skip first time
				log.Infof(
					"Read all messages until %s (total %d, saved %d)",
					lastMsg.Format("2006-01-02 15:04:05"),
					s.counter.total,
					s.counter.saved,
				)
			}
			lastLog = time.Now()
		}

		if !s.check(msg.data) {
			continue
		}

		if err := s.save(msg.data); err != nil {
			return fmt.Errorf("save message: %v", err)
		}
		s.counter.saved++
	}
	return nil
}

func (s *storage) check(m map[string]interface{}) bool {
	for k, v := range s.filter {
		data, ok := m[k]
		if !ok || !equal(v, data) {
			return false
		}
	}
	return true
}

func (s *storage) save(m map[string]interface{}) error {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}
	_, err = s.storage.Write(append(data, []byte("\n")...))
	return err
}
