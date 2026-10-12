# Agent Guidelines & Repository Directives

Welcome to **IDXStock**! This file serves as the primary orientation and operating guide for AI agents and pair-programmers working on this repository.

---

## 1. Project Memory & Architecture Reference

All comprehensive, ground-truth architectural documentation is maintained in the [`.context/`](file:///.context/) directory:

- [`.context/PROJECT_CONTEXT.md`](file:///.context/PROJECT_CONTEXT.md) — System purpose, business boundaries, domain glossary, and runtime constraints.
- [`.context/ARCHITECTURE.md`](file:///.context/ARCHITECTURE.md) — System component diagrams, request lifecycles, concurrency models, and fault tolerance.
- [`.context/CODEBASE_MAP.md`](file:///.context/CODEBASE_MAP.md) — Module navigation, responsibilities, key files, and package dependencies.
- [`.context/DECISIONS.md`](file:///.context/DECISIONS.md) — Architecture Decision Records (ADRs) capturing design choices (append-only).
- [`.context/TODO.md`](file:///.context/TODO.md) — Inventory of technical debt, immediate tasks, and refactoring items.

> [!IMPORTANT]
> Always inspect the corresponding documents in [`.context/`](file:///.context/) before implementing architectural modifications or adding new modules.

---

## 2. Technology Stack & Runtime

- **Language / Runtime**: Go 1.26.3
- **Web Framework**: Fiber v3 (`github.com/gofiber/fiber/v3`)
- **Database Driver / Pool**: `pgx/v5` (`github.com/jackc/pgx/v5/pgxpool`)
- **Schema Migrations**: Managed internally via `internal/database/database.go` and `migrations/*.sql`
- **Config Management**: Spring Cloud Config client via `github.com/KAnggara75/scc2go` + Viper
- **Logging**: Logrus JSON formatted logger (`internal/config/logger.go`)

---

## 3. Development Workflow & Commands

### Running Locally
```bash
# Set required configuration variables
export SCC_URL="http://your-config-server:8888"
export SCC_AUTH="Basic <credentials>"
export DATABASE_URL="postgres://user:password@localhost:5432/dbname"

# Run with hot reloading
air

# Or build & execute directly
go build -o idxstock ./cmd/api/main.go
./idxstock
```

### Testing & Validation
```bash
# Run unit and repository tests
go test -v ./...

# Run test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## 4. Coding & Delivery Standards

- **Strict Evidence & Non-Destructive**: Do not assume unconfirmed business logic or delete existing architectural records.
- **Conventional Commits**: Format commit messages according to the Conventional Commits specification (`feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`).
- **Token Efficiency**: Use `rtk` when running shell commands (`rtk git status`, `rtk git diff`) to minimize token consumption.
- **Error Handling**: Propagate errors cleanly with context and log using `logrus` with structured fields.
