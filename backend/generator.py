"""
CMake project generator module.
"""

import io
import zipfile
from typing import List, Optional
from libraries import Library, get_library_by_id


def generate_cmake_lists(
    project_name: str,
    cpp_standard: int,
    libraries: List[Library],
    include_tests: bool = True,
) -> str:
    """Generate the main CMakeLists.txt content."""
    
    # Find maximum required C++ standard
    max_standard = cpp_standard
    for lib in libraries:
        if lib["cpp_standard"] > max_standard:
            max_standard = lib["cpp_standard"]
    
    # Separate test libraries from main libraries
    test_libraries = [lib for lib in libraries if lib["category"] == "testing"]
    main_libraries = [lib for lib in libraries if lib["category"] != "testing"]
    
    cmake_content = f"""cmake_minimum_required(VERSION 3.20)
project({project_name} VERSION 1.0.0 LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD {max_standard})
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for IDE support
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Include FetchContent module
include(FetchContent)

"""

    # Add FetchContent declarations for main libraries
    if main_libraries:
        cmake_content += "# ============================================================================\n"
        cmake_content += "# Dependencies\n"
        cmake_content += "# ============================================================================\n\n"
        
        for lib in main_libraries:
            cmake_content += f"# {lib['name']}\n"
            cmake_content += lib["fetch_content"] + "\n\n"

    # Main library target
    cmake_content += f"""# ============================================================================
# Main Library
# ============================================================================

add_library({project_name}_lib
    src/{project_name}.cpp
)

target_include_directories({project_name}_lib
    PUBLIC
        $<BUILD_INTERFACE:${{CMAKE_CURRENT_SOURCE_DIR}}/include>
        $<INSTALL_INTERFACE:include>
)

"""

    # Link main libraries
    if main_libraries:
        link_libs = []
        for lib in main_libraries:
            link_libs.extend(lib["link_libraries"])
        
        if link_libs:
            cmake_content += f"target_link_libraries({project_name}_lib\n"
            cmake_content += "    PUBLIC\n"
            for link_lib in link_libs:
                cmake_content += f"        {link_lib}\n"
            cmake_content += ")\n\n"

    # Main executable
    cmake_content += f"""# ============================================================================
# Main Executable
# ============================================================================

add_executable({project_name} src/main.cpp)
target_link_libraries({project_name} PRIVATE {project_name}_lib)

"""

    # Test configuration
    if include_tests and test_libraries:
        cmake_content += """# ============================================================================
# Testing
# ============================================================================

enable_testing()

"""
        # Add test library FetchContent
        for lib in test_libraries:
            cmake_content += f"# {lib['name']}\n"
            cmake_content += lib["fetch_content"] + "\n\n"

        cmake_content += "add_subdirectory(tests)\n"

    return cmake_content


def generate_test_cmake(
    project_name: str,
    test_libraries: List[Library],
) -> str:
    """Generate the tests/CMakeLists.txt content."""
    
    cmake_content = f"""# Test configuration for {project_name}

add_executable({project_name}_tests
    test_main.cpp
)

target_link_libraries({project_name}_tests
    PRIVATE
        {project_name}_lib
"""
    
    # Add test framework link libraries
    for lib in test_libraries:
        for link_lib in lib["link_libraries"]:
            cmake_content += f"        {link_lib}\n"
    
    cmake_content += ")\n\n"

    # Add test discovery based on framework
    if any(lib["id"] == "googletest" for lib in test_libraries):
        cmake_content += """include(GoogleTest)
gtest_discover_tests({}_tests)
""".format(project_name)
    elif any(lib["id"] == "catch2" for lib in test_libraries):
        cmake_content += """include(CTest)
include(Catch)
catch_discover_tests({}_tests)
""".format(project_name)
    else:
        cmake_content += f"""add_test(NAME {project_name}_tests COMMAND {project_name}_tests)
"""

    return cmake_content


def generate_main_cpp(project_name: str, libraries: List[Library]) -> str:
    """Generate the main.cpp file."""
    
    includes = []
    
    # Add relevant includes based on selected libraries
    for lib in libraries:
        if lib["id"] == "nlohmann_json":
            includes.append('#include <nlohmann/json.hpp>')
        elif lib["id"] == "spdlog":
            includes.append('#include <spdlog/spdlog.h>')
        elif lib["id"] == "fmt":
            includes.append('#include <fmt/format.h>')
        elif lib["id"] == "cli11":
            includes.append('#include <CLI/CLI.hpp>')
        elif lib["id"] == "argparse":
            includes.append('#include <argparse/argparse.hpp>')
    
    includes_str = "\n".join(includes) if includes else ""
    
    # Check if spdlog is included for logging
    has_spdlog = any(lib["id"] == "spdlog" for lib in libraries)
    has_cli11 = any(lib["id"] == "cli11" for lib in libraries)
    has_argparse = any(lib["id"] == "argparse" for lib in libraries)
    
    main_content = f"""#include <{project_name}/{project_name}.hpp>

#include <iostream>
{includes_str}

int main(int argc, char* argv[]) {{
"""

    if has_spdlog:
        main_content += """    spdlog::info("Starting {} v1.0.0");
""".format(project_name)
    
    if has_cli11:
        main_content += f"""
    CLI::App app{{"{project_name} application"}};
    
    std::string name = "World";
    app.add_option("-n,--name", name, "Name to greet");
    
    CLI11_PARSE(app, argc, argv);
"""
    elif has_argparse:
        main_content += f"""
    argparse::ArgumentParser program("{project_name}");
    
    program.add_argument("-n", "--name")
        .default_value(std::string("World"))
        .help("Name to greet");
    
    try {{
        program.parse_args(argc, argv);
    }} catch (const std::exception& err) {{
        std::cerr << err.what() << std::endl;
        std::cerr << program;
        return 1;
    }}
    
    auto name = program.get<std::string>("--name");
"""
    else:
        main_content += """    (void)argc;
    (void)argv;
"""

    main_content += f"""
    {project_name}::greet();
    
    return 0;
}}
"""

    return main_content


def generate_lib_header(project_name: str) -> str:
    """Generate the main library header file."""
    
    guard = f"{project_name.upper()}_HPP"
    
    return f"""#ifndef {guard}
#define {guard}

#include <string>

namespace {project_name} {{

/**
 * @brief Greet function
 */
void greet();

/**
 * @brief Get the library version
 * @return Version string
 */
std::string version();

}}  // namespace {project_name}

#endif  // {guard}
"""


def generate_lib_source(project_name: str, libraries: List[Library]) -> str:
    """Generate the main library source file."""
    
    has_spdlog = any(lib["id"] == "spdlog" for lib in libraries)
    has_fmt = any(lib["id"] == "fmt" for lib in libraries)
    
    includes = [f'#include <{project_name}/{project_name}.hpp>']
    
    if has_spdlog:
        includes.append('#include <spdlog/spdlog.h>')
    if has_fmt:
        includes.append('#include <fmt/format.h>')
    
    includes.append('#include <iostream>')
    
    source = "\n".join(includes) + "\n\n"
    
    source += f"""namespace {project_name} {{

void greet() {{
"""
    
    if has_spdlog:
        source += '    spdlog::info("Hello from {}!");\n'.format(project_name)
    else:
        source += '    std::cout << "Hello from {}!" << std::endl;\n'.format(project_name)
    
    source += """}

std::string version() {
    return "1.0.0";
}

}  // namespace """ + project_name + "\n"

    return source


def generate_test_main(project_name: str, test_libraries: List[Library]) -> str:
    """Generate the test main file."""
    
    if any(lib["id"] == "googletest" for lib in test_libraries):
        return f"""#include <gtest/gtest.h>
#include <{project_name}/{project_name}.hpp>

TEST({project_name.capitalize()}Test, VersionTest) {{
    EXPECT_EQ({project_name}::version(), "1.0.0");
}}

TEST({project_name.capitalize()}Test, GreetTest) {{
    // Should not throw
    EXPECT_NO_THROW({project_name}::greet());
}}
"""
    elif any(lib["id"] == "catch2" for lib in test_libraries):
        return f"""#include <catch2/catch_test_macros.hpp>
#include <{project_name}/{project_name}.hpp>

TEST_CASE("{project_name}::version returns correct version", "[version]") {{
    REQUIRE({project_name}::version() == "1.0.0");
}}

TEST_CASE("{project_name}::greet does not throw", "[greet]") {{
    REQUIRE_NOTHROW({project_name}::greet());
}}
"""
    elif any(lib["id"] == "doctest" for lib in test_libraries):
        return f"""#define DOCTEST_CONFIG_IMPLEMENT_WITH_MAIN
#include <doctest/doctest.h>
#include <{project_name}/{project_name}.hpp>

TEST_CASE("testing version") {{
    CHECK({project_name}::version() == "1.0.0");
}}

TEST_CASE("testing greet") {{
    CHECK_NOTHROW({project_name}::greet());
}}
"""
    else:
        return f"""// Basic test file - add a test framework for better testing support
#include <{project_name}/{project_name}.hpp>
#include <cassert>
#include <iostream>

int main() {{
    assert({project_name}::version() == "1.0.0");
    {project_name}::greet();
    std::cout << "All tests passed!" << std::endl;
    return 0;
}}
"""


