#!/bin/sh
set -e

if ! docker network inspect blockchain_net >/dev/null 2>&1; then
  echo "🔗 Creating docker network: blockchain_net"
  docker network create blockchain_net
fi

if [ "$1" != "" ]; then
  echo "⚙️  Running custom command: ./app $@"
  exec ./app "$@" --InstanceId "$INSTANCE_ID"
else
  echo "🔧 Initializing blockchain instance..."
  ./app init --Address "$WALLET_ADDRESS" --InstanceId "$INSTANCE_ID"

  echo "🚀 Starting blockchain node..."
  exec ./app startNode --Port "$PORT" --InstanceId "$INSTANCE_ID" $START_NODE_FLAGS
fi
