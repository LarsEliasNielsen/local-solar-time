# Deployment

The app runs on [Fly.io](https://fly.io) at `https://local-solar-time.fly.dev`. Fly builds the root `Dockerfile` and runs it as a single Machine (a small VM) that stops when nobody is connected and starts again on the next request.

## Before you deploy

Follow Fly's [Install flyctl](https://docs.fly.io/flyctl/install/) and [Create an app with Fly Launch](https://docs.fly.io/launch/create/) guides. When you're done, you should have:

- A Fly.io account with a credit card added. Fly does not deploy without one.
- `flyctl` installed and signed in (`fly auth login`). It installs as both `fly` and `flyctl`; this guide uses `fly`.
- An app created from the repo root with `fly launch --no-deploy`, using the root `Dockerfile` and no databases or other add-ons.
- `app` and `primary_region` in `fly.toml` set to your app name and region (`fly platform regions` lists them). The rest of the committed `fly.toml` stays as it is.

`fly launch` can also generate a GitHub Actions workflow under `.github/`. Delete it unless you want CI deploys (see [Later](#later)).

## `fly.toml` settings

| Setting | Value | Effect |
|---|---|---|
| `primary_region` | `arn` | Region the Machine runs in |
| `kill_timeout` | `10` | Seconds between Fly's `SIGINT` and a forced kill; matches the server's 10 s shutdown timeout |
| `internal_port` | `8080` | Port the server listens on inside the Machine |
| `force_https` | `true` | Redirects `http://` to `https://` |
| `auto_stop_machines` | `'stop'` | Stops the Machine when it has no connections |
| `auto_start_machines` | `true` | Starts the Machine on the next request |
| `min_machines_running` | `0` | Lets the app scale to zero |
| `concurrency` | `connections`, soft 50, hard 100 | Each open tab holds one WebSocket connection |
| `checks` | `GET /healthz` every 30 s, 5 s timeout | Health check; does not keep a stopped Machine awake |
| `vm` | `shared-cpu-1x`, `256mb` | Smallest Machine size |

## Deploy and redeploy

The first deploy creates one Machine instead of Fly's default of two:

```
fly deploy --ha=false
```

Every later deploy is plain `fly deploy`. Fly builds the image on its remote builder, so Docker does not need to run locally.

Check the result:

```
fly status        # Machine state: started or stopped
fly logs          # server output, including "starting server"
fly ips list      # should show a shared IPv4 and an IPv6
```

If `fly ips list` shows a dedicated IPv4, release it with `fly ips release <ip>`. It costs $2/month and the app does not need it.

## Cost

Fly has no spending cap. These settings bound the bill:

- One Machine with no autoscaling. Running all month, `shared-cpu-1x` with 256 MB costs about $2.19, which is the ceiling. A stopped Machine bills only a few cents for its disk.
- A shared IPv4 instead of a dedicated one.
- `SOLAR_MAX_SESSION` (default `2h`) closes any WebSocket after that long.
- A connection that sends no valid subscribe within 10 seconds is closed.
- The frontend closes its socket when the tab has been hidden for 5 minutes.
- WebSocket upgrades from other origins are rejected, so other sites cannot hold the Machine awake.
- Egress is about $0.02/GB. One viewer for one hour uses about 3.6 MB.

Fly does not offer billing alerts on every account. Without one, watch usage and cost under Billing in the Fly dashboard.

## Cold start

The first visit after the Machine has stopped waits for it to start, page included. The drift probe measured a wake time of about 1.7 s, from opening the connection to the WebSocket upgrade completing. Fly stops the Machine a few minutes after the last connection closes.

## Clock accuracy

Solar time is computed from the server's wall clock, so the host must be NTP-synchronized. Fly does not document its hosts' clock accuracy, so the drift probe measures it.

The probe in `tools/drift-probe/` measures how far the server's clock is from true UTC:

1. It queries an NTP server to find the local clock's offset.
2. It opens `/ws` on a stopped Machine, timing the wake, and subscribes.
3. For each of 10 updates, it times a ping/pong round trip and reads the next update. The offset is the update's `utc` minus the NTP-corrected moment the server sent it, estimated as the receive time minus half the round trip.

Run it after `fly status` shows `stopped`:

```
cd tools/drift-probe
go run . -url wss://local-solar-time.fly.dev/ws -out drift.csv
```

Pass criterion: across at least 10 cold starts spread over one week, median offset under 50 ms and worst case under 100 ms.

## Later

- Custom domain: `fly certs add <domain>`, then point the domain's DNS at the app as the command instructs.
- CI deploys: a GitHub Actions workflow that runs `flyctl deploy --remote-only` with a deploy token stored as a repository secret (`fly tokens create deploy`).
