package solar

import (
	"math"
	"testing"
	"time"
)

func TestComputeMidLatitude(t *testing.T) {
	res := Compute(time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC), 51.5, 0)

	if res.SolarTime == nil {
		t.Fatal("SolarTime = nil, want a value at mid-latitude")
	}
	if res.AzimuthDeg == nil {
		t.Fatal("AzimuthDeg = nil, want a value at mid-latitude")
	}
	if res.UTCOffsetSeconds == nil {
		t.Fatal("UTCOffsetSeconds = nil, want a value at mid-latitude")
	}
	if res.SolarNoon == nil || res.SolarNoon.SolarTime != "12:00:00" {
		t.Fatalf("SolarNoon = %v, want solar time 12:00:00", res.SolarNoon)
	}
	if res.TodaySunrise == nil || res.TodaySunset == nil {
		t.Fatal("TodaySunrise/TodaySunset = nil, want both populated at mid-latitude in April")
	}
	if res.PreviousSunrise == nil || res.NextSunrise == nil || res.PreviousSunset == nil || res.NextSunset == nil {
		t.Fatal("previous/next sunrise/sunset = nil, want all populated at mid-latitude")
	}
	if res.PolarCapReason != "" {
		t.Errorf("PolarCapReason = %q, want empty below the polar cap", res.PolarCapReason)
	}
}

func TestComputePolarCap(t *testing.T) {
	res := Compute(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC), 89.9, 0)

	if res.PolarCapReason == "" {
		t.Fatal("PolarCapReason = \"\", want a reason above the polar cap")
	}
	if res.SolarTime != nil {
		t.Errorf("SolarTime = %v, want nil above the polar cap", *res.SolarTime)
	}
	if res.AzimuthDeg != nil {
		t.Errorf("AzimuthDeg = %v, want nil above the polar cap", *res.AzimuthDeg)
	}
	if res.UTCOffsetSeconds != nil {
		t.Errorf("UTCOffsetSeconds = %v, want nil above the polar cap", *res.UTCOffsetSeconds)
	}
	if res.SolarNoon != nil {
		t.Errorf("SolarNoon = %v, want nil above the polar cap", res.SolarNoon)
	}
	if res.TodaySunrise != nil || res.TodaySunset != nil {
		t.Error("TodaySunrise/TodaySunset != nil, want nil above the polar cap")
	}
	if math.IsNaN(res.AltitudeDeg) {
		t.Error("AltitudeDeg is NaN above the polar cap")
	}
}

func TestComputeExactPole(t *testing.T) {
	res := Compute(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC), 90, 0)

	if math.IsNaN(res.AltitudeDeg) {
		t.Error("AltitudeDeg is NaN at the exact pole")
	}
	if res.AzimuthDeg != nil {
		t.Errorf("AzimuthDeg = %v, want nil at the exact pole", *res.AzimuthDeg)
	}
}

func TestComputeArcticSummerPolarDay(t *testing.T) {
	res := Compute(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC), 78, 0)

	if res.PolarCapReason != "" {
		t.Fatalf("PolarCapReason = %q, want empty below the polar cap", res.PolarCapReason)
	}
	if res.TodaySunrise != nil || res.TodaySunset != nil {
		t.Error("TodaySunrise/TodaySunset != nil, want both nil during polar day")
	}
	if res.PreviousSunrise != nil || res.NextSunrise != nil {
		t.Error("previous/next sunrise != nil, want both nil deep in polar day")
	}
}
