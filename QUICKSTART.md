# ⚡ Quick Start - Test the Turbo Service

## 🚀 **1-Minute Setup**

```bash
# Clone and enter the directory
cd /path/to/bangchu

# Start the turbo service
go run cmd/turbo/main.go
```

**That's it!** 🎉

## 🌐 **Open the Web Interface**

**In your browser, go to:** `http://localhost:8080`

You'll see a beautiful web interface where you can:
- ✨ Type numbers and see **instant** Vietnamese conversion
- 📊 Watch real-time latency metrics (usually 100-500μs)
- 🔢 Track request count and average performance

## 🧪 **Try These Test Numbers**

| Type This | See This Vietnamese |
|-----------|-------------------|
| `123` | một trăm hai mươi ba đồng |
| `1001` | một nghìn lẻ một đồng |
| `123456789` | một trăm hai mươi ba triệu bốn trăm năm mươi sáu nghìn bảy trăm tám mươi chín đồng |

## 🔧 **Test the API Directly**

```bash
# Health check
curl http://localhost:8080/health

# Convert a number
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 123456789}'

# View metrics
curl http://localhost:8080/metrics
```

## 📊 **Expected Performance**

- **Latency**: 40-100μs (P95)
- **Throughput**: 1000+ RPS
- **Memory**: <5MB total usage
- **Response**: Instant conversion as you type

**If you see these numbers, the turbo service is working perfectly!** ⚡

---

For detailed testing and deployment options, see [TESTING_GUIDE.md](TESTING_GUIDE.md)