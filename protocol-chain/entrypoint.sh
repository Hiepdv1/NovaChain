#!/bin/sh
set -e

# Nếu có tham số được truyền vào container → chạy command đó thay vì startNode
if [ "$1" != "" ]; then
  echo "⚙️  Running custom command: ./app $@"
  exec ./app "$@" --InstanceId "$INSTANCE_ID"
else
  echo "🔧 Initializing blockchain instance..."
  ./app init --Address "$WALLET_ADDRESS" --InstanceId "$INSTANCE_ID"

  echo "🚀 Starting blockchain node..."
  exec ./app startNode --Port "$PORT" --InstanceId "$INSTANCE_ID" $START_NODE_FLAGS
fi
