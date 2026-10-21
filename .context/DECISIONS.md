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

---

## ADR-014: Targeted Sector and Industry Taxonomy Synchronization for Unmapped Equities

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Codebase implementation (`internal/routes/routes.go`, `internal/handlers/stock_handler.go`, `internal/usecases/stock_usecase.go`)
- **Context**: Ingested stock entries frequently lack sector, sub-sector, industry, and sub-industry relational mappings (`sector_id`, `sector_sub_id`, `industry_id`, `industry_sub_id`). Querying all 1,100+ stocks sequentially or querying external endpoints indiscriminately creates unnecessary latency and rate-limit risks.
- **Decision**: 
  - Provide a dedicated synchronization endpoint: `PUT /api/v1/stocks/sectors-industries/sync`.
  - Filter target stock codes at the repository level to only retrieve records where any of the four taxonomy keys are `NULL` or empty strings (`TRIM(CAST(... AS TEXT)) = ''`).
  - Run concurrent worker pools (20 workers) to fetch individual stock profiles via `https://www.pasardana.id/api/Stock/GetByCode?code={code}`.
  - Extract and update `fkNewSectorId`, `fkNewSubSectorId`, `fkNewIndustryId`, and `fkNewSubIndustryId` only when values change.
  - Run `go generate ./...` automatically via `.pre-commit-config.yaml` on both `pre-commit` and `pre-push` hooks.
- **Consequences**:
  - Dramatically reduces sync execution time from full iterations to only unmapped equities (~104 stocks).
  - Keeps SQLite relational taxonomy updated cleanly without blocking standard HTTP request threads.
  - Guarantees OpenAPI documentation remains continuously in sync with codebase modifications.

---

## ADR-015: Master Listing Board REST Endpoint and Column Cleanup

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/routes/routes.go`)
- **Context**: The `listing_board` table in `database/master_stocks.db` contained a redundant timestamp column `last_modified` that provided no operational value for static board definitions (`Main`, `Development`, `Acceleration`, `Ekonomi Baru`, `Watchlist`). Additionally, clients required a dedicated REST endpoint to retrieve full master board metadata (`code`, `name`, `name_id`, `description`) distinct from `/api/v1/stocks/boards` which only returns string codes.
- **Decision**:
  - Drop column `last_modified` from table `listing_board` via SQLite schema migration.
  - Introduce `models.ListingBoard` entity and expose `GET /api/v1/stocks/listing-boards` returning all 5 master board metadata entries sorted by `code ASC`.
  - Retain `GET /api/v1/stocks/boards` for backwards compatibility with consumers expecting a flat `[]string`.
- **Consequences**:
  - Simplifies `listing_board` table structure to pure static metadata.
  - Exposes localized Indonesian and English display names along with official descriptive blurbs for frontend UI selectors.

---

## ADR-016: Foreign Key Enforcement & Direct Taxonomic Resolution in Stock Details

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/repositories/stock_repository.go`, `internal/models/stock.go`)
- **Context**: Columns `sector_id`, `sector_sub_id`, `industry_id`, and `industry_sub_id` in the `stocks` table previously lacked explicit Foreign Key constraints to `sector(id)`, `sector_sub(id)`, `industry(id)`, and `industry_sub(id)`. Furthermore, clients requesting single stock details via `GET /api/v1/stocks/:code` or stock listings received raw foreign key integer IDs rather than meaningful, human-readable taxonomy names (`sector`, `sub_sector`, `industry`, `sub_industry`).
- **Decision**:
  - Recreate the `stocks` table schema with explicit `REFERENCES` foreign keys on `sector_id`, `sector_sub_id`, `industry_id`, and `industry_sub_id`. Ensure zero foreign key violations via `PRAGMA foreign_key_check`.
  - Populate missing taxonomy definitions (`industry` and `industry_sub`) via `PUT /api/v1/industries/sync`.
  - Refactor `StockDetailResponse` in `internal/models/stock.go` to return concise named fields: `sector`, `sub_sector`, `industry`, and `sub_industry` (`*string`).
  - Update `GetAllStocks` and `GetStockByCode` queries in `internal/repositories/stock_repository.go` to join `sector`, `sector_sub`, `industry`, and `industry_sub` via `LEFT JOIN`.
