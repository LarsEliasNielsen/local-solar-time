package solar

import (
	"fmt"
	"time"
)

// apparentSolarTime returns the time of day, since midnight UTC, that a
// sundial would read at the given longitude for the given UTC instant. It
// combines the equation of time with the standard 4-minutes-per-degree
// longitude correction relative to the UTC (0 degrees) meridian.
func apparentSolarTime(t time.Time, lon float64) time.Duration {
	t = t.UTC()
	utcTimeOfDay := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())

	offsetMinutes := lon*4 + equationOfTimeMinutes(t)
	offset := time.Duration(offsetMinutes * float64(time.Minute))

	const day = 24 * time.Hour
	solar := (utcTimeOfDay + offset) % day
	if solar < 0 {
		solar += day
	}
	return solar
}

// formatTimeOfDay formats d, a duration since midnight, as HH:MM:SS.
func formatTimeOfDay(d time.Duration) string {
	const day = 24 * time.Hour
	d = ((d % day) + day) % day
	s := int(d.Round(time.Second).Seconds())
	return fmt.Sprintf("%02d:%02d:%02d", s/3600, (s%3600)/60, s%60)
}
