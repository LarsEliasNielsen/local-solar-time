package config

import "time"

type Config struct {
	Listen  string
	Cadence time.Duration
}

func Load() (Config, error) {
	panic("not implemented")
}
