param(
    [string]$cmd = "run",
    [int]$steps = 1,
    [string]$name = ""
)

$ErrorActionPreference = "Stop"

# -------------------------
# Configs
# -------------------------
$env:DATABASE_URL="postgresql://postgres:1@postgres:5432/chain?sslmode=disable"
$migratePath = "migrations"
$sqlcConfig = "sqlc.yaml"

# -------------------------
# Methods
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

function migrateDownAll {
    Write-Host "⚠️ Rolling back ALL migrations..." -ForegroundColor Red
    migrate -path $migratePath -database $env:POSTGRES_URL -verbose down -all
}

function migrateNew {
    param([string]$name)
    if (-not $name) {
        Write-Host "❌ Migration name is required! Use -name <migration_name>" -ForegroundColor Red
        return
    }
    Write-Host "📦 Creating new migration: $name" -ForegroundColor Cyan
    migrate create -ext sql -dir $migratePath -seq $name
}

function sqlcGen {
    Write-Host "🛠 Generating SQLC code..." -ForegroundColor Cyan
    sqlc generate -f $sqlcConfig
}

function goRun {
    Write-Host "🚀 Running Go application..." -ForegroundColor Green
    go run ./cmd/api/main.go
}

function goBuild {
    Write-Host "🔨 Building Go application..." -ForegroundColor Yellow
    go build -o bin/app ./cmd/api
}

function goTest {
    Write-Host "🧪 Running tests..." -ForegroundColor Cyan
    go test ./...
}

function dockerBuild {
    Write-Host "🐳 Building Docker image fiber_api:latest (linux/amd64)..." -ForegroundColor Yellow
    docker build --platform=linux/amd64 -t fiber_api:latest . --progress=plain
}

# -------------------------
# Command dispatcher
# -------------------------
function dispatch($cmd, $steps, $name) {
    switch ($cmd) {
        "run"       { goRun }
        "build"     { goBuild }
        "test"      { goTest }
        "mup"       { migrateUp }
        "mdown"     { migrateDown $steps }
        "mdownall"  { migrateDownAll }
        "mnew"      { migrateNew $name }
        "sqlc"      { sqlcGen }
        "docker"    { dockerBuild }
        default {
            Write-Host "❌ Unknown command: $cmd" -ForegroundColor Red
            Write-Host "Available commands:" -ForegroundColor Yellow
            Write-Host "  run       -> run application"
            Write-Host "  build     -> build application"
            Write-Host "  test      -> run tests"
            Write-Host "  mup       -> migrate up"
            Write-Host "  mdown     -> migrate down <steps> (default 1)"
            Write-Host "  mdownall  -> rollback ALL migrations (drop all tables)"
            Write-Host "  mnew      -> create new migration <-name required>"
            Write-Host "  sqlc      -> sqlc generate"
            Write-Host "  docker    -> build Docker image fiber_api:latest"
        }
    }
}
# -------------------------
# Run
# -------------------------
dispatch $cmd $steps $name
