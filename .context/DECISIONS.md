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
- **Decision**: Initialize remote configuration retrieval using `github.com/KAnggara75/scc2go` during package `init()`, fetching settings dynamically based on `SCC_URL` and `SCC_AUTH`.
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
