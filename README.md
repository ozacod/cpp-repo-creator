# Forge

A C++ dependency manager and project generator - Forge Your Code!

## Features

- **Full Cargo-like CLI**: `generate`, `build`, `run`, `test`, `add`, `remove`, `update`, `fmt`, `lint`, `doc` commands
- **60+ Libraries**: Curated collection of popular C++ libraries (5000+ GitHub stars)
- **Recipe System**: Libraries defined in YAML files - easy to customize and extend
- **Web UI**: Browse and select libraries visually
- **Lock File**: `forge.lock` for reproducible builds
- **Code Quality Tools**: Built-in clang-format and clang-tidy integration
- **Documentation**: Doxygen integration with `forge doc`

## Quick Install

Install with a single command (auto-detects your OS and architecture):

```bash
# via curl
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/forge/master/install.sh)"

# via wget
sh -c "$(wget -qO- https://raw.githubusercontent.com/ozacod/forge/master/install.sh)"
```

## Quick Start

```bash
# Create and run a new project
forge new my_app
cd my_app
forge generate      # Generate CMake project from yaml
forge build         # Compile with CMake
forge run           # Run the executable

# Add dependencies
forge add spdlog
forge add --dev catch2
forge generate      # Regenerate to include new deps

# Run tests
forge test

# Format and lint
forge fmt
forge lint
```

## Server Setup (for self-hosting)

If you want to run your own server:

```bash
make setup-server
make run-server
```

## forge.yaml Format

```yaml
package:
  name: my_project
  version: "0.1.0"
  cpp_standard: 17
  authors: ["Your Name"]
  description: "My awesome project"

build:
  shared_libs: false
  clang_format: Google
  build_type: Debug  # Debug, Release, RelWithDebInfo

testing:
  framework: googletest  # googletest, catch2, doctest, none

dependencies:
  spdlog:
    spdlog_header_only: true
  nlohmann_json: {}
  fmt: {}
  cli11: {}

dev-dependencies:
  catch2: {}
```

## CLI Commands

### Project Management
```bash
forge new <name>              # Create new project directory
forge new <name> --lib        # Create library project
forge init                    # Create forge.yaml in current dir
forge init -t <template>      # Use template (minimal, web-server, game, cli-tool, networking, data-processing)
```

### Generate & Build
```bash
forge generate                # Generate CMake project from forge.yaml (alias: gen)
forge generate -o ./output    # Output to specific directory
forge build                   # Compile the project (Debug mode)
forge build --release         # Build in release mode (O2)
forge build -O3               # Build with O3 optimization
forge build -Os               # Optimize for size
forge build --clean           # Clean and rebuild
forge build -j 8              # Use 8 parallel jobs
forge run                     # Build and run executable
forge run --release           # Run in release mode
forge run -- arg1 arg2        # Pass arguments to executable
forge test                    # Build and run tests
forge test -v                 # Verbose test output
forge check                   # Check code compiles
forge clean                   # Remove build artifacts
forge clean --all             # Also remove generated files
```

### Dependency Management
```bash
forge add <library>           # Add dependency
forge add --dev <library>     # Add dev dependency
forge remove <library>        # Remove dependency
forge update                  # Update all dependencies
forge update <library>        # Update specific dependency
forge list                    # List available libraries
forge search <query>          # Search for libraries
forge info <library>          # Show library details
```

### Code Quality
```bash
forge fmt                     # Format code with clang-format
forge fmt --check             # Check formatting without modifying
forge lint                    # Run clang-tidy static analysis
forge lint --fix              # Auto-fix lint issues
```

### Documentation
```bash
forge doc                     # Generate Doxygen documentation
forge doc --open              # Open docs in browser
```

### Versioning
```bash
forge release patch           # Bump 0.1.0 → 0.1.1
forge release minor           # Bump 0.1.0 → 0.2.0
forge release major           # Bump 0.1.0 → 1.0.0
```

## Project Structure

```
forge/
├── forge-client/            # Go CLI tool (statically compiled)
│   ├── main.go
│   └── go.mod
├── forge-server/            # FastAPI server + recipes
│   ├── main.py              # API endpoints
│   ├── generator.py         # CMake/project generator
│   ├── recipe_loader.py     # YAML recipe loader
│   ├── recipes/             # Library recipe files
│   │   ├── spdlog.yaml
│   │   ├── fmt.yaml
│   │   └── ...
│   └── requirements.txt
├── frontend/                # React web UI
├── Makefile
└── README.md
```

## Building from Source

### Prerequisites

- Go 1.21+ (for CLI client)
- Python 3.9+ (for server)
- Node.js 18+ (for web UI, optional)

### Build CLI Client

```bash
# Build for current platform
make build-client

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Install to /usr/local/bin
make install
```

### Run Server

```bash
make setup-server
make run-server
```

### Run Web UI (Optional)

```bash
cd frontend
npm install
npm run dev
```

## Available Libraries (60+)

| Category | Libraries |
|----------|-----------|
| **Serialization** | nlohmann/json, json11, RapidJSON, simdjson, cereal |
| **Logging** | spdlog, Google glog, plog |
| **Testing** | Google Test, Catch2, doctest, Google Benchmark |
| **Networking** | Asio, CPR, cpp-httplib, Crow, Drogon, WebSocket++, POCO, libevent, libcurl |
| **CLI** | CLI11, argparse, cxxopts, indicators, tabulate |
| **GUI/Graphics** | Dear ImGui, SFML, raylib, GLFW |
| **Utility** | Abseil, fmt, range-v3, magic_enum, EnTT, stb, xxHash, mimalloc, pybind11, backward-cpp |
| **Database** | hiredis, sqlite_modern_cpp |
| **Compression** | zlib, zstd, LZ4 |
| **Cryptography** | OpenSSL, Mbed TLS |
| **Math** | Eigen, GLM |

Run `forge list` to see all available libraries.

## Adding New Libraries

Create a new YAML file in `forge-server/recipes/`:

```yaml
id: mylib
name: MyLib
description: My awesome library
category: utility

github_url: https://github.com/user/mylib
cpp_standard: 17
header_only: true
tags:
  - awesome
  - header-only

fetch_content:
  repository: https://github.com/user/mylib.git
  tag: v1.0.0

link_libraries:
  - mylib::mylib

options:
  - id: mylib_option
    name: Enable Feature
    description: Enable some feature
    type: boolean
    default: false
    cmake_var: MYLIB_ENABLE_FEATURE
```

Recipes are hot-reloaded - no server restart needed.

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/libraries` | GET | Get all libraries |
| `/api/libraries/{id}` | GET | Get library with options |
| `/api/categories` | GET | Get categories |
| `/api/cargo` | POST | Generate from forge.yaml |
| `/api/cargo/template` | GET | Get template |
| `/api/cargo/example/{name}` | GET | Get example template |
| `/api/generate` | POST | Generate project (JSON) |
| `/api/preview` | POST | Preview CMakeLists.txt |

## Generated Project Structure

```
my_project/
├── CMakeLists.txt          # CMake with FetchContent
├── include/
│   └── my_project/
│       └── my_project.hpp
├── src/
│   ├── main.cpp
│   └── my_project.cpp
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
├── .gitignore
├── .clang-format
└── README.md
```

## License

MIT