def generate_readme(project_name: str, libraries: List[Library], cpp_standard: int) -> str:
    """Generate the README.md file."""
    
    lib_list = "\n".join([f"- [{lib['name']}]({lib['github_url']}) - {lib['description']}" for lib in libraries])
    
    return f"""# {project_name}

A C++ project using modern CMake and FetchContent for dependency management.

## Requirements

- CMake 3.20 or higher
- C++{cpp_standard} compatible compiler

## Dependencies

{lib_list if lib_list else "No external dependencies."}

## Building

```bash
mkdir build && cd build
cmake ..
cmake --build .
```

## Running

```bash
./build/{project_name}
```

## Testing

```bash
cd build
ctest --output-on-failure
```

## Project Structure

```
{project_name}/
├── CMakeLists.txt
├── include/
│   └── {project_name}/
│       └── {project_name}.hpp
├── src/
│   ├── main.cpp
│   └── {project_name}.cpp
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
└── README.md
```

## License

MIT License
"""


def generate_gitignore() -> str:
    """Generate .gitignore file."""
    return """# Build directories
build/
cmake-build-*/
out/

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# Compiled files
*.o
*.obj
*.a
*.lib
*.so
*.dylib
*.dll

# CMake
CMakeFiles/
CMakeCache.txt
cmake_install.cmake
Makefile
compile_commands.json

# Testing
Testing/

# Package
*.zip
*.tar.gz
"""


def generate_clang_format() -> str:
    """Generate .clang-format file (Google style)."""
    return """BasedOnStyle: Google
IndentWidth: 4
ColumnLimit: 100
AllowShortFunctionsOnASingleLine: Empty
AllowShortIfStatementsOnASingleLine: Never
AllowShortLoopsOnASingleLine: false
BreakBeforeBraces: Attach
PointerAlignment: Left
SpaceAfterCStyleCast: false
SpaceBeforeParens: ControlStatements
"""


def create_project_zip(
    project_name: str,
    cpp_standard: int,
    library_ids: List[str],
    include_tests: bool = True,
) -> bytes:
    """Create a ZIP file containing the complete project."""
    
    # Get library objects
    libraries = []
    for lib_id in library_ids:
        lib = get_library_by_id(lib_id)
        if lib:
            libraries.append(lib)
    
    # Separate test libraries
    test_libraries = [lib for lib in libraries if lib["category"] == "testing"]
    
    # Create in-memory ZIP file
    zip_buffer = io.BytesIO()
    
    with zipfile.ZipFile(zip_buffer, 'w', zipfile.ZIP_DEFLATED) as zf:
        base_path = project_name
        
        # CMakeLists.txt
        zf.writestr(
            f"{base_path}/CMakeLists.txt",
            generate_cmake_lists(project_name, cpp_standard, libraries, include_tests)
        )
        
        # README.md
        zf.writestr(
            f"{base_path}/README.md",
            generate_readme(project_name, libraries, cpp_standard)
        )
        
        # .gitignore
        zf.writestr(
            f"{base_path}/.gitignore",
            generate_gitignore()
        )
        
        # .clang-format
        zf.writestr(
            f"{base_path}/.clang-format",
            generate_clang_format()
        )
        
        # Include directory
        zf.writestr(
            f"{base_path}/include/{project_name}/{project_name}.hpp",
            generate_lib_header(project_name)
        )
        
        # Source directory
        zf.writestr(
            f"{base_path}/src/main.cpp",
            generate_main_cpp(project_name, libraries)
        )
        zf.writestr(
            f"{base_path}/src/{project_name}.cpp",
            generate_lib_source(project_name, libraries)
        )
        
        # Tests directory
        if include_tests:
            zf.writestr(
                f"{base_path}/tests/CMakeLists.txt",
                generate_test_cmake(project_name, test_libraries)
            )
            zf.writestr(
                f"{base_path}/tests/test_main.cpp",
                generate_test_main(project_name, test_libraries)
            )
    
    zip_buffer.seek(0)
    return zip_buffer.getvalue()

