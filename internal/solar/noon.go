package solar

import "time"

// dateOf truncates t to midnight UTC on its calendar date.
func dateOf(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// solarNoon returns the UTC instant of the Sun's meridian crossing for the
// UTC calendar date of date at the given longitude. Apparent solar time at
// that instant is always exactly 12:00:00 by definition.
func solarNoon(date time.Time, lon float64) time.Time {
	midnight := dateOf(date)
	noonApprox := midnight.Add(12 * time.Hour)
	eot := equationOfTimeMinutes(noonApprox)
	offset := time.Duration((lon*4 + eot) * float64(time.Minute))
	return midnight.Add(12*time.Hour - offset)
}
