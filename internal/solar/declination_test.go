package solar

import (
	"math"
	"testing"
	"time"
)

func TestDeclination(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want float64
		tol  float64
	}{
		{
			name: "March equinox",
			time: time.Date(2026, 3, 20, 13, 0, 0, 0, time.UTC),
			want: 0,
			tol:  1.0,
		},
		{
			name: "June solstice",
			time: time.Date(2026, 6, 21, 4, 0, 0, 0, time.UTC),
			want: 23.44,
			tol:  0.5,
		},
		{
			name: "December solstice",
			time: time.Date(2026, 12, 21, 21, 0, 0, 0, time.UTC),
			want: -23.44,
			tol:  0.5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := declination(tc.time)
			if math.Abs(got-tc.want) > tc.tol {
				t.Errorf("declination(%v) = %v, want within %v of %v", tc.time, got, tc.tol, tc.want)
			}
		})
	}
}

func TestEquationOfTimeMinutes(t *testing.T) {
	// The equation of time stays within roughly +-17 minutes year-round.
	tests := []struct {
		name string
		time time.Time
	}{
		{name: "January", time: time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)},
		{name: "April", time: time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)},
		{name: "July", time: time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)},
		{name: "November", time: time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := equationOfTimeMinutes(tc.time)
			if math.Abs(got) > 17 {
				t.Errorf("equationOfTimeMinutes(%v) = %v, want within +-17 minutes", tc.time, got)
			}
		})
	}
}
