# cargo-cpp

A C++ dependency manager and project generator - like Cargo for Rust, but for C++!

## Features

- **CLI Tool**: `cargo-cpp` command to create projects from `cpp-cargo.yaml`
- **60+ Libraries**: Curated collection of popular C++ libraries (5000+ GitHub stars)
- **Recipe System**: Libraries defined in YAML files - easy to customize and extend
- **Web UI**: Browse and select libraries visually
- **Testing Framework Selector**: GoogleTest, Catch2, doctest, or none
- **Clang-Format Styles**: Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU

## Quick Start

### 1. Setup Server

```bash
make setup-server
make run-server
```

### 2. Build CLI Client

```bash
make build-client
```

### 3. Create a Project

```bash
# Initialize a new cpp-cargo.yaml
./bin/cargo-cpp init

# Or use a template
./bin/cargo-cpp init -t web-server

# Build the project
./bin/cargo-cpp build

# Build and run
cd my_project
cmake -B build
cmake --build build
./build/my_project
```

## cpp-cargo.yaml Format

```yaml
package:
  name: my_project
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest  # googletest, catch2, doctest, none

dependencies:
  spdlog:
    spdlog_header_only: true
  nlohmann_json: {}
  fmt: {}
  cli11: {}
```

## CLI Commands

```bash
cargo-cpp build                    # Build from cpp-cargo.yaml
cargo-cpp build -c myconfig.yaml   # Build from specific config
cargo-cpp build -o ./output        # Output to specific directory
cargo-cpp init                     # Create cpp-cargo.yaml template
cargo-cpp init -t game             # Create from template (minimal, web-server, game, cli-tool, networking, data-processing)
cargo-cpp list                     # Show available libraries
cargo-cpp -s http://server:8000    # Use custom server
cargo-cpp --help                   # Show help
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
