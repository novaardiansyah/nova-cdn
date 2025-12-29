.PHONY: run build build-linux dev install migrate clean help

# Run the application
run:
	go run cmd/api/main.go

# Build the application (Windows)
build:
	set GOOS=windows&& set GOARCH=amd64&& go build -o ./runner-app.exe cmd/api/main.go

# Build for Linux (production)
build-linux:
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; go build -ldflags='-s -w' -o ./runner-app ./cmd/api"

# Run with hot reload (requires air)
dev:
	air

# Install dependencies
install:
	go mod download
	go mod tidy

# Clean build artifacts (Windows compatible)
clean:
	-cmd /c if exist tmp rmdir /s /q tmp
	-cmd /c if exist runner-app del /q runner-app
	-cmd /c if exist runner-app.exe del /q runner-app.exe

# Install air for hot reload
install-air:
	go install github.com/air-verse/air@latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go --parseDependency

# Show help
help:
	@echo "Available commands:"
	@echo "  make run         - Run the application"
	@echo "  make build       - Build the application (Windows)"
	@echo "  make build-linux - Build for Linux production"
	@echo "  make dev         - Run with hot reload (requires air)"
	@echo "  make swagger     - Generate Swagger documentation"
	@echo "  make install     - Install dependencies"
	@echo "  make install-air - Install air for hot reload"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make help        - Show this help message"

