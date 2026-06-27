# local-solar-time

A Go service that streams apparent (true) solar time, Sun position, and rise/set events over WebSocket for any latitude/longitude, computed continuously from the server's own NTP-disciplined clock. Ships with a Vite + React + TypeScript web frontend served by nginx that visualizes solar time on an SVG half-circle clock, with browser geolocation, manual coordinate input, and automatic reconnection.

## Prerequisites

- Docker with the Compose plugin
- Go 1.23+ — only needed for local builds outside Docker
- An NTP-synchronized host clock — solar time accuracy depends entirely on it

## Quick start

```sh
git clone <repo-url>
cd local-solar-time
make docker-up   # builds images and starts both services
```

Open `http://localhost` for the web UI, or connect a headless client directly to the WebSocket endpoint (see [Backend](#backend)).

## Frontend

A Vite + React + TypeScript single-page app served by nginx on port 80.

- Connects immediately with the default location (Copenhagen, `lat: 55.6761, lon: 12.5683`) so the clock is visible before any geolocation prompt.
- Requests browser geolocation in the background; on grant, transitions to the real location with a ~0.4s tween.
- Manual lat/lon inputs with inline validation let you override the location at any time.
- Reconnects automatically with exponential backoff (1s initial, 30s cap) on connection loss.

## Backend

### WebSocket API

Connect, then send a subscribe message once with your location:

```json
{ "lat": 55.6761, "lon": 12.5683 }
```

The server then streams JSON updates at the configured cadence:

```json
{
  "solar_time": "12:34:56",
  "equation_of_time_minutes": -3.2,
  "utc_offset_seconds": -192,
  "altitude_deg": 42.1,
  "azimuth_deg": 195.3,
  "solar_noon": {
    "solar_time": "12:00:00",
    "utc": "2026-06-14T10:23:00Z"
  },
  "today": {
    "sunrise": { "solar_time": "04:51:00", "utc": "2026-06-14T02:54:00Z" },
    "sunset":  { "solar_time": "19:09:00", "utc": "2026-06-14T17:12:00Z" }
  },
  "previous_sunrise": { "solar_time": "04:51:00", "utc": "2026-06-13T02:54:00Z" },
  "next_sunrise":     { "solar_time": "04:51:00", "utc": "2026-06-14T02:54:00Z" },
  "previous_sunset":  { "solar_time": "19:09:00", "utc": "2026-06-13T17:12:00Z" },
  "next_sunset":      { "solar_time": "19:09:00", "utc": "2026-06-14T17:12:00Z" }
}
```

Above latitude ~89.4°, `solar_time`, `azimuth_deg`, `solar_noon`, `today`, and the four `*_sunrise`/`*_sunset` fields are `null`, and a `polar_cap` object with a `reason` string is added instead. During polar day or polar night at lower latitudes, `today.sunrise`/`today.sunset` are `null` but the bracketing `previous_sunrise`/`next_sunrise`/`previous_sunset`/`next_sunset` fields are omitted entirely.

Invalid coordinates (e.g. `lat` outside ±90) get `{ "error": "..." }` instead of an update, and the connection stays open for another subscribe attempt.

### Configuration

| Variable | Flag | Default | Description |
|---|---|---|---|
| `SOLAR_PORT` | `--port` | `8080` | Backend WebSocket listen port (internal to the Compose network; the nginx frontend is always on port 80) |
| `SOLAR_CADENCE` | `--cadence` | `1s` | Update push interval (e.g. `1s`, `500ms`) |

Precedence is flag > environment variable > default. Copy `.env.example` to `.env` to override locally.

The process logs structured JSON to stdout (via `zerolog`) and shuts down gracefully on `SIGINT`/`SIGTERM`, closing in-flight WebSocket connections before exiting.

## Development

```sh
make build   # compile the local binary to ./local-solar-time
```

## Tests and lint

```sh
make test   # Go package tests
make lint   # Go package linting (golangci-lint)
```

## License

Licensed under the [PolyForm Noncommercial License 1.0.0](LICENSE). Noncommercial use is free; commercial use requires a separate agreement.
