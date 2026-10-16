package config

import (
	"os"
	"testing"
	"time"
)

func withArgs(t *testing.T, args ...string) {
	t.Helper()
	orig := os.Args
	os.Args = append([]string{"local-solar-time"}, args...)
	t.Cleanup(func() { os.Args = orig })
}

func TestLoadPrecedence(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		env            map[string]string
		wantPort       int
		wantCadence    time.Duration
		wantStaticDir  string
		wantMaxSession time.Duration
	}{
		{
			name:           "defaults",
			wantPort:       8080,
			wantCadence:    time.Second,
			wantMaxSession: 2 * time.Hour,
		},
		{
			name:           "env overrides default",
			env:            map[string]string{"SOLAR_PORT": "9001", "SOLAR_CADENCE": "500ms", "SOLAR_STATIC_DIR": "/env", "SOLAR_MAX_SESSION": "30m"},
			wantPort:       9001,
			wantCadence:    500 * time.Millisecond,
			wantStaticDir:  "/env",
			wantMaxSession: 30 * time.Minute,
		},
		{
			name:           "flag overrides env and default",
			args:           []string{"-port=7000", "-cadence=2s", "-static-dir=/flag", "-max-session=5s"},
			env:            map[string]string{"SOLAR_PORT": "9001", "SOLAR_CADENCE": "500ms", "SOLAR_STATIC_DIR": "/env", "SOLAR_MAX_SESSION": "30m"},
			wantPort:       7000,
			wantCadence:    2 * time.Second,
			wantStaticDir:  "/flag",
			wantMaxSession: 5 * time.Second,
		},
		{
			name:           "zero max session disables the limit",
			args:           []string{"-max-session=0"},
			wantPort:       8080,
			wantCadence:    time.Second,
			wantMaxSession: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			withArgs(t, tt.args...)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", cfg.Port, tt.wantPort)
			}
			if cfg.Cadence != tt.wantCadence {
				t.Errorf("Cadence = %v, want %v", cfg.Cadence, tt.wantCadence)
			}
			if cfg.StaticDir != tt.wantStaticDir {
				t.Errorf("StaticDir = %q, want %q", cfg.StaticDir, tt.wantStaticDir)
			}
			if cfg.MaxSession != tt.wantMaxSession {
				t.Errorf("MaxSession = %v, want %v", cfg.MaxSession, tt.wantMaxSession)
			}
		})
	}
}

func TestLoadInvalidCadence(t *testing.T) {
	withArgs(t, "-cadence=not-a-duration")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an invalid cadence")
	}
}

func TestLoadInvalidMaxSession(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "not a duration", arg: "-max-session=forever"},
		{name: "negative", arg: "-max-session=-1m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withArgs(t, tt.arg)

			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil, want an error for %s", tt.arg)
			}
		})
	}
}

func TestLoadInvalidPort(t *testing.T) {
	withArgs(t, "-port=not-a-port")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an invalid port")
	}
}
