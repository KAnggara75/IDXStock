# Codebase Navigation Map

## `cmd/api`
- **Responsibility**: Application entry point and runtime bootstrap.
- **Entry / Key Files**:
  - `cmd/api/main.go` — Service initialization, configuration loading (`scc2go`), logger setup, DB connection with auto-migrations, broker seeding, and Fiber HTTP web server startup.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `github.com/KAnggara75/scc2go`, `github.com/sirupsen/logrus`, `internal/config`, `internal/database`, `internal/repositories`, `internal/routes`, `internal/services`.
- **Consumers**: Executable container / process entry point.
- **External Integrations**: Spring Cloud Config server (via `scc2go`), PostgreSQL database.
- **Key Notes**: Reads `SCC_URL` and `SCC_AUTH` during package `init()` to inject remote configurations into Viper before main runtime boots.

---

## `internal/routes`
- **Responsibility**: HTTP route registration, middleware binding, and dependency injection wiring.
- **Entry / Key Files**:
  - `internal/routes/routes.go` — Instantiates repositories, services, use cases, handlers, and mounts routes onto `/health` and `/api/v1/*`.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `internal/handlers`, `internal/usecases`, `internal/services`, `internal/repositories`, `internal/database`.
- **Consumers**: `cmd/api/main.go`.
- **External Integrations**: None directly.
- **Key Notes**: Pure manual DI setup connecting PostgreSQL connection pool (`database.Pool`) down to REST handlers.

---

## `internal/handlers`
- **Responsibility**: HTTP transport layer, request validation, query param binding, and JSON response formatting.
- **Entry / Key Files**:
  - `internal/handlers/stock_handler.go` — Preview and upload Excel stock lists, sync company details, and fetch delisted stocks.
  - `internal/handlers/history_handler.go` — Query stock historical prices (`/api/v1/stocks/:code/history`), sync daily EOD history, and execute stock split adjustments.
  - `internal/handlers/broker_handler.go` — Sync broker activities and maintain weekly PostgreSQL partitions for `broker_activities`.
  - `internal/handlers/sector_handler.go` — Synchronize sector and sub-sector master data.
  - `internal/handlers/industry_handler.go` — Synchronize industry and sub-industry master data.
  - `internal/handlers/asset_handler.go` — Proxy company logos from external sources with HTTP stream piping and client caching.
- **Dependencies**: `github.com/gofiber/fiber/v3`, `internal/usecases`, `internal/models`.
- **Consumers**: Fiber router in `internal/routes/routes.go`.
- **External Integrations**: None directly (delegated to use cases).
- **Key Notes**: Handlers strictly handle HTTP-level concerns and return standardized JSON payload responses.

---

## `internal/usecases`
- **Responsibility**: Core application business logic, orchestration of repositories, and external upstream APIs.
- **Entry / Key Files**:
  - `internal/usecases/stock_usecase.go` — Parsing Excel lists, validating stock data, and syncing company profiles.
  - `internal/usecases/history_usecase.go` — Processing historical market pricing, applying mathematical ratios for stock splits (`close_price`, `open_price`, `high_price`, `low_price`, and `volume`), and backfilling missing days.
  - `internal/usecases/broker_usecase.go` — Orchestrates broker activity synchronization and dynamic partition provisioning.
  - `internal/usecases/sector_usecase.go` — Fetching and updating sector taxonomies.
  - `internal/usecases/industry_usecase.go` — Fetching and updating industry classifications.
- **Dependencies**: `internal/repositories`, `internal/services`, `internal/models`.
- **Consumers**: `internal/handlers`.
- **External Integrations**: Orchestrates `PasardanaService`, `IdxService`, and `BrokerService`.
- **Key Notes**: Business transaction boundaries are encapsulated here.

---

## `internal/services`
- **Responsibility**: Outbound HTTP integrations with external financial market data providers and file parsing.
- **Entry / Key Files**:
  - `internal/services/idx_service.go` — Direct integration with Indonesia Stock Exchange (`idx.co.id`) with custom stealth headers and cookie management to bypass anti-scraping / 403 filters.
  - `internal/services/pasardana_service.go` — Integration with `pasardana.id` REST endpoints (stocks, sectors, historical quotes).
  - `internal/services/broker_service.go` — Fetching broker summary transactions from external market data endpoints.
  - `internal/services/stock_service.go` — Parsing Microsoft Excel (`.xlsx`) files containing official IDX listing reports using `excelize/v2`.
  - `internal/services/seeder_service.go` — Seeding master broker codes and broker names into PostgreSQL on startup.
- **Dependencies**: `net/http`, `github.com/xuri/excelize/v2`, `github.com/sirupsen/logrus`.
- **Consumers**: `internal/usecases`, `cmd/api/main.go`.
- **External Integrations**: `idx.co.id`, `pasardana.id`, external broker API.
- **Key Notes**: Implements resilient HTTP clients with timeout controls and custom response decoders (including gzip compression support).

---

## `internal/repositories`
- **Responsibility**: Persistence layer, executing SQL queries, transactional batches, and partition management on PostgreSQL.
- **Entry / Key Files**:
  - `internal/repositories/stock_repository.go` — CRUD and batch upserts for company profiles, board types, and delisting statuses.
  - `internal/repositories/history_repository.go` — Historical EOD price querying, bulk insert, and SQL-based stock split adjustments.
  - `internal/repositories/broker_activity_repository.go` — Ingestion of broker transactions and dynamic DDL execution for weekly partition tables.
  - `internal/repositories/broker_repository.go` — Query and upsert broker master records.
  - `internal/repositories/sector_search_repository.go` — Upsert sector & sub-sector records.
  - `internal/repositories/industry_repository.go` — Upsert industry & sub-industry records.
- **Dependencies**: `github.com/jackc/pgx/v5/pgxpool`, `internal/models`.
- **Consumers**: `internal/usecases`, `internal/services/seeder_service`.
- **External Integrations**: PostgreSQL (schema `idxstock`).
- **Key Notes**: Uses `pgx/v5` connection pool for high-performance parameterized SQL executions and explicit transaction blocks.

---

## `internal/models`
- **Responsibility**: Data transfer objects (DTOs), domain models, and external API payload contracts.
- **Entry / Key Files**:
  - `internal/models/stock.go` — Stock profile entity and Excel upload schema.
  - `internal/models/history.go` — Stock price history entity, query filters, and stock split request models.
  - `internal/models/broker_activity.go` — Broker transaction entity and sync parameters.
  - `internal/models/broker.go` — Broker master entity.
  - `internal/models/sector_new.go` & `internal/models/industry.go` — Classification hierarchy models.
- **Dependencies**: Standard library.
- **Consumers**: All layers across `internal/*`.
- **External Integrations**: None.
- **Key Notes**: Shared core types with JSON struct tags.

---

## `internal/database` & `internal/config`
- **Responsibility**: Database connection management, schema migration orchestration, and structured logging.
- **Entry / Key Files**:
  - `internal/database/database.go` — Initializes `pgxpool.Pool` from Viper configuration (`db.idxstock.dsn`) and executes automatic schema migrations against `idxstock.schema_migrations`.
  - `internal/config/logger.go` — Configures `logrus` JSON logging and log levels based on `LOG_LEVEL` environment variable.
- **Dependencies**: `github.com/jackc/pgx/v5`, `github.com/spf13/viper`, `github.com/sirupsen/logrus`.
- **Consumers**: `cmd/api/main.go`, `internal/routes/routes.go`.
- **External Integrations**: PostgreSQL.
- **Key Notes**: Includes auto-migration executor applying ordered SQL scripts from `migrations/` at application boot.
