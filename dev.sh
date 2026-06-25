#!/usr/bin/env bash
# dev.sh — start all services in parallel in a single terminal
# Usage: ./dev.sh
# Press Ctrl+C to stop everything.

set -euo pipefail

BLUE='\033[1;34m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RESET='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Kill all background jobs when the script exits (Ctrl+C or error)
trap 'echo ""; echo "Stopping all services..."; kill $(jobs -p) 2>/dev/null; wait' EXIT

# ── Postgres (Docker Compose) ──────────────────────────────────────────────────
echo -e "${BLUE}[postgres] Starting Docker Compose (Postgres)...${RESET}"
(cd "$SCRIPT_DIR/backend" && make docker-run) &

# ── Backend (Air hot-reload) — wait for Postgres to be ready ──────────────────
(
  echo -e "${GREEN}[backend]  Waiting 5s for Postgres to be healthy...${RESET}"
  sleep 5
  echo -e "${GREEN}[backend]  Starting Go backend with Air (hot-reload) on :8080...${RESET}"
  cd "$SCRIPT_DIR/backend" && make watch
) &

# ── Web (Next.js) ─────────────────────────────────────────────────────────────
echo -e "${YELLOW}[web]      Starting Next.js dev server on :3000...${RESET}"
(cd "$SCRIPT_DIR/web" && pnpm dev) &

# Block until Ctrl+C triggers the trap above
wait
