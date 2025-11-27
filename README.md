# cargo-cpp

A C++ dependency manager and project generator - like Cargo for Rust, but for C++!

## Features

- **Full Cargo-like CLI**: `build`, `run`, `test`, `add`, `remove`, `update`, `fmt`, `lint`, `doc` commands
- **60+ Libraries**: Curated collection of popular C++ libraries (5000+ GitHub stars)
- **Recipe System**: Libraries defined in YAML files - easy to customize and extend
- **Web UI**: Browse and select libraries visually
- **Lock File**: `cpp-cargo.lock` for reproducible builds
- **Code Quality Tools**: Built-in clang-format and clang-tidy integration
- **Documentation**: Doxygen integration with `cargo-cpp doc`

## Quick Install

Install with a single command (auto-detects your OS and architecture):

```bash
# via curl
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/cpp-repo-creator/master/install.sh)"

# via wget
sh -c "$(wget -qO- https://raw.githubusercontent.com/ozacod/cpp-repo-creator/master/install.sh)"
```

## Quick Start

```bash
# Create and run a new project
cargo-cpp new my_app
cd my_app
cargo-cpp build
cargo-cpp run

# Add dependencies
cargo-cpp add spdlog
cargo-cpp add --dev catch2

# Run tests
cargo-cpp test

# Format and lint
cargo-cpp fmt
cargo-cpp lint
```

## Server Setup (for self-hosting)

If you want to run your own server:

```bash
make setup-server
make run-server
```

## cpp-cargo.yaml Format

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
cargo-cpp new <name>              # Create new project directory
cargo-cpp new <name> --lib        # Create library project
cargo-cpp init                    # Create cpp-cargo.yaml in current dir
cargo-cpp init -t <template>      # Use template (minimal, web-server, game, cli-tool, networking, data-processing)
```

### Build & Run
```bash
cargo-cpp build                   # Generate CMake project from cpp-cargo.yaml
cargo-cpp build --release         # Build in release mode
cargo-cpp run                     # Build and run executable
cargo-cpp run --release           # Run in release mode
cargo-cpp run -- arg1 arg2        # Pass arguments to executable
cargo-cpp test                    # Build and run tests
cargo-cpp test -v                 # Verbose test output
cargo-cpp check                   # Check code compiles
cargo-cpp clean                   # Remove build artifacts
cargo-cpp clean --all             # Also remove generated files
```

### Dependency Management
```bash
cargo-cpp add <library>           # Add dependency
cargo-cpp add --dev <library>     # Add dev dependency
cargo-cpp remove <library>        # Remove dependency
cargo-cpp update                  # Update all dependencies
cargo-cpp update <library>        # Update specific dependency
cargo-cpp list                    # List available libraries
cargo-cpp search <query>          # Search for libraries
cargo-cpp info <library>          # Show library details
```

### Code Quality
```bash
cargo-cpp fmt                     # Format code with clang-format
cargo-cpp fmt --check             # Check formatting without modifying
cargo-cpp lint                    # Run clang-tidy static analysis
cargo-cpp lint --fix              # Auto-fix lint issues
```

### Documentation
```bash
cargo-cpp doc                     # Generate Doxygen documentation
cargo-cpp doc --open              # Open docs in browser
```

### Versioning
```bash
cargo-cpp release patch           # Bump 0.1.0 → 0.1.1
cargo-cpp release minor           # Bump 0.1.0 → 0.2.0
cargo-cpp release major           # Bump 0.1.0 → 1.0.0
```

## Project Structure

```
cargo-cpp/
├── cargo-cpp-client/        # Go CLI tool (statically compiled)
│   ├── main.go
│   └── go.mod
├── cargo-cpp-server/        # FastAPI server + recipes
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

Run `cargo-cpp list` to see all available libraries.

## Adding New Libraries

Create a new YAML file in `cargo-cpp-server/recipes/`:

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
| `/api/cargo` | POST | Generate from cpp-cargo.yaml |
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
