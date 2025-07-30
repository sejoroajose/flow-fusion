.PHONY: setup build test deploy clean dev

# =============================================================================
# FLOW FUSION - ONE DAY DEVELOPMENT MAKEFILE
# =============================================================================

# Load environment variables
include .env
export

# Default target
all: setup dev

# =============================================================================
# SETUP PHASE (0-30 minutes)
# =============================================================================

setup: setup-env setup-deps
	@echo "✅ Flow Fusion development environment ready!"

setup-env:
	@echo "🔧 Setting up environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "⚠️  Please update .env with your keys"; fi
	@chmod +x scripts/*.sh

setup-deps:
	@echo "📦 Installing dependencies..."
	@yarn install
	@echo "✅ Dependencies installed"

# =============================================================================
# DEVELOPMENT PHASE (30+ minutes)
# =============================================================================

dev: dev-chains
	@echo "🚀 Development environment started!"
	@echo "📊 Ethereum RPC: http://localhost:8545"
	@echo "🌌 Cosmos RPC: http://localhost:26657"

dev-chains:
	@echo "🐳 Starting blockchain networks..."
	@docker-compose up -d ethereum cosmos redis postgres
	@echo "⏳ Waiting for chains to initialize..."
	@sleep 30
	@./scripts/verify-chains.sh

dev-relayer:
	@echo "⚙️ Starting relayer service..."
	@docker-compose --profile relayer up -d relayer

dev-frontend:
	@echo "🎨 Starting frontend..."
	@docker-compose --profile frontend up -d frontend

# =============================================================================
# BUILD PHASE
# =============================================================================

build: build-contracts build-relayer-binary

build-contracts: build-ethereum build-cosmos
	@echo "✅ All contracts built successfully"

build-ethereum:
	@echo "🔨 Building Ethereum contracts..."
	@docker-compose exec ethereum forge build
	@echo "✅ Ethereum contracts built"

build-cosmos:
	@echo "🦀 Building CosmWasm contracts..."
	@docker-compose exec cosmos bash -c "cd /contracts && cargo build --release --target wasm32-unknown-unknown"
	@echo "✅ CosmWasm contracts built"

build-relayer-binary:
	@echo "⚙️ Building Go relayer..."
	@cd relayer && go build -o bin/relayer ./cmd/relayer
	@echo "✅ Relayer binary built"

# =============================================================================
# DEPLOYMENT PHASE
# =============================================================================

deploy: deploy-ethereum deploy-cosmos
	@echo "✅ All contracts deployed"

deploy-ethereum:
	@echo "🚀 Deploying Ethereum contracts..."
	@docker-compose exec ethereum forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --private-key $(DEPLOYER_PRIVATE_KEY) --broadcast
	@echo "✅ Ethereum contracts deployed"

deploy-cosmos:
	@echo "🚀 Deploying CosmWasm contracts..."
	@./scripts/deploy-cosmos.sh
	@echo "✅ CosmWasm contracts deployed"

# =============================================================================
# TESTING PHASE
# =============================================================================

test: test-contracts test-integration
	@echo "✅ All tests passed"

test-contracts:
	@echo "🧪 Testing smart contracts..."
	@docker-compose exec ethereum forge test -vv
	@docker-compose exec cosmos bash -c "cd /contracts && cargo test"

test-integration:
	@echo "🔗 Running integration tests..."
	@cd relayer && go test ./tests/integration/... -v

test-e2e:
	@echo "🎭 Running end-to-end tests..."
	@./scripts/e2e-test.sh

# =============================================================================
# DEMO PREPARATION
# =============================================================================

demo: prepare-demo run-demo

prepare-demo:
	@echo "🎬 Preparing demo environment..."
	@./scripts/prepare-demo.sh

run-demo:
	@echo "🎪 Running demo scenario..."
	@./scripts/demo-scenario.sh

# =============================================================================
# UTILITY COMMANDS
# =============================================================================

logs:
	@docker-compose logs -f

status:
	@docker-compose ps
	@echo ""
	@echo "🌐 Service Status:"
	@echo "  Ethereum RPC: http://localhost:8545"
	@echo "  Cosmos RPC: http://localhost:26657"
	@echo "  Relayer API: http://localhost:8080"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Redis: localhost:6379"
	@echo "  PostgreSQL: localhost:5432"

clean:
	@echo "🧹 Cleaning up..."
	@docker-compose down -v
	@docker system prune -f
	@rm -rf contracts/ethereum/out contracts/ethereum/cache
	@rm -rf contracts/cosmos/target
	@rm -rf relayer/bin

restart:
	@docker-compose restart

# =============================================================================
# DEVELOPMENT SHORTCUTS
# =============================================================================

quick-start: setup dev deploy
	@echo "🚀 Flow Fusion is running!"

ethereum-shell:
	@docker-compose exec ethereum bash

cosmos-shell:
	@docker-compose exec cosmos bash

redis-cli:
	@docker-compose exec redis redis-cli

db-shell:
	@docker-compose exec postgres psql -U flow -d flowfusion

# =============================================================================
# HACKATHON HELPERS
# =============================================================================

check-requirements:
	@echo "📋 Checking hackathon requirements..."
	@./scripts/check-requirements.sh

package-submission:
	@echo "📦 Packaging hackathon submission..."
	@./scripts/package-submission.sh

# =============================================================================
# HOUR-BY-HOUR CHECKPOINTS
# =============================================================================

checkpoint-1: # Hour 1: Environment ready
	@echo "⏰ Hour 1 Checkpoint: Environment Setup"
	@make status
	@./scripts/verify-chains.sh

checkpoint-2: # Hour 2: Contracts building
	@echo "⏰ Hour 2 Checkpoint: Contract Development"
	@make build-contracts

checkpoint-4: # Hour 4: Cross-chain swaps working
	@echo "⏰ Hour 4 Checkpoint: Basic Bridge Working"
	@make test-integration

checkpoint-6: # Hour 6: TWAP implementation
	@echo "⏰ Hour 6 Checkpoint: TWAP Integration"
	@./scripts/test-twap.sh

checkpoint-8: # Hour 8: End-to-end working
	@echo "⏰ Hour 8 Checkpoint: Complete System"
	@make test-e2e

checkpoint-10: # Hour 10: Demo ready
	@echo "⏰ Hour 10 Checkpoint: Demo Ready"
	@make demo