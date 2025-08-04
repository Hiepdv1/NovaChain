#!/bin/sh

echo "🔧 Initializing blockchain instance..."
./app init --Address "$WALLET_ADDRESS" --InstanceId "$INSTANCE_ID"

echo "🚀 Starting blockchain node..."
./app startNode --Port "$PORT" --InstanceId "$INSTANCE_ID" $START_NODE_FLAGS