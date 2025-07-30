# Flow Fusion

**TWAP-Enabled Cross-Chain Bridge: Ethereum ↔ Cosmos**

Flow Fusion combines 1inch Fusion+ with Time-Weighted Average Price (TWAP) execution for efficient cross-chain swaps between Ethereum and Cosmos.

## ✨ Key Features

- **Bidirectional Bridge**: Seamless ETH ↔ Cosmos atomic swaps
- **TWAP Execution**: Minimize price impact for large trades
- **1inch Integration**: Leverage best-in-class liquidity aggregation
- **Professional UI**: Advanced trading interface

## 🚀 Quick Start

```bash
# Setup environment
make setup

# Start development environment
make dev

# Run tests
make test

# Start demo
make demo
```

## 🏗️ Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│  Ethereum   │◄──►│ Flow Fusion  │◄──►│   Cosmos    │
│   (EVM)     │    │   Relayer    │    │   (SDK)     │
└─────────────┘    └──────────────┘    └─────────────┘
```

## 📄 License

MIT - Built for ETHGlobal Unite 2025