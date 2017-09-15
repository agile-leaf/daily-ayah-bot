package main

import (
	"github.com/jinzhu/configor"
)

type config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func loadConfig() *config {
	cfg := &config{}
	configor.Load(cfg, "./config.yaml")

	return cfg
}
