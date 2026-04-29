# Vietnamese Number Converter — Deployment Guide

## Option A: Standard Service (Recommended for most use cases)

### Prerequisites

- Go 1.21+
- Docker (optional, for container deployment)

### Build

```bash
make build
# Output: bin/server (CGO-disabled, Linux-compatible binary)
```

### Run Directly

```bash
make run
# or
go run cmd/server/main.go
```

### Docker

```bash
# Build image
make docker
# Tag: vietnamese-converter:latest

# Run container
make docker-run
# Maps container port 8080 to host 8080
```

---

## Option B: Turbo Service (High-throughput workloads)

Use when you need 1000+ RPS with sub-100μs latency.

```bash
make turbo-build
make turbo-run
```

**Docker:**
```bash
make turbo-docker
make turbo-deploy
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `DISABLE_GC` | `false` | Turbo only. Disable GC for max performance |
| `GOMAXPROCS` | CPU count | Turbo only. Number of goroutines |

---

## Internal Network Deployment

### Behind Nginx Reverse Proxy

```nginx
upstream bangchu_backend {
    server 10.0.1.10:8080;
    server 10.0.1.11:8080;
    # Add more instances as needed
}

server {
    listen 80;
    server_name bangchu.internal;

    location / {
        proxy_pass http://bangchu_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_connect_timeout 1s;
        proxy_read_timeout 5s;
    }
}
```

### Docker Compose (Multi-instance)

```yaml
version: "3.8"
services:
  converter-1:
    image: vietnamese-converter:latest
    ports: ["8081:8080"]
    environment: ["LOG_LEVEL=info"]
    restart: always
  converter-2:
    image: vietnamese-converter:latest
    ports: ["8082:8080"]
    environment: ["LOG_LEVEL=info"]
    restart: always
  converter-3:
    image: vietnamese-converter:latest
    ports: ["8083:8080"]
    environment: ["LOG_LEVEL=info"]
    restart: always
```

Place Nginx or your internal load balancer in front of the three instances.

---

## Health Check Configuration

**Kubernetes liveness/readiness probe:**
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
```

**Simple TCP ping for monitoring tools:**
```bash
curl -sf http://bangchu.internal:8080/ping
# Returns "pong" on success, exits non-zero on failure
```

---

## Scaling

| Load Level | Instances | Notes |
|------------|-----------|-------|
| <500 RPS | 1 instance | Standard service is sufficient |
| 500–5000 RPS | 2–3 instances | Standard service behind load balancer |
| 5000–10000 RPS | 3+ instances | Switch to turbo service |
| 10000+ RPS | Scale horizontally | Turbo service, one per CPU core |

No shared state. Scale by adding instances behind your load balancer.

---

## System Tuning (Production)

```bash
# Increase file descriptor limits
ulimit -n 65536

# TCP optimizations (Linux)
sysctl -w net.core.somaxconn=32768
sysctl -w net.ipv4.tcp_tw_reuse=1
```

---

## Security Notes

- No authentication is configured by default. Deploy behind your internal firewall/VPC.
- Rate limiting (10,000 req/s) protects against accidental overload.
- No CORS headers — API is internal-only. Add a CORS middleware if external access is needed.
