package solar

import (
	"math"
	"time"
)

// fractionalYear returns the Spencer (1971) fractional-year angle gamma, in
// radians, used by both declination and the equation of time. It depends on
// the day of year and the fraction of the day elapsed, per the NOAA solar
// calculator's low-order series.
func fractionalYear(t time.Time) float64 {
	t = t.UTC()
	dayOfYear := t.YearDay()
	hour := float64(t.Hour()) + float64(t.Minute())/60 + float64(t.Second())/3600
	daysInYear := 365.0
	if isLeapYear(t.Year()) {
		daysInYear = 366.0
	}
	return 2 * math.Pi / daysInYear * (float64(dayOfYear) + (hour-12)/24)
}

// isLeapYear reports whether year is a leap year in the Gregorian calendar.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// declination returns the Sun's geometric declination in degrees for the
// given UTC instant, using the Spencer (1971) low-order Fourier series as
// reproduced by the NOAA solar calculator. Accurate to within ~0.05 degrees.
func declination(t time.Time) float64 {
	gamma := fractionalYear(t)
	decl := 0.006918 -
		0.399912*math.Cos(gamma) +
		0.070257*math.Sin(gamma) -
		0.006758*math.Cos(2*gamma) +
		0.000907*math.Sin(2*gamma) -
		0.002697*math.Cos(3*gamma) +
		0.00148*math.Sin(3*gamma)
	return decl * 180 / math.Pi
}

// equationOfTimeMinutes returns the equation of time in minutes for the
// given UTC instant, using the same Spencer (1971) series family as
// declination.
func equationOfTimeMinutes(t time.Time) float64 {
	gamma := fractionalYear(t)
	return 229.18 * (0.000075 +
		0.001868*math.Cos(gamma) -
		0.032077*math.Sin(gamma) -
		0.014615*math.Cos(2*gamma) -
		0.040849*math.Sin(2*gamma))
}
