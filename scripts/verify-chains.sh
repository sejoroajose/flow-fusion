#!/bin/bash
set -e

echo "🔍 Verifying blockchain networks..."

# Check Ethereum
echo "📡 Checking Ethereum fork..."
ETH_BLOCK=$(curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545 | jq -r '.result')

if [ "$ETH_BLOCK" != "null" ] && [ -n "$ETH_BLOCK" ]; then
    echo "✅ Ethereum fork running - Block: $ETH_BLOCK"
else
    echo "❌ Ethereum fork not responding"
    exit 1
fi

# Check Cosmos
echo "🌌 Checking Cosmos chain..."
COSMOS_STATUS=$(curl -s http://localhost:26657/status | jq -r '.result.sync_info.catching_up')

if [ "$COSMOS_STATUS" = "false" ]; then
    echo "✅ Cosmos chain running and synced"
else
    echo "⏳ Cosmos chain starting up..."
    sleep 10
fi

# Check Redis
echo "🔴 Checking Redis..."
if docker-compose exec redis redis-cli ping | grep -q PONG; then
    echo "✅ Redis running"
else
    echo "❌ Redis not responding"
    exit 1
fi

# Check PostgreSQL
echo "🐘 Checking PostgreSQL..."
if docker-compose exec postgres pg_isready -U flow | grep -q "accepting connections"; then
    echo "✅ PostgreSQL running"
else
    echo "❌ PostgreSQL not responding"
    exit 1
fi

echo "🎉 All blockchain networks are ready!"