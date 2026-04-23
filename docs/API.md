# GeoCalc API

Base URL: `http://<host>:<APP_PORT>`

All responses have `Content-Type: application/json`.

---

## Response formats

### Success

```json
{
  "success": true,
  "data": { ... }
}
```

### Error

```json
{
  "success": false,
  "message": "error description"
}
```

---

## Endpoints

### GET /

Health check. Always returns 200 with empty data.

**Response `200`**
```json
{
  "success": true,
  "data": {}
}
```

---

### GET /status

Returns application state.

**Response `200`**
```json
{
  "version": "v1.0",
  "deploy_flag": "blue",
  "polygon": "exists",
  "uptime": "2006-01-02 15:04:05",
  "last_updated_at": "2006-01-02 15:04:05"
}
```

| Field | Type | Description |
|---|---|---|
| `version` | string | Value of `APP_VERSION` env variable |
| `deploy_flag` | string | Deployment label set via `APP_DEPLOY_FLAG` (e.g. `"blue"`, `"green"`, `"v2"`) — useful for identifying which instance is responding |
| `polygon` | string | `"exists"` or `"not exists"` |
| `uptime` | string | Server start time, format `YYYY-MM-DD HH:MM:SS` |
| `last_updated_at` | string | Last polygon load time; empty string if polygon was never loaded |

---

### POST /load-polygon

Loads a GeoJSON FeatureCollection as the active polygon. Replaces any previously loaded polygon.

**Request**

`Content-Type: application/x-www-form-urlencoded`

| Parameter | Required | Description |
|---|---|---|
| `polygon` | yes | GeoJSON FeatureCollection string |

**Response `200`**
```json
{
  "success": true,
  "data": {}
}
```

**Errors**

| Status | Message | Reason |
|---|---|---|
| `400` | `Polygon is required` | Parameter `polygon` is missing or blank |
| `400` | `Invalid GeoJSON` | Value is not a valid GeoJSON FeatureCollection |
| `405` | — | Method other than POST used |

---

### GET /show-polygon

Returns the currently loaded polygon.

**Response `200` — polygon loaded**
```json
{
  "success": true,
  "data": {
    "polygon": "{\"type\":\"FeatureCollection\", ...}"
  }
}
```

**Response `200` — no polygon loaded**
```json
{
  "success": true,
  "data": {
    "polygon": null
  }
}
```

---

### GET /check-point

Checks whether a geographic point falls inside the loaded polygon. Supports both `Polygon` and `MultiPolygon` geometry types.

**Query parameters**

| Parameter | Required | Description |
|---|---|---|
| `lat` | yes | Latitude, decimal degrees (e.g. `55.7558`) |
| `lon` | yes | Longitude, decimal degrees (e.g. `37.6173`) |

**Response `200`**
```json
{
  "success": true,
  "data": {
    "point_status": "inside polygon"
  }
}
```

`point_status` is either `"inside polygon"` or `"out of polygon"`.

**Errors**

| Status | Message | Reason |
|---|---|---|
| `400` | `Polygon not found` | No polygon has been loaded yet |
| `400` | `Invalid lat value` | `lat` is missing or not a valid number |
| `400` | `Invalid lon value` | `lon` is missing or not a valid number |
| `500` | `Invalid polygon GeoJSON` | Stored polygon failed to parse (should not occur after validation on load) |
| `405` | — | Method other than GET used |

---

## Blue/Green deployment

The service runs as two identical instances (`app_blue` and `app_green`) behind an nginx reverse proxy. At any moment only one instance receives traffic. `deploy_flag` in `/status` shows which instance is currently active.

**Start the stack**

```bash
docker compose up -d
```

Blue is active by default.

**Switch traffic to green**

```bash
./switch.sh green
```

**Switch traffic back to blue**

```bash
./switch.sh blue
```

**Verify the active instance**

```bash
curl http://localhost:<APP_PORT>/status
# "deploy_flag": "green"
```

The switch is zero-downtime: nginx reloads its config while both containers keep running.

---

## Example flow

```bash
# 1. Load a polygon
curl -X POST http://localhost:8080/load-polygon \
  -d 'polygon={"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[37.0,55.0],[38.0,55.0],[38.0,56.0],[37.0,56.0],[37.0,55.0]]]},"properties":{}}]}'

# 2. Check a point
curl "http://localhost:8080/check-point?lat=55.5&lon=37.5"

# 3. Check status
curl http://localhost:8080/status
```
