//go:build meeus

// This file validates same-lineage agreement with an independent Go
// implementation of the same NOAA/Meeus algorithm family. It is not run by
// default; absolute accuracy is covered by golden_test.go against
// independently published constants. Run explicitly with:
//
//	go test -tags meeus ./internal/solar/...
package solar

import (
	"math"
	"testing"
	"time"

	"github.com/soniakeys/meeus/v3/julian"
	"github.com/soniakeys/meeus/v3/solar"
)

func TestDeclinationAgreesWithMeeus(t *testing.T) {
	tests := []time.Time{
		time.Date(2026, 3, 20, 13, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 21, 4, 0, 0, 0, time.UTC),
		time.Date(2026, 9, 22, 22, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 21, 21, 0, 0, 0, time.UTC),
	}

	for _, tm := range tests {
		t.Run(tm.Format(time.RFC3339), func(t *testing.T) {
			got := declination(tm)

			jde := julian.TimeToJD(tm)
			_, declRad := solar.ApparentEquatorial(jde)
			want := declRad.Deg()

			const tol = 0.06 // degrees; Spencer (1971)'s own documented accuracy bound (~0.034deg), worst case near the equinoxes
			if diff := math.Abs(got - want); diff > tol {
				t.Errorf("declination(%v) = %v, meeus = %v, diff %v exceeds tolerance %v", tm, got, want, diff, tol)
			}
		})
	}
}
