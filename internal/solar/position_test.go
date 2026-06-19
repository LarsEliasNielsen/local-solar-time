package solar

import (
	"math"
	"testing"
	"time"
)

func TestAltitudeAzimuth(t *testing.T) {
	tests := []struct {
		name    string
		time    time.Time
		lat     float64
		lon     float64
		wantAlt float64
		altTol  float64
		wantAz  *float64
		azTol   float64
	}{
		{
			// Near the June solstice, solar noon at the equator: the Sun
			// sits a bit north of the zenith and roughly due north.
			name:    "equator near June solstice solar noon",
			time:    time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC),
			lat:     0,
			lon:     0,
			wantAlt: 66.5,
			altTol:  1.5,
			wantAz:  ptr(0),
			azTol:   15,
		},
		{
			name:    "midnight altitude is negative",
			time:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			lat:     50,
			lon:     0,
			wantAlt: -60,
			altTol:  20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			alt, az := altitudeAzimuth(tc.time, tc.lat, tc.lon)
			if math.Abs(alt-tc.wantAlt) > tc.altTol {
				t.Errorf("altitude = %v, want within %v of %v", alt, tc.altTol, tc.wantAlt)
			}
			if tc.wantAz != nil {
				if az == nil {
					t.Fatalf("azimuth = nil, want ~%v", *tc.wantAz)
				}
				diff := math.Abs(*az - *tc.wantAz)
				if diff > 180 {
					diff = 360 - diff
				}
				if diff > tc.azTol {
					t.Errorf("azimuth = %v, want within %v of %v", *az, tc.azTol, *tc.wantAz)
				}
			}
		})
	}
}

func TestAltitudeAzimuthPoleGuard(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
	}{
		{name: "north pole", lat: 90},
		{name: "south pole", lat: -90},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tm := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
			alt, az := altitudeAzimuth(tm, tc.lat, 0)
			if az != nil {
				t.Errorf("azimuth = %v, want nil at exact pole", *az)
			}
			if math.IsNaN(alt) {
				t.Errorf("altitude is NaN at exact pole")
			}
		})
	}
}

func ptr(f float64) *float64 {
	return &f
}
