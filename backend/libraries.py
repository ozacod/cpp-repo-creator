"""
C++ Library catalog with FetchContent configurations.
"""

from typing import TypedDict, List, Optional


class Library(TypedDict):
    id: str
    name: str
    description: str
    category: str
    tags: List[str]
    github_url: str
    fetch_content: str
    link_libraries: List[str]
    header_only: bool
    cpp_standard: int
    alternatives: List[str]


LIBRARIES: List[Library] = [
    # JSON Libraries
    {
        "id": "nlohmann_json",
        "name": "nlohmann/json",
        "description": "JSON for Modern C++ - A header-only library with intuitive syntax",
        "category": "serialization",
        "tags": ["json", "serialization", "header-only"],
        "github_url": "https://github.com/nlohmann/json",
        "fetch_content": """FetchContent_Declare(
    json
    GIT_REPOSITORY https://github.com/nlohmann/json.git
    GIT_TAG v3.11.3
)
FetchContent_MakeAvailable(json)""",
        "link_libraries": ["nlohmann_json::nlohmann_json"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["json11", "rapidjson", "simdjson"],
    },
    {
        "id": "json11",
        "name": "json11",
        "description": "A tiny JSON library for C++11 by Dropbox",
        "category": "serialization",
        "tags": ["json", "serialization", "lightweight"],
        "github_url": "https://github.com/dropbox/json11",
        "fetch_content": """FetchContent_Declare(
    json11
    GIT_REPOSITORY https://github.com/dropbox/json11.git
    GIT_TAG v1.0.0
)
FetchContent_MakeAvailable(json11)""",
        "link_libraries": ["json11"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["nlohmann_json", "rapidjson"],
    },
    {
        "id": "rapidjson",
        "name": "RapidJSON",
        "description": "A fast JSON parser/generator with SAX/DOM style API",
        "category": "serialization",
        "tags": ["json", "serialization", "high-performance", "header-only"],
        "github_url": "https://github.com/Tencent/rapidjson",
        "fetch_content": """FetchContent_Declare(
    rapidjson
    GIT_REPOSITORY https://github.com/Tencent/rapidjson.git
    GIT_TAG v1.1.0
)
FetchContent_MakeAvailable(rapidjson)""",
        "link_libraries": ["rapidjson"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["nlohmann_json", "simdjson"],
    },
    {
        "id": "simdjson",
        "name": "simdjson",
        "description": "Parsing gigabytes of JSON per second using SIMD instructions",
        "category": "serialization",
        "tags": ["json", "serialization", "simd", "high-performance"],
        "github_url": "https://github.com/simdjson/simdjson",
        "fetch_content": """FetchContent_Declare(
    simdjson
    GIT_REPOSITORY https://github.com/simdjson/simdjson.git
    GIT_TAG v3.6.3
)
FetchContent_MakeAvailable(simdjson)""",
        "link_libraries": ["simdjson::simdjson"],
        "header_only": False,
        "cpp_standard": 17,
        "alternatives": ["nlohmann_json", "rapidjson"],
    },
    # Logging Libraries
    {
        "id": "spdlog",
        "name": "spdlog",
        "description": "Fast C++ logging library with support for multiple sinks",
        "category": "logging",
        "tags": ["logging", "header-only", "fast"],
        "github_url": "https://github.com/gabime/spdlog",
        "fetch_content": """FetchContent_Declare(
    spdlog
    GIT_REPOSITORY https://github.com/gabime/spdlog.git
    GIT_TAG v1.12.0
)
FetchContent_MakeAvailable(spdlog)""",
        "link_libraries": ["spdlog::spdlog"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["glog", "plog"],
    },
    {
        "id": "glog",
        "name": "Google glog",
        "description": "Google's C++ logging library with severity levels",
        "category": "logging",
        "tags": ["logging", "google"],
        "github_url": "https://github.com/google/glog",
        "fetch_content": """FetchContent_Declare(
    glog
    GIT_REPOSITORY https://github.com/google/glog.git
    GIT_TAG v0.6.0
)
FetchContent_MakeAvailable(glog)""",
        "link_libraries": ["glog::glog"],
        "header_only": False,
        "cpp_standard": 14,
        "alternatives": ["spdlog", "plog"],
    },
    {
        "id": "plog",
        "name": "plog",
        "description": "Portable, simple and extensible C++ logging library",
        "category": "logging",
        "tags": ["logging", "header-only", "portable"],
        "github_url": "https://github.com/SergiusTheBest/plog",
        "fetch_content": """FetchContent_Declare(
    plog
    GIT_REPOSITORY https://github.com/SergiusTheBest/plog.git
    GIT_TAG 1.1.10
)
FetchContent_MakeAvailable(plog)""",
        "link_libraries": ["plog::plog"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["spdlog", "glog"],
    },
    # Testing Libraries
    {
        "id": "googletest",
        "name": "Google Test",
        "description": "Google's C++ testing and mocking framework",
        "category": "testing",
        "tags": ["testing", "mocking", "google"],
        "github_url": "https://github.com/google/googletest",
        "fetch_content": """FetchContent_Declare(
    googletest
    GIT_REPOSITORY https://github.com/google/googletest.git
    GIT_TAG v1.14.0
)
FetchContent_MakeAvailable(googletest)""",
        "link_libraries": ["GTest::gtest", "GTest::gtest_main", "GTest::gmock"],
        "header_only": False,
        "cpp_standard": 14,
        "alternatives": ["catch2", "doctest"],
    },
    {
        "id": "catch2",
        "name": "Catch2",
        "description": "Modern, C++-native test framework with natural language syntax",
        "category": "testing",
        "tags": ["testing", "header-only", "bdd"],
        "github_url": "https://github.com/catchorg/Catch2",
        "fetch_content": """FetchContent_Declare(
    Catch2
    GIT_REPOSITORY https://github.com/catchorg/Catch2.git
    GIT_TAG v3.5.0
)
FetchContent_MakeAvailable(Catch2)""",
        "link_libraries": ["Catch2::Catch2WithMain"],
        "header_only": False,
        "cpp_standard": 14,
        "alternatives": ["googletest", "doctest"],
    },
    {
        "id": "doctest",
        "name": "doctest",
        "description": "Fastest feature-rich C++ single-header testing framework",
        "category": "testing",
        "tags": ["testing", "header-only", "fast"],
        "github_url": "https://github.com/doctest/doctest",
        "fetch_content": """FetchContent_Declare(
    doctest
    GIT_REPOSITORY https://github.com/doctest/doctest.git
    GIT_TAG v2.4.11
)
FetchContent_MakeAvailable(doctest)""",
        "link_libraries": ["doctest::doctest"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["googletest", "catch2"],
    },
    # Networking Libraries
    {
        "id": "asio",
        "name": "Asio",
        "description": "Cross-platform C++ library for network and I/O programming",
        "category": "networking",
        "tags": ["networking", "async", "header-only"],
        "github_url": "https://github.com/chriskohlhoff/asio",
        "fetch_content": """FetchContent_Declare(
    asio
    GIT_REPOSITORY https://github.com/chriskohlhoff/asio.git
    GIT_TAG asio-1-28-2
)
FetchContent_MakeAvailable(asio)
add_library(asio INTERFACE)
target_include_directories(asio INTERFACE ${asio_SOURCE_DIR}/asio/include)
target_compile_definitions(asio INTERFACE ASIO_STANDALONE)""",
        "link_libraries": ["asio"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["boost_asio", "libuv"],
    },
    {
        "id": "cpr",
        "name": "CPR",
        "description": "C++ Requests - A simple wrapper around libcurl inspired by Python Requests",
        "category": "networking",
        "tags": ["http", "networking", "curl"],
        "github_url": "https://github.com/libcpr/cpr",
        "fetch_content": """FetchContent_Declare(
    cpr
    GIT_REPOSITORY https://github.com/libcpr/cpr.git
    GIT_TAG 1.10.5
)
FetchContent_MakeAvailable(cpr)""",
        "link_libraries": ["cpr::cpr"],
        "header_only": False,
        "cpp_standard": 17,
        "alternatives": ["httplib"],
    },
    {
        "id": "httplib",
        "name": "cpp-httplib",
        "description": "A C++ header-only HTTP/HTTPS server and client library",
        "category": "networking",
        "tags": ["http", "networking", "header-only"],
        "github_url": "https://github.com/yhirose/cpp-httplib",
        "fetch_content": """FetchContent_Declare(
    httplib
    GIT_REPOSITORY https://github.com/yhirose/cpp-httplib.git
    GIT_TAG v0.14.3
)
FetchContent_MakeAvailable(httplib)""",
        "link_libraries": ["httplib::httplib"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["cpr"],
    },
    # CLI Libraries
    {
        "id": "cli11",
        "name": "CLI11",
        "description": "Command line parser for C++11 and beyond",
        "category": "cli",
        "tags": ["cli", "argument-parser", "header-only"],
        "github_url": "https://github.com/CLIUtils/CLI11",
        "fetch_content": """FetchContent_Declare(
    CLI11
    GIT_REPOSITORY https://github.com/CLIUtils/CLI11.git
    GIT_TAG v2.3.2
)
FetchContent_MakeAvailable(CLI11)""",
        "link_libraries": ["CLI11::CLI11"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["argparse", "cxxopts"],
    },
    {
        "id": "argparse",
        "name": "argparse",
        "description": "Argument parser for modern C++ inspired by Python",
        "category": "cli",
        "tags": ["cli", "argument-parser", "header-only"],
        "github_url": "https://github.com/p-ranav/argparse",
        "fetch_content": """FetchContent_Declare(
    argparse
    GIT_REPOSITORY https://github.com/p-ranav/argparse.git
    GIT_TAG v3.0
)
FetchContent_MakeAvailable(argparse)""",
        "link_libraries": ["argparse::argparse"],
        "header_only": True,
        "cpp_standard": 17,
        "alternatives": ["cli11", "cxxopts"],
    },
    {
        "id": "cxxopts",
        "name": "cxxopts",
        "description": "Lightweight C++ command line option parser",
        "category": "cli",
        "tags": ["cli", "argument-parser", "header-only", "lightweight"],
        "github_url": "https://github.com/jarro2783/cxxopts",
        "fetch_content": """FetchContent_Declare(
    cxxopts
    GIT_REPOSITORY https://github.com/jarro2783/cxxopts.git
    GIT_TAG v3.1.1
)
FetchContent_MakeAvailable(cxxopts)""",
        "link_libraries": ["cxxopts::cxxopts"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["cli11", "argparse"],
    },
    # Configuration Libraries
    {
        "id": "yaml_cpp",
        "name": "yaml-cpp",
        "description": "YAML parser and emitter in C++",
        "category": "configuration",
        "tags": ["yaml", "configuration", "serialization"],
        "github_url": "https://github.com/jbeder/yaml-cpp",
        "fetch_content": """FetchContent_Declare(
    yaml-cpp
    GIT_REPOSITORY https://github.com/jbeder/yaml-cpp.git
    GIT_TAG 0.8.0
)
FetchContent_MakeAvailable(yaml-cpp)""",
        "link_libraries": ["yaml-cpp::yaml-cpp"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["toml11", "tomlplusplus"],
    },
    {
        "id": "toml11",
        "name": "toml11",
        "description": "TOML for Modern C++ - A header-only TOML library",
        "category": "configuration",
        "tags": ["toml", "configuration", "header-only"],
        "github_url": "https://github.com/ToruNiina/toml11",
        "fetch_content": """FetchContent_Declare(
    toml11
    GIT_REPOSITORY https://github.com/ToruNiina/toml11.git
    GIT_TAG v3.8.1
)
FetchContent_MakeAvailable(toml11)""",
        "link_libraries": ["toml11::toml11"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["yaml_cpp", "tomlplusplus"],
    },
    {
        "id": "tomlplusplus",
        "name": "toml++",
        "description": "Header-only TOML config file parser and serializer for C++17",
        "category": "configuration",
        "tags": ["toml", "configuration", "header-only"],
        "github_url": "https://github.com/marzer/tomlplusplus",
        "fetch_content": """FetchContent_Declare(
    tomlplusplus
    GIT_REPOSITORY https://github.com/marzer/tomlplusplus.git
    GIT_TAG v3.4.0
)
FetchContent_MakeAvailable(tomlplusplus)""",
        "link_libraries": ["tomlplusplus::tomlplusplus"],
        "header_only": True,
        "cpp_standard": 17,
        "alternatives": ["yaml_cpp", "toml11"],
    },
    # GUI Libraries
    {
        "id": "imgui",
        "name": "Dear ImGui",
        "description": "Immediate mode graphical user interface for C++",
        "category": "gui",
        "tags": ["gui", "immediate-mode", "graphics"],
        "github_url": "https://github.com/ocornut/imgui",
        "fetch_content": """FetchContent_Declare(
    imgui
    GIT_REPOSITORY https://github.com/ocornut/imgui.git
    GIT_TAG v1.90.1
)
FetchContent_MakeAvailable(imgui)
add_library(imgui STATIC
    ${imgui_SOURCE_DIR}/imgui.cpp
    ${imgui_SOURCE_DIR}/imgui_demo.cpp
    ${imgui_SOURCE_DIR}/imgui_draw.cpp
    ${imgui_SOURCE_DIR}/imgui_tables.cpp
    ${imgui_SOURCE_DIR}/imgui_widgets.cpp
)
target_include_directories(imgui PUBLIC ${imgui_SOURCE_DIR})""",
        "link_libraries": ["imgui"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["qt"],
    },
    # Formatting Libraries
    {
        "id": "fmt",
        "name": "fmt",
        "description": "A modern formatting library - faster and safer alternative to printf",
        "category": "formatting",
        "tags": ["formatting", "string", "header-only"],
        "github_url": "https://github.com/fmtlib/fmt",
        "fetch_content": """FetchContent_Declare(
    fmt
    GIT_REPOSITORY https://github.com/fmtlib/fmt.git
    GIT_TAG 10.1.1
)
FetchContent_MakeAvailable(fmt)""",
        "link_libraries": ["fmt::fmt"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": [],
    },
    # Concurrency Libraries
    {
        "id": "concurrentqueue",
        "name": "moodycamel::ConcurrentQueue",
        "description": "Fast multi-producer, multi-consumer lock-free concurrent queue",
        "category": "concurrency",
        "tags": ["concurrency", "lock-free", "queue", "header-only"],
        "github_url": "https://github.com/cameron314/concurrentqueue",
        "fetch_content": """FetchContent_Declare(
    concurrentqueue
    GIT_REPOSITORY https://github.com/cameron314/concurrentqueue.git
    GIT_TAG v1.0.4
)
FetchContent_MakeAvailable(concurrentqueue)""",
        "link_libraries": ["concurrentqueue"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["readerwriterqueue"],
    },
    {
        "id": "readerwriterqueue",
        "name": "moodycamel::ReaderWriterQueue",
        "description": "Fast single-producer, single-consumer lock-free queue",
        "category": "concurrency",
        "tags": ["concurrency", "lock-free", "queue", "header-only"],
        "github_url": "https://github.com/cameron314/readerwriterqueue",
        "fetch_content": """FetchContent_Declare(
    readerwriterqueue
    GIT_REPOSITORY https://github.com/cameron314/readerwriterqueue.git
    GIT_TAG v1.0.6
)
FetchContent_MakeAvailable(readerwriterqueue)""",
        "link_libraries": ["readerwriterqueue"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["concurrentqueue"],
    },
    {
        "id": "transwarp",
        "name": "transwarp",
        "description": "Task-based parallelism library with dependency graphs",
        "category": "concurrency",
        "tags": ["concurrency", "task-based", "header-only"],
        "github_url": "https://github.com/bloomen/transwarp",
        "fetch_content": """FetchContent_Declare(
    transwarp
    GIT_REPOSITORY https://github.com/bloomen/transwarp.git
    GIT_TAG 2.2.3
)
FetchContent_MakeAvailable(transwarp)""",
        "link_libraries": ["transwarp::transwarp"],
        "header_only": True,
        "cpp_standard": 17,
        "alternatives": ["taskflow"],
    },
    {
        "id": "taskflow",
        "name": "Taskflow",
        "description": "Parallel and heterogeneous task programming system",
        "category": "concurrency",
        "tags": ["concurrency", "task-based", "header-only", "gpu"],
        "github_url": "https://github.com/taskflow/taskflow",
        "fetch_content": """FetchContent_Declare(
    taskflow
    GIT_REPOSITORY https://github.com/taskflow/taskflow.git
    GIT_TAG v3.6.0
)
FetchContent_MakeAvailable(taskflow)""",
        "link_libraries": ["Taskflow::Taskflow"],
        "header_only": True,
        "cpp_standard": 17,
        "alternatives": ["transwarp"],
    },
    # Utility Libraries
    {
        "id": "expected_lite",
        "name": "expected-lite",
        "description": "Expected objects for C++11 and later (single-file header-only)",
        "category": "utility",
        "tags": ["utility", "error-handling", "header-only"],
        "github_url": "https://github.com/martinmoene/expected-lite",
        "fetch_content": """FetchContent_Declare(
    expected-lite
    GIT_REPOSITORY https://github.com/martinmoene/expected-lite.git
    GIT_TAG v0.6.3
)
FetchContent_MakeAvailable(expected-lite)""",
        "link_libraries": ["nonstd::expected-lite"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["tl_expected"],
    },
    {
        "id": "tl_expected",
        "name": "tl::expected",
        "description": "C++11/14/17 std::expected with functional-style extensions",
        "category": "utility",
        "tags": ["utility", "error-handling", "header-only"],
        "github_url": "https://github.com/TartanLlama/expected",
        "fetch_content": """FetchContent_Declare(
    tl_expected
    GIT_REPOSITORY https://github.com/TartanLlama/expected.git
    GIT_TAG v1.1.0
)
FetchContent_MakeAvailable(tl_expected)""",
        "link_libraries": ["tl::expected"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["expected_lite"],
    },
    {
        "id": "span_lite",
        "name": "span-lite",
        "description": "A C++20-like span for C++98, C++11 and later",
        "category": "utility",
        "tags": ["utility", "container", "header-only"],
        "github_url": "https://github.com/martinmoene/span-lite",
        "fetch_content": """FetchContent_Declare(
    span-lite
    GIT_REPOSITORY https://github.com/martinmoene/span-lite.git
    GIT_TAG v0.10.3
)
FetchContent_MakeAvailable(span-lite)""",
        "link_libraries": ["nonstd::span-lite"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": [],
    },
    {
        "id": "range_v3",
        "name": "range-v3",
        "description": "Range library for C++14/17/20",
        "category": "utility",
        "tags": ["ranges", "functional", "header-only"],
        "github_url": "https://github.com/ericniebler/range-v3",
        "fetch_content": """FetchContent_Declare(
    range-v3
    GIT_REPOSITORY https://github.com/ericniebler/range-v3.git
    GIT_TAG 0.12.0
)
FetchContent_MakeAvailable(range-v3)""",
        "link_libraries": ["range-v3::range-v3"],
        "header_only": True,
        "cpp_standard": 14,
        "alternatives": [],
    },
    # Database Libraries
    {
        "id": "sqlite_modern_cpp",
        "name": "sqlite_modern_cpp",
        "description": "Modern C++ wrapper around sqlite library",
        "category": "database",
        "tags": ["database", "sqlite", "header-only"],
        "github_url": "https://github.com/SqliteModernCpp/sqlite_modern_cpp",
        "fetch_content": """FetchContent_Declare(
    sqlite_modern_cpp
    GIT_REPOSITORY https://github.com/SqliteModernCpp/sqlite_modern_cpp.git
    GIT_TAG v3.2
)
FetchContent_MakeAvailable(sqlite_modern_cpp)""",
        "link_libraries": ["sqlite_modern_cpp"],
        "header_only": True,
        "cpp_standard": 14,
        "alternatives": ["sqlitecpp"],
    },
    {
        "id": "sqlitecpp",
        "name": "SQLiteCpp",
        "description": "Smart and easy to use C++ SQLite3 wrapper",
        "category": "database",
        "tags": ["database", "sqlite"],
        "github_url": "https://github.com/SRombauts/SQLiteCpp",
        "fetch_content": """FetchContent_Declare(
    SQLiteCpp
    GIT_REPOSITORY https://github.com/SRombauts/SQLiteCpp.git
    GIT_TAG 3.3.1
)
FetchContent_MakeAvailable(SQLiteCpp)""",
        "link_libraries": ["SQLiteCpp"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["sqlite_modern_cpp"],
    },
    # Compression Libraries
    {
        "id": "zlib",
        "name": "zlib",
        "description": "General purpose compression library",
        "category": "compression",
        "tags": ["compression", "zlib"],
        "github_url": "https://github.com/madler/zlib",
        "fetch_content": """FetchContent_Declare(
    zlib
    GIT_REPOSITORY https://github.com/madler/zlib.git
    GIT_TAG v1.3
)
FetchContent_MakeAvailable(zlib)""",
        "link_libraries": ["ZLIB::ZLIB"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["lz4", "zstd"],
    },
    {
        "id": "lz4",
        "name": "LZ4",
        "description": "Extremely fast compression algorithm",
        "category": "compression",
        "tags": ["compression", "fast"],
        "github_url": "https://github.com/lz4/lz4",
        "fetch_content": """FetchContent_Declare(
    lz4
    GIT_REPOSITORY https://github.com/lz4/lz4.git
    GIT_TAG v1.9.4
    SOURCE_SUBDIR build/cmake
)
FetchContent_MakeAvailable(lz4)""",
        "link_libraries": ["lz4_static"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["zlib", "zstd"],
    },
    {
        "id": "zstd",
        "name": "Zstandard",
        "description": "Fast real-time compression algorithm",
        "category": "compression",
        "tags": ["compression", "fast"],
        "github_url": "https://github.com/facebook/zstd",
        "fetch_content": """FetchContent_Declare(
    zstd
    GIT_REPOSITORY https://github.com/facebook/zstd.git
    GIT_TAG v1.5.5
    SOURCE_SUBDIR build/cmake
)
set(ZSTD_BUILD_PROGRAMS OFF)
set(ZSTD_BUILD_TESTS OFF)
FetchContent_MakeAvailable(zstd)""",
        "link_libraries": ["libzstd_static"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["zlib", "lz4"],
    },
    # Math Libraries
    {
        "id": "eigen",
        "name": "Eigen",
        "description": "C++ template library for linear algebra",
        "category": "math",
        "tags": ["math", "linear-algebra", "header-only"],
        "github_url": "https://gitlab.com/libeigen/eigen",
        "fetch_content": """FetchContent_Declare(
    Eigen
    GIT_REPOSITORY https://gitlab.com/libeigen/eigen.git
    GIT_TAG 3.4.0
)
FetchContent_MakeAvailable(Eigen)""",
        "link_libraries": ["Eigen3::Eigen"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["glm"],
    },
    {
        "id": "glm",
        "name": "GLM",
        "description": "OpenGL Mathematics - Header only C++ mathematics library for graphics",
        "category": "math",
        "tags": ["math", "graphics", "header-only"],
        "github_url": "https://github.com/g-truc/glm",
        "fetch_content": """FetchContent_Declare(
    glm
    GIT_REPOSITORY https://github.com/g-truc/glm.git
    GIT_TAG 1.0.0
)
FetchContent_MakeAvailable(glm)""",
        "link_libraries": ["glm::glm"],
        "header_only": True,
        "cpp_standard": 11,
        "alternatives": ["eigen"],
    },
    # Cryptography Libraries
    {
        "id": "openssl",
        "name": "OpenSSL",
        "description": "TLS/SSL and crypto library",
        "category": "cryptography",
        "tags": ["crypto", "ssl", "tls"],
        "github_url": "https://github.com/openssl/openssl",
        "fetch_content": """# OpenSSL is typically found as a system library
find_package(OpenSSL REQUIRED)""",
        "link_libraries": ["OpenSSL::SSL", "OpenSSL::Crypto"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["botan", "cryptopp"],
    },
    {
        "id": "cryptopp",
        "name": "Crypto++",
        "description": "Free C++ class library of cryptographic schemes",
        "category": "cryptography",
        "tags": ["crypto", "encryption"],
        "github_url": "https://github.com/weidai11/cryptopp",
        "fetch_content": """FetchContent_Declare(
    cryptopp
    GIT_REPOSITORY https://github.com/weidai11/cryptopp.git
    GIT_TAG CRYPTOPP_8_9_0
)
FetchContent_MakeAvailable(cryptopp)""",
        "link_libraries": ["cryptopp"],
        "header_only": False,
        "cpp_standard": 11,
        "alternatives": ["openssl", "botan"],
    },
]

CATEGORIES = [
    {"id": "serialization", "name": "Serialization", "icon": "📦", "description": "JSON, XML, Binary serialization"},
    {"id": "logging", "name": "Logging", "icon": "📝", "description": "Logging and diagnostics"},
    {"id": "testing", "name": "Testing", "icon": "🧪", "description": "Unit testing and mocking frameworks"},
    {"id": "networking", "name": "Networking", "icon": "🌐", "description": "HTTP, TCP/UDP, async I/O"},
    {"id": "cli", "name": "CLI", "icon": "💻", "description": "Command line argument parsing"},
    {"id": "configuration", "name": "Configuration", "icon": "⚙️", "description": "Config file parsing (YAML, TOML)"},
    {"id": "gui", "name": "GUI", "icon": "🖼️", "description": "Graphical user interfaces"},
    {"id": "formatting", "name": "Formatting", "icon": "✨", "description": "String formatting and text processing"},
    {"id": "concurrency", "name": "Concurrency", "icon": "⚡", "description": "Threading, async, lock-free structures"},
    {"id": "utility", "name": "Utility", "icon": "🔧", "description": "General utilities and helpers"},
    {"id": "database", "name": "Database", "icon": "💾", "description": "Database clients and ORMs"},
    {"id": "compression", "name": "Compression", "icon": "🗜️", "description": "Data compression libraries"},
    {"id": "math", "name": "Math", "icon": "📐", "description": "Mathematics and linear algebra"},
    {"id": "cryptography", "name": "Cryptography", "icon": "🔐", "description": "Encryption and cryptographic functions"},
]


def get_library_by_id(library_id: str) -> Optional[Library]:
    """Get a library by its ID."""
    for lib in LIBRARIES:
        if lib["id"] == library_id:
            return lib
    return None


def get_libraries_by_category(category: str) -> List[Library]:
    """Get all libraries in a specific category."""
    return [lib for lib in LIBRARIES if lib["category"] == category]


def search_libraries(query: str) -> List[Library]:
    """Search libraries by name, description, or tags."""
    query = query.lower()
    results = []
    for lib in LIBRARIES:
        if (
            query in lib["name"].lower()
            or query in lib["description"].lower()
            or any(query in tag.lower() for tag in lib["tags"])
        ):
            results.append(lib)
    return results

