# C++ Project Creator

A web application for generating modern C++ project templates with selectable libraries and automatic CMake configuration using FetchContent.

## Features

- **Library Catalog**: Browse and select from a curated collection of popular C++ libraries
- **Category Filtering**: Filter libraries by category (JSON, logging, testing, networking, etc.)
- **Search**: Find libraries by name, description, or tags
- **CMake Preview**: Real-time preview of generated CMakeLists.txt
- **Download ZIP**: Generate complete project structure with all configurations

## Tech Stack

- **Backend**: FastAPI (Python)
- **Frontend**: React + TypeScript + TailwindCSS + Vite

## Getting Started

### Prerequisites

- Python 3.9+
- Node.js 18+
- npm or yarn

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

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/libraries` | GET | Get all available libraries |
| `/api/libraries/{id}` | GET | Get a specific library |
| `/api/categories` | GET | Get all library categories |
| `/api/categories/{id}/libraries` | GET | Get libraries in a category |
| `/api/search?q={query}` | GET | Search libraries |
| `/api/preview` | GET | Preview generated CMakeLists.txt |
| `/api/generate` | POST | Generate and download project ZIP |

## Project Structure

```
cpp-repo-creator/
├── backend/
│   ├── main.py              # FastAPI application
│   ├── libraries.py         # Library catalog
│   ├── generator.py         # CMake/project generator
│   └── requirements.txt
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── App.tsx          # Main application
│   │   ├── api.ts           # API client
│   │   └── types.ts         # TypeScript types
│   ├── index.html
│   └── package.json
└── README.md
```

## Available Libraries

### Serialization
- nlohmann/json
- json11
- RapidJSON
- simdjson

### Logging
- spdlog
- Google glog
- plog

### Testing
- Google Test
- Catch2
- doctest

### Networking
- Asio
- CPR
- cpp-httplib

### CLI
- CLI11
- argparse
- cxxopts

### Configuration
- yaml-cpp
- toml11
- toml++

### And more!

## Generated Project Structure

When you download a project, you get:

```
project_name/
├── CMakeLists.txt          # Main CMake configuration
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

