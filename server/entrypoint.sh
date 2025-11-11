#!/bin/sh
set -e

# -----------------------
# Kiểm tra DATABASE_URL
# -----------------------
if [ -z "$DATABASE_URL" ]; then
  echo "❌ DATABASE_URL is not set"
  exit 1
fi

# -----------------------
# Chạy migrations
# -----------------------
echo "🚀 Running migrations UP..."
migrate -path ./migrations -database "$DATABASE_URL" up

exec "./app"