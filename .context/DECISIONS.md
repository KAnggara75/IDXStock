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

- **Status**: Deprecated & Removed (Superseded by ADR-013)
- **Date**: 2026-10-12
- **Source**: Codebase evidence (`migrations/006_create_broker_activity_table.sql`, `internal/repositories/broker_activity_repository.go`)
- **Context**: Broker transactions (`broker_activity`) accumulate tens of thousands of rows daily across dozens of broker codes, causing performance bottlenecks and index bloat over time.
- **Decision**: Partition `idxstock.broker_activity` by range on `date` into weekly ISO partitions (`broker_activity_YYYY_WW`) and provide an administrative maintenance endpoint (`PUT /api/v1/partition/broker-activity`) to provision upcoming partitions ahead of time.
- **Consequences**:
  - Enables efficient query pruning by date and easy retention management.
  - Requires automated or periodic execution of partition creation before new weekly data arrives.

---

## ADR-003: Centralize Remote Configuration via Spring Cloud Config (scc2go)

- **Status**: Superseded by ADR-011
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

---

## ADR-008: Deprecate and Remove Schema Migrations in Favor of Pre-Packaged Master Database

- **Status**: Accepted
- **Date**: 2026-10-19
- **Source**: Developer request & Codebase refactor
- **Context**: With the adoption of SQLite and the distribution of pre-populated master and historical databases (`database/master_stocks.db`), runtime SQL migrations (`migrations/*.sql`) became obsolete, redundant, and introduced unnecessary startup file inspection and execution overhead.
- **Decision**: 
  - Delete `migrations/*.sql` and remove the `migrations/` directory entirely.
  - Remove `Migrate()` and migration table tracking logic (`schema_migrations`) from `internal/database/database.go`.
  - Clean up build and container configurations (`.air.toml`, `Dockerfile`, `.github/workflows/pipeline.yaml`, `.github/CODEOWNERS`) to eliminate references to migrations.
- **Consequences**:
  - Faster, cleaner application bootstrap without filesystem scanning for SQL files.
  - Reduced codebase footprint and simplified operational maintenance.
  - Future schema evolution will be handled directly in the master database schema or through explicit database update tools.

---

## ADR-009: Pasardana Stock ID Converter Mapping Table

- **Status**: Accepted
- **Date**: 2026-10-19
- **Source**: Developer request & Database schema update (`database/master_stocks.db`)
- **Context**: Upstream Pasardana endpoints (`pasardana.id`) utilize internal integer identifiers (`Id`) to reference listed equities and historical records, whereas IDX official standards identify securities primarily by ticker symbol (`code`). A deterministic lookup table was needed inside the master database to convert between Pasardana stock IDs and ticker symbols without repeated outbound API calls.
- **Decision**: 
  - Create table `pasardana_stock_mapping` in `database/master_stocks.db` with columns `id` (INTEGER PRIMARY KEY), `code` (TEXT NOT NULL UNIQUE), `name` (TEXT NOT NULL), and `last_modified` (TEXT NOT NULL DEFAULT datetime('now')).
  - Add index `idx_pasardana_code` on `code` for constant-time ticker-to-ID resolution.
  - Ingest and populate all 979 equity mappings directly from `https://www.pasardana.id/api/Stock/GetAllSimpleStocks`.
- **Consequences**:
  - Provides instant local mapping between Pasardana IDs and ticker codes for background synchronization jobs.
  - Prevents redundant network calls to `GetAllSimpleStocks` during historical price fetching and profile sync.

---

## ADR-010: Consolidation of Master Stocks & Historical Price Dataset via Zstandard Release Asset

- **Status**: Accepted
- **Date**: 2026-10-20
- **Source**: Developer request & Database consolidation (`database/master_stocks.db`, GitHub Release `v2026-10`)
- **Context**: Historical EOD price records (1.34+ million rows spanning 2019 to 2026) were previously fragmented across 83 monthly SQLite databases (`database/history/YYYY-MM.db`). Querying across multiple months required managing multiple SQLite connections or external file switching. Furthermore, storing the compiled master SQLite database (~437 MB uncompressed) directly in Git would severely bloat repository clone size and breach version control best practices.
- **Decision**: 
  - Consolidate all historical EOD prices into a unified table `stocks_history` inside `database/master_stocks.db`.
  - Implement high-performance B-tree indexes (`idx_stocks_history_code_date` composite, `idx_stocks_history_date`, `idx_stocks_history_code`) to guarantee sub-millisecond query execution over 1.34M rows.
  - Exclude raw binary database files (`database/master_stocks.db`, `database/master_stocks.db-shm`, `database/master_stocks.db-wal`, and `database/master_stocks.db.zst`) from Git tracking via `.gitignore`.
  - Compress the unified master database with high-ratio Zstandard compression (`zstd -19`) producing `database/master_stocks.db.zst` (~58 MB, ~86% reduction).
  - Clean up fragmented monthly draft releases and publish a consolidated draft release (`v2026-10`) tagged as `--latest` with `master_stocks.db.zst` attached as the primary distribution asset.