- **Consequences**:
  - Enforces strict relational integrity across the SQLite master stocks database.
  - Delivers complete, self-describing JSON responses directly to clients without requiring downstream consumers to perform manual ID lookups.

---

## ADR-017: Brokers Table Schema Simplification to Pure Master Metadata

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/models/broker.go`, `internal/repositories/broker_repository.go`)
- **Context**: The `brokers` table previously retained legacy transactional and activity columns (`investor_type`, `total_value`, `net_value`, `buy_value`, `sell_value`, `total_volume`, `total_frequency`, `created_at`) from earlier broker activity tracking prototypes. Since broker activity tracking was deprecated and removed (ADR-013), `brokers` now functions strictly as a master reference entity for broker codes, names, and groups.
- **Decision**:
  - Migrate table `brokers` in `database/master_stocks.db` to only include columns `code` (TEXT PRIMARY KEY), `name` (TEXT NOT NULL), and `broker_group` (TEXT).
  - Drop columns `investor_type`, `total_value`, `net_value`, `buy_value`, `sell_value`, `total_volume`, `total_frequency`, and `created_at`.
  - Simplify `models.Broker` struct to `Code`, `Name`, and `Group`.
  - Update `BatchInsertBrokers` query and unit tests in `internal/repositories/broker_repository.go` and `internal/repositories/broker_repository_test.go`.
- **Consequences**:
  - Eliminates unused columns and storage footprint from the master database.
  - Aligns database entity definitions directly with application models and domain boundaries.

---

## ADR-018: Elimination of Redundant `last_modified` Timestamps in Taxonomy Tables

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/repositories/industry_repository.go`, `internal/repositories/sector_search_repository.go`)
- **Context**: Taxonomy tables (`industry`, `industry_sub`, `sector`, `sector_sub`) previously contained `last_modified` timestamp columns. Because these taxonomic classifications represent static financial classifications whose timestamps are never consumed by client APIs or cache invalidators, retaining them introduced unnecessary schema overhead and redundant SQL writes on every upsert.
- **Decision**:
  - Drop column `last_modified` from tables `industry`, `industry_sub`, `sector`, and `sector_sub` in `database/master_stocks.db`.
  - Update SQL queries in `internal/repositories/industry_repository.go` and `internal/repositories/sector_search_repository.go` to remove references to `last_modified`.
  - Confirm table `stocks_history` retains its operational `last_modified` column.
- **Consequences**:
  - Simplifies taxonomy table DDL and eliminates unnecessary timestamp writes during upserts.
  - Preserves relational integrity verified via `PRAGMA foreign_key_check`.

---

## ADR-019: Broker Group Normalization to Clean Categorical Literals

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database migration (`database/master_stocks.db`, `internal/services/seeder_service.go`)
- **Context**: The `broker_group` column in the `brokers` master table used verbose string constants prefixed with `BROKER_GROUP_` (e.g. `BROKER_GROUP_FOREIGN`, `BROKER_GROUP_LOCAL`, `BROKER_GROUP_GOVERNMENT`). Clients and internal logic require cleaner, concise categorical literals without technical enum prefixes.
- **Decision**:
  - Normalize existing values in `database/master_stocks.db`:
    - `BROKER_GROUP_FOREIGN` → `FOREIGN`
    - `BROKER_GROUP_LOCAL` → `LOCAL`
    - `BROKER_GROUP_GOVERNMENT` → `GOVERNMENT`
  - Enhance `SeederService.SeedBrokersData` in `internal/services/seeder_service.go` to automatically strip `BROKER_GROUP_` prefixes when seeding data from upstream JSON.
  - Update repository and seeder unit test suites.
- **Consequences**:
  - Standardizes broker categorization values across database and API payloads.
  - Ensures newly seeded broker data adheres to normalized conventions.

---

