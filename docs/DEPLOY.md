# Deploy

## Prerequisites

- Docker
- Docker Compose v2

## Configuration

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

| Variable | Description | Example |
|---|---|---|
| `APP_NAME` | Application name | `GeoCalc` |
| `APP_PORT` | Public port exposed by nginx | `8080` |
| `APP_VERSION` | Application version label | `v1.0` |
| `POLYGON_SOURCE_URL` | URL to fetch GeoJSON polygon on startup (optional) | `https://example.com/polygon.json` |

> `APP_DEPLOY_FLAG` is set automatically per instance (`blue` / `green`) in `docker-compose.yml` and should not be set in `.env`.

## Start

```bash
docker compose up -d
```

Builds both images and starts three containers: `geo_calc_blue`, `geo_calc_green`, `nginx`. Blue is active by default.

## Blue/Green switch

Traffic is routed by nginx to one active instance at a time. The other instance keeps running in the background.

**Switch to green:**

```bash
./switch.sh green
```

**Switch back to blue:**

```bash
./switch.sh blue
```

The script patches `nginx/nginx.conf` and reloads nginx with zero downtime. Confirm the active instance:

```bash
curl http://localhost:<APP_PORT>/status
# { "deploy_flag": "green", ... }
```

### Typical deploy flow

```
1. New version is ready
2. Stack is running on blue
3. Rebuild green with the new image:
       docker compose build app_green
       docker compose up -d app_green
4. Verify green is healthy:
       docker compose logs app_green
       curl http://localhost:<APP_PORT>/status   # still blue
4. Switch traffic:
       ./switch.sh green
5. Verify:
       curl http://localhost:<APP_PORT>/status   # deploy_flag: green
6. Blue stays running as fallback. Roll back if needed:
       ./switch.sh blue
```

## Stop

```bash
docker compose down
```

## Logs

```bash
# All containers
docker compose logs -f

# Specific instance
docker compose logs -f app_green
```

## Rebuild

```bash
# Rebuild both images and restart
docker compose up -d --build
```
