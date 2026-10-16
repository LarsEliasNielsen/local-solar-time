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
	defaultPort       = 8080
	defaultCadence    = time.Second
	defaultMaxSession = 2 * time.Hour
)

type Config struct {
	Port       int
	Cadence    time.Duration
	StaticDir  string
	MaxSession time.Duration
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
	maxSessionStr := envOrDefault("SOLAR_MAX_SESSION", defaultMaxSession.String())

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	portFlag := fs.String("port", portStr, "WebSocket listen port")
	cadenceFlag := fs.String("cadence", cadenceStr, "update tick cadence (e.g. 1s, 500ms)")
	maxSessionFlag := fs.String("max-session", maxSessionStr, "maximum WebSocket session length (e.g. 2h); 0 disables")
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

	maxSession, err := time.ParseDuration(*maxSessionFlag)
	if err != nil {
		return Config{}, fmt.Errorf("parse max session %q: %w", *maxSessionFlag, err)
	}
	if maxSession < 0 {
		return Config{}, fmt.Errorf("max session %q must not be negative", *maxSessionFlag)
	}

	return Config{Port: port, Cadence: cadence, StaticDir: *staticDirFlag, MaxSession: maxSession}, nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
