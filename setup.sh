#!/usr/bin/env bash
# setup.sh — first-run setup for new contributors
# Usage: ./setup.sh
# Run once after cloning the repo.

set -euo pipefail

BLUE='\033[1;34m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
RESET='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ERRORS=0

# ── Banner ─────────────────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}${BLUE}╔═══════════════════════════════════════════════════════════╗${RESET}"
echo -e "${BOLD}${BLUE}║         Fullstack Template — First-Run Setup              ║${RESET}"
echo -e "${BOLD}${BLUE}╚═══════════════════════════════════════════════════════════╝${RESET}"
echo ""

# ── Helper: pass / fail printing ──────────────────────────────────────────────
pass() { echo -e "  ${GREEN}✓${RESET} $*"; }
fail() { echo -e "  ${RED}✗${RESET} $*"; ERRORS=$((ERRORS + 1)); }
info() { echo -e "  ${YELLOW}→${RESET} $*"; }
header() { echo -e "\n${BOLD}$*${RESET}"; }

# ── 1. Prerequisite checks ─────────────────────────────────────────────────────
header "Checking prerequisites..."

# go ≥ 1.25
if command -v go &>/dev/null; then
  GO_VERSION=$(go version | grep -oP '\d+\.\d+' | head -1)
  GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
  GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
  if [[ "$GO_MAJOR" -gt 1 ]] || [[ "$GO_MAJOR" -eq 1 && "$GO_MINOR" -ge 25 ]]; then
    pass "go $GO_VERSION (≥ 1.25)"
  else
    fail "go $GO_VERSION found — need ≥ 1.25. Install from https://go.dev/dl/"
  fi
else
  fail "go not found. Install from https://go.dev/dl/"
fi

# node ≥ 22
if command -v node &>/dev/null; then
  NODE_VERSION=$(node --version | sed 's/v//')
  NODE_MAJOR=$(echo "$NODE_VERSION" | cut -d. -f1)
  if [[ "$NODE_MAJOR" -ge 22 ]]; then
    pass "node v$NODE_VERSION (≥ 22)"
  else
    fail "node v$NODE_VERSION found — need ≥ 22. Install from https://nodejs.org/"
  fi
else
  fail "node not found. Install from https://nodejs.org/"
fi

# pnpm (any version)
if command -v pnpm &>/dev/null; then
  PNPM_VERSION=$(pnpm --version)
  pass "pnpm $PNPM_VERSION"
else
  fail "pnpm not found. Install with: npm install -g pnpm"
fi

# docker (running)
if command -v docker &>/dev/null; then
  if docker info &>/dev/null 2>&1; then
    pass "docker (running)"
  else
    fail "docker is installed but not running. Start Docker Desktop and retry."
  fi
else
  fail "docker not found. Install Docker Desktop from https://www.docker.com/products/docker-desktop/"
fi

# java ≥ 17
if command -v java &>/dev/null; then
  JAVA_VERSION=$(java -version 2>&1 | grep -oP '(?<=version ")\d+' | head -1)
  if [[ "$JAVA_VERSION" -ge 17 ]]; then
    pass "java $JAVA_VERSION (≥ 17)"
  else
    fail "java $JAVA_VERSION found — need ≥ 17. Install from https://adoptium.net/"
  fi
else
  fail "java not found. Install from https://adoptium.net/"
fi

# Android SDK (ANDROID_HOME)
if [[ -n "${ANDROID_HOME:-}" ]] && [[ -d "$ANDROID_HOME" ]]; then
  pass "Android SDK at \$ANDROID_HOME=$ANDROID_HOME"
else
  fail "ANDROID_HOME is not set or points to a missing directory."
  info "Install Android Studio and set ANDROID_HOME to your SDK path."
  info "Typical path: \$HOME/Android/Sdk  (Linux/macOS)"
fi

# Abort if any hard prerequisite failed
if [[ "$ERRORS" -gt 0 ]]; then
  echo ""
  echo -e "${RED}${BOLD}$ERRORS prerequisite(s) failed. Fix the issues above and re-run ./setup.sh.${RESET}"
  echo ""
  exit 1
fi

# ── 2. Install web dependencies ────────────────────────────────────────────────
header "Installing web dependencies (pnpm install)..."
(cd "$SCRIPT_DIR/web" && pnpm install)
pass "web dependencies installed"

# ── 3. Copy env files (if not already present) ────────────────────────────────
header "Copying environment files..."

if [[ -f "$SCRIPT_DIR/backend/.env" ]]; then
  info "backend/.env already exists — skipping"
else
  cp "$SCRIPT_DIR/backend/.env.example" "$SCRIPT_DIR/backend/.env"
  pass "backend/.env.example → backend/.env"
fi

if [[ -f "$SCRIPT_DIR/web/.env.local" ]]; then
  info "web/.env.local already exists — skipping"
else
  cp "$SCRIPT_DIR/web/.env.example" "$SCRIPT_DIR/web/.env.local"
  pass "web/.env.example → web/.env.local"
fi

# ── 4. Start Postgres and run migrations ──────────────────────────────────────
header "Starting Postgres (Docker Compose)..."
(cd "$SCRIPT_DIR/backend" && make docker-run) &
DOCKER_PID=$!
info "Waiting 5s for Postgres to be healthy..."
sleep 5

header "Running database migrations..."
(cd "$SCRIPT_DIR/backend" && make migrate-up)
pass "Migrations applied"

# ── 5. Ready summary ──────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}${GREEN}╔═══════════════════════════════════════════════════════════╗${RESET}"
echo -e "${BOLD}${GREEN}║   You're all set! Run these commands to start developing: ║${RESET}"
echo -e "${BOLD}${GREEN}╚═══════════════════════════════════════════════════════════╝${RESET}"
echo ""
echo -e "  ${BOLD}Quickstart (all services in one terminal):${RESET}"
echo -e "    ${YELLOW}./dev.sh${RESET}"
echo ""
echo -e "  ${BOLD}Or start each service individually:${RESET}"
echo -e "    ${BLUE}cd backend && make docker-run${RESET}   # Postgres"
echo -e "    ${GREEN}cd backend && make watch${RESET}        # Go backend → http://localhost:8080"
echo -e "    ${YELLOW}cd web && pnpm dev${RESET}              # Next.js web → http://localhost:3000"
echo ""
echo -e "  ${BOLD}Edit your env files before starting:${RESET}"
echo -e "    backend/.env    — add DB credentials, Firebase config, etc."
echo -e "    web/.env.local  — add Firebase client config, AUTH_SECRET, etc."
echo ""
