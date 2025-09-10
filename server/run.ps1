param(
    [string]$cmd = "run",
    [int]$steps = 1
)

$ErrorActionPreference = "Stop"

# -------------------------
# Configs
# -------------------------
$env:POSTGRES_URL = "postgres://postgres:1@localhost:5433/chain?sslmode=disable"
$migratePath = "migrations"
$sqlcConfig = "sqlc.yaml"

# -------------------------
# Method
# -------------------------
function migrateUp {
    Write-Host "🚀 Running migrations UP..." -ForegroundColor Green
    migrate -path $migratePath -database $env:POSTGRES_URL -verbose up
}

function migrateDown {
    param([int]$steps)
    Write-Host "⚠️ Rolling back $steps migration(s)..." -ForegroundColor Yellow
    migrate -path $migratePath -database $env:POSTGRES_URL -verbose down $steps
}

function sqlcGen {
    Write-Host "🛠 Generating SQLC code..." -ForegroundColor Cyan
    sqlc generate -f $sqlcConfig
}

# -------------------------
# Switch command
# -------------------------
switch ($cmd) {
    "run" {
        Write-Host "🚀 Running Go application..." -ForegroundColor Green
        go run ./cmd/api/main.go
    }
    "build" {
        Write-Host "🔨 Building Go application..." -ForegroundColor Yellow
        go build -o bin/app ./cmd/api
    }
    "test" {
        Write-Host "🧪 Running tests..." -ForegroundColor Cyan
        go test ./...
    }
    "mup" {
        migrateUp
    }
    "mdown" {
        migrateDown $steps
    }
    "sqlc" {
        sqlcGen
    }
    default {
        Write-Host "❌ Unknown command: $cmd" -ForegroundColor Red
        Write-Host "Available commands:" -ForegroundColor Yellow
        Write-Host "  run       -> run application"
        Write-Host "  build     -> build application"
        Write-Host "  test      -> run tests"
        Write-Host "  mup       -> migrate up"
        Write-Host "  mdown     -> migrate down <steps> (default 1)"
        Write-Host "  sqlc      -> sqlc generate"
    }
}
