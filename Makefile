run:
	go run cmd/web/main.go

test:
	go run gotest.tools/gotestsum@latest --hide-summary=skipped ./...

test-e2e:
	npm run test:e2e

test-e2e-headed:
	npm run test:e2e:headed

install-e2e:
	npm install

setup-e2e: install-e2e
	npx cypress install

# Versioning commands
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_HASH ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build with version information
build:
	@echo "Building SimpleWebServer v$(VERSION)..."
	go build -ldflags "-X github.com/anglesson/simple-web-server/internal/config.Version=$(VERSION) -X github.com/anglesson/simple-web-server/internal/config.CommitHash=$(COMMIT_HASH) -X github.com/anglesson/simple-web-server/internal/config.BuildTime=$(BUILD_TIME)" -o bin/simple-web-server cmd/web/main.go

# Build for production
build-prod:
	@echo "Building SimpleWebServer for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X github.com/anglesson/simple-web-server/internal/config.Version=$(VERSION) -X github.com/anglesson/simple-web-server/internal/config.CommitHash=$(COMMIT_HASH) -X github.com/anglesson/simple-web-server/internal/config.BuildTime=$(BUILD_TIME)" -o bin/simple-web-server-linux-amd64 cmd/web/main.go

# Build for macOS
build-mac:
	@echo "Building SimpleWebServer for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X github.com/anglesson/simple-web-server/internal/config.Version=$(VERSION) -X github.com/anglesson/simple-web-server/internal/config.CommitHash=$(COMMIT_HASH) -X github.com/anglesson/simple-web-server/internal/config.BuildTime=$(BUILD_TIME)" -o bin/simple-web-server-darwin-amd64 cmd/web/main.go

# Build for Windows
build-windows:
	@echo "Building SimpleWebServer for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X github.com/anglesson/simple-web-server/internal/config.Version=$(VERSION) -X github.com/anglesson/simple-web-server/internal/config.CommitHash=$(COMMIT_HASH) -X github.com/anglesson/simple-web-server/internal/config.BuildTime=$(BUILD_TIME)" -o bin/simple-web-server-windows-amd64.exe cmd/web/main.go

# Build all platforms
build-all: build-prod build-mac build-windows

# Create a new version tag
tag:
	@echo "Creating tag v$(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Show current version info
version:
	@echo "Current version: $(VERSION)"
	@echo "Commit hash: $(COMMIT_HASH)"
	@echo "Build time: $(BUILD_TIME)"

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with hot reload (requires air)
dev:
	air

# Docker build
docker-build:
	docker build -t simple-web-server:$(VERSION) .

# Docker run
docker-run:
	docker run -p 8080:8080 simple-web-server:$(VERSION)
