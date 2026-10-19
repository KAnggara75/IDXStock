# Project Context & Domain Overview

## 1. Project Purpose
IDXStock is a specialized Golang backend microservice engineered to ingest, normalize, persist, and serve Indonesian financial market data (Indonesia Stock Exchange / IDX). It bridges disparate financial data feeds (official IDX endpoints, Pasardana, and listing spreadsheets) into a centralized, queryable SQLite data store with interactive OpenAPI documentation.

---

## 2. System Boundaries
- **In-Scope**:
  - Ingesting and validating official IDX listing reports from Excel files (`.xlsx`).
  - Serving master listed equities via `GET /api/v1/stocks` with board filtering and text search.
  - Serving stock detail by ticker code (`/api/v1/stocks/:code`), distinct trading boards (`/api/v1/stocks/boards`), and stocks filtered by board (`/api/v1/stocks/boards/:board`).
  - Maintaining Pasardana ID to ticker code mapping (`pasardana_stock_mapping`).
  - Synchronizing company profiles, sectors, sub-sectors, industries, and sub-industries.
  - Tracking delisted securities.
  - Ingesting and querying end-of-day (EOD) historical stock prices with support for date range filters and field selections.
  - Adjusting historical price series for corporate actions such as stock splits.
  - Streaming and caching company logo images.
  - Ingesting and persisting daily broker transactions and summaries.
  - Serving interactive OpenAPI/Swagger documentation at `/swagger/*`.
  - Monthly historical price dataset packaging via Zstandard-compressed SQLite releases (`YYYY-MM.db.zst`).
- **Out-of-Scope**:
  - Direct order routing or real-time trading execution.
  - User authentication and multi-tenant authorization (delegated to upstream API gateway/proxy).
  - Portfolio management or predictive financial analysis.

---

## 3. Main Actors & Consumers
- **Automated Schedulers / Cron Jobs**: Invoke periodic synchronization endpoints (`/api/v1/stocks/history/sync`, `/api/v1/broker/sync`, `/api/v1/partition/broker-activity`).
- **Internal Analysts / Developers**: Query master stocks (`/api/v1/stocks`), upload listing spreadsheets (`/api/v1/stocks/upload`), trigger corporate stock splits, and explore API documentation via Swagger UI (`/swagger/*`).
- **Client Applications / Dashboard Services**: Query historical prices (`/api/v1/stocks/:code/history`) and company logos.

---

## 4. Domain Glossary
- **Stock Code / Ticker (`code`)**: 4-letter symbol representing an IDX listed security (e.g., `BBCA`, `BBRI`, `TLKM`).
- **Listing Board (`board`)**: Trading board category on the IDX:
  - `UTAMA` (Main Board)
  - `PENGEMBANGAN` (Development Board)
  - `AKSELERASI` (Acceleration Board)
  - `EKONOMI BARU` (New Economy Board)
  - `PEMANTAUAN KHUSUS` (Watchlist Board)
- **Delisted Stock**: A security that has been removed from active trading on the exchange.
- **Stock Split**: A corporate action where existing shares are divided into multiple shares, requiring historical pricing data (open, high, low, close) to be divided by the split ratio and volume multiplied.
- **Broker Code (`broker_code`)**: 2-letter identifier for exchange broker firms (e.g., `YP`, `CC`, `PD`, `NI`).
- **Broker Activity**: Net buy/sell volume, transaction value, and frequency categorized by investor type (domestic vs foreign).

---

## 5. External Integrations
- **Indonesia Stock Exchange (`idx.co.id`)**: Upstream data source for delisted stocks and daily summaries. Accessed with custom HTTP headers and cookie management.
- **Pasardana (`pasardana.id`)**: Upstream data source for sector, industry, company profile, and quote history feeds.
- **Spring Cloud Config (`scc2go`)**: External configuration server providing database connection strings and environment parameters.
- **External Logo CDNs**: Upstream image servers for company logos streamed via `/api/v1/assets/logos/companies/:code`.
- **GitHub Releases (`gh`)**: Distribution channel for monthly historical SQLite databases (`database/history/YYYY-MM.db.zst`).

---

## 6. Runtime Environment & Constraints
- **Language / Runtime**: Go 1.26.3
- **Web Framework**: Fiber v3 (`github.com/gofiber/fiber/v3`)
- **Database Engine**: Pure Go SQLite 3 via `modernc.org/sqlite` (CGO-free, file-based default: `database/master_stocks.db`).
- **Containerization**: Alpine Linux-based Docker image (`Dockerfile`), Podman build support (`build_local.sh`).
- **Target OS / Architecture**: Linux `x86_64` / `arm64`, macOS.
- **Timezone Requirement**: `Asia/Jakarta` (WIB) strictly enforced for stock exchange trading day calculations.

---

## 7. Coding Conventions & Standards
- **Linter & Formatting**: Pre-commit hooks (`.pre-commit-config.yaml`) using standard `gofmt`, `goimports`, and `golangci-lint`.
- **Pre-Push Quality Gate**: Automated `go generate ./...` (Swagger sync), `gotest` (>= 80% coverage), `gobuild`, and `govulncheck`.
- **Commit Style**: Conventional Commits (`feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `perf`).
- **Build Command**: `go build -o idxstock ./cmd/api/main.go` or `./build_local.sh`.
