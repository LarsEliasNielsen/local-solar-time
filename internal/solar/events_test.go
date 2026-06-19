package solar

import (
	"testing"
	"time"
)

func TestPreviousNextEvent(t *testing.T) {
	// Mid-afternoon at a mid-latitude: today's sunrise is in the past,
	// today's sunset is still ahead.
	ref := time.Date(2026, 4, 15, 14, 0, 0, 0, time.UTC)
	lat, lon := 51.5, 0.0

	prevSunrise, nextSunrise := previousNextEvent(ref, lat, lon, true)
	if prevSunrise == nil || nextSunrise == nil {
		t.Fatalf("previousNextEvent(sunrise) = %v, %v, want both non-nil", prevSunrise, nextSunrise)
	}
	if !prevSunrise.UTC.Before(ref) {
		t.Errorf("previous sunrise %v is not before reference %v", prevSunrise.UTC, ref)
	}
	if !nextSunrise.UTC.After(ref) {
		t.Errorf("next sunrise %v is not after reference %v", nextSunrise.UTC, ref)
	}
	if nextSunrise.UTC.Sub(prevSunrise.UTC) > 25*time.Hour {
		t.Errorf("previous/next sunrise %v / %v are more than a day apart", prevSunrise.UTC, nextSunrise.UTC)
	}

	prevSunset, nextSunset := previousNextEvent(ref, lat, lon, false)
	if prevSunset == nil || nextSunset == nil {
		t.Fatalf("previousNextEvent(sunset) = %v, %v, want both non-nil", prevSunset, nextSunset)
	}
	if !prevSunset.UTC.Before(ref) {
		t.Errorf("previous sunset %v is not before reference %v", prevSunset.UTC, ref)
	}
	if !nextSunset.UTC.After(ref) {
		t.Errorf("next sunset %v is not after reference %v", nextSunset.UTC, ref)
	}
}

func TestPreviousNextEventPolarNight(t *testing.T) {
	// Deep into Arctic winter, no sunrise should be found within the
	// search window in either direction.
	ref := time.Date(2026, 12, 21, 12, 0, 0, 0, time.UTC)
	prev, next := previousNextEvent(ref, 78, 0, true)
	if prev != nil || next != nil {
		t.Errorf("previousNextEvent during polar night = %v, %v, want both nil", prev, next)
	}
}

func TestPolarPhase(t *testing.T) {
	tests := []struct {
		name          string
		date          time.Time
		lat           float64
		wantPolarDay  bool
		wantPolarNite bool
	}{
		{
			name:          "mid-latitude has neither",
			date:          time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			lat:           51.5,
			wantPolarDay:  false,
			wantPolarNite: false,
		},
		{
			name:          "Arctic summer is polar day",
			date:          time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC),
			lat:           78,
			wantPolarDay:  true,
			wantPolarNite: false,
		},
		{
			name:          "Arctic winter is polar night",
			date:          time.Date(2026, 12, 21, 0, 0, 0, 0, time.UTC),
			lat:           78,
			wantPolarDay:  false,
			wantPolarNite: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isDay, isNight := polarPhase(tc.date, tc.lat)
			if isDay != tc.wantPolarDay || isNight != tc.wantPolarNite {
				t.Errorf("polarPhase(%v, %v) = (%v, %v), want (%v, %v)",
					tc.date, tc.lat, isDay, isNight, tc.wantPolarDay, tc.wantPolarNite)
			}
		})
	}
}
