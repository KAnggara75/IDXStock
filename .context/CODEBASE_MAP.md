# Codebase Navigation Map

## `cmd/api`
- **Responsibility**: Application entry point and runtime bootstrap.
- **Entry / Key Files**:
  - `cmd/api/main.go` — Service initialization, package `init()` lifecycle (logger initialization via `config.InitLogger()` and automated snapshot download via `ensureMasterDatabase()`), DB connection, broker seeding, Fiber HTTP web server startup, and Swagger annotation declarations (`//go:generate swag init ...`).
  - `cmd/api/main_test.go` — Unit tests for port binding, Fiber app construction, seeder execution, runtime start, and `ensureMasterDatabase` snapshot retrieval.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `github.com/klauspost/compress/zstd`, `github.com/sirupsen/logrus`, `internal/config`, `internal/database`, `internal/repositories`, `internal/routes`, `internal/services`.
- **Consumers**: Executable container / process entry point.
- **External Integrations**: GitHub Releases (via HTTP for initial `master_stocks.db.zst` snapshot retrieval), SQLite database (`database/master_stocks.db`).
- **Key Notes**: In `init()`, verifies presence of local database at `DATABASE_URL` / `database/master_stocks.db`. If absent, automatically pulls `master_stocks.db.zst` from GitHub releases and decompresses it into place before Fiber boots.

---

## `internal/routes`
- **Responsibility**: HTTP route registration, middleware binding, and dependency injection wiring.
- **Entry / Key Files**:
  - `internal/routes/routes.go` — Instantiates repositories, services, use cases, handlers, mounts Swagger UI at `/swagger/*`, and registers routes on `/health` and `/api/v1/*`.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `github.com/gofiber/contrib/v3/swaggo`, `docs`, `internal/handlers`, `internal/usecases`, `internal/services`, `internal/repositories`, `internal/database`.
- **Consumers**: `cmd/api/main.go`.
- **External Integrations**: None directly.
- **Key Notes**: Pure manual DI setup connecting SQLite connection executor (`database.Pool`) down to REST handlers.

---

## `internal/handlers`
- **Responsibility**: HTTP transport layer, request validation, query param binding, JSON response formatting, and OpenAPI Swagger documentation annotations.
- **Entry / Key Files**:
  - `internal/handlers/stock_handler.go` — Get all stocks (`/api/v1/stocks`), get stock by ticker code (`/api/v1/stocks/:code`), get distinct trading boards (`/api/v1/stocks/boards`), get stocks filtered by board (`/api/v1/stocks/boards/:board`), preview and upload Excel stock lists (`/api/v1/stocks/upload`), sync company details (`/api/v1/stocks/sync`), and fetch delisted stocks (`/api/v1/stocks/delisting/sync`).
  - `internal/handlers/history_handler.go` — Query stock historical prices (`/api/v1/stocks/:code/history`), sync daily EOD history (`/api/v1/stocks/history/sync`), and execute stock split adjustments (`/api/v1/stocks/split`).
  - `internal/handlers/sector_handler.go` — Synchronize sector and sub-sector master data.
  - `internal/handlers/industry_handler.go` — Synchronize industry and sub-industry master data.
  - `internal/handlers/asset_handler.go` — Proxy company logos from external sources with HTTP stream piping and client caching.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `internal/usecases`, `internal/models`.
- **Consumers**: Fiber router in `internal/routes/routes.go`.
- **External Integrations**: None directly (delegated to use cases).
- **Key Notes**: Handlers annotate Swagger specs and return standardized JSON payload responses.

---

## `internal/usecases`
- **Responsibility**: Core application business logic, orchestration of repositories, and external upstream APIs.
- **Entry / Key Files**:
  - `internal/usecases/stock_usecase.go` — Retrieving filtered stocks (`GetAllStocks`), lookup stock by code (`GetStockByCode`), listing distinct trading boards (`GetDistinctBoards`), query stocks by board (`GetStocksByBoard`), parsing Excel lists, validating stock data, and syncing company profiles.
  - `internal/usecases/history_usecase.go` — Processing historical market pricing, applying mathematical ratios for stock splits (`close`, `open_price`, `high`, `low`, and `volume`), and backfilling missing days.
  - `internal/usecases/sector_usecase.go` — Fetching and updating sector taxonomies.
  - `internal/usecases/industry_usecase.go` — Fetching and updating industry classifications.
