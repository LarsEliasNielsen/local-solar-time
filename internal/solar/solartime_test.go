package solar

import (
	"testing"
	"time"
)

func TestApparentSolarTime(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		lon  float64
		want time.Duration
		tol  time.Duration
	}{
		{
			name: "Greenwich meridian, equation of time near zero",
			time: time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC),
			lon:  0,
			want: 12 * time.Hour,
			tol:  time.Minute,
		},
		{
			name: "15 degrees east is one hour ahead of UTC",
			time: time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC),
			lon:  15,
			want: 13 * time.Hour,
			tol:  time.Minute,
		},
		{
			name: "15 degrees west is one hour behind UTC",
			time: time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC),
			lon:  -15,
			want: 11 * time.Hour,
			tol:  time.Minute,
		},
		{
			name: "wraps past midnight",
			time: time.Date(2026, 4, 15, 23, 0, 0, 0, time.UTC),
			lon:  90,
			want: 5 * time.Hour,
			tol:  time.Minute,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := apparentSolarTime(tc.time, tc.lon)
			diff := got - tc.want
			if diff < 0 {
				diff = -diff
			}
			if diff > tc.tol {
				t.Errorf("apparentSolarTime(%v, %v) = %v, want within %v of %v", tc.time, tc.lon, got, tc.tol, tc.want)
			}
		})
	}
}
