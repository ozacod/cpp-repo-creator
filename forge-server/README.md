# Forge Server

Go implementation of the Forge backend server using Gin.

## Features

- Recipe loading from YAML files
- CMake project generation
- ZIP file creation
- Static file serving for frontend
- CORS support

## Building

```bash
cd forge-server
go build ./cmd/server
```

## Running

```bash
./server
```

Or set environment variables:

```bash
PORT=8000 FORGE_RECIPES_DIR=recipes ./server
```

## API Endpoints

- `GET /api` - API root
- `GET /api/version` - Version information
- `GET /api/libraries` - Get all libraries
- `GET /api/libraries/:id` - Get specific library
- `GET /api/categories` - Get all categories
- `GET /api/categories/:id/libraries` - Get libraries by category
- `GET /api/search?q=query` - Search libraries
- `POST /api/reload-recipes` - Reload recipes
- `POST /api/generate` - Generate project ZIP
- `POST /api/preview` - Preview CMakeLists.txt
- `POST /api/forge` - Generate from forge.yaml
- `POST /api/forge/dependencies` - Generate dependencies.cmake only
- `GET /api/forge/template` - Get forge.yaml template
- `GET /api/forge/example/:template` - Get example templates

## Structure

```
forge-server/
├── cmd/
│   └── server/
│       └── main.go          # Main server entry point
├── internal/
│   ├── recipe/
│   │   └── loader.go       # Recipe loader (YAML parsing)
│   └── generator/
│       ├── generator.go    # CMake generation
│       ├── files.go        # File generation (main.cpp, tests, etc.)
│       └── zip.go          # ZIP file creation
├── recipes/                 # YAML recipe files
└── go.mod                   # Go module definition
```

## Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `gopkg.in/yaml.v3` - YAML parsing

