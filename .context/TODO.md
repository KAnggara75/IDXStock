# Project TODO & Technical Debt

## 1. Immediate Tasks
_Urgent tasks affecting reliability, security, or immediate execution._
- [ ] `cmd/api/main.go`: Add graceful shutdown handling (`app.ShutdownWithContext`) on SIGINT/SIGTERM signals to cleanly drain active DB connections and in-flight HTTP requests.
- [ ] `internal/database/database.go:86`: Clean up legacy fallback handling for `migrations/000_board_type.sql` and `migrations/001_create_stocks_table.sql` once SQLite migration baseline is firmly established across all runtime environments.

## 2. Existing Code Annotations (TODO / FIXME)
_Annotations found directly in source code._
- None currently present in application source code (`CONFIRMED`).

## 3. Technical Debt & Structural Improvements
_Long-term architectural and refactoring initiatives._
- [x] **Test Coverage Expansion**: Project test coverage reached **90.3%** across all packages (`cmd/api`, `internal/handlers`, `internal/repositories`, `internal/usecases`, `internal/services`, `internal/routes`, `internal/utils`), surpassing the 80% threshold.
- [ ] **Standardized Error Handling**: Unify API error response schema across all handlers (some return `{"error": ...}`, others return `{"status": "error", "message": ...}`).
- [ ] **Distributed Configuration Fallback**: If Spring Cloud Config (`scc2go`) is unreachable during `init()`, ensure deterministic fallback or informative warning logging rather than silent defaults.
- [ ] **Rate Limiting & Anti-Scraping Resilience**: Outbound calls to `idx.co.id` and `pasardana.id` rely on hardcoded HTTP headers and cookie bypasses; consider implementing exponential backoff and rotating proxy/token support.
