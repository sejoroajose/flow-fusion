#!/bin/bash
set -e

echo "Building Solidity contracts..."
cd contracts/ethereum
forge build
cd ../..

echo "Generating Go bindings..."
mkdir -p relayer/internal/ethereum/generated

# Generate all contract bindings
for contract in FlowFusionEscrowFactory IEscrowFactory IBaseEscrow; do
    echo "Generating binding for $contract..."
    
    # Extract just the ABI from the Foundry artifact
    jq '.abi' contracts/ethereum/out/$contract.sol/$contract.json > /tmp/${contract}_abi.json
    
    abigen --abi /tmp/${contract}_abi.json \
           --pkg generated \
           --type $contract \
           --out relayer/internal/ethereum/generated/${contract,,}.go
    
    # Clean up temporary file
    rm /tmp/${contract}_abi.json
done

echo "Go bindings generated successfully!"