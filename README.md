# IDXStock

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?style=flat&logo=go)](https://golang.org)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v3-00ACD7?style=flat&logo=go)](https://gofiber.io)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)](https://www.sqlite.org)
[![Code Coverage](https://img.shields.io/badge/Coverage-90.0%25-brightgreen?style=flat)](https://github.com/KAnggara75/IDXStock)
[![codecov](https://codecov.io/gh/KAnggara75/IDXStock/branch/main/graph/badge.svg)](https://codecov.io/gh/KAnggara75/IDXStock)

A high-performance Golang backend microservice engineered to ingest, normalize, persist, and serve Indonesian financial market data (Indonesia Stock Exchange / IDX).

---

## 📊 Data Sources

This application interfaces with reliable financial and exchange data sources:
- **IDX (Indonesia Stock Exchange)**: [idx.co.id](https://www.idx.co.id) — Listed/delisted equities and daily trading statistics.
- **Pasardana**: [pasardana.id](https://pasardana.id) — Sector, sub-sector, industry classifications, and historical quote summaries.
- **Official Listing Spreadsheets**: Microsoft Excel (`.xlsx`) listing reports published by the exchange.

---

## 🚀 Key Features

- **Historical Stock Pricing**: Retrieve end-of-day (EOD) historical quotes with flexible date filtering (`start_date`, `end_date`) and sparse field selection (`fields`).
- **Corporate Action Adjustments**: Dynamic stock split calculation endpoint (`PUT /api/v1/stocks/split`) that adjusts historical open, high, low, close prices and recalculates trading volume.
- **Listing & Delisting Synchronization**: Ingest and track active listing profiles, sector mapping, and delisted stocks.
- **High-Volume Broker Activities**: Automated weekly ISO partition management (`broker_activity_YYYY_WW`) for scalable broker summary transaction ingestion.
- **Media Asset Streaming**: In-memory proxy and caching for company logo assets (`/api/v1/assets/logos/companies/:code`).
- **Decoupled Architecture**: Spring Cloud Config remote configuration, zero-downtime database auto-migrations, and structured JSON telemetry.

---

## 🛠 Tech Stack

- **Language / Runtime**: Go 1.26.3
- **HTTP Engine**: Fiber v3 (`github.com/gofiber/fiber/v3`)
- **Database Driver / Engine**: SQLite 3 via `modernc.org/sqlite` (pure Go, CGO-free)
- **Database**: SQLite 3 (single-file or in-memory DB with weekly table partitions)
- **Spreadsheet Processing**: Excelize v2 (`github.com/xuri/excelize/v2`)
- **Configuration Management**: Spring Cloud Config client (`scc2go`) + Viper
- **Logging**: Logrus JSON structured logger (`sirupsen/logrus`)

---

## 📁 Project Architecture & Memory

Architectural documentation and project memory are centralized in [`.context/`](.context/):

- [`.context/PROJECT_CONTEXT.md`](.context/PROJECT_CONTEXT.md) — System boundaries, business glossary, and runtime constraints.
- [`.context/ARCHITECTURE.md`](.context/ARCHITECTURE.md) — Architectural diagrams, data lifecycle, concurrency model, and fault tolerance.
- [`.context/CODEBASE_MAP.md`](.context/CODEBASE_MAP.md) — Module navigation, subsystem responsibilities, and consumers.
- [`.context/DECISIONS.md`](.context/DECISIONS.md) — Architecture Decision Records (ADRs: Fiber v3, ISO-week partitioning, scc2go, pre-commit standards).
- [`.context/TODO.md`](.context/TODO.md) — Technical debt register, immediate tasks, and coverage targets.
- [`AGENTS.md`](AGENTS.md) — Primary orientation guide for AI agents and pair-programmers.

---

## ⚙️ Getting Started

### Prerequisites
- Go 1.26+ installed
- PostgreSQL 15+ instance
- (Optional) Podman / Docker

### Environment Variables
Configure the following environment variables:

| Variable | Description | Example |
| :--- | :--- | :--- |
| `IDXSTOCK_CONFIG_URL` | Spring Cloud Config server URL | `http://config-server:8888` |
| `SCC_AUTH` | HTTP Basic Auth for Config server | `Basic <base64-credentials>` |
| `DATABASE_URL` | PostgreSQL connection DSN (if not supplied via SCC) | `postgres://user:password@localhost:5432/idxstock?sslmode=disable` |
| `PORT` | HTTP server listen port (default: `3000`) | `3000` |
| `LOG_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |

### Running Locally

```bash
# Clone the repository
git clone git@github.com:KAnggara75/IDXStock.git
cd IDXStock

# Run hot-reloading development server
air

# Or compile and run directly
go build -o /dev/null ./cmd/api/main.go
go run ./cmd/api/main.go
```

### Running with Docker / Podman

```bash
# Build container image
podman build -t idxstock:latest .

# Run container
podman run -d -p 3000:3000 --env-file .env idxstock:latest
```

---

## 🧪 Testing & Quality Gates

The project uses `pre-commit` hooks for format and quality compliance:

```bash
# Install git hooks
pre-commit install --hook-type pre-commit --hook-type pre-push

# Run all pre-commit checks manually
pre-commit run --all-files

# Run test suite
go test -v ./...

# Run test coverage calculation
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## 📚 API Reference

### System & Health
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Verify HTTP service and PostgreSQL connection pool health |

### Stocks & Equities
| Method | Endpoint | Description | Query / Body Params |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/stocks/upload` | Preview parsed records from an Excel `.xlsx` file | Multipart form: `file` |
| `PATCH` | `/api/v1/stocks/upload` | Persist and upsert equities parsed from `.xlsx` | Multipart form: `file` |
| `PUT` | `/api/v1/stocks/sync` | Sync detailed company profiles & statistics | None |
| `PUT` | `/api/v1/stocks/delisting/sync` | Sync delisted equities from official sources | None |
| `PUT` | `/api/v1/stocks/history/sync` | Ingest end-of-day (EOD) quote history | None |
| `GET` | `/api/v1/stocks/:code/history` | Retrieve historical price data by stock ticker | `start_date`, `end_date`, `fields` |
| `PUT` | `/api/v1/stocks/split` | Adjust historical price & volume for stock splits | JSON body: `stock_code`, `date`, `from_ratio`, `to_ratio` |

### Master Data (Sectors & Industries)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `PUT` | `/api/v1/sectors/sync` | Fetch and update sector and sub-sector taxonomies |
| `PUT` | `/api/v1/industries/sync` | Fetch and update industry and sub-industry classifications |

### Broker Activities
| Method | Endpoint | Description | Query / Header Params |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/broker/sync` | Synchronize daily broker transaction summaries | Header: `Authorization: Bearer <token>`, Query: `broker_code`, `from`, `to`, `transaction_type`, `market_board`, `investor_type` |
| `PUT` | `/api/v1/partition/broker-activity` | Pre-allocate weekly PostgreSQL partitions (next 2 ISO weeks) | None |

### Assets & Media
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/assets/logos/companies/:code` | Stream company logo (e.g., `BBCA.png`) with HTTP caching |
