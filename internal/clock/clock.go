package clock

import "time"

type Clock interface {
	Now() time.Time
}

type WallClock struct{}

func (WallClock) Now() time.Time {
	panic("not implemented")
}

type FixedClock struct {
	Time time.Time
}

func (f FixedClock) Now() time.Time {
	panic("not implemented")
}
