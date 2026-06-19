package solar

import (
	"testing"
	"time"
)

func TestRiseSetEstimate(t *testing.T) {
	tests := []struct {
		name    string
		date    time.Time
		lat     float64
		lon     float64
		sunrise bool
		wantOK  bool
	}{
		{
			name:    "mid-latitude sunrise has a solution",
			date:    time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			lat:     51.5,
			lon:     0,
			sunrise: true,
			wantOK:  true,
		},
		{
			name:    "mid-latitude sunset has a solution",
			date:    time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			lat:     51.5,
			lon:     0,
			sunrise: false,
			wantOK:  true,
		},
		{
			name:    "Arctic summer has no sunset (polar day)",
			date:    time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC),
			lat:     78,
			lon:     0,
			sunrise: false,
			wantOK:  false,
		},
		{
			name:    "Arctic winter has no sunrise (polar night)",
			date:    time.Date(2026, 12, 21, 0, 0, 0, 0, time.UTC),
			lat:     78,
			lon:     0,
			sunrise: true,
			wantOK:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := riseSetEstimate(tc.date, tc.lat, tc.lon, tc.sunrise)
			if ok != tc.wantOK {
				t.Fatalf("riseSetEstimate(...) ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got.IsZero() {
				t.Errorf("riseSetEstimate(...) returned zero time with ok=true")
			}
		})
	}
}

func TestRefine(t *testing.T) {
	date := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)
	lat, lon := 51.5, 0.0

	estimate, ok := riseSetEstimate(date, lat, lon, true)
	if !ok {
		t.Fatal("riseSetEstimate returned ok=false, want true")
	}

	refined := refine(estimate, lat, lon, true)

	alt, _ := altitudeAzimuth(refined, lat, lon)
	const tol = 0.05 // degrees; well within the +-1 minute time accuracy target
	if diff := alt - riseSetThresholdDeg; diff > tol || diff < -tol {
		t.Errorf("altitude at refined sunrise = %v, want within %v of threshold %v", alt, tol, riseSetThresholdDeg)
	}
}
