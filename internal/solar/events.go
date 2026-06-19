package solar

import (
	"math"
	"time"
)

// maxSearchDays caps how far previousNextEvent walks in each direction
// before treating the absence of a crossing as polar day or polar night.
const maxSearchDays = 7

// eventOnDate runs riseSetEstimate then refine for the UTC calendar date
// of date, returning ok=false if that date has no crossing.
func eventOnDate(date time.Time, lat, lon float64, sunrise bool) (event *Event, ok bool) {
	estimate, ok := riseSetEstimate(date, lat, lon, sunrise)
	if !ok {
		return nil, false
	}
	refined := refine(estimate, lat, lon, sunrise)
	return &Event{
		SolarTime: formatTimeOfDay(apparentSolarTime(refined, lon)),
		UTC:       refined,
	}, true
}

// previousNextEvent returns the nearest sunrise (or sunset, if sunrise is
// false) before and after t, searching backward and forward by calendar
// day up to maxSearchDays. Either side is nil if no crossing is found
// within that window (polar day or polar night).
func previousNextEvent(t time.Time, lat, lon float64, sunrise bool) (prev, next *Event) {
	date := dateOf(t)

	for i := 0; i <= maxSearchDays; i++ {
		if event, ok := eventOnDate(date.AddDate(0, 0, -i), lat, lon, sunrise); ok && event.UTC.Before(t) {
			prev = event
			break
		}
	}

	for i := 0; i <= maxSearchDays; i++ {
		if event, ok := eventOnDate(date.AddDate(0, 0, i), lat, lon, sunrise); ok && event.UTC.After(t) {
			next = event
			break
		}
	}

	return prev, next
}

// polarPhase reports whether the UTC calendar date of date has no
// sunrise/sunset at all at the given latitude, and if so whether the Sun
// stays above (polar day) or below (polar night) the rise/set threshold
// all day.
func polarPhase(date time.Time, lat float64) (isPolarDay, isPolarNight bool) {
	if _, ok := riseSetEstimate(date, lat, 0, true); ok {
		return false, false
	}

	noon := dateOf(date).Add(12 * time.Hour)
	declRad := declination(noon) * math.Pi / 180
	latRad := lat * math.Pi / 180

	maxAltDeg := math.Asin(math.Sin(latRad)*math.Sin(declRad)+math.Cos(latRad)*math.Cos(declRad)) * 180 / math.Pi
	if maxAltDeg > riseSetThresholdDeg {
		return true, false
	}
	return false, true
}
