# C++ Project Creator

A web application for generating modern C++ project templates with selectable libraries and automatic CMake configuration using FetchContent. Works like Cargo for Rust, but for C++!

## Features

- **60+ Libraries**: Browse and select from popular C++ libraries with 5000+ GitHub stars
- **Recipe-Based Configuration**: Libraries defined in YAML files - easy to customize and extend
- **Category Filtering**: Filter by serialization, logging, testing, networking, CLI, GUI, etc.
- **Library Options**: Configure library-specific build options (header-only, SSL, etc.)
- **Testing Framework Selector**: Choose GoogleTest, Catch2, doctest, or none
- **Clang-Format Styles**: Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU
- **CMake Preview**: Real-time preview of generated CMakeLists.txt
- **cpp-cargo.yaml**: Cargo.toml-like dependency management via curl

## Tech Stack

- **Backend**: FastAPI (Python) with YAML recipe loader
- **Frontend**: React + TypeScript + TailwindCSS + Vite

## Getting Started

### Prerequisites

- Python 3.9+
- Node.js 18+

### Backend Setup

```bash
cd backend
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
uvicorn main:app --reload
```

The API will be available at `http://localhost:8000`

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

## cpp-cargo.yaml (CLI Usage)

Generate projects from the command line using a YAML dependency file:

```yaml
# cpp-cargo.yaml
package:
  name: my_project
  cpp_standard: 17

build:
  clang_format: Google

testing:
  framework: googletest

dependencies:
  spdlog:
    spdlog_header_only: true
  nlohmann_json: {}
  fmt: {}
  cli11: {}
```

```bash
# Generate project from cpp-cargo.yaml
curl -X POST -F "file=@cpp-cargo.yaml" http://localhost:8000/api/cargo -o project.zip

# Get a sample template
curl http://localhost:8000/api/cargo/template > cpp-cargo.yaml

# Get example templates
curl http://localhost:8000/api/cargo/example/minimal > cpp-cargo.yaml
curl http://localhost:8000/api/cargo/example/web-server > cpp-cargo.yaml
curl http://localhost:8000/api/cargo/example/game > cpp-cargo.yaml
curl http://localhost:8000/api/cargo/example/cli-tool > cpp-cargo.yaml
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/libraries` | GET | Get all available libraries |
| `/api/libraries/{id}` | GET | Get a specific library with options |
| `/api/categories` | GET | Get all library categories |
| `/api/search?q={query}` | GET | Search libraries |
| `/api/preview` | POST | Preview generated CMakeLists.txt |
| `/api/generate` | POST | Generate and download project ZIP |
| `/api/cargo` | POST | Generate from cpp-cargo.yaml file |
| `/api/cargo/template` | GET | Get cpp-cargo.yaml template |
| `/api/cargo/example/{name}` | GET | Get example templates |
| `/api/reload-recipes` | POST | Hot-reload recipe files |

## Project Structure

```
cpp-repo-creator/
├── backend/
│   ├── main.py              # FastAPI application
│   ├── recipe_loader.py     # YAML recipe loader
│   ├── generator.py         # CMake/project generator
│   ├── recipes/             # Library recipe YAML files
│   │   ├── spdlog.yaml
│   │   ├── fmt.yaml
│   │   └── ...
│   └── requirements.txt
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── App.tsx          # Main application
│   │   ├── api.ts           # API client
│   │   └── types.ts         # TypeScript types
│   └── package.json
├── cpp-cargo.yaml           # Example dependency file
└── README.md
```

## Available Libraries (60+)

### Serialization
- nlohmann/json, json11, RapidJSON, simdjson, cereal

### Logging
- spdlog, Google glog, plog

### Testing
- Google Test, Catch2, doctest, Google Benchmark

### Networking
- Asio, CPR, cpp-httplib, Crow, Drogon, WebSocket++, POCO, libevent, libcurl

### CLI
- CLI11, argparse, cxxopts, indicators, tabulate

### GUI/Graphics
- Dear ImGui, SFML, raylib, GLFW

### Utility
- Abseil, fmt, range-v3, magic_enum, EnTT, stb, xxHash, mimalloc, pybind11, backward-cpp

### And more!

## Adding New Libraries

Create a new YAML file in `backend/recipes/`:

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

Then reload recipes: `curl -X POST http://localhost:8000/api/reload-recipes`

## Generated Project Structure

```
project_name/
├── CMakeLists.txt          # Main CMake with FetchContent
├── include/
│   └── project_name/
│       └── project_name.hpp
├── src/
│   ├── main.cpp
│   └── project_name.cpp
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
├── .gitignore
├── .clang-format
└── README.md
```

## License

MIT
