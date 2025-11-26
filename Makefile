# cargo-cpp Makefile

.PHONY: all build-client build-server install clean setup-server run-server run-frontend help

# Default target
all: build-client

# Build the Go CLI client (statically linked)
build-client:
	@echo "🔨 Building cargo-cpp client..."
	cd cargo-cpp-client && \
		CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp .
	@echo "✅ Built: bin/cargo-cpp"

# Build for all platforms
build-all: build-client
	@echo "🔨 Building for all platforms..."
	@mkdir -p bin
	cd cargo-cpp-client && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp-linux-amd64 . && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp-linux-arm64 . && \
		GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp-darwin-amd64 . && \
		GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp-darwin-arm64 . && \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cargo-cpp-windows-amd64.exe .
	@echo "✅ Built binaries for all platforms in bin/"

# Install the client to /usr/local/bin
install: build-client
	@echo "📦 Installing cargo-cpp to /usr/local/bin..."
	sudo cp bin/cargo-cpp /usr/local/bin/
	@echo "✅ Installed! Run 'cargo-cpp --help' to get started"

# Setup the Python server environment
setup-server:
	@echo "🐍 Setting up Python server..."
	cd cargo-cpp-server && \
		python3 -m venv venv && \
		./venv/bin/pip install -r requirements.txt
	@echo "✅ Server setup complete"

# Run the server
run-server:
	@echo "🚀 Starting cargo-cpp server on http://localhost:8000..."
	cd cargo-cpp-server && \
		./venv/bin/uvicorn main:app --reload --port 8000

# Run the frontend
run-frontend:
	@echo "🚀 Starting frontend on http://localhost:5173..."
	cd frontend && npm run dev

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf cargo-cpp-client/cargo-cpp
	rm -rf cargo-cpp-server/__pycache__
	rm -rf cargo-cpp-server/venv
	@echo "✅ Cleaned build artifacts"

# Download Go dependencies
deps:
	cd cargo-cpp-client && go mod tidy

# Help
help:
	@echo "cargo-cpp - C++ Project Generator"
	@echo ""
	@echo "Usage:"
	@echo "  make build-client   Build the Go CLI client"
	@echo "  make build-all      Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install        Install cargo-cpp to /usr/local/bin"
	@echo "  make setup-server   Setup Python virtual environment for server"
	@echo "  make run-server     Start the FastAPI server"
	@echo "  make run-frontend   Start the React frontend"
	@echo "  make clean          Remove build artifacts"
	@echo "  make deps           Download Go dependencies"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make setup-server"
	@echo "  2. make run-server    (in one terminal)"
	@echo "  3. make build-client"
	@echo "  4. ./bin/cargo-cpp init"
	@echo "  5. ./bin/cargo-cpp build"

