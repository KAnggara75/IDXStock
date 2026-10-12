# Project Context & Domain Overview

## 1. Project Purpose
IDXStock is a specialized Golang backend microservice engineered to ingest, normalize, persist, and serve Indonesian financial market data (Indonesia Stock Exchange / IDX). It bridges disparate financial data feeds (official IDX endpoints, Pasardana, and listing spreadsheets) into a centralized, queryable PostgreSQL data store.

---

## 2. System Boundaries
- **In-Scope**:
  - Ingesting and validating official IDX listing reports from Excel files (`.xlsx`).
  - Synchronizing company profiles, sectors, sub-sectors, industries, and sub-industries.
  - Tracking delisted securities.
  - Ingesting and querying end-of-day (EOD) historical stock prices with support for date range filters and field selections.
  - Adjusting historical price series for corporate actions such as stock splits.
  - Streaming and caching company logo images.
  - Ingesting and partitioning daily broker transactions and summaries.
- **Out-of-Scope**:
  - Direct order routing or real-time trading execution.
  - User authentication and multi-tenant authorization (delegated to upstream API gateway/proxy).
  - Portfolio management or predictive financial analysis.

---

## 3. Main Actors & Consumers
- **Automated Schedulers / Cron Jobs**: Invoke periodic synchronization endpoints (`/api/v1/stocks/history/sync`, `/api/v1/broker/sync`, `/api/v1/partition/broker-activity`).
- **Internal Analysts / Admins**: Upload listing spreadsheets (`/api/v1/stocks/upload`) and trigger corporate stock splits.
- **Client Applications / Dashboard Services**: Query historical prices (`/api/v1/stocks/:code/history`) and company logos.

---

## 4. Domain Glossary
- **Stock Code / Ticker (`code`)**: 4-letter symbol representing an IDX listed security (e.g., `BBCA`, `BBRI`, `TLKM`).
- **Board Type (`board_type`)**: Trading board category on the IDX:
  - `MAIN` (Papan Utama)
  - `DEVELOPMENT` (Papan Pengembangan)
  - `ACCELERATION` (Papan Akselerasi)
  - `NEW_ECONOMY` (Papan Ekonomi Baru)
  - `WATCHLIST` (Papan Pemantauan Khusus)
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

---

## 6. Runtime Environment & Constraints
- **Language / Runtime**: Go 1.26.3
- **Containerization**: Alpine Linux-based Docker image (`Dockerfile`), Podman build support (`build_local.sh`).
- **Target OS / Architecture**: Linux `x86_64` / `arm64` container runtime.
- **Database Engine**: PostgreSQL 15+ (requires table partitioning support).
- **Timezone Requirement**: `Asia/Jakarta` (WIB) strictly enforced for stock exchange trading day calculations.

---

## 7. Coding Conventions & Standards
- **Linter & Formatting**: Pre-commit hooks (`.pre-commit-config.yaml`) using standard `gofmt`, `goimports`, and `golangci-lint`.
- **Commit Style**: Conventional Commits (`feat`, `fix`, `chore`, `docs`, `refactor`).
- **Build Command**: `go build -o idxstock ./cmd/api/main.go` or `./build_local.sh`.
