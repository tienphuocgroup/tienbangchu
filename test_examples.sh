#!/bin/bash

echo "🧪 Testing Vietnamese Number Converter..."

# Wait for server to start
sleep 2

echo "Testing key Vietnamese language rules:"

echo -n "4 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 4}' | jq -r '.vietnamese'

echo -n "14 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 14}' | jq -r '.vietnamese'

echo -n "24 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 24}' | jq -r '.vietnamese'

echo -n "40 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 40}' | jq -r '.vietnamese'

echo -n "34000 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 34000}' | jq -r '.vietnamese'

echo -n "21 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 21}' | jq -r '.vietnamese'

echo -n "101 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 101}' | jq -r '.vietnamese'

echo -n "50050050 → "
curl -s -X POST http://localhost:8080/api/v1/convert -H "Content-Type: application/json" -d '{"number": 50050050}' | jq -r '.vietnamese'

echo ""
echo "✅ Test complete!"
