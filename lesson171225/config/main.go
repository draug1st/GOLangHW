package config

import (
	"os"
	"strconv"
)

var Api *Config

type Counter struct {
	Interval int
}

type Config struct {
	Counter *Counter
}

func init() {
	c := os.Getenv("INTERVAL")
	i := 0
	if c != "" {
		if v, err := strconv.Atoi(c); err == nil {
			i = v
		} else {
			i = 1
		}
	}
	Api = &Config{
		Counter: &Counter{
			Interval: i,
		},
	}
}
