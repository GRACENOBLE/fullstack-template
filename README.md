# Fullstack Template

A production-ready fullstack starter. Clone it, rename things, and focus on your business logic — the infrastructure is already wired. Ships with a Go + Gin backend, Next.js 16 web app, Android mobile app (Kotlin + Compose), PostgreSQL, Docker Compose, hot reload, integration testing, and a full agentic development setup for AI coding assistants.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the App](#running-the-app)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [Development](#development)
- [Testing](#testing)
- [Working with AI Agents](#working-with-ai-agents)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Go backend** with Gin, structured into clean `cmd/` and `internal/` layers
- **Next.js 16 web app** with React 19, TypeScript, and Tailwind CSS 4
- **Android mobile app** with Kotlin 2.2, Jetpack Compose BOM 2026.02, and Material3
- **PostgreSQL** database managed via Docker Compose
- **Hot reload** on both web (`next dev`) and backend ([Air](https://github.com/air-verse/air))
- **Integration tests** using [Testcontainers](https://testcontainers.com/) — no mocks, real DB
- **CORS** pre-configured between web and backend
- **`.env` support** via `godotenv`
- **Makefile** for common backend tasks
- **Agentic infrastructure** — AGENTS.md, CLAUDE.md, topic docs, subagents, hooks, and slash commands ready out of the box for all three layers

## Tech Stack

| Layer     | Technology                                          |
|-----------|-----------------------------------------------------|
| Web       | Next.js 16, React 19, TypeScript                    |
| Styling   | Tailwind CSS 4                                      |
| Backend   | Go, Gin                                             |
| Database  | PostgreSQL 16 (via Docker)                          |
| DB Driver | pgx v5 (standard library style)                     |
| Mobile    | Android, Kotlin 2.2, Jetpack Compose, Material3     |
| Dev Tools | Air (hot reload), pnpm, Docker, Gradle 9.4          |
| Testing   | Go test, Testcontainers, JUnit 4, Espresso          |

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/) and [pnpm](https://pnpm.io/installation)
- [Docker](https://www.docker.com/) and Docker Compose
- [Android Studio Meerkat (2024.3+)](https://developer.android.com/studio) with Android SDK API 36 and JDK 17 (for mobile)

### Installation

```bash
git clone https://github.com/your-username/fullstack-template.git
cd fullstack-template
```

**Backend:**

```bash
cd backend
cp .env .env.local   # adjust values as needed
go mod download
```

**Web:**

```bash
cd web
pnpm install
```

**Mobile:**

Open `mobile/` in Android Studio. The Gradle wrapper handles all SDK downloads. Alternatively:

```bash
cd mobile && ./gradlew assembleDebug
```

### Running the App

**1. Start the database:**

```bash
cd backend && make docker-run
```

**2. Start the backend** (new terminal):

```bash
cd backend && make watch   # hot reload via Air
# or: make run             # run once, no reload
```

**3. Start the web app** (new terminal):

```bash
cd web && pnpm dev
```

**4. Run the mobile app** — connect a device or start an emulator, then:

```bash
cd mobile && ./gradlew installDebug
```

Or run directly from Android Studio.

Web: `http://localhost:3000` — Backend API: `http://localhost:8080`

## Project Structure

```
fullstack-template/
├── AGENTS.md                    # AI agent instructions (all agents)
├── CLAUDE.md                    # Claude Code workflow and conventions
├── .claude/
│   ├── agents/                  # Specialized Claude subagents
│   ├── commands/                # Custom slash commands
│   ├── hooks/                   # Auto-format + guard hooks
│   └── settings.json            # Hook configuration
├── backend/
│   ├── AGENTS.md                # Backend agent instructions
│   ├── docs/                    # Topic docs: database, routing, testing, errors, env
│   ├── cmd/api/main.go          # Entry point
│   ├── internal/
│   │   ├── database/
│   │   │   ├── database.go      # DB connection + queries
│   │   │   └── database_test.go # Integration tests
│   │   └── server/
│   │       ├── server.go        # HTTP server setup
│   │       └── routes.go        # Route definitions
│   ├── .air.toml
│   ├── .env
│   ├── docker-compose.yml
│   └── Makefile
├── web/
│   ├── AGENTS.md                # Web agent instructions
│   ├── docs/                    # Topic docs: routing, data-fetching, styling, components
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── globals.css
│   ├── public/
│   ├── next.config.ts
│   ├── package.json
│   └── tsconfig.json
└── mobile/
    ├── AGENTS.md                # Mobile agent instructions
    ├── docs/                    # Topic docs: compose-conventions, architecture, testing
    ├── app/
    │   ├── build.gradle.kts     # App dependencies and build config
    │   └── src/main/java/com/company/template/
    │       ├── MainActivity.kt  # Single entry point
    │       └── ui/theme/        # Color, Theme, Type (Material3)
    ├── gradle/
    │   └── libs.versions.toml   # Version catalog
    ├── build.gradle.kts
    └── settings.gradle.kts
```

## Environment Variables

Copy `backend/.env` and fill in your values. Never commit secrets to source control.

| Variable                | Description                  | Default     |
|-------------------------|------------------------------|-------------|
| `PORT`                  | Backend server port          | `8080`      |
| `APP_ENV`               | Environment (`local`/`prod`) | `local`     |
| `BLUEPRINT_DB_HOST`     | Postgres host                | `localhost` |
| `BLUEPRINT_DB_PORT`     | Postgres port                | `5432`      |
| `BLUEPRINT_DB_DATABASE` | Database name                | `blueprint` |
| `BLUEPRINT_DB_USERNAME` | Database user                | —           |
| `BLUEPRINT_DB_PASSWORD` | Database password            | —           |
| `BLUEPRINT_DB_SCHEMA`   | Postgres schema              | `public`    |

## Development

### Backend Makefile commands

```bash
make build       # compile the binary
make run         # run without hot reload
make watch       # run with Air hot reload
make docker-run  # start the Postgres container
make docker-down # stop the Postgres container
make test        # run all tests
make itest       # run integration tests only
make clean       # remove compiled binary
```

### Web commands

```bash
pnpm dev     # start dev server with hot reload
pnpm build   # production build
pnpm start   # serve production build
pnpm lint    # run ESLint
```

### Mobile Gradle commands

```bash
./gradlew assembleDebug         # compile debug APK
./gradlew installDebug          # build and install on connected device/emulator
./gradlew lint                  # run Android lint
./gradlew test                  # run unit tests (no device needed)
./gradlew connectedAndroidTest  # run instrumented tests (device/emulator required)
./gradlew clean                 # clean build outputs
```

On Windows outside Git Bash, use `.\gradlew.bat` instead of `./gradlew`.

## Testing

Backend tests use [Testcontainers](https://testcontainers.com/) to spin up a real PostgreSQL instance — no mocking the database layer.

```bash
cd backend
make test    # unit + integration tests
make itest   # integration tests only
```

Docker must be running for integration tests.

Mobile has two test tiers:

```bash
cd mobile
./gradlew test                   # unit tests — runs on JVM, no device needed
./gradlew connectedAndroidTest   # instrumented tests — requires emulator or device
```

## Working with AI Agents

This template ships with a complete agentic development setup so AI assistants have the context they need to work on your project accurately and consistently.

### For any AI coding agent

A layered `AGENTS.md` system follows the [AGENTS.md open standard](https://agents.md). The closest file to the code you are editing takes precedence:

| File | Covers |
|---|---|
| [`AGENTS.md`](AGENTS.md) | Project overview, setup, cross-cutting conventions, security |
| [`backend/AGENTS.md`](backend/AGENTS.md) | Go commands, project structure, links to topic docs |
| [`web/AGENTS.md`](web/AGENTS.md) | pnpm commands, Next.js conventions, links to topic docs |
| [`mobile/AGENTS.md`](mobile/AGENTS.md) | Gradle commands, Android conventions, links to topic docs |

Topic-specific documentation lives in `backend/docs/`, `web/docs/`, and `mobile/docs/`. Each file is kept in sync with the source code it describes and includes `last_verified` metadata so agents can detect when it may be stale.

### For Claude Code

Additional infrastructure in `.claude/` provides a deeper integration:

| Path | Purpose |
|---|---|
| [`CLAUDE.md`](CLAUDE.md) | Feature development workflow, all conventions |
| `.claude/agents/` | Specialized subagents: `backend`, `web`, `mobile`, `reviewer`, `db-explorer`, `docs` |
| `.claude/commands/` | Slash commands: `/project:implement`, `/project:check`, `/project:test`, `/project:new-route` |
| `.claude/hooks/` | Auto-formats Go and TypeScript files on save; blocks dangerous commands |

The recommended workflow for any implementation:

1. Check the relevant topic doc in `backend/docs/`, `web/docs/`, or `mobile/docs/` before writing code
2. Implement against documented patterns rather than general training data
3. Update the doc file after implementation so the next agent session starts with accurate context

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
