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
		name        string
		args        []string
		env         map[string]string
		wantPort    int
		wantCadence time.Duration
	}{
		{
			name:        "defaults",
			wantPort:    8080,
			wantCadence: time.Second,
		},
		{
			name:        "env overrides default",
			env:         map[string]string{"SOLAR_PORT": "9001", "SOLAR_CADENCE": "500ms"},
			wantPort:    9001,
			wantCadence: 500 * time.Millisecond,
		},
		{
			name:        "flag overrides env and default",
			args:        []string{"-port=7000", "-cadence=2s"},
			env:         map[string]string{"SOLAR_PORT": "9001", "SOLAR_CADENCE": "500ms"},
			wantPort:    7000,
			wantCadence: 2 * time.Second,
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
		})
	}
}

func TestLoadInvalidCadence(t *testing.T) {
	withArgs(t, "-cadence=not-a-duration")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an invalid cadence")
	}
}

func TestLoadInvalidPort(t *testing.T) {
	withArgs(t, "-port=not-a-port")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error for an invalid port")
	}
}
