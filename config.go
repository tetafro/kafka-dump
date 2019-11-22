package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-yaml/yaml"
)

const defaultLogPeriod = 10 * time.Second

type config struct {
	File      string                 `yaml:"file"`
	LogLevel  string                 `yaml:"logLevel"`
	LogPeriod time.Duration          `yaml:"logPeriod"`
	Brokers   []string               `yaml:"brokers"`
	Topic     string                 `yaml:"topic"`
	Offset    int64                  `yaml:"offset"`
	GroupID   string                 `yaml:"groupID"`
	Filter    map[string]interface{} `yaml:"filter"`
}

func readConfig(file string) (config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return config{}, fmt.Errorf("read file: %v", err)
	}
	var conf config
	if err := yaml.Unmarshal(f, &conf); err != nil {
		return config{}, fmt.Errorf("unmarshal yaml: %v", err)
	}
	if conf.LogPeriod == 0 {
		conf.LogPeriod = defaultLogPeriod
	}
	return conf, nil
}