## ADR-020: Listing Board Column Internationalization and Description Removal

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/models/stock.go`, `internal/repositories/stock_repository.go`)
- **Context**: In table `listing_board`, columns were inconsistently named compared to `sector` and `sector_sub` (`name` contained English, `name_id` contained Indonesian). To align with the existing schema pattern across the repository, the Indonesian name should be `name` and English name should be `name_en`. Additionally, the redundant `description` column provided no operational value for board lookups and was dropped.
- **Decision**:
  - Restructure `listing_board` table schema to: `code` (TEXT PRIMARY KEY), `name` (TEXT NOT NULL), and `name_en` (TEXT).
  - Migrate data: copy previous `name_id` to `name`, copy previous `name` to `name_en`, and drop `description`.
  - Update `models.ListingBoard` struct and `GetListingBoards` repository query.
  - Regenerate Swagger OpenAPI specifications.
- **Consequences**:
  - Unifies multi-language column conventions (`name`, `name_en`) across all reference metadata tables.
  - Streamlines `listing_board` schema to essential identifying attributes.

---

## ADR-021: Code Column Removal from Sector and Sub-Sector Reference Tables

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`, `internal/models/sector_new.go`, `internal/repositories/sector_search_repository.go`)
- **Context**: The `code` column in `sector` and `sector_sub` tables was deprecated by upstream data feeds and is no longer used for securities classification. Primary keys (`id`) and textual names (`name`, `name_en`) provide sufficient taxonomic identity.
- **Decision**:
  - Drop column `code` from tables `sector` and `sector_sub` in `database/master_stocks.db`.
  - Remove `Code` field from `models.SectorNew` and `models.SectorSub`.
  - Adjust SQL `INSERT/UPDATE` queries in `internal/repositories/sector_search_repository.go` and `internal/usecases/sector_usecase.go`.
- **Consequences**:
  - Eliminates unused column storage in classification tables.
  - Retains full foreign key reference integrity (`stocks.sector_id` -> `sector.id`, `stocks.sector_sub_id` -> `sector_sub.id`).

---

## ADR-022: Relational Foreign Key Binding Between `pasardana_stock_mapping` and `stocks`

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`)
- **Context**: Table `pasardana_stock_mapping` lacked a formal database Foreign Key constraint to table `stocks` on column `code`. As a strict constraint, the developer mandated adding the foreign key constraint directly on `pasardana_stock_mapping` without modifying the `stocks` table.
- **Decision**:
  - Recreate table `pasardana_stock_mapping` with column definition `code TEXT NOT NULL UNIQUE REFERENCES stocks(code)`.
  - Purge orphaned non-equity entries (`XSPI`, `GOTOM`) from `pasardana_stock_mapping` that do not exist in `stocks`.
  - Confirm table `stocks` schema remains completely unaltered.
  - Verify relational integrity via `PRAGMA foreign_key_check`.
- **Consequences**:
  - Prevents orphaned stock mapping entries from entering the database.
  - Leaves the `stocks` primary entity untouched as required.

---

## ADR-023: Deletion of Legacy `sector_old` Table

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database cleanup (`database/master_stocks.db`)
- **Context**: Table `sector_old` was an obsolete artifact from earlier migration stages and was neither referenced by foreign keys nor used by any active application code.
- **Decision**:
  - Drop table `sector_old` from `database/master_stocks.db`.
  - Update `.context/ARCHITECTURE.md` table inventory.
- **Consequences**:
  - Removes dead database schema artifacts and prevents developer confusion.

---

## ADR-024: Foreign Key Enforcement Between `stocks_history` and `stocks`

- **Status**: Accepted
- **Date**: 2026-10-21
- **Source**: Developer request & Database schema update (`database/master_stocks.db`)
- **Context**: High-volume historical EOD pricing table `stocks_history` (1.34M+ records) identified equities solely by plain text `code` without a formal database foreign key pointing to `stocks(code)`. Enforcing relational integrity ensures that price histories can only belong to valid master listed equities.
- **Decision**:
  - Recreate table `stocks_history` with `code TEXT NOT NULL REFERENCES stocks(code)` while preserving composite primary key `(code, date)` and B-tree indexes (`idx_stocks_history_code`, `idx_stocks_history_date`, `idx_stocks_history_code_date`).
  - Clean up 347 orphaned records belonging to unlisted symbol `GOTOM` before migrating.
  - Verify zero foreign key violations across all 1.34+ million rows via `PRAGMA foreign_key_check`.
- **Consequences**:
  - Guarantees strict database-level referential integrity between master equity definitions and EOD historical price series.
  - Sub-millisecond composite index query performance remains preserved.
