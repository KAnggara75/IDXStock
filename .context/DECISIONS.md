# Architecture Decision Records (ADRs)

## ADR-001: Adopt Fiber v3 Web Framework for High-Throughput REST API

- **Status**: Accepted
- **Date**: 2026-10-12
- **Source**: Codebase evidence (`go.mod`, `cmd/api/main.go`, `internal/routes/routes.go`)
- **Context**: The service requires high-performance HTTP routing, JSON serialization, and memory-efficient request parsing for high-volume market quote data and Excel file uploads.
- **Decision**: Use `github.com/gofiber/fiber/v3` as the HTTP application framework, configured with recovery middleware and structured logging.
- **Consequences**:
  - Leverages Fasthttp's zero-allocation principles and high-throughput concurrency.
  - Requires maintaining Fiber v3-compatible handler signatures and middlewares.

---

## ADR-002: PostgreSQL ISO-Week Partitioning for Broker Activity Records

- **Status**: Accepted
- **Date**: 2026-10-12
- **Source**: Codebase evidence (`migrations/006_create_broker_activity_table.sql`, `internal/repositories/broker_activity_repository.go`)
- **Context**: Broker transactions (`broker_activity`) accumulate tens of thousands of rows daily across dozens of broker codes, causing performance bottlenecks and index bloat over time.
- **Decision**: Partition `idxstock.broker_activity` by range on `date` into weekly ISO partitions (`broker_activity_YYYY_WW`) and provide an administrative maintenance endpoint (`PUT /api/v1/partition/broker-activity`) to provision upcoming partitions ahead of time.
- **Consequences**:
  - Enables efficient query pruning by date and easy retention management.
  - Requires automated or periodic execution of partition creation before new weekly data arrives.

---

## ADR-003: Centralize Remote Configuration via Spring Cloud Config (scc2go)

- **Status**: Accepted
- **Date**: 2026-10-12
- **Source**: Codebase evidence (`cmd/api/main.go#L20-L22`, `go.mod`)
- **Context**: Application deployments need decoupled configuration without hardcoding sensitive database credentials and API endpoints inside local repositories or image builds.
- **Decision**: Initialize remote configuration retrieval using `github.com/KAnggara75/scc2go` during package `init()`, fetching settings dynamically based on `IDXSTOCK_CONFIG_URL` and `SCC_AUTH`.
- **Consequences**:
  - Centralizes credential and environment management across microservices.
  - Requires the Spring Cloud Config service to be reachable at boot time.

---

## ADR-004: Standardize Local Quality Gates & Pre-Commit Hook Architecture

- **Status**: Accepted
- **Date**: 2026-10-13
- **Source**: Codebase evidence (`.pre-commit-config.yaml`, git hooks setup)
- **Context**: The repository needs automated pre-commit and pre-push quality gates that enforce local import ordering, code cleanliness, build sanity, and vulnerability scanning without introducing irrelevant tool overhead or build failures.
- **Decision**: 
  - Standardize pre-commit with `pre-commit-hooks` (whitespace, EOF, yaml, large files, merge conflicts) and `tekwizely/pre-commit-golang` configured specifically with `-local github.com/KAnggara75/IDXStock`.
  - Shift heavier validation checks (`gotest`, `gobuild`, `govulncheck`) to the `pre-push` stage using native system executables.
  - Eliminate dead framework hooks (e.g. `zerolog-use-stringer`) since the project standardizes on `logrus`.
- **Consequences**:
  - Ensures clean, consistent import grouping across the entire codebase prior to commit.
  - Prevents broken builds and known vulnerabilities from reaching remote branches while keeping commit latency near-instant.

---

## ADR-005: Pure-Go SQLite Driver & Schema Migration

- **Status**: Accepted
- **Date**: 2026-10-19
- **Source**: Codebase evidence (`internal/database/database.go`, `migrations/*.sql`, PR #51)
- **Context**: The service previously required an external PostgreSQL database container/instance (`pgx/v5`) with weekly table partitioning, increasing operational complexity and deployment footprint for embedded, local, and portable setups.
- **Decision**: 
  - Migrate the primary database engine to SQLite using `modernc.org/sqlite` (pure Go implementation, zero CGO dependency).
  - Adopt file-based persistence defaulting to `database/master_stocks.db` or environment `DATABASE_URL`.
  - Abstract SQL operations through `database.DBExecutor` interface wrapping standard `*sql.DB` / `*sql.Tx` to preserve transaction atomicity and testability via `go-sqlmock`.
  - Adapt migrations (`migrations/*.sql`) from PostgreSQL syntax (`SERIAL`, `ON CONFLICT (...) DO UPDATE`) to standard SQLite syntax (`INTEGER PRIMARY KEY AUTOINCREMENT`, `ON CONFLICT (...) DO UPDATE` with `EXCLUDED`).
- **Consequences**:
  - Eliminates external PostgreSQL service dependencies, simplifying local development, CI/CD, and embedded deployments.
  - Retains cross-platform portability without requiring C compilers (CGO-free).
  - Concurrency is bounded by SQLite write-locking semantics; mitigated through connection pooling constraints and fast transactional operations.

---

## ADR-006: Automated OpenAPI Documentation Generation via Swaggo

- **Status**: Accepted
- **Date**: 2026-10-19
- **Source**: Codebase evidence (`cmd/api/main.go`, `docs/*`, `.pre-commit-config.yaml`)
- **Context**: REST API endpoints, DTO contracts, query filters, and Swagger UI lacked an automated, synchronized specification, risking documentation drift as endpoints evolve.
- **Decision**: 
  - Integrate `github.com/swaggo/swag` using declarative code annotations on API entry points and HTTP handlers.
  - Add `//go:generate swag init ...` directive in `cmd/api/main.go` generating `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`.
  - Serve interactive Swagger UI directly via `github.com/gofiber/contrib/v3/swaggo` at `/swagger/*`.
  - Enforce automated execution of `go generate ./...` in the `pre-push` git quality gate hook to ensure documentation is always in lockstep with handler definitions.
- **Consequences**:
  - Automatically provides interactive API documentation without manual YAML maintenance.
  - Ensures changes to API models and parameters are immediately reflected in generated Swagger assets before push.

---

## ADR-007: Monthly Historical Market Data Distribution via Zstandard-Compressed GitHub Releases

- **Status**: Accepted
- **Date**: 2026-10-19
- **Source**: Codebase evidence (`database/history/`, `scripts/create_history_draft_releases.sh`, GitHub Releases)
- **Context**: End-of-day historical stock price datasets spanning 2019 to 2026 accumulate large storage volume (~405 MB across 83 monthly SQLite databases). Storing these directly in git bloats repository clone size.
- **Decision**: 
  - Partition historical stock database by calendar month into standalone SQLite files (`database/history/YYYY-MM.db`).
  - Compress each monthly database using high-ratio Zstandard compression (`zstd -19`) producing `YYYY-MM.db.zst` (~86% compression ratio, reducing total payload from ~405 MB to ~56 MB).
  - Automate draft release creation on GitHub via `scripts/create_history_draft_releases.sh` utilizing GitHub CLI (`gh release create vYYYY-MM YYYY-MM.db.zst --draft`).
  - Add `database/history/*.db*` to `.gitignore` to prevent tracking binary artifacts in the git tree.
- **Consequences**:
  - Dramatically reduces storage requirements and bandwidth consumption for consumer distribution.
  - Keeps the core repository lightweight while providing versioned monthly release assets for historical backtesting and analysis.
