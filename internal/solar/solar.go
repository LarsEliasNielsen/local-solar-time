package solar

import "time"

type Event struct {
	SolarTime string
	UTC       time.Time
}

type Result struct {
	SolarTime             *string
	EquationOfTimeMinutes float64
	UTCOffsetSeconds      *int
	AltitudeDeg           float64
	AzimuthDeg            *float64
	SolarNoon             *Event
	TodaySunrise          *Event
	TodaySunset           *Event
	PreviousSunrise       *Event
	NextSunrise           *Event
	PreviousSunset        *Event
	NextSunset            *Event
	PolarCapReason        string
}

func Compute(t time.Time, lat, lon float64) Result {
	panic("not implemented")
}
