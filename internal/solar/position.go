package solar

import (
	"math"
	"time"
)

// altitudeAzimuth returns the Sun's geometric altitude in degrees and its
// azimuth in degrees, measured clockwise from true North (N=0, E=90,
// S=180, W=270), for the given UTC instant and location.
//
// At the exact pole (|lat| == 90) azimuth is undefined: every direction
// points either due north or due south, and the atan2 numerator and
// denominator below are both zero, which would otherwise produce NaN. That
// case is guarded explicitly and reports a nil azimuth.
func altitudeAzimuth(t time.Time, lat, lon float64) (alt float64, az *float64) {
	decl := declination(t) * math.Pi / 180
	latRad := lat * math.Pi / 180

	solarTime := apparentSolarTime(t, lon)
	hourAngleDeg := solarTime.Minutes()/4 - 180
	hourAngle := hourAngleDeg * math.Pi / 180

	sinAlt := math.Sin(latRad)*math.Sin(decl) + math.Cos(latRad)*math.Cos(decl)*math.Cos(hourAngle)
	altRad := math.Asin(sinAlt)
	alt = altRad * 180 / math.Pi

	if math.Abs(lat) == 90 {
		return alt, nil
	}

	cosAlt := math.Cos(altRad)
	sinAz := -math.Sin(hourAngle) * math.Cos(decl) / cosAlt
	cosAz := (math.Sin(decl) - sinAlt*math.Sin(latRad)) / (cosAlt * math.Cos(latRad))
	azRad := math.Atan2(sinAz, cosAz)
	azDeg := azRad * 180 / math.Pi
	if azDeg < 0 {
		azDeg += 360
	}
	return alt, &azDeg
}
