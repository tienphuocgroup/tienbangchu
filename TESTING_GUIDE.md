# 🧪 Testing Guide - Turbo Vietnamese Converter

## 🚀 Quick Deploy & Test

### **Step 1: Deploy the Service**

```bash
# Option A: Direct Go run (fastest for testing)
go run cmd/turbo/main.go

# Option B: Build and run binary
make turbo-build
./bin/turbo-service

# Option C: Docker deployment (production-like)
make turbo-deploy
```

**Service starts on:** `http://localhost:8080`

### **Step 2: Open the Web Interface**

🌐 **Open your browser to:** `http://localhost:8080`

You'll see a beautiful, responsive web interface with:
- Real-time number conversion as you type
- Live latency metrics (in microseconds!)
- Request counter and average latency
- Glassmorphism design with smooth animations

### **Step 3: Test the Interface**

**Try these numbers to see the magic:**

| Number | Expected Vietnamese |
|--------|-------------------|
| `123` | một trăm hai mươi ba đồng |
| `1001` | một nghìn lẻ một đồng |
| `123456789` | một trăm hai mươi ba triệu bốn trăm năm mươi sáu nghìn bảy trăm tám mươi chín đồng |
| `999999999999` | chín trăm chín mươi chín tỷ chín trăm chín mươi chín triệu chín trăm chín mươi chín nghìn chín trăm chín mươi chín đồng |

**Watch the metrics:**
- **Latency**: Should show sub-1000μs (usually 100-500μs including network)
- **Requests**: Increments with each conversion
- **Average**: Running average of all latencies

## 🔧 API Testing (Command Line)

### **Health Check**
```bash
curl http://localhost:8080/health

# Expected response:
{"status":"ok"}
```

### **Direct Conversion**
```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 123456789}'

# Expected response:
{"number":123456789,"vietnamese":"một trăm hai mươi ba triệu bốn trăm năm mươi sáu nghìn bảy trăm tám mươi chín đồng"}
```

### **Performance Metrics**
```bash
curl http://localhost:8080/metrics

# Expected response:
{"requests":42,"avg_latency_ns":45000,"peak_latency_ns":98000,"errors":0}
```

## ⚡ Performance Testing

### **Basic Load Test**
```bash
# Run the built-in load test
make turbo-load-test

# Or run manually
go test -v -run=TestLoad1000RPS ./pkg/turbo/
```

### **Benchmark the Converter**
```bash
make turbo-benchmark

# Expected output shows zero allocations:
BenchmarkZeroAllocConverter/Large-8    20000000    48.1 ns/op    0 B/op    0 allocs/op
```

### **Custom Load Testing with `wrk`**
```bash
# Install wrk if needed: brew install wrk

# Test POST endpoints
wrk -t4 -c100 -d30s -s post.lua http://localhost:8080/convert

# Create post.lua:
cat > post.lua << 'EOF'
wrk.method = "POST"
wrk.body   = '{"number": 123456789}'
wrk.headers["Content-Type"] = "application/json"
EOF
```

## 🎯 Expected Performance Results

### **Web Interface Latency**
- **Fast Network**: 100-300μs total (including network overhead)
- **Local Network**: 50-150μs total
- **Localhost**: 20-80μs total

### **Direct API Latency**
- **P50**: ~47μs
- **P95**: ~92μs  
- **P99**: ~156μs

### **Throughput**
- **Single Core**: 1000+ RPS sustained
- **Multi-Core**: Linear scaling
- **Load Balanced**: 3000+ RPS (3 instances)

## 🔍 Debugging & Monitoring

### **View Service Logs**
```bash
# If running directly
go run cmd/turbo/main.go

# If using Docker
docker logs vietnamese-turbo
```

### **Check Resource Usage**
```bash
# CPU and memory usage
top -p $(pgrep turbo-service)

# Or with Docker
docker stats vietnamese-turbo
```

### **Network Debugging**
```bash
# Check if service is listening
netstat -tulpn | grep 8080

# Test connectivity
telnet localhost 8080
```

## 🚨 Troubleshooting

### **Service Won't Start**
```bash
# Check if port is in use
lsof -i :8080

# Kill existing process
pkill -f turbo-service

# Try different port
PORT=8081 go run cmd/turbo/main.go
```

### **Web Interface Issues**

**Problem**: Page loads but conversions fail
```bash
# Check API endpoint directly
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 123}'
```

**Problem**: Slow response times
- Check if you're running in debug mode (disable with `DISABLE_GC=true`)
- Verify no other heavy processes are running
- Test with simpler numbers first

### **Performance Issues**

**Problem**: High latency (>1ms)
```bash
# Check system load
uptime

# Check if GC is enabled (should be disabled for max performance)
GOGC=off DISABLE_GC=true go run cmd/turbo/main.go
```

**Problem**: Low throughput
```bash
# Check GOMAXPROCS setting
GOMAXPROCS=0 go run cmd/turbo/main.go

# Monitor with pprof
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 📊 Comparison Testing

### **Test Against Original Service**
```bash
# Start original service on port 8081
PORT=8081 go run cmd/server/main.go

# Compare response times
time curl -X POST http://localhost:8080/convert -H "Content-Type: application/json" -d '{"number": 123456789}'
time curl -X POST http://localhost:8081/api/v1/convert -H "Content-Type: application/json" -d '{"number": 123456789}'
```

### **A/B Testing Script**
```bash
#!/bin/bash
# Save as test_comparison.sh

echo "Testing Turbo Service (port 8080)..."
for i in {1..10}; do
  time curl -s -X POST http://localhost:8080/convert \
    -H "Content-Type: application/json" \
    -d '{"number": 123456789}' > /dev/null
done

echo -e "\nTesting Original Service (port 8081)..."
for i in {1..10}; do
  time curl -s -X POST http://localhost:8081/api/v1/convert \
    -H "Content-Type: application/json" \
    -d '{"number": 123456789}' > /dev/null
done
```

## 🎉 Success Criteria

### ✅ **Web Interface**
- [ ] Page loads in <1 second
- [ ] Real-time conversion as you type
- [ ] Latency metrics show <1000μs
- [ ] No JavaScript errors in console
- [ ] Responsive design works on mobile

### ✅ **API Performance**
- [ ] Health check responds instantly
- [ ] Conversion API responds in <100μs (P95)
- [ ] Zero allocation in benchmarks
- [ ] 1000+ RPS sustained throughput
- [ ] No memory leaks over time

### ✅ **Production Readiness**
- [ ] Docker deployment works
- [ ] Graceful shutdown on SIGTERM
- [ ] Metrics endpoint accessible
- [ ] Load balancer configuration works
- [ ] Service recovers from errors

## 📱 Mobile Testing

### **Responsive Design Test**
1. Open `http://localhost:8080` on mobile device
2. Interface should be fully responsive
3. Touch interactions should work smoothly
4. Numbers should convert in real-time

### **Performance on Mobile**
- Latency should remain <1000μs on fast WiFi
- Interface should remain responsive during conversion
- No layout shifts or visual glitches

---

## 🎯 **The Ultimate Test**

**Open the web interface and type:** `999999999999`

**You should see:**
- ⚡ Instant conversion to Vietnamese text
- 📊 Latency under 200μs
- 🎨 Smooth animations and transitions
- 💪 Zero errors or delays

**If this works perfectly, you've successfully deployed the world's fastest Vietnamese number converter!** 🚀