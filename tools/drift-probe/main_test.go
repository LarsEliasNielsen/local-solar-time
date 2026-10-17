package main

import (
	"testing"
	"time"
)

func TestOffset(t *testing.T) {
	recv := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		utc       time.Time
		rtt       time.Duration
		ntpOffset time.Duration
		want      time.Duration
	}{
		{"in sync", recv.Add(-20 * time.Millisecond), 40 * time.Millisecond, 0, 0},
		{"server ahead", recv.Add(30 * time.Millisecond), 40 * time.Millisecond, 0, 50 * time.Millisecond},
		{"local clock slow", recv.Add(80 * time.Millisecond), 40 * time.Millisecond, 100 * time.Millisecond, 0},
		{"server behind", recv.Add(-70 * time.Millisecond), 0, 0, -70 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := offset(tt.utc, recv, tt.rtt, tt.ntpOffset); got != tt.want {
				t.Errorf("offset = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedianAndMaxAbs(t *testing.T) {
	tests := []struct {
		name       string
		in         []time.Duration
		wantMedian time.Duration
		wantMax    time.Duration
	}{
		{"empty", nil, 0, 0},
		{"odd", []time.Duration{3, -9, 1}, 3, 9},
		{"even", []time.Duration{4, -2, 8, 6}, 5, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := median(absAll(tt.in)); got != tt.wantMedian {
				t.Errorf("median = %v, want %v", got, tt.wantMedian)
			}
			if got := maxAbs(tt.in); got != tt.wantMax {
				t.Errorf("maxAbs = %v, want %v", got, tt.wantMax)
			}
		})
	}
}
