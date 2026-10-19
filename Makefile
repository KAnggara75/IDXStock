.PHONY: generate swagger test build run

# Generate OpenAPI / Swagger documentation
swagger:
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -d ./cmd/api,./internal/handlers,./internal/models -g main.go -o docs


# Run go generate
generate:
	go generate ./...

# Run all unit tests with coverage
test:
	go test -v -coverprofile=coverage.out ./...

# Build binary
build:
	go build -o idxstock ./cmd/api/main.go

# Run locally
run:
	go run ./cmd/api/main.go
