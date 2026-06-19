package server

import (
	"encoding/json"
	"time"

	"local-solar-time/internal/solar"
)

// ErrorResponse is sent in place of an update when a subscribe message fails validation.
type ErrorResponse struct {
	Error string `json:"error"`
}

// EventPayload is the wire representation of a solar.Event.
type EventPayload struct {
	SolarTime string    `json:"solar_time"`
	UTC       time.Time `json:"utc"`
}

// TodayPayload holds today's sunrise/sunset, each null during polar day or polar night.
type TodayPayload struct {
	Sunrise *EventPayload `json:"sunrise"`
	Sunset  *EventPayload `json:"sunset"`
}

// PolarCapPayload is present only when the subscribed latitude exceeds the polar cap threshold.
type PolarCapPayload struct {
	Reason string `json:"reason"`
}

// UpdateResponse is the canonical streaming update payload. See
// local-solar-time-epics.md for the nullability rules MarshalJSON enforces:
// previous/next sunrise and sunset are omitted during polar day/night but
// null above the polar cap, which a plain `omitempty` tag can't express
// since both cases are a nil pointer in Go.
type UpdateResponse struct {
	SolarTime             *string          `json:"solar_time"`
	EquationOfTimeMinutes float64          `json:"equation_of_time_minutes"`
	UTCOffsetSeconds      *int             `json:"utc_offset_seconds"`
	AltitudeDeg           float64          `json:"altitude_deg"`
	AzimuthDeg            *float64         `json:"azimuth_deg"`
	SolarNoon             *EventPayload    `json:"solar_noon"`
	Today                 *TodayPayload    `json:"today"`
	PreviousSunrise       *EventPayload    `json:"previous_sunrise"`
	NextSunrise           *EventPayload    `json:"next_sunrise"`
	PreviousSunset        *EventPayload    `json:"previous_sunset"`
	NextSunset            *EventPayload    `json:"next_sunset"`
	PolarCap              *PolarCapPayload `json:"polar_cap"`
}

func (u UpdateResponse) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"solar_time":               u.SolarTime,
		"equation_of_time_minutes": u.EquationOfTimeMinutes,
		"utc_offset_seconds":       u.UTCOffsetSeconds,
		"altitude_deg":             u.AltitudeDeg,
		"azimuth_deg":              u.AzimuthDeg,
		"solar_noon":               u.SolarNoon,
		"today":                    u.Today,
	}

	if u.PolarCap != nil {
		m["previous_sunrise"] = nil
		m["next_sunrise"] = nil
		m["previous_sunset"] = nil
		m["next_sunset"] = nil
		m["polar_cap"] = u.PolarCap
	} else {
		if u.PreviousSunrise != nil {
			m["previous_sunrise"] = u.PreviousSunrise
		}
		if u.NextSunrise != nil {
			m["next_sunrise"] = u.NextSunrise
		}
		if u.PreviousSunset != nil {
			m["previous_sunset"] = u.PreviousSunset
		}
		if u.NextSunset != nil {
			m["next_sunset"] = u.NextSunset
		}
	}

	return json.Marshal(m)
}

func toEventPayload(e *solar.Event) *EventPayload {
	if e == nil {
		return nil
	}
	return &EventPayload{SolarTime: e.SolarTime, UTC: e.UTC}
}

// toUpdateResponse maps the solar engine's flat Result onto the canonical nested wire schema.
func toUpdateResponse(r solar.Result) UpdateResponse {
	resp := UpdateResponse{
		SolarTime:             r.SolarTime,
		EquationOfTimeMinutes: r.EquationOfTimeMinutes,
		UTCOffsetSeconds:      r.UTCOffsetSeconds,
		AltitudeDeg:           r.AltitudeDeg,
		AzimuthDeg:            r.AzimuthDeg,
		SolarNoon:             toEventPayload(r.SolarNoon),
		PreviousSunrise:       toEventPayload(r.PreviousSunrise),
		NextSunrise:           toEventPayload(r.NextSunrise),
		PreviousSunset:        toEventPayload(r.PreviousSunset),
		NextSunset:            toEventPayload(r.NextSunset),
	}

	if r.PolarCapReason != "" {
		resp.PolarCap = &PolarCapPayload{Reason: r.PolarCapReason}
	} else {
		resp.Today = &TodayPayload{
			Sunrise: toEventPayload(r.TodaySunrise),
			Sunset:  toEventPayload(r.TodaySunset),
		}
	}

	return resp
}
