# drift-probe

Measures how far a deployed server's clock is from true UTC. It is its own Go module, so `make test` and `make lint` at the repo root never touch it. Delete the folder to remove it.

## Method

1. Query an NTP server to find the local clock's offset, trying up to 3 times. Each pool lookup can return a different server, so a retry skips an unresponsive one. This runs before the dial, so a failed run doesn't wake the Machine.
2. Dial `/ws` and time the handshake. On a stopped Machine this is the cold start wake latency.
3. Subscribe, then for each sample: send a ping, time the pong for the round trip, and read the next update.
4. Each sample's offset is `utc - (recv + ntpOffset - RTT/2)`: the server's timestamp minus the NTP-corrected moment it sent the update.

## Usage

```
cd tools/drift-probe
go run . -url wss://<app>.fly.dev/ws -out drift.csv
```

| Flag | Default | Meaning |
|---|---|---|
| `-url` | (required) | WebSocket URL, ending in `/ws` |
| `-ntp` | `pool.ntp.org` | NTP server |
| `-samples` | `10` | Updates measured per run |
| `-out` | (none) | CSV file to append the run to; the header is written when the file is new |

Each run prints a summary and appends one row:

```
run_at,wake_ms,ntp_offset_ms,median_rtt_ms,median_abs_offset_ms,max_abs_offset_ms,samples
```

## Clock gate

Run it at least 10 times over one week, each time after `fly status` shows the Machine as `stopped`. Pass: median `|offset|` under 50 ms and worst case under 100 ms. During the same runs, confirm the Machine still stops with the `/healthz` check enabled.
