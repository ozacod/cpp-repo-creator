package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateVersionHpp generates version.hpp directly from project name and version
func generateVersionHpp(projectName, projectVersion string) string {
	if projectVersion == "" {
		projectVersion = "1.0.0"
	}

	// Parse version components
	parts := strings.Split(projectVersion, ".")
	major := "0"
	minor := "0"
	patch := "0"
	if len(parts) > 0 {
		major = parts[0]
	}
	if len(parts) > 1 {
		minor = parts[1]
	}
	if len(parts) > 2 {
		patch = parts[2]
	}

	projectNameUpper := strings.ToUpper(projectName)
	guard := projectNameUpper + "_VERSION_H_"

	return fmt.Sprintf(`#ifndef %s
#define %s

#define %s_VERSION "%s"
#define %s_MAJOR_VERSION %s
#define %s_MINOR_VERSION %s
#define %s_PATCH_VERSION %s

#endif  // %s
`, guard, guard, projectNameUpper, projectVersion, projectNameUpper, major, projectNameUpper, minor, projectNameUpper, patch, guard)
}

// generateProjectFiles generates all project files locally (except dependencies.cmake)
func generateProjectFiles(config ForgeConfig, outputDir string, dependenciesCMake string) error {
	projectName := config.Package.Name
	if projectName == "" {
		projectName = "my_project"
	}

	projectVersion := config.Package.Version
	if projectVersion == "" {
		projectVersion = "1.0.0"
	}

	cppStandard := config.Package.CppStandard
	if cppStandard == 0 {
		cppStandard = 17
	}

	projectType := "exe"
	if config.Build.SharedLibs {
		projectType = "lib"
	}

	includeTests := config.Testing.Framework != "" && config.Testing.Framework != "none"
	testingFramework := config.Testing.Framework
	if testingFramework == "" {
		testingFramework = "none"
	}

	buildShared := config.Build.SharedLibs

	// Get library IDs from dependencies
	libraryIDs := make([]string, 0, len(config.Dependencies))
	for libID := range config.Dependencies {
		libraryIDs = append(libraryIDs, libID)
	}

	// Create directories
	dirs := []string{
		".cmake/forge",
		"include/" + projectName,
		"src",
		"tests",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(outputDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Write dependencies.cmake (from server)
	if err := os.WriteFile(
		filepath.Join(outputDir, ".cmake/forge/dependencies.cmake"),
		[]byte(dependenciesCMake),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write dependencies.cmake: %w", err)
	}

	// Generate and write version.hpp directly (no CMake pipeline needed)
	versionHpp := generateVersionHpp(projectName, projectVersion)
	if err := os.WriteFile(
		filepath.Join(outputDir, "include/"+projectName+"/version.hpp"),
		[]byte(versionHpp),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write version.hpp: %w", err)
	}

	// Generate and write CMakeLists.txt
	cmakeLists, err := generateCMakeLists(projectName, cppStandard, libraryIDs, includeTests, testingFramework, buildShared, projectType, projectVersion)
	if err != nil {
		return fmt.Errorf("failed to generate CMakeLists.txt: %w", err)
	}
	if err := os.WriteFile(
		filepath.Join(outputDir, "CMakeLists.txt"),
		[]byte(cmakeLists),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write CMakeLists.txt: %w", err)
	}

	// Generate and write header file (always generated for both exe and lib)
	libHeader := generateLibHeader(projectName)
	if err := os.WriteFile(
		filepath.Join(outputDir, "include/"+projectName+"/"+projectName+".hpp"),
		[]byte(libHeader),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate and write main.cpp for executable projects
	if projectType == "exe" {
		mainCpp := generateMainCpp(projectName, libraryIDs)
		if err := os.WriteFile(
			filepath.Join(outputDir, "src/main.cpp"),
			[]byte(mainCpp),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write main.cpp: %w", err)
		}
	}

	// Generate and write project source file (always generated, uses libSource which includes version())
	libSource := generateLibSource(projectName, libraryIDs)
	if err := os.WriteFile(
		filepath.Join(outputDir, "src/"+projectName+".cpp"),
		[]byte(libSource),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write project source: %w", err)
	}

	// Generate and write README.md
	readme := generateReadme(projectName, libraryIDs, cppStandard, projectType)
	if err := os.WriteFile(
		filepath.Join(outputDir, "README.md"),
		[]byte(readme),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Generate and write .gitignore
	gitignore := generateGitignore()
	if err := os.WriteFile(
		filepath.Join(outputDir, ".gitignore"),
		[]byte(gitignore),
		0644,
	); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Generate test files if needed
	if includeTests {
		testCMake := generateTestCMake(projectName, libraryIDs, testingFramework)
		if err := os.WriteFile(
			filepath.Join(outputDir, "tests/CMakeLists.txt"),
			[]byte(testCMake),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write tests/CMakeLists.txt: %w", err)
		}

		testMain := generateTestMain(projectName, libraryIDs, testingFramework)
		if err := os.WriteFile(
			filepath.Join(outputDir, "tests/test_main.cpp"),
			[]byte(testMain),
			0644,
		); err != nil {
			return fmt.Errorf("failed to write tests/test_main.cpp: %w", err)
		}
	}

	return nil
}

// Generation functions (simplified versions that work with library IDs only)

func generateVersionCMake(projectVersion string) string {
	if projectVersion == "" {
		projectVersion = "1.0.0"
	}
	return fmt.Sprintf(`# =============================================================================
# Version from forge.yaml
# =============================================================================
# This file is auto-generated by Forge from forge.yaml package.version
# Do not edit manually. Regenerate with 'forge generate' to update.

set(FORGE_PROJECT_VERSION "%s")
`, projectVersion)
}

func generateConfigureVersionCMake() string {
	return `# =============================================================================
# Configure Version Header Script
# =============================================================================
# This script is used to regenerate version.hpp when version.cmake changes

# Get the source directory (where version.cmake is located)
# FORGE_SOURCE_DIR is passed as -DFORGE_SOURCE_DIR=<path>
if(NOT DEFINED FORGE_SOURCE_DIR)
    set(FORGE_SOURCE_DIR "")
endif()

# CRITICAL FIX: Remove literal quotes if they were passed in the argument.
# This converts "/path/to/src" (relative string) back to /path/to/src (absolute path)
string(REPLACE "\"" "" FORGE_SOURCE_DIR "${FORGE_SOURCE_DIR}")

if("${FORGE_SOURCE_DIR}" STREQUAL "")
    if(DEFINED CMAKE_SCRIPT_MODE_FILE)
        get_filename_component(SCRIPT_DIR "${CMAKE_SCRIPT_MODE_FILE}" DIRECTORY)
        get_filename_component(FORGE_SOURCE_DIR "${SCRIPT_DIR}/../.." ABSOLUTE)
    else()
        message(FATAL_ERROR "FORGE_SOURCE_DIR must be set")
    endif()
endif()

# Ensure we have an absolute path
get_filename_component(FORGE_SOURCE_DIR "${FORGE_SOURCE_DIR}" ABSOLUTE)

# Include version from forge.yaml
include("${FORGE_SOURCE_DIR}/.cmake/forge/version.cmake")

# Parse version components
string(REGEX REPLACE "^([0-9]+)\\..*" "\\1" PROJECT_VERSION_MAJOR "${FORGE_PROJECT_VERSION}")
string(REGEX REPLACE "^[0-9]+\\.([0-9]+).*" "\\1" PROJECT_VERSION_MINOR "${FORGE_PROJECT_VERSION}")
string(REGEX REPLACE "^[0-9]+\\.[0-9]+\\.([0-9]+).*" "\\1" PROJECT_VERSION_PATCH "${FORGE_PROJECT_VERSION}")

# Set default values if parsing failed
if("${PROJECT_VERSION_MAJOR}" STREQUAL "")
    set(PROJECT_VERSION_MAJOR "0")
endif()
if("${PROJECT_VERSION_MINOR}" STREQUAL "")
    set(PROJECT_VERSION_MINOR "0")
endif()
if("${PROJECT_VERSION_PATCH}" STREQUAL "")
    set(PROJECT_VERSION_PATCH "0")
endif()

# Set PROJECT_VERSION for template
set(PROJECT_VERSION "${FORGE_PROJECT_VERSION}")

# PROJECT_NAME and PROJECT_NAME_UPPERCASE should be passed as -D parameters
if(NOT DEFINED PROJECT_NAME)
    message(FATAL_ERROR "PROJECT_NAME must be set via -DPROJECT_NAME=<name>")
endif()

if(NOT DEFINED PROJECT_NAME_UPPERCASE)
    string(TOUPPER "${PROJECT_NAME}" PROJECT_NAME_UPPERCASE)
endif()

# Remove old file to force regeneration (ensures content is updated)
file(REMOVE "${FORGE_SOURCE_DIR}/include/${PROJECT_NAME}/version.hpp")

# Configure the file - always overwrite to ensure it's updated
configure_file(
    "${FORGE_SOURCE_DIR}/.cmake/forge/version.hpp.in"
    "${FORGE_SOURCE_DIR}/include/${PROJECT_NAME}/version.hpp"
    @ONLY
    NEWLINE_STYLE UNIX
)
`
}

func generateUtilsCMake() string {
	return `# =============================================================================
# Forge Utils - Version Header Generation
# =============================================================================
# This file is auto-generated by Forge. Do not edit manually.
# Regenerate with 'forge generate' to update.

function(forge_configure_version_header PROJECT_NAME)
    # Include version from forge.yaml (generated by forge generate)
    include(${CMAKE_CURRENT_SOURCE_DIR}/.cmake/forge/version.cmake)
    
    string(TOUPPER "${PROJECT_NAME}" PROJECT_NAME_UPPERCASE)
    
    # Parse version components from FORGE_PROJECT_VERSION
    string(REGEX REPLACE "^([0-9]+)\\..*" "\\1" PROJECT_VERSION_MAJOR "${FORGE_PROJECT_VERSION}")
    string(REGEX REPLACE "^[0-9]+\\.([0-9]+).*" "\\1" PROJECT_VERSION_MINOR "${FORGE_PROJECT_VERSION}")
    string(REGEX REPLACE "^[0-9]+\\.[0-9]+\\.([0-9]+).*" "\\1" PROJECT_VERSION_PATCH "${FORGE_PROJECT_VERSION}")
    
    # Set default values if parsing failed
    if("${PROJECT_VERSION_MAJOR}" STREQUAL "")
        set(PROJECT_VERSION_MAJOR "0")
    endif()
    if("${PROJECT_VERSION_MINOR}" STREQUAL "")
        set(PROJECT_VERSION_MINOR "0")
    endif()
    if("${PROJECT_VERSION_PATCH}" STREQUAL "")
        set(PROJECT_VERSION_PATCH "0")
    endif()
    
    # Set PROJECT_VERSION for template substitution
    set(PROJECT_VERSION "${FORGE_PROJECT_VERSION}")
    
    # Add custom command to generate/regenerate version.hpp when version.cmake changes
    # This ensures the header is regenerated at build time if version changes
    # The OUTPUT must be a file that will be used by the target
    # Store source directory in a variable for the custom command
    set(FORGE_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}")
    
    add_custom_command(
        OUTPUT "${CMAKE_CURRENT_SOURCE_DIR}/include/${PROJECT_NAME}/version.hpp"
        COMMAND ${CMAKE_COMMAND}
            -DFORGE_SOURCE_DIR="${FORGE_SOURCE_DIR}"
            -DPROJECT_NAME="${PROJECT_NAME}"
            -DPROJECT_NAME_UPPERCASE="${PROJECT_NAME_UPPERCASE}"
            -P "${FORGE_SOURCE_DIR}/.cmake/forge/configure_version.cmake"
        DEPENDS
            "${FORGE_SOURCE_DIR}/.cmake/forge/version.cmake"
            "${FORGE_SOURCE_DIR}/.cmake/forge/version.hpp.in"
        COMMENT "Regenerating version.hpp from forge.yaml"
        VERBATIM
    )
    
    # Also configure at configure time for initial generation
    # This ensures the file exists even if custom command hasn't run yet
    if(NOT EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/include/${PROJECT_NAME}/version.hpp")
        configure_file(
            "${CMAKE_CURRENT_SOURCE_DIR}/.cmake/forge/version.hpp.in"
            "${CMAKE_CURRENT_SOURCE_DIR}/include/${PROJECT_NAME}/version.hpp"
            @ONLY
        )
    endif()
endfunction()
`
}

func generateVersionHppIn() string {
	return `#ifndef @PROJECT_NAME_UPPERCASE@_VERSION_H_
#define @PROJECT_NAME_UPPERCASE@_VERSION_H_

#define @PROJECT_NAME_UPPERCASE@_VERSION "@PROJECT_VERSION@"
#define @PROJECT_NAME_UPPERCASE@_MAJOR_VERSION @PROJECT_VERSION_MAJOR@
#define @PROJECT_NAME_UPPERCASE@_MINOR_VERSION @PROJECT_VERSION_MINOR@
#define @PROJECT_NAME_UPPERCASE@_PATCH_VERSION @PROJECT_VERSION_PATCH@

#endif  // @PROJECT_NAME_UPPERCASE@_VERSION_H_
`
}

func generateCMakeLists(projectName string, cppStandard int, libraryIDs []string, includeTests bool, testingFramework string, buildShared bool, projectType string, projectVersion string) (string, error) {
	buildSharedStr := "OFF"
	if buildShared {
		buildSharedStr = "ON"
	}

	if projectVersion == "" {
		projectVersion = "1.0.0"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`cmake_minimum_required(VERSION 3.20)
project(%s VERSION %s LANGUAGES CXX)

# Set C++ standard
set(CMAKE_CXX_STANDARD %d)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for IDE support
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Build options
option(BUILD_SHARED_LIBS "Build shared libraries" %s)

# =============================================================================
# Dependencies (managed by Forge - regenerate with 'forge generate')
# =============================================================================
include(${CMAKE_CURRENT_SOURCE_DIR}/.cmake/forge/dependencies.cmake)

`, projectName, projectVersion, cppStandard, buildSharedStr, projectName))

	if projectType == "exe" {
		sb.WriteString(fmt.Sprintf(`# =============================================================================
# Main Executable
# =============================================================================

add_executable(%s
    src/main.cpp
    src/%s.cpp
)

target_include_directories(%s
    PRIVATE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
)

target_link_libraries(%s
    PRIVATE
        ${FORGE_LINK_LIBRARIES}
)

`, projectName, projectName, projectName, projectName, projectName))
	} else {
		sb.WriteString(fmt.Sprintf(`# =============================================================================
# Main Library
# =============================================================================

add_library(%s
    src/%s.cpp
)

target_include_directories(%s
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
)

target_link_libraries(%s
    PUBLIC
        ${FORGE_LINK_LIBRARIES}
)

# =============================================================================
# Installation
# =============================================================================

install(TARGETS %s
    EXPORT %sTargets
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
    INCLUDES DESTINATION include
)

install(DIRECTORY include/ DESTINATION include)

`, projectName, projectName, projectName, projectName, projectName, projectName, projectName))
	}

	// Test configuration
	if includeTests {
		sb.WriteString(`# =============================================================================
# Testing
# =============================================================================

enable_testing()

add_subdirectory(tests)
`)
	}

	return sb.String(), nil
}

func generateMainCpp(projectName string, libraryIDs []string) string {
	var includes []string
	hasSpdlog := false
	hasCLI11 := false
	hasArgparse := false

	for _, libID := range libraryIDs {
		switch libID {
		case "nlohmann_json":
			includes = append(includes, "#include <nlohmann/json.hpp>")
		case "spdlog":
			includes = append(includes, "#include <spdlog/spdlog.h>")
			hasSpdlog = true
		case "fmt":
			includes = append(includes, "#include <fmt/format.h>")
		case "cli11":
			includes = append(includes, "#include <CLI/CLI.hpp>")
			hasCLI11 = true
		case "argparse":
			includes = append(includes, "#include <argparse/argparse.hpp>")
			hasArgparse = true
		}
	}

	includesStr := strings.Join(includes, "\n")
	if includesStr != "" {
		includesStr = "\n" + includesStr
	}

	var sb strings.Builder
	projectNameUpper := strings.ToUpper(projectName)
	versionMacro := projectNameUpper + "_VERSION"
	sb.WriteString(fmt.Sprintf(`#include <%s/%s.hpp>
#include <%s/version.hpp>
#include <iostream>%s

int main(int argc, char* argv[]) {
`, projectName, projectName, projectName, includesStr))

	if hasSpdlog {
		sb.WriteString(fmt.Sprintf(`    spdlog::info("Starting %s {}", %s);
`, projectName, versionMacro))
	} else {
		sb.WriteString(fmt.Sprintf(`    std::cout << "Starting %s " << %s << std::endl;
`, projectName, versionMacro))
	}

	if hasCLI11 {
		sb.WriteString(fmt.Sprintf(`
    CLI::App app{"%s application"};
    
    std::string name = "World";
    app.add_option("-n,--name", name, "Name to greet");
    
    CLI11_PARSE(app, argc, argv);
`, projectName))
	} else if hasArgparse {
		sb.WriteString(fmt.Sprintf(`
    argparse::ArgumentParser program("%s");
    
    program.add_argument("-n", "--name")
        .default_value(std::string("World"))
        .help("Name to greet");
    
    try {
        program.parse_args(argc, argv);
    } catch (const std::exception& err) {
        std::cerr << err.what() << std::endl;
        std::cerr << program;
        return 1;
    }
    
    auto name = program.get<std::string>("--name");
`, projectName))
	} else {
		sb.WriteString(`    (void)argc;
    (void)argv;
`)
	}

	sb.WriteString(fmt.Sprintf(`
    %s::greet();
    
    return 0;
}
`, projectName))

	return sb.String()
}

func generateLibHeader(projectName string) string {
	guard := strings.ToUpper(projectName) + "_HPP"
	return fmt.Sprintf(`#ifndef %s
#define %s

#include <string>

namespace %s {

/**
 * @brief Greet function
 */
void greet();

/**
 * @brief Get the library version
 * @return Version string
 */
std::string version();

}  // namespace %s

#endif  // %s
`, guard, guard, projectName, projectName, guard)
}

func generateLibSource(projectName string, libraryIDs []string) string {
	hasSpdlog := false
	hasFmt := false

	for _, libID := range libraryIDs {
		switch libID {
		case "spdlog":
			hasSpdlog = true
		case "fmt":
			hasFmt = true
		}
	}

	var includes []string
	includes = append(includes, fmt.Sprintf("#include <%s/%s.hpp>", projectName, projectName))

	if hasSpdlog {
		includes = append(includes, "#include <spdlog/spdlog.h>")
	}
	if hasFmt {
		includes = append(includes, "#include <fmt/format.h>")
	}
	includes = append(includes, "#include <iostream>")

	var sb strings.Builder
	sb.WriteString(strings.Join(includes, "\n"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("namespace %s {\n\n", projectName))
	sb.WriteString("void greet() {\n")

	if hasSpdlog {
		sb.WriteString(fmt.Sprintf(`    spdlog::info("Hello from %s!");
`, projectName))
	} else {
		sb.WriteString(fmt.Sprintf(`    std::cout << "Hello from %s!" << std::endl;
`, projectName))
	}

	sb.WriteString(`}

std::string version() {
    return "1.0.0";
}

}  // namespace ` + projectName + "\n")

	return sb.String()
}

func generateProjectCpp(projectName string, libraryIDs []string) string {
	hasSpdlog := false
	hasFmt := false

	for _, libID := range libraryIDs {
		switch libID {
		case "spdlog":
			hasSpdlog = true
		case "fmt":
			hasFmt = true
		}
	}

	var includes []string
	includes = append(includes, fmt.Sprintf("#include <%s/%s.hpp>", projectName, projectName))

	if hasSpdlog {
		includes = append(includes, "#include <spdlog/spdlog.h>")
	}
	if hasFmt {
		includes = append(includes, "#include <fmt/format.h>")
	}
	includes = append(includes, "#include <iostream>")

	var sb strings.Builder
	sb.WriteString(strings.Join(includes, "\n"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("namespace %s {\n\n", projectName))
	sb.WriteString("void greet() {\n")

	if hasSpdlog {
		sb.WriteString(fmt.Sprintf(`    spdlog::info("Hello from %s!");
`, projectName))
	} else {
		sb.WriteString(fmt.Sprintf(`    std::cout << "Hello from %s!" << std::endl;
`, projectName))
	}

	sb.WriteString(`}

}  // namespace ` + projectName + "\n")

	return sb.String()
}

func generateTestCMake(projectName string, libraryIDs []string, testingFramework string) string {
	hasGtest := false
	hasCatch2 := false

	for _, libID := range libraryIDs {
		if libID == "googletest" {
			hasGtest = true
		}
		if libID == "catch2" {
			hasCatch2 = true
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`# Test configuration for %s

add_executable(%s_tests
    test_main.cpp
    ${CMAKE_CURRENT_SOURCE_DIR}/../src/%s.cpp
)

target_include_directories(%s_tests
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/../include
)

# Link libraries from dependencies.cmake (FORGE_LINK_LIBRARIES + FORGE_TEST_LINK_LIBRARIES)
target_link_libraries(%s_tests
    PRIVATE
        ${FORGE_LINK_LIBRARIES}
        ${FORGE_TEST_LINK_LIBRARIES}
)

`, projectName, projectName, projectName, projectName, projectName))

	if hasGtest {
		sb.WriteString(fmt.Sprintf(`include(GoogleTest)
gtest_discover_tests(%s_tests)
`, projectName))
	} else if hasCatch2 {
		sb.WriteString(fmt.Sprintf(`include(CTest)
include(Catch)
catch_discover_tests(%s_tests)
`, projectName))
	} else {
		sb.WriteString(fmt.Sprintf(`add_test(NAME %s_tests COMMAND %s_tests)
`, projectName, projectName))
	}

	return sb.String()
}

func generateTestMain(projectName string, libraryIDs []string, testingFramework string) string {
	hasGtest := false
	hasCatch2 := false
	hasDoctest := false

	for _, libID := range libraryIDs {
		switch libID {
		case "googletest":
			hasGtest = true
		case "catch2":
			hasCatch2 = true
		case "doctest":
			hasDoctest = true
		}
	}

	if hasGtest {
		capName := projectName
		if len(projectName) > 0 {
			capName = strings.ToUpper(projectName[:1]) + projectName[1:]
		}
		return fmt.Sprintf(`#include <gtest/gtest.h>
#include <%s/%s.hpp>

TEST(%sTest, VersionTest) {
    EXPECT_EQ(%s::version(), "1.0.0");
}

TEST(%sTest, GreetTest) {
    // Should not throw
    EXPECT_NO_THROW(%s::greet());
}
`, projectName, projectName, capName, projectName, capName, projectName)
	} else if hasCatch2 {
		return fmt.Sprintf(`#include <catch2/catch_test_macros.hpp>
#include <%s/%s.hpp>

TEST_CASE("%s::version returns correct version", "[version]") {
    REQUIRE(%s::version() == "1.0.0");
}

TEST_CASE("%s::greet does not throw", "[greet]") {
    REQUIRE_NOTHROW(%s::greet());
}
`, projectName, projectName, projectName, projectName, projectName, projectName)
	} else if hasDoctest {
		return fmt.Sprintf(`#define DOCTEST_CONFIG_IMPLEMENT_WITH_MAIN
#include <doctest/doctest.h>
#include <%s/%s.hpp>

TEST_CASE("testing version") {
    CHECK(%s::version() == "1.0.0");
}

TEST_CASE("testing greet") {
    CHECK_NOTHROW(%s::greet());
}
`, projectName, projectName, projectName, projectName)
	} else {
		return fmt.Sprintf(`// Basic test file - add a test framework for better testing support
#include <%s/%s.hpp>
#include <cassert>
#include <iostream>

int main() {
    assert(%s::version() == "1.0.0");
    %s::greet();
    std::cout << "All tests passed!" << std::endl;
    return 0;
}
`, projectName, projectName, projectName, projectName)
	}
}

func generateReadme(projectName string, libraryIDs []string, cppStandard int, projectType string) string {
	var libList strings.Builder
	if len(libraryIDs) > 0 {
		for _, libID := range libraryIDs {
			libList.WriteString(fmt.Sprintf("- %s\n", libID))
		}
	} else {
		libList.WriteString("No external dependencies.")
	}

	if projectType == "lib" {
		return fmt.Sprintf(`# %s

A C++ library using modern CMake and FetchContent for dependency management.

## Requirements

- CMake 3.20 or higher
- C++%d compatible compiler

## Dependencies

%s

## Building

`+"```bash\nmkdir build && cd build\ncmake ..\ncmake --build .\n```"+`

## Installation

`+"```bash\ncd build\ncmake --install . --prefix /usr/local\n```"+`

## Usage

`+"```cmake\nfind_package(%s REQUIRED)\ntarget_link_libraries(your_target PRIVATE %s)\n```"+`

## Testing

`+"```bash\ncd build\nctest --output-on-failure\n```"+`

## Project Structure

`+"```\n%s/\n├── .cmake/\n│   └── forge/\n│       └── dependencies.cmake  # Managed by Forge - regenerate to update\n├── CMakeLists.txt\n├── include/\n│   └── %s/\n│       └── %s.hpp\n├── src/\n│   └── %s.cpp\n├── tests/\n│   ├── CMakeLists.txt\n│   └── test_main.cpp\n└── README.md\n```"+`

## Updating Dependencies

To update dependencies, edit `+"`forge.yaml`"+` and run:
`+"```bash\nforge generate\n```"+`

This regenerates .cmake/forge/dependencies.cmake without modifying your CMakeLists.txt.

## License

MIT License
`, projectName, cppStandard, libList.String(), projectName, projectName, projectName, projectName, projectName, projectName)
	} else {
		return fmt.Sprintf(`# %s

A C++ project using modern CMake and FetchContent for dependency management.

## Requirements

- CMake 3.20 or higher
- C++%d compatible compiler

## Dependencies

%s

## Building

`+"```bash\nmkdir build && cd build\ncmake ..\ncmake --build .\n```"+`

## Running

`+"```bash\n./build/%s\n```"+`

## Testing

`+"```bash\ncd build\nctest --output-on-failure\n```"+`

## Project Structure

`+"```\n%s/\n├── .cmake/\n│   └── forge/\n│       └── dependencies.cmake  # Managed by Forge - regenerate to update\n├── CMakeLists.txt\n├── include/\n│   └── %s/\n│       └── %s.hpp\n├── src/\n│   ├── main.cpp\n│   └── %s.cpp\n├── tests/\n│   ├── CMakeLists.txt\n│   └── test_main.cpp\n└── README.md\n```"+`

## Updating Dependencies

To update dependencies, edit `+"`forge.yaml`"+` and run:
`+"```bash\nforge generate\n```"+`

This regenerates .cmake/forge/dependencies.cmake without modifying your CMakeLists.txt.

## License

MIT License
`, projectName, cppStandard, libList.String(), projectName, projectName, projectName, projectName, projectName)
	}
}

func generateGitignore() string {
	return `# Build directories
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
`
}