- **Dependencies**: `internal/repositories`, `internal/services`, `internal/models`.
- **Consumers**: `internal/handlers`.
- **External Integrations**: Orchestrates `PasardanaService` and `IdxService`.
- **Key Notes**: Encapsulates business transaction boundaries.

---

## `internal/services`
- **Responsibility**: Outbound HTTP integrations with external financial market data providers and file parsing.
- **Entry / Key Files**:
  - `internal/services/idx_service.go` — Direct integration with Indonesia Stock Exchange (`idx.co.id`) with custom stealth headers and cookie management to bypass anti-scraping / 403 filters.
  - `internal/services/pasardana_service.go` — Integration with `pasardana.id` REST endpoints (stocks, sectors, historical quotes).
  - `internal/services/stock_service.go` — Parsing Microsoft Excel (`.xlsx`) files containing official IDX listing reports using `excelize/v2`.
  - `internal/services/seeder_service.go` — Seeding master broker codes and broker names into SQLite on startup.
- **Dependencies**: `net/http`, `github.com/xuri/excelize/v2`, `github.com/sirupsen/logrus`.
- **Consumers**: `internal/usecases`, `cmd/api/main.go`.
- **External Integrations**: `idx.co.id`, `pasardana.id`.
- **Key Notes**: Implements resilient HTTP clients with timeout controls and custom response decoders.

---

## `internal/repositories`
- **Responsibility**: Persistence layer, executing SQL queries, transactional batches, and data operations on SQLite.
- **Entry / Key Files**:
  - `internal/repositories/stock_repository.go` — Master stock retrieval (`GetAllStocks`), stock detail by code (`GetStockByCode`), distinct boards retrieval (`GetDistinctBoards`), CRUD, and batch upserts for company profiles, board types, and delisting statuses.
  - `internal/repositories/history_repository.go` — Historical EOD price querying, bulk insert, and SQL-based stock split adjustments.
  - `internal/repositories/broker_repository.go` — Query and upsert broker master records.
  - `internal/repositories/sector_search_repository.go` — Upsert sector & sector_sub records.
  - `internal/repositories/industry_repository.go` — Upsert industry & industry_sub records.
  - `internal/repositories/db_executor.go` — Database executor type alias.
- **Dependencies**: `modernc.org/sqlite`, `internal/database`, `internal/models`.
- **Consumers**: `internal/usecases`, `internal/services/seeder_service`.
- **External Integrations**: SQLite 3 database engine.
- **Key Notes**: Relies on `database.DBExecutor` interface to allow seamless transaction execution and unit testing via `go-sqlmock`.

---

## `internal/models`
- **Responsibility**: Data transfer objects (DTOs), domain models, and Swagger API payload contracts.
- **Entry / Key Files**:
  - `internal/models/stock.go` — Stock profile entity, `GetAllStocksFilter`, `StockDetailResponse`, and Excel upload schemas.
  - `internal/models/history.go` — Stock price history entity, query filters, and stock split request models.
  - `internal/models/broker.go` — Broker master entity.
  - `internal/models/sector_new.go` & `internal/models/industry.go` — Classification hierarchy models.
- **Dependencies**: Standard library.
- **Consumers**: All layers across `internal/*`.
- **External Integrations**: None.
- **Key Notes**: Shared core types with JSON and query struct tags.

---

## `internal/database` & `internal/config`
- **Responsibility**: Database connection management and structured logging.
- **Entry / Key Files**:
  - `internal/database/database.go` — Connects to SQLite (`modernc.org/sqlite`) using DSN from `DATABASE_URL` (fallback `database/master_stocks.db`), defines `DBExecutor`, and manages the thread-safe connection pool.
  - `internal/config/logger.go` — Configures `logrus` JSON logging and log levels based on `LOG_LEVEL` environment variable.
- **Dependencies**: `modernc.org/sqlite`, `github.com/sirupsen/logrus`.
- **Consumers**: `cmd/api/main.go`, `internal/routes/routes.go`.
- **External Integrations**: SQLite.
- **Key Notes**: Pure Go CGO-free database driver with connection pooling. Configuration is loaded directly via environment variables (`DATABASE_URL`, `LOG_LEVEL`).

---

## `scripts`
- **Responsibility**: Build, dataset management, and automated release utilities.
- **Entry / Key Files**:
  - `scripts/create_history_draft_releases.sh` — Compresses 83 monthly SQLite historical databases (`database/history/YYYY-MM.db`) using `zstd -19` and creates GitHub Draft Releases (`vYYYY-MM`) with attached `YYYY-MM.db.zst` assets.
  - `build_local.sh` — Local container image build script using Podman.
