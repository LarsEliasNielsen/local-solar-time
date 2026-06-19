package server

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"local-solar-time/internal/solar"
)

func marshalToMap(t *testing.T, resp UpdateResponse) map[string]any {
	t.Helper()

	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(raw), "NaN") {
		t.Fatalf("response contains NaN: %s", raw)
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

func TestToUpdateResponseNormalDay(t *testing.T) {
	when := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	m := marshalToMap(t, toUpdateResponse(solar.Compute(when, 55.6761, 12.5683)))

	for _, key := range []string{
		"solar_time", "equation_of_time_minutes", "utc_offset_seconds", "altitude_deg",
		"azimuth_deg", "solar_noon", "today",
		"previous_sunrise", "next_sunrise", "previous_sunset", "next_sunset",
	} {
		if _, ok := m[key]; !ok {
			t.Errorf("expected key %q to be present", key)
		}
	}
	if _, ok := m["polar_cap"]; ok {
		t.Error("polar_cap should be absent for a normal mid-latitude day")
	}
}

func TestToUpdateResponsePolarCap(t *testing.T) {
	when := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	m := marshalToMap(t, toUpdateResponse(solar.Compute(when, 90, 0)))

	for _, key := range []string{"solar_time", "azimuth_deg", "solar_noon", "today",
		"previous_sunrise", "next_sunrise", "previous_sunset", "next_sunset"} {
		v, ok := m[key]
		if !ok {
			t.Errorf("expected key %q to be present (null) above the polar cap", key)
			continue
		}
		if v != nil {
			t.Errorf("expected %q to be null above the polar cap, got %v", key, v)
		}
	}
	if m["polar_cap"] == nil {
		t.Error("expected polar_cap to be present at the exact pole")
	}
}

func TestToUpdateResponsePolarDay(t *testing.T) {
	when := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	m := marshalToMap(t, toUpdateResponse(solar.Compute(when, 78, 0)))

	today, ok := m["today"].(map[string]any)
	if !ok {
		t.Fatalf("expected today to be present, got %v", m["today"])
	}
	if today["sunrise"] != nil || today["sunset"] != nil {
		t.Errorf("expected today.sunrise and today.sunset to be null during polar day, got %+v", today)
	}

	for _, key := range []string{"previous_sunrise", "next_sunrise", "previous_sunset", "next_sunset"} {
		if v, ok := m[key]; ok {
			t.Errorf("expected key %q to be omitted during polar day, got %v", key, v)
		}
	}
}

func TestTodayConsistentWithPreviousNext(t *testing.T) {
	tests := []struct {
		name string
		when time.Time
	}{
		{"daytime", time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC)},
		{"evening", time.Date(2026, 3, 20, 20, 0, 0, 0, time.UTC)},
		{"after midnight", time.Date(2026, 3, 21, 2, 0, 0, 0, time.UTC)},
	}

	const lat, lon = 0.0, 0.0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := toUpdateResponse(solar.Compute(tt.when, lat, lon))

			if resp.Today == nil {
				t.Fatal("expected today to be present")
			}

			switch {
			case resp.Today.Sunrise != nil && resp.NextSunrise != nil && resp.Today.Sunrise.UTC.Equal(resp.NextSunrise.UTC):
			case resp.Today.Sunrise != nil && resp.PreviousSunrise != nil && resp.Today.Sunrise.UTC.Equal(resp.PreviousSunrise.UTC):
			default:
				t.Errorf("today.sunrise %+v not consistent with previous/next sunrise (prev=%+v, next=%+v)",
					resp.Today.Sunrise, resp.PreviousSunrise, resp.NextSunrise)
			}

			switch {
			case resp.Today.Sunset != nil && resp.NextSunset != nil && resp.Today.Sunset.UTC.Equal(resp.NextSunset.UTC):
			case resp.Today.Sunset != nil && resp.PreviousSunset != nil && resp.Today.Sunset.UTC.Equal(resp.PreviousSunset.UTC):
			default:
				t.Errorf("today.sunset %+v not consistent with previous/next sunset (prev=%+v, next=%+v)",
					resp.Today.Sunset, resp.PreviousSunset, resp.NextSunset)
			}
		})
	}
}
