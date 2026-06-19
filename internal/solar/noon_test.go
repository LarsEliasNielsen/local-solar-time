package solar

import (
	"testing"
	"time"
)

func TestSolarNoon(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		lon  float64
		want time.Time
		tol  time.Duration
	}{
		{
			name: "Greenwich meridian",
			date: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			lon:  0,
			want: time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC),
			tol:  time.Minute,
		},
		{
			name: "45 degrees east is three hours before UTC noon",
			date: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			lon:  45,
			want: time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC),
			tol:  time.Minute,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := solarNoon(tc.date, tc.lon)
			diff := got.Sub(tc.want)
			if diff < 0 {
				diff = -diff
			}
			if diff > tc.tol {
				t.Errorf("solarNoon(%v, %v) = %v, want within %v of %v", tc.date, tc.lon, got, tc.tol, tc.want)
			}

			// By definition, apparent solar time at solar noon is 12:00:00,
			// within the package's ~1-2s solar time-of-day accuracy budget.
			solarTOD := apparentSolarTime(got, tc.lon)
			noonDiff := solarTOD - 12*time.Hour
			if noonDiff < 0 {
				noonDiff = -noonDiff
			}
			if noonDiff > 2*time.Second {
				t.Errorf("apparent solar time at solar noon = %v, want within 2s of 12:00:00", formatTimeOfDay(solarTOD))
			}
		})
	}
}
