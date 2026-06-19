package clock

import "time"

// Clock provides the current time, so callers never depend on time.Now directly.
type Clock interface {
	Now() time.Time
}

// WallClock is the production Clock, backed by the OS wall clock.
type WallClock struct{}

func (WallClock) Now() time.Time {
	return time.Now().UTC()
}

// FixedClock is a test Clock that always returns Time.
type FixedClock struct {
	Time time.Time
}

func (f FixedClock) Now() time.Time {
	return f.Time
}
