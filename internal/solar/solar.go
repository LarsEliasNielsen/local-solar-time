// Package solar is a pure, deterministic solar-position engine: no I/O,
// no imports outside the standard library.
package solar

import (
	"math"
	"time"
)

// polarCapLatitude is the threshold above which solar-time accuracy is
// not maintained; only AltitudeDeg and EquationOfTimeMinutes stay valid.
const polarCapLatitude = 89.4

const polarCapReason = "latitude exceeds solar-time accuracy threshold (~89.4°)"

// Event is a solar event such as solar noon, sunrise, or sunset, reported
// in both apparent solar time and UTC.
type Event struct {
	SolarTime string
	UTC       time.Time
}

// Result is the solar payload for one UTC instant and location. Above the
// polar cap latitude, only AltitudeDeg and EquationOfTimeMinutes are
// valid; every other field is its zero value and PolarCapReason explains
// why.
type Result struct {
	SolarTime             *string
	EquationOfTimeMinutes float64
	UTCOffsetSeconds      *int
	AltitudeDeg           float64
	AzimuthDeg            *float64 // nil at the exact pole and above the polar cap.
	SolarNoon             *Event
	TodaySunrise          *Event // nil during polar day or polar night.
	TodaySunset           *Event // nil during polar day or polar night.
	PreviousSunrise       *Event
	NextSunrise           *Event
	PreviousSunset        *Event
	NextSunset            *Event
	PolarCapReason        string // non-empty only above ~89.4 degrees latitude.
}

// Compute returns the solar payload for the UTC instant t at the given
// latitude and longitude in degrees. It never clamps an out-of-range
// latitude; validating |lat| <= 90 is the caller's responsibility.
func Compute(t time.Time, lat, lon float64) Result {
	eot := equationOfTimeMinutes(t)
	alt, az := altitudeAzimuth(t, lat, lon)

	if math.Abs(lat) > polarCapLatitude {
		return Result{
			AltitudeDeg:           alt,
			EquationOfTimeMinutes: eot,
			PolarCapReason:        polarCapReason,
		}
	}

	solarTime := formatTimeOfDay(apparentSolarTime(t, lon))
	utcOffsetSeconds := int(math.Round(lon*240 + eot*60))
	noonEvent := &Event{SolarTime: "12:00:00", UTC: solarNoon(t, lon)}

	var todaySunrise, todaySunset *Event
	if isPolarDay, isPolarNight := polarPhase(t, lat); !isPolarDay && !isPolarNight {
		todaySunrise, _ = eventOnDate(t, lat, lon, true)
		todaySunset, _ = eventOnDate(t, lat, lon, false)
	}

	prevSunrise, nextSunrise := previousNextEvent(t, lat, lon, true)
	prevSunset, nextSunset := previousNextEvent(t, lat, lon, false)

	return Result{
		SolarTime:             &solarTime,
		EquationOfTimeMinutes: eot,
		UTCOffsetSeconds:      &utcOffsetSeconds,
		AltitudeDeg:           alt,
		AzimuthDeg:            az,
		SolarNoon:             noonEvent,
		TodaySunrise:          todaySunrise,
		TodaySunset:           todaySunset,
		PreviousSunrise:       prevSunrise,
		NextSunrise:           nextSunrise,
		PreviousSunset:        prevSunset,
		NextSunset:            nextSunset,
	}
}