- **Consequences**:
  - Provides a single, unified database schema for master company metadata, taxonomic classifications, broker data, and comprehensive historical price series.
  - Eliminates large binary artifacts from git history while providing automated, versioned release distribution on GitHub.
  - Consumers download and decompress a single portable artifact (`master_stocks.db.zst`) for immediate queryability.

---

## ADR-011: Direct Environment Configuration & Complete Viper Removal

- **Status**: Accepted
- **Date**: 2026-10-20
- **Source**: Developer request & Codebase refactor (`internal/database/database.go`, `go.mod`)
- **Context**: The application previously used `github.com/spf13/viper` (and legacy `scc2go`) for reading database DSN and configuration parameters. Since configuration needs were simplified to standard environment variables (`DATABASE_URL`, `PORT`, `LOG_LEVEL`, `MASTER_STOCKS_DOWNLOAD_URL`), retaining Viper introduced unnecessary transitive dependencies and binary bloat without operational benefit.
- **Decision**: 
  - Completely remove `github.com/spf13/viper` from `internal/database/database.go` and `internal/database/database_test.go`.
  - Read configuration values directly via `os.Getenv(...)` with safe local defaults (e.g. `DATABASE_URL` defaulting to `database/master_stocks.db`).
  - Purge `github.com/spf13/viper` and its unused transitive packages from `go.mod` and `go.sum` via `go mod tidy`.
- **Consequences**:
  - Drastically trimmed `go.mod` and `go.sum` dependencies, reducing compile times and attack surfaces.
  - Configuration logic is explicit, transparent, and idiomatically aligned with standard Go patterns.
  - Eliminates runtime dependency on external Spring Cloud Config server.

---

## ADR-012: Startup Initialization & Automated Database Snapshot Provisioning in `func init()`

- **Status**: Accepted
- **Date**: 2026-10-20
- **Source**: Developer request & Codebase implementation (`cmd/api/main.go`)
- **Context**: Newly cloned or containerized deployments lack the pre-populated SQLite database (`database/master_stocks.db`) because large database files are excluded from git. Previously, bootstrapping required manually executing scripts or downloading assets before starting the service. Furthermore, database download checks belonged in early process lifecycle initialization before application components boot.
- **Decision**: 
  - Implement process lifecycle initialization inside `func init()` in `cmd/api/main.go`:
    1. Initialize structured logging (`config.InitLogger()`).
    2. Check if local database (`DATABASE_URL` or `database/master_stocks.db`) exists via `ensureMasterDatabase()`.
  - If missing (and not in-memory mode `:memory:`), automatically download `master_stocks.db.zst` from GitHub Releases (`https://github.com/KAnggara75/IDXStock/releases/latest/download/master_stocks.db.zst` or `MASTER_STOCKS_DOWNLOAD_URL`).
  - Decompress the snapshot stream on-the-fly using `github.com/klauspost/compress/zstd` into a temporary file (`.tmp`) and finalize atomically via `os.Rename`.
- **Consequences**:
  - Zero-touch bootstrap for local developers and container deployments; running `go run ./cmd/api` automatically fetches the latest consolidated dataset if absent.
  - Ensures atomic database file creation via temporary file renaming, preventing partial/corrupted database reads.
  - Clean separation of concerns keeping `internal/database` pure and free of circular dependencies.

---

## ADR-013: Deprecation and Complete Removal of Broker Activity Tracking

- **Status**: Accepted
- **Date**: 2026-10-20
- **Source**: Developer request ("hapus broker_activity") & Codebase refactor
- **Context**: The `broker_activity` tracking mechanism ingested external high-volume broker transaction summaries and managed weekly table partitions (`broker_activity_p_YYYY_WW`). It introduced substantial maintenance overhead, complex partition logic, and was no longer needed for the core equities and taxonomy service. The master broker metadata (`brokers` table seeded via `seederService`) is retained for reference.
- **Decision**: 
  - Drop the `broker_activity` table from `database/master_stocks.db` and execute `VACUUM`.
  - Remove all broker activity REST routes (`GET /api/v1/broker/sync` and `PUT /api/v1/partition/broker-activity`).
  - Delete broker activity handlers, use cases, services, repositories, and models (`broker_handler.go`, `broker_usecase.go`, `broker_service.go`, `broker_activity_repository.go`, `broker_activity.go`).
  - Retain `brokers` master table, `BrokerRepository`, and `SeederService.SeedBrokersData`.
  - Recompress `database/master_stocks.db` to `database/master_stocks.db.zst` (`zstd -19`) and update the GitHub Release `v2026-10` asset.
  - Regenerate Swagger OpenAPI specifications (`docs/`).
- **Consequences**:
  - Eliminates unnecessary upstream dependency on Exodus broker activity endpoints.
  - Reduces codebase size, eliminates table partition complexity, and simplifies database maintenance.
  - Releases distributed database assets without bloat.
