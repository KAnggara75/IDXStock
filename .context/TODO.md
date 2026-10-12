# Project TODO & Technical Debt

## 1. Immediate Tasks
_Urgent tasks affecting reliability, security, or immediate execution._
- [ ] `internal/database/database.go:86`: Remove or clean legacy fallback hardcoding for `migrations/000_board_type.sql` and `migrations/001_create_stocks_table.sql` once historical migrations are completely standardized across all deployment environments.
- [ ] `internal/routes/routes.go:38`: Add graceful shutdown handling (`app.ShutdownWithContext`) on SIGINT/SIGTERM signals in `cmd/api/main.go` to cleanly drain active DB connections and HTTP requests.

## 2. Existing Code Annotations (TODO / FIXME)
_Annotations found directly in source code._
- None currently present in application source code (`CONFIRMED`).

## 3. Technical Debt & Structural Improvements
_Long-term architectural and refactoring initiatives._
- [ ] **Test Coverage Expansion**: Current test coverage is ~7.4% (target: 75% as noted in `README.md`). Repository lacks comprehensive integration and unit tests for `StockUsecase`, `HistoryUsecase`, `StockRepository`, `HistoryRepository`, and external HTTP clients (`IdxService`, `PasardanaService`).
- [ ] **Standardized Error Handling**: Unify API error response schema across all handlers (some return `{"error": ...}`, others return `{"status": "error", "message": ...}`).
- [ ] **Distributed Configuration Fallback**: If Spring Cloud Config (`scc2go`) is unreachable during `init()`, ensure deterministic fallback or informative error logging rather than silent defaults.
- [ ] **Rate Limiting & Anti-Scraping Resilience**: Outbound calls to `idx.co.id` and `pasardana.id` rely on hardcoded HTTP headers and cookie bypasses; consider implementing exponential backoff and rotating proxy/token support.
