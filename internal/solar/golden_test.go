package solar

import (
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"
)

// goldenFixture is one row of testdata/golden.json. Optional fields are
// pointers; a nil field means that assertion is skipped for this row.
type goldenFixture struct {
	Name              string    `json:"name"`
	Source            string    `json:"source"`
	Time              time.Time `json:"time"`
	Lat               float64   `json:"lat"`
	Lon               float64   `json:"lon"`
	AltitudeDeg       *float64  `json:"altitude_deg,omitempty"`
	AltitudeTolDeg    float64   `json:"altitude_tol_deg,omitempty"`
	SunriseSolarTime  *string   `json:"sunrise_solar_time,omitempty"`
	SunsetSolarTime   *string   `json:"sunset_solar_time,omitempty"`
	RiseSetTolMinutes float64   `json:"rise_set_tol_minutes,omitempty"`
	PolarDay          *bool     `json:"polar_day,omitempty"`
	PolarNight        *bool     `json:"polar_night,omitempty"`
}

func loadGoldenFixtures(t *testing.T) []goldenFixture {
	data, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatalf("reading golden fixtures: %v", err)
	}
	var fixtures []goldenFixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatalf("parsing golden fixtures: %v", err)
	}
	return fixtures
}

// parseSolarTime parses an "HH:MM:SS" string into a duration since midnight.
func parseSolarTime(t *testing.T, s string) time.Duration {
	parsed, err := time.Parse("15:04:05", s)
	if err != nil {
		t.Fatalf("parsing solar time %q: %v", s, err)
	}
	return time.Duration(parsed.Hour())*time.Hour +
		time.Duration(parsed.Minute())*time.Minute +
		time.Duration(parsed.Second())*time.Second
}

func TestGoldenFixtures(t *testing.T) {
	for _, fx := range loadGoldenFixtures(t) {
		t.Run(fx.Name, func(t *testing.T) {
			res := Compute(fx.Time, fx.Lat, fx.Lon)

			if fx.AltitudeDeg != nil {
				if diff := math.Abs(res.AltitudeDeg - *fx.AltitudeDeg); diff > fx.AltitudeTolDeg {
					t.Errorf("AltitudeDeg = %v, want within %v of %v (%s)", res.AltitudeDeg, fx.AltitudeTolDeg, *fx.AltitudeDeg, fx.Source)
				}
			}

			if fx.SunriseSolarTime != nil {
				if res.TodaySunrise == nil {
					t.Fatalf("TodaySunrise = nil, want %v (%s)", *fx.SunriseSolarTime, fx.Source)
				}
				want := parseSolarTime(t, *fx.SunriseSolarTime)
				got := parseSolarTime(t, res.TodaySunrise.SolarTime)
				if diff := (got - want).Abs(); diff > time.Duration(fx.RiseSetTolMinutes*float64(time.Minute)) {
					t.Errorf("TodaySunrise.SolarTime = %v, want within %vmin of %v (%s)", res.TodaySunrise.SolarTime, fx.RiseSetTolMinutes, *fx.SunriseSolarTime, fx.Source)
				}
			}

			if fx.SunsetSolarTime != nil {
				if res.TodaySunset == nil {
					t.Fatalf("TodaySunset = nil, want %v (%s)", *fx.SunsetSolarTime, fx.Source)
				}
				want := parseSolarTime(t, *fx.SunsetSolarTime)
				got := parseSolarTime(t, res.TodaySunset.SolarTime)
				if diff := (got - want).Abs(); diff > time.Duration(fx.RiseSetTolMinutes*float64(time.Minute)) {
					t.Errorf("TodaySunset.SolarTime = %v, want within %vmin of %v (%s)", res.TodaySunset.SolarTime, fx.RiseSetTolMinutes, *fx.SunsetSolarTime, fx.Source)
				}
			}

			if fx.PolarDay != nil || fx.PolarNight != nil {
				isDay, isNight := polarPhase(fx.Time, fx.Lat)
				if fx.PolarDay != nil && isDay != *fx.PolarDay {
					t.Errorf("polarPhase isPolarDay = %v, want %v (%s)", isDay, *fx.PolarDay, fx.Source)
				}
				if fx.PolarNight != nil && isNight != *fx.PolarNight {
					t.Errorf("polarPhase isPolarNight = %v, want %v (%s)", isNight, *fx.PolarNight, fx.Source)
				}
			}
		})
	}
}
