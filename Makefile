.PHONY: setup dev \
        dev-backend dev-web \
        test test-backend test-web test-e2e \
        build build-backend build-web \
        lint lint-backend lint-web lint-mobile \
        tidy itest-backend

# ── Setup ──────────────────────────────────────────────────────────────────────

setup:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File setup.ps1
else
	bash setup.sh
endif

# ── Development ────────────────────────────────────────────────────────────────

dev:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File dev.ps1
else
	bash dev.sh
endif

dev-backend:
	cd backend && make watch

dev-web:
	cd web && pnpm dev

# ── Testing ────────────────────────────────────────────────────────────────────

test-backend:
	cd backend && go test ./...

test-web:
	cd web && pnpm test:run

test-e2e:
	cd web && pnpm test:e2e

test: test-backend test-web test-e2e

# ── Build ──────────────────────────────────────────────────────────────────────

build-backend:
	cd backend && make build

build-web:
	cd web && pnpm build

build: build-backend build-web

# ── Utilities ──────────────────────────────────────────────────────────────────

itest-backend:
	cd backend && make itest

lint-backend:
	cd backend && go vet ./...

lint-web:
	cd web && pnpm lint

GRADLEW := $(if $(filter Windows_NT,$(OS)),gradlew.bat,./gradlew)

lint-mobile:
	cd mobile && $(GRADLEW) lint

lint: lint-backend lint-web lint-mobile

tidy:
	cd backend && go mod tidy
	cd web && pnpm install
