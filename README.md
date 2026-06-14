# Fullstack Template

A production-ready fullstack starter. Clone it, rename things, and focus on your business logic — the infrastructure is already wired. Ships with a Go + Gin backend, Next.js 16 frontend, PostgreSQL, Docker Compose, hot reload, integration testing, and a full agentic development setup for AI coding assistants.

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
- **Next.js 16 frontend** with React 19, TypeScript, and Tailwind CSS 4
- **PostgreSQL** database managed via Docker Compose
- **Hot reload** on both frontend (`next dev`) and backend ([Air](https://github.com/air-verse/air))
- **Integration tests** using [Testcontainers](https://testcontainers.com/) — no mocks, real DB
- **CORS** pre-configured between frontend and backend
- **`.env` support** via `godotenv`
- **Makefile** for common backend tasks
- **Agentic infrastructure** — AGENTS.md, CLAUDE.md, topic docs, subagents, hooks, and slash commands ready out of the box

## Tech Stack

| Layer     | Technology                          |
|-----------|-------------------------------------|
| Frontend  | Next.js 16, React 19, TypeScript    |
| Styling   | Tailwind CSS 4                      |
| Backend   | Go, Gin                             |
| Database  | PostgreSQL 16 (via Docker)          |
| DB Driver | pgx v5 (standard library style)     |
| Dev Tools | Air (hot reload), pnpm, Docker      |
| Testing   | Go test, Testcontainers             |

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/) and [pnpm](https://pnpm.io/installation)
- [Docker](https://www.docker.com/) and Docker Compose

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

**Frontend:**

```bash
cd frontend
pnpm install
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

**3. Start the frontend** (new terminal):

```bash
cd frontend && pnpm dev
```

Frontend: `http://localhost:3000` — Backend API: `http://localhost:8080`

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
└── frontend/
    ├── AGENTS.md                # Frontend agent instructions
    ├── docs/                    # Topic docs: routing, data-fetching, styling, components
    ├── app/
    │   ├── layout.tsx
    │   ├── page.tsx
    │   └── globals.css
    ├── public/
    ├── next.config.ts
    ├── package.json
    └── tsconfig.json
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

### Frontend commands

```bash
pnpm dev     # start dev server with hot reload
pnpm build   # production build
pnpm start   # serve production build
pnpm lint    # run ESLint
```

## Testing

Backend tests use [Testcontainers](https://testcontainers.com/) to spin up a real PostgreSQL instance — no mocking the database layer.

```bash
cd backend
make test    # unit + integration tests
make itest   # integration tests only
```

Docker must be running for integration tests.

## Working with AI Agents

This template ships with a complete agentic development setup so AI assistants have the context they need to work on your project accurately and consistently.

### For any AI coding agent

A layered `AGENTS.md` system follows the [AGENTS.md open standard](https://agents.md). The closest file to the code you are editing takes precedence:

| File | Covers |
|---|---|
| [`AGENTS.md`](AGENTS.md) | Project overview, setup, cross-cutting conventions, security |
| [`backend/AGENTS.md`](backend/AGENTS.md) | Go commands, project structure, links to topic docs |
| [`frontend/AGENTS.md`](frontend/AGENTS.md) | pnpm commands, Next.js conventions, links to topic docs |

Topic-specific documentation lives in `backend/docs/` and `frontend/docs/`. Each file is kept in sync with the source code it describes and includes `last_verified` metadata so agents can detect when it may be stale.

### For Claude Code

Additional infrastructure in `.claude/` provides a deeper integration:

| Path | Purpose |
|---|---|
| [`CLAUDE.md`](CLAUDE.md) | Feature development workflow, all conventions |
| `.claude/agents/` | Specialized subagents: `backend`, `frontend`, `reviewer`, `db-explorer`, `docs` |
| `.claude/commands/` | Slash commands: `/project:implement`, `/project:check`, `/project:test`, `/project:new-route` |
| `.claude/hooks/` | Auto-formats Go and TypeScript files on save; blocks dangerous commands |

The recommended workflow for any implementation:

1. Check the relevant topic doc in `backend/docs/` or `frontend/docs/` before writing code
2. Implement against documented patterns rather than general training data
3. Update the doc file after implementation so the next agent session starts with accurate context

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
