package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-yaml/yaml"
)

const (
	configFile       = "./config.yaml"
	defaultLogPeriod = 10 * time.Second
)

// Config is a main app configuration.
type Config struct {
	File   string                 `yaml:"file"`
	Mongo  MongoConf              `yaml:"mongo"`
	Kafka  KafkaConf              `yaml:"kafka"`
	Filter map[string]interface{} `yaml:"filter"`
	Logs   LogsConf               `yaml:"logs"`
}

// MongoConf is a set of mongodb parameters.
type MongoConf struct {
	Addr       string `yaml:"addr"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// KafkaConf is a set of kafka parameters.
type KafkaConf struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	Offset  int64    `yaml:"offset"`
	GroupID string   `yaml:"group_id"`
}

// LogsConf is a logging configuration.
type LogsConf struct {
	Level  string        `yaml:"level"`
	Period time.Duration `yaml:"period"`
}

// ReadConfig reads configuration from a file.
func ReadConfig() (Config, error) {
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("read file: %v", err)
	}
	var conf Config
	if err := yaml.Unmarshal(f, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal yaml: %v", err)
	}
	if conf.File != "" && conf.Mongo.Addr != "" {
		return Config{}, fmt.Errorf("only one storage should be specified: file or mongodb")
	}
	if conf.File == "" && conf.Mongo.Addr == "" {
		return Config{}, fmt.Errorf("no storage specified")
	}
	if conf.Logs.Period == 0 {
		conf.Logs.Period = defaultLogPeriod
	}
	return conf, nil
}
