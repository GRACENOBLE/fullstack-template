# setup.ps1 — first-run setup for new contributors (Windows / PowerShell)
# Usage: .\setup.ps1
# Run once after cloning the repo.

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$Errors = 0

# ── Banner ─────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║         Fullstack Template — First-Run Setup              ║" -ForegroundColor Blue
Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# ── Helpers ────────────────────────────────────────────────────────────────────
function Pass   { param($msg) Write-Host "  ✓ $msg" -ForegroundColor Green }
function Fail   { param($msg) Write-Host "  ✗ $msg" -ForegroundColor Red; $script:Errors++ }
function Info   { param($msg) Write-Host "  → $msg" -ForegroundColor Yellow }
function Header { param($msg) Write-Host "`n$msg" -ForegroundColor White }

# ── 1. Prerequisite checks ─────────────────────────────────────────────────────
Header "Checking prerequisites..."

# go ≥ 1.25
try {
    $goVersion = (go version 2>&1) -replace '.*go(\d+\.\d+).*','$1'
    $parts = $goVersion -split '\.'
    $major = [int]$parts[0]; $minor = [int]$parts[1]
    if ($major -gt 1 -or ($major -eq 1 -and $minor -ge 25)) {
        Pass "go $goVersion (≥ 1.25)"
    } else {
        Fail "go $goVersion found — need ≥ 1.25. Install from https://go.dev/dl/"
    }
} catch {
    Fail "go not found. Install from https://go.dev/dl/"
}

# node ≥ 22
try {
    $nodeVersion = (node --version 2>&1).TrimStart('v')
    $nodeMajor = [int]($nodeVersion -split '\.')[0]
    if ($nodeMajor -ge 22) {
        Pass "node v$nodeVersion (≥ 22)"
    } else {
        Fail "node v$nodeVersion found — need ≥ 22. Install from https://nodejs.org/"
    }
} catch {
    Fail "node not found. Install from https://nodejs.org/"
}

# pnpm (any version)
try {
    $pnpmVersion = (pnpm --version 2>&1)
    Pass "pnpm $pnpmVersion"
} catch {
    Fail "pnpm not found. Install with: npm install -g pnpm"
}

# docker (running)
try {
    $null = docker info 2>&1
    if ($LASTEXITCODE -eq 0) {
        Pass "docker (running)"
    } else {
        Fail "docker is installed but not running. Start Docker Desktop and retry."
    }
} catch {
    Fail "docker not found. Install Docker Desktop from https://www.docker.com/products/docker-desktop/"
}

# java ≥ 17
try {
    $javaOut = (java -version 2>&1) | Select-Object -First 1
    if ($javaOut -match '"(\d+)') {
        $javaMajor = [int]$Matches[1]
        if ($javaMajor -ge 17) {
            Pass "java $javaMajor (≥ 17)"
        } else {
            Fail "java $javaMajor found — need ≥ 17. Install from https://adoptium.net/"
        }
    } else {
        Fail "Could not parse java version. Install from https://adoptium.net/"
    }
} catch {
    Fail "java not found. Install from https://adoptium.net/"
}

# Android SDK (ANDROID_HOME or well-known default)
$androidHome = $env:ANDROID_HOME
if (-not $androidHome -or -not (Test-Path $androidHome)) {
    $androidHome = "$env:LOCALAPPDATA\Android\Sdk"
}
if (Test-Path $androidHome) {
    Pass "Android SDK at $androidHome"
    if (-not $env:ANDROID_HOME) {
        Info "ANDROID_HOME is not set as an environment variable — tools like the Gradle wrapper will still work, but you may want to add it:"
        Info "  [System.Environment]::SetEnvironmentVariable('ANDROID_HOME', '$androidHome', 'User')"
    }
} else {
    Fail "Android SDK not found at ANDROID_HOME or $env:LOCALAPPDATA\Android\Sdk"
    Info  "Install Android Studio — it sets up the SDK automatically."
    Info  "Then optionally set ANDROID_HOME to: $env:LOCALAPPDATA\Android\Sdk"
}

# Abort on failures
if ($Errors -gt 0) {
    Write-Host ""
    Write-Host "$Errors prerequisite(s) failed. Fix the issues above and re-run .\setup.ps1." -ForegroundColor Red
    Write-Host ""
    exit 1
}

# ── 2. Install web dependencies ────────────────────────────────────────────────
Header "Installing web dependencies (pnpm install)..."
Set-Location "$ScriptDir\web"
pnpm install
Set-Location $ScriptDir
Pass "web dependencies installed"

# ── 3. Copy env files (if not already present) ────────────────────────────────
Header "Copying environment files..."

$backendEnv = "$ScriptDir\backend\.env"
$backendExample = "$ScriptDir\backend\.env.example"
if (Test-Path $backendEnv) {
    Info "backend\.env already exists — skipping"
} else {
    Copy-Item $backendExample $backendEnv
    Pass "backend\.env.example → backend\.env"
}

$webEnvLocal = "$ScriptDir\web\.env.local"
$webExample  = "$ScriptDir\web\.env.example"
if (Test-Path $webEnvLocal) {
    Info "web\.env.local already exists — skipping"
} else {
    Copy-Item $webExample $webEnvLocal
    Pass "web\.env.example → web\.env.local"
}

# ── 4. Start Postgres and run migrations ──────────────────────────────────────
Header "Starting Postgres (Docker Compose) in a background window..."
Start-Process powershell -ArgumentList "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", "Set-Location '$ScriptDir\backend'; make docker-run" `
    -WindowStyle Minimized

Info "Waiting 5s for Postgres to be healthy..."
Start-Sleep 5

Header "Running database migrations..."
Set-Location "$ScriptDir\backend"
make migrate-up
Set-Location $ScriptDir
Pass "Migrations applied"

# ── 5. Ready summary ──────────────────────────────────────────────────────────
Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║   You're all set! Run these commands to start developing: ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "  Quickstart (all services in separate windows):" -ForegroundColor White
Write-Host "    .\dev.ps1" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Or start each service individually:" -ForegroundColor White
Write-Host "    cd backend; make docker-run   # Postgres" -ForegroundColor Blue
Write-Host "    cd backend; make watch        # Go backend → http://localhost:8080" -ForegroundColor Green
Write-Host "    cd web; pnpm dev              # Next.js web → http://localhost:3000" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Edit your env files before starting:" -ForegroundColor White
Write-Host "    backend\.env    — add DB credentials, Firebase config, etc."
Write-Host "    web\.env.local  — add Firebase client config, AUTH_SECRET, etc."
Write-Host ""
