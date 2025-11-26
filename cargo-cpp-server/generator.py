"""
CMake project generator module.
"""

import io
import zipfile
from typing import List, Dict, Any, Tuple, Optional
from recipe_loader import Library, get_library_by_id


def generate_cmake_lists(
    project_name: str,
    cpp_standard: int,
    libraries_with_options: List[Tuple[Library, Dict[str, Any]]],
    include_tests: bool = True,
    testing_framework: str = "googletest",
    build_shared: bool = False,
) -> str:
    """Generate the main CMakeLists.txt content.
    
    Args:
        project_name: Name of the project.
        cpp_standard: C++ standard version (11, 14, 17, 20, 23).
        libraries_with_options: List of (Library, options_dict) tuples.
        include_tests: Whether to include test configuration.
        testing_framework: Testing framework to use (none, googletest, catch2, doctest).
        build_shared: Whether to build shared libraries.
    
    Returns:
        Generated CMakeLists.txt content.
    """
    
    # Find maximum required C++ standard
    max_standard = cpp_standard
    for lib, _ in libraries_with_options:
        if lib.get("cpp_standard", 11) > max_standard:
            max_standard = lib["cpp_standard"]
    
    # Separate test libraries from main libraries
    test_libraries = [(lib, opts) for lib, opts in libraries_with_options if lib["category"] == "testing"]
    main_libraries = [(lib, opts) for lib, opts in libraries_with_options if lib["category"] != "testing"]
    
    # Add selected testing framework if not already present
    if include_tests and testing_framework and testing_framework != "none":
        existing_test_ids = [lib["id"] for lib, _ in test_libraries]
        if testing_framework not in existing_test_ids:
            test_lib = get_library_by_id(testing_framework)
            if test_lib:
                test_libraries.insert(0, (test_lib, {}))
                # Update max standard if needed
                if test_lib.get("cpp_standard", 11) > max_standard:
                    max_standard = test_lib["cpp_standard"]
    
    cmake_content = f"""cmake_minimum_required(VERSION 3.20)
project({project_name} VERSION 1.0.0 LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD {max_standard})
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for IDE support
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Build options
option(BUILD_SHARED_LIBS "Build shared libraries" {"ON" if build_shared else "OFF"})

# Include FetchContent module
include(FetchContent)

"""

    # Add FetchContent declarations for main libraries
    if main_libraries:
        cmake_content += "# ============================================================================\n"
        cmake_content += "# Dependencies\n"
        cmake_content += "# ============================================================================\n\n"
        
        for lib, options in main_libraries:
            cmake_content += generate_library_cmake(lib, options)
            cmake_content += "\n"

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
    link_libs = collect_link_libraries(main_libraries)
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
        for lib, options in test_libraries:
            cmake_content += generate_library_cmake(lib, options)
            cmake_content += "\n"

        cmake_content += "add_subdirectory(tests)\n"

    return cmake_content


def generate_library_cmake(lib: Library, options: Dict[str, Any]) -> str:
    """Generate CMake code for a single library.
    
    Args:
        lib: Library definition.
        options: User-selected options for this library.
        
    Returns:
        CMake code string.
    """
    result = f"# {lib['name']}\n"
    
    # Generate CMake variables from options
    for opt in lib.get("options", []):
        opt_id = opt["id"]
        opt_value = options.get(opt_id, opt.get("default"))
        
        if opt_value is None:
            continue
            
        # Handle cmake_var
        if "cmake_var" in opt:
            if opt["type"] == "boolean":
                cmake_val = "ON" if opt_value else "OFF"
                result += f"set({opt['cmake_var']} {cmake_val})\n"
            elif opt["type"] == "string" and opt_value:
                result += f'set({opt["cmake_var"]} "{opt_value}")\n'
            elif opt["type"] == "integer":
                result += f"set({opt['cmake_var']} {opt_value})\n"
            elif opt["type"] == "choice":
                result += f'set({opt["cmake_var"]} "{opt_value}")\n'
    
    # Add cmake_pre if present
    if "cmake_pre" in lib:
        result += lib["cmake_pre"].strip() + "\n"
    
    # System package (find_package)
    if lib.get("system_package", False):
        pkg_name = lib.get("find_package_name", lib["name"])
        result += f"find_package({pkg_name} REQUIRED)\n"
    else:
        # FetchContent
        fc = lib.get("fetch_content", {})
        if fc:
            result += "FetchContent_Declare(\n"
            result += f"    {lib['id']}\n"
            result += f"    GIT_REPOSITORY {fc['repository']}\n"
            result += f"    GIT_TAG {fc['tag']}\n"
            if "source_subdir" in fc:
                result += f"    SOURCE_SUBDIR {fc['source_subdir']}\n"
            result += ")\n"
            result += f"FetchContent_MakeAvailable({lib['id']})\n"
    
    # Add cmake_post if present
    if "cmake_post" in lib:
        result += lib["cmake_post"].strip() + "\n"
    
    # Generate compile definitions from options
    for opt in lib.get("options", []):
        opt_id = opt["id"]
        opt_value = options.get(opt_id, opt.get("default"))
        
        if "cmake_define" in opt and opt_value:
            if opt["type"] == "boolean" and opt_value:
                result += f"add_compile_definitions({opt['cmake_define']})\n"
            elif opt["type"] == "integer" and opt_value:
                result += f"add_compile_definitions({opt['cmake_define']}={opt_value})\n"
            elif opt["type"] in ("string", "choice") and opt_value:
                result += f'add_compile_definitions({opt["cmake_define"]}={opt_value})\n'
    
    return result


