# local-solar-time

A Go service that streams apparent (true) solar time, Sun position, and rise/set events over WebSocket for any latitude/longitude, computed continuously from the server's own NTP-disciplined clock.

## Prerequisites

- Go 1.23+
- Docker (with the Compose plugin)
- An NTP-synchronized host clock - solar time accuracy depends entirely on it

## Quick start

```sh
git clone <repo-url>
cd local-solar-time
make build       # local binary at ./local-solar-time
make docker-up   # builds and runs the service in Docker
```

## Connecting

The service exposes a WebSocket endpoint. Connect, then send a subscribe message once with your location:

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

## Environment variables

| Variable | Flag | Default | Description |
|---|---|---|---|
| `SOLAR_PORT` | `--port` | `8080` | WebSocket listen port (also used by `docker-compose.yml`'s published port mapping) |
| `SOLAR_CADENCE` | `--cadence` | `1s` | Update push interval (e.g. `1s`, `500ms`) |

Precedence is flag > environment variable > default. Copy `.env.example` to `.env` to override locally.

The process logs structured JSON to stdout (via `zerolog`) and shuts down gracefully on `SIGINT`/`SIGTERM`, closing in-flight WebSocket connections before exiting.

## Tests and lint

```sh
make test
make lint
```

## License

Licensed under the [PolyForm Noncommercial License 1.0.0](LICENSE). Noncommercial use is free; commercial use requires a separate agreement.
