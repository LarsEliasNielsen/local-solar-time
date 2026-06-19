package clock

import (
	"testing"
	"time"
)

var (
	_ Clock = WallClock{}
	_ Clock = FixedClock{}
)

func TestWallClockNowIsUTC(t *testing.T) {
	got := WallClock{}.Now()
	if got.Location() != time.UTC {
		t.Errorf("WallClock.Now() location = %v, want UTC", got.Location())
	}
}

func TestFixedClockNowReturnsConfiguredTime(t *testing.T) {
	want := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	got := FixedClock{Time: want}.Now()
	if !got.Equal(want) {
		t.Errorf("FixedClock.Now() = %v, want %v", got, want)
	}
}