def collect_link_libraries(libraries_with_options: List[Tuple[Library, Dict[str, Any]]]) -> List[str]:
    """Collect all link libraries from library selections.
    
    Args:
        libraries_with_options: List of (Library, options_dict) tuples.
        
    Returns:
        List of library names to link.
    """
    link_libs = []
    
    for lib, options in libraries_with_options:
        # Base link libraries
        link_libs.extend(lib.get("link_libraries", []))
        
        # Check options that affect linking
        for opt in lib.get("options", []):
            opt_id = opt["id"]
            opt_value = options.get(opt_id, opt.get("default"))
            
            if opt.get("affects_link", False) and opt_value:
                additional_libs = opt.get("link_libraries_when_enabled", [])
                link_libs.extend(additional_libs)
    
    # Remove duplicates while preserving order
    seen = set()
    unique_libs = []
    for lib in link_libs:
        if lib not in seen:
            seen.add(lib)
            unique_libs.append(lib)
    
    return unique_libs


def generate_test_cmake(
    project_name: str,
    test_libraries: List[Tuple[Library, Dict[str, Any]]],
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
    link_libs = collect_link_libraries(test_libraries)
    for lib in link_libs:
        cmake_content += f"        {lib}\n"
    
    cmake_content += ")\n\n"

    # Add test discovery based on framework
    has_gtest = any(lib["id"] == "googletest" for lib, _ in test_libraries)
    has_catch2 = any(lib["id"] == "catch2" for lib, _ in test_libraries)
    
    if has_gtest:
        cmake_content += f"""include(GoogleTest)
gtest_discover_tests({project_name}_tests)
"""
    elif has_catch2:
        cmake_content += f"""include(CTest)
include(Catch)
catch_discover_tests({project_name}_tests)
"""
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
        main_content += f'    spdlog::info("Starting {project_name} v1.0.0");\n'
    
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
        source += f'    spdlog::info("Hello from {project_name}!");\n'
    else:
        source += f'    std::cout << "Hello from {project_name}!" << std::endl;\n'
    
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


CLANG_FORMAT_STYLES = {
    "Google": """BasedOnStyle: Google
IndentWidth: 4
ColumnLimit: 100
AllowShortFunctionsOnASingleLine: Empty
AllowShortIfStatementsOnASingleLine: Never
AllowShortLoopsOnASingleLine: false
BreakBeforeBraces: Attach
PointerAlignment: Left
SpaceAfterCStyleCast: false
SpaceBeforeParens: ControlStatements
""",
    "LLVM": """BasedOnStyle: LLVM
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: All
AllowShortIfStatementsOnASingleLine: Never
BreakBeforeBraces: Attach
PointerAlignment: Right
SpaceBeforeParens: ControlStatements
""",
    "Chromium": """BasedOnStyle: Chromium
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: Inline
AllowShortIfStatementsOnASingleLine: Never
BreakBeforeBraces: Attach
PointerAlignment: Left
DerivePointerAlignment: false
""",
    "Mozilla": """BasedOnStyle: Mozilla
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: Inline
BreakBeforeBraces: Mozilla
PointerAlignment: Left
AlwaysBreakAfterDefinitionReturnType: TopLevel
""",
    "WebKit": """BasedOnStyle: WebKit
IndentWidth: 4
ColumnLimit: 0
AllowShortFunctionsOnASingleLine: All
BreakBeforeBraces: WebKit
PointerAlignment: Left
NamespaceIndentation: Inner
""",
    "Microsoft": """BasedOnStyle: Microsoft
IndentWidth: 4
ColumnLimit: 120
AllowShortFunctionsOnASingleLine: None
BreakBeforeBraces: Allman
PointerAlignment: Left
AccessModifierOffset: -4
AlignAfterOpenBracket: Align
""",
    "GNU": """BasedOnStyle: GNU
IndentWidth: 2
ColumnLimit: 79
AllowShortFunctionsOnASingleLine: None
BreakBeforeBraces: GNU
PointerAlignment: Right
SpaceBeforeParens: Always
""",
}


def generate_clang_format(style: str = "Google") -> str:
    """Generate .clang-format file with specified style.
    
    Args:
        style: One of Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU
        
    Returns:
        .clang-format file content
    """
    return CLANG_FORMAT_STYLES.get(style, CLANG_FORMAT_STYLES["Google"])


def create_project_zip(
    project_name: str,
    cpp_standard: int,
    library_selections: List[Any],  # List of LibrarySelection from pydantic
    include_tests: bool = True,
    testing_framework: str = "googletest",
    build_shared: bool = False,
    clang_format_style: str = "Google",
) -> bytes:
    """Create a ZIP file containing the complete project.
    
    Args:
        project_name: Name of the project.
        cpp_standard: C++ standard version.
        library_selections: List of library selections with options.
        include_tests: Whether to include test configuration.
        testing_framework: Testing framework (none, googletest, catch2, doctest).
        build_shared: Whether to build shared libraries.
        clang_format_style: Clang-format style (Google, LLVM, etc.).
        
    Returns:
        ZIP file content as bytes.
    """
    
    # Get library objects with their options
    libraries_with_options: List[Tuple[Library, Dict[str, Any]]] = []
    all_libraries: List[Library] = []
    
    for selection in library_selections:
        lib = get_library_by_id(selection.library_id)
        if lib:
            libraries_with_options.append((lib, selection.options))
            all_libraries.append(lib)
    
    # Separate test libraries
    test_libraries = [(lib, opts) for lib, opts in libraries_with_options if lib["category"] == "testing"]
    
    # Add selected testing framework if not already present
    if include_tests and testing_framework and testing_framework != "none":
        existing_test_ids = [lib["id"] for lib, _ in test_libraries]
        if testing_framework not in existing_test_ids:
            test_lib = get_library_by_id(testing_framework)
            if test_lib:
                test_libraries.insert(0, (test_lib, {}))
    
    test_libs_only = [lib for lib, _ in test_libraries]
    
    # Create in-memory ZIP file
    zip_buffer = io.BytesIO()
    
    with zipfile.ZipFile(zip_buffer, 'w', zipfile.ZIP_DEFLATED) as zf:
        base_path = project_name
        
        # CMakeLists.txt
        zf.writestr(
            f"{base_path}/CMakeLists.txt",
            generate_cmake_lists(project_name, cpp_standard, libraries_with_options, include_tests, testing_framework, build_shared)
        )
        
        # README.md
        zf.writestr(
            f"{base_path}/README.md",
            generate_readme(project_name, all_libraries, cpp_standard)
        )
        
        # .gitignore
        zf.writestr(
            f"{base_path}/.gitignore",
            generate_gitignore()
        )
        
        # .clang-format
        zf.writestr(
            f"{base_path}/.clang-format",
            generate_clang_format(clang_format_style)
        )
        
        # Include directory
        zf.writestr(
            f"{base_path}/include/{project_name}/{project_name}.hpp",
            generate_lib_header(project_name)
        )
        
        # Source directory
        zf.writestr(
            f"{base_path}/src/main.cpp",
            generate_main_cpp(project_name, all_libraries)
        )
        zf.writestr(
            f"{base_path}/src/{project_name}.cpp",
            generate_lib_source(project_name, all_libraries)
        )
        
        # Tests directory
        if include_tests:
            zf.writestr(
                f"{base_path}/tests/CMakeLists.txt",
                generate_test_cmake(project_name, test_libraries)
            )
            zf.writestr(
                f"{base_path}/tests/test_main.cpp",
                generate_test_main(project_name, test_libs_only)
            )
    
    zip_buffer.seek(0)
    return zip_buffer.getvalue()
