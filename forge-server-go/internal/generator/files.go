package generator

import (
	"fmt"
	"strings"

	"github.com/ozacod/forge/forge-server-go/internal/recipe"
)

func GenerateTestCMake(
	projectName string,
	testLibraries []LibraryWithOptions,
	mainLibraries []LibraryWithOptions,
	projectType string,
) string {
	hasGtest := false
	hasCatch2 := false

	for _, lwo := range testLibraries {
		if lwo.Lib.ID == "googletest" {
			hasGtest = true
		}
		if lwo.Lib.ID == "catch2" {
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

func GenerateMainCpp(projectName string, libraries []*recipe.Library) string {
	var includes []string

	// Add relevant includes based on selected libraries
	for _, lib := range libraries {
		switch lib.ID {
		case "nlohmann_json":
			includes = append(includes, "#include <nlohmann/json.hpp>")
		case "spdlog":
			includes = append(includes, "#include <spdlog/spdlog.h>")
		case "fmt":
			includes = append(includes, "#include <fmt/format.h>")
		case "cli11":
			includes = append(includes, "#include <CLI/CLI.hpp>")
		case "argparse":
			includes = append(includes, "#include <argparse/argparse.hpp>")
		}
	}

	includesStr := strings.Join(includes, "\n")
	if includesStr != "" {
		includesStr = "\n" + includesStr
	}

	hasSpdlog := false
	hasCLI11 := false
	hasArgparse := false

	for _, lib := range libraries {
		switch lib.ID {
		case "spdlog":
			hasSpdlog = true
		case "cli11":
			hasCLI11 = true
		case "argparse":
			hasArgparse = true
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`#include <%s/%s.hpp>
#include <iostream>%s

int main(int argc, char* argv[]) {
`, projectName, projectName, includesStr))

	if hasSpdlog {
		sb.WriteString(fmt.Sprintf(`    spdlog::info("Starting %s v1.0.0");
`, projectName))
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

func GenerateLibHeader(projectName string) string {
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

func GenerateLibSource(projectName string, libraries []*recipe.Library) string {
	hasSpdlog := false
	hasFmt := false

	for _, lib := range libraries {
		switch lib.ID {
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

func GenerateTestMain(projectName string, testLibraries []*recipe.Library) string {
	hasGtest := false
	hasCatch2 := false
	hasDoctest := false

	for _, lib := range testLibraries {
		switch lib.ID {
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

func GenerateReadme(projectName string, libraries []*recipe.Library, cppStandard int, projectType string) string {
	var libList strings.Builder
	if len(libraries) > 0 {
		for _, lib := range libraries {
			libList.WriteString(fmt.Sprintf("- [%s](%s) - %s\n", lib.Name, lib.GitHubURL, lib.Description))
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
`, projectName, cppStandard, libList.String(), projectName, projectName, projectName, projectName, projectName, projectName)
	}
}

func GenerateGitignore() string {
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

var clangFormatStyles = map[string]string{
	"Google": `BasedOnStyle: Google
IndentWidth: 4
ColumnLimit: 100
AllowShortFunctionsOnASingleLine: Empty
AllowShortIfStatementsOnASingleLine: Never
AllowShortLoopsOnASingleLine: false
BreakBeforeBraces: Attach
PointerAlignment: Left
SpaceAfterCStyleCast: false
SpaceBeforeParens: ControlStatements
`,
	"LLVM": `BasedOnStyle: LLVM
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: All
AllowShortIfStatementsOnASingleLine: Never
BreakBeforeBraces: Attach
PointerAlignment: Right
SpaceBeforeParens: ControlStatements
`,
	"Chromium": `BasedOnStyle: Chromium
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: Inline
AllowShortIfStatementsOnASingleLine: Never
BreakBeforeBraces: Attach
PointerAlignment: Left
DerivePointerAlignment: false
`,
	"Mozilla": `BasedOnStyle: Mozilla
IndentWidth: 2
ColumnLimit: 80
AllowShortFunctionsOnASingleLine: Inline
BreakBeforeBraces: Mozilla
PointerAlignment: Left
AlwaysBreakAfterDefinitionReturnType: TopLevel
`,
	"WebKit": `BasedOnStyle: WebKit
IndentWidth: 4
ColumnLimit: 0
AllowShortFunctionsOnASingleLine: All
BreakBeforeBraces: WebKit
PointerAlignment: Left
NamespaceIndentation: Inner
`,
	"Microsoft": `BasedOnStyle: Microsoft
IndentWidth: 4
ColumnLimit: 120
AllowShortFunctionsOnASingleLine: None
BreakBeforeBraces: Allman
PointerAlignment: Left
AccessModifierOffset: -4
AlignAfterOpenBracket: Align
`,
	"GNU": `BasedOnStyle: GNU
IndentWidth: 2
ColumnLimit: 79
AllowShortFunctionsOnASingleLine: None
BreakBeforeBraces: GNU
PointerAlignment: Right
SpaceBeforeParens: Always
`,
}

func GenerateClangFormat(style string) string {
	if s, ok := clangFormatStyles[style]; ok {
		return s
	}
	return clangFormatStyles["Google"]
}

