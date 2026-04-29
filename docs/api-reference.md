# Vietnamese Number Converter — API Reference

**Base URL:** `http://localhost:8080` (development) | `http://bangchu.internal:8080` (production)

## Endpoints

### `POST /api/v1/convert` — Convert Number (Primary)

Converts an integer to Vietnamese text with optional currency suffix.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 1433433225}'
```

**Response (200):**
```json
{
  "number": 1433433225,
  "vietnamese": "một tỷ bốn trăm ba mươi ba triệu bốn trăm ba mươi ba nghìn hai trăm hai mươi lăm đồng",
  "processing_time_ms": 0.024084
}
```

**Request Schema:**
| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `number` | int64 | Yes | — | Value to convert (0–999,999,999,999,999) |
| `currency` | string | No | `"đồng"` | Currency suffix. Empty string `""` omits currency |

**Error Responses:**
| Code | Error | Cause |
|------|-------|-------|
| 400 | `"Number must be non-negative"` | `number` < 0 |
| 400 | `"Number too large"` | `number` > 999,999,999,999,999 |
| 400 | `"Invalid request body"` | Malformed JSON or missing `number` field |
| 429 | `"Rate limit exceeded"` | Exceeded 10,000 req/s limit |
| 500 | `"Internal Server Error"` | Unexpected server error |

---

### `GET /api/v1/convert` — Convert via Query Params

```bash
curl "http://localhost:8080/api/v1/convert?number=123456789&currency=VNĐ"
```

**Parameters:**
| Param | Required | Default | Description |
|-------|----------|---------|-------------|
| `number` | Yes | — | Same constraints as POST |
| `currency` | No | `"đồng"` | Same behavior as POST |

---

### `GET /health` — Health Check

```bash
curl http://localhost:8080/health
```

**Response (200):**
```json
{ "status": "healthy" }
```

Use for load balancer health checks and Kubernetes liveness probes.

---

### `GET /ping` — Connectivity Test

```bash
curl http://localhost:8080/ping
# pong
```

Plain text response. Lower overhead than `/health` for TCP-level monitoring.

---

## Client Code Examples

### Node.js / TypeScript

```typescript
interface ConvertResponse {
  number: number;
  vietnamese: string;
  processing_time_ms: number;
}

async function convertNumber(number: number, currency = "đồng"): Promise<ConvertResponse> {
  const res = await fetch("http://bangchu.internal:8080/api/v1/convert", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ number, currency }),
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`);
  return res.json();
}

// Usage:
const result = await convertNumber(1433433225);
console.log(result.vietnamese);
// "một tỷ bốn trăm ba mươi ba triệu bốn trăm ba mươi ba nghìn hai trăm hai mươi lăm đồng"
```

### Python

```python
import requests

def convert_number(number: int, currency: str = "đồng") -> dict:
    resp = requests.post(
        "http://bangchu.internal:8080/api/v1/convert",
        json={"number": number, "currency": currency},
        timeout=5,
    )
    resp.raise_for_status()
    return resp.json()

# Usage:
result = convert_number(1433433225)
print(result["vietnamese"])
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type ConvertRequest struct {
    Number   int64  `json:"number"`
    Currency string `json:"currency,omitempty"`
}

type ConvertResponse struct {
    Number           int64   `json:"number"`
    Vietnamese       string  `json:"vietnamese"`
    ProcessingTimeMs float64 `json:"processing_time_ms"`
}

func ConvertNumber(number int64, currency string) (*ConvertResponse, error) {
    req := ConvertRequest{Number: number, Currency: currency}
    body, _ := json.Marshal(req)

    resp, err := http.Post(
        "http://bangchu.internal:8080/api/v1/convert",
        "application/json",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result ConvertResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

---

## Vietnamese Language Rules

The converter handles these linguistic exceptions automatically:

| Rule | Example | Output |
|------|---------|--------|
| `1` → `"mốt"` in tens position | 11, 21, 31 | `"mười mốt"`, `"hai mươi mốt"` |
| `4` → `"tư"` in tens position | 14, 24, 34 | `"mười bốn"`, `"hai mươi tư"` |
| `5` → `"lăm"` in tens position | 15, 25, 35 | `"mười lăm"`, `"hai mươi lăm"` |
| `0` → `"lẻ"` for gaps | 101, 1001 | `"một trăm lẻ một"`, `"nghìn lẻ một"` |
| Scale transitions | 1000, 1000000 | `"nghìn"`, `"triệu"` |

---

## Rate Limiting

- **Default:** 10,000 requests per second
- **Burst:** Up to 10,000 requests before throttling
- Response `429` with `{"error":"Rate limit exceeded","code":429}` when exceeded

---

## Performance Characteristics

| Metric | Standard Service | Turbo Service |
|--------|-----------------|---------------|
| P50 Latency | <25μs | <50μs |
| P95 Latency | <100μs | <100μs |
| Throughput | ~8,000 RPS | 10,000+ RPS |
| Memory/Request | Minimal (buffer pooled) | Zero allocation |

---

## Headers

| Header | Direction | Description |
|--------|-----------|-------------|
| `X-Request-ID` | Response | Unique UUID for each request, used for log correlation |
| `Content-Type: application/json` | Request | Required for `POST /api/v1/convert` |

---

## Go Library Alternative

If your service is written in Go, import the converter directly instead of calling the API:

```go
import "vietnamese-converter/pkg/converter"

conv := converter.NewVietnameseConverter()
text, _ := conv.Convert(1433433225)
// "một tỷ bốn trăm ba mươi ba triệu bốn trăm ba mươi ba nghìn hai trăm hai mươi lăm đồng"
```

This eliminates network latency entirely. See `README.md` for module setup.
