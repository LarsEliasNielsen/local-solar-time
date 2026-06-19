package solar

import (
	"math"
	"time"
)

// riseSetThresholdDeg is the standard refraction + solar upper-limb
// correction applied to rise/set events only; all other quantities use
// the true geometric Sun.
const riseSetThresholdDeg = -0.833

// riseSetEstimate returns a first-pass estimate of sunrise (or sunset, if
// sunrise is false) on the UTC calendar date of date, using the standard
// hour-angle formula against riseSetThresholdDeg. ok is false when the
// date has no such crossing (polar day or polar night).
func riseSetEstimate(date time.Time, lat, lon float64, sunrise bool) (t time.Time, ok bool) {
	noon := solarNoon(date, lon)
	declRad := declination(noon) * math.Pi / 180
	latRad := lat * math.Pi / 180
	thresholdRad := riseSetThresholdDeg * math.Pi / 180

	cosH0 := (math.Sin(thresholdRad) - math.Sin(latRad)*math.Sin(declRad)) / (math.Cos(latRad) * math.Cos(declRad))
	if cosH0 < -1 || cosH0 > 1 {
		return time.Time{}, false
	}

	hourAngleHours := math.Acos(cosH0) * 180 / math.Pi / 15
	if sunrise {
		return noon.Add(-time.Duration(hourAngleHours * float64(time.Hour))), true
	}
	return noon.Add(time.Duration(hourAngleHours * float64(time.Hour))), true
}

// refine applies one Newton-Raphson iteration to estimate, the result of
// riseSetEstimate, against the same riseSetThresholdDeg crossing.
func refine(estimate time.Time, lat, lon float64, sunrise bool) time.Time {
	const step = time.Minute

	altitudeAt := func(t time.Time) float64 {
		alt, _ := altitudeAzimuth(t, lat, lon)
		return alt
	}

	f := altitudeAt(estimate) - riseSetThresholdDeg
	derivative := (altitudeAt(estimate.Add(step)) - altitudeAt(estimate.Add(-step))) / (2 * step.Seconds())
	if derivative == 0 {
		return estimate
	}

	deltaSeconds := -f / derivative
	return estimate.Add(time.Duration(deltaSeconds * float64(time.Second)))
}
