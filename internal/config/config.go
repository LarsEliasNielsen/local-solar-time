package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultPort    = 8080
	defaultCadence = time.Second
)

type Config struct {
	Port      int
	Cadence   time.Duration
	StaticDir string
}

// Load resolves Config from, in order of precedence, command-line flags,
// then SOLAR_* environment variables (auto-loaded from a local .env file
// if present), then defaults.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	portStr := envOrDefault("SOLAR_PORT", strconv.Itoa(defaultPort))
	cadenceStr := envOrDefault("SOLAR_CADENCE", defaultCadence.String())
	staticDirStr := os.Getenv("SOLAR_STATIC_DIR")

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	portFlag := fs.String("port", portStr, "WebSocket listen port")
	cadenceFlag := fs.String("cadence", cadenceStr, "update tick cadence (e.g. 1s, 500ms)")
	staticDirFlag := fs.String("static-dir", staticDirStr, "directory of frontend files to serve; empty serves only /ws")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return Config{}, fmt.Errorf("parse flags: %w", err)
	}

	port, err := strconv.Atoi(*portFlag)
	if err != nil {
		return Config{}, fmt.Errorf("parse port %q: %w", *portFlag, err)
	}

	cadence, err := time.ParseDuration(*cadenceFlag)
	if err != nil {
		return Config{}, fmt.Errorf("parse cadence %q: %w", *cadenceFlag, err)
	}

	return Config{Port: port, Cadence: cadence, StaticDir: *staticDirFlag}, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
