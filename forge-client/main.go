package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Version        = "1.0.7"
	DefaultServer  = "http://localhost:8000"
	DefaultCfgFile = "forge.yaml"
	LockFile       = "forge.lock"
)

// Colors for terminal output
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

// ForgeConfig represents the forge.yaml structure
type ForgeConfig struct {
	Package struct {
		Name        string   `yaml:"name"`
		Version     string   `yaml:"version"`
		CppStandard int      `yaml:"cpp_standard"`
		Authors     []string `yaml:"authors,omitempty"`
		Description string   `yaml:"description,omitempty"`
	} `yaml:"package"`
	Build struct {
		SharedLibs  bool   `yaml:"shared_libs"`
		ClangFormat string `yaml:"clang_format"`
		BuildType   string `yaml:"build_type,omitempty"`
		CxxFlags    string `yaml:"cxx_flags,omitempty"`
	} `yaml:"build"`
	Testing struct {
		Framework string `yaml:"framework"`
	} `yaml:"testing"`
	Features        map[string]FeatureConfig          `yaml:"features,omitempty"`
	Dependencies    map[string]map[string]interface{} `yaml:"dependencies"`
	DevDependencies map[string]map[string]interface{} `yaml:"dev-dependencies,omitempty"`
}

type FeatureConfig struct {
	Dependencies map[string]map[string]interface{} `yaml:"dependencies,omitempty"`
}

// LockConfig represents the forge.lock structure
type LockConfig struct {
	Version      int                  `yaml:"version"`
	Dependencies map[string]LockEntry `yaml:"dependencies"`
}

type LockEntry struct {
	Git    string `yaml:"git"`
	Tag    string `yaml:"tag"`
	Commit string `yaml:"commit,omitempty"`
}

// Library represents a library from the server
type Library struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Category     string            `json:"category"`
	HeaderOnly   bool              `json:"header_only"`
	CppStandard  int               `json:"cpp_standard"`
	GithubURL    string            `json:"github_url"`
	Tags         []string          `json:"tags"`
	Options      []LibraryOption   `json:"options"`
	FetchContent map[string]string `json:"fetch_content"`
}

type LibraryOption struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	CMakeVar    string      `json:"cmake_var"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	// Handle global flags
	if command == "-v" || command == "--version" || command == "version" {
		fmt.Printf("%sforge%s version %s%s%s\n", Bold, Reset, Cyan, Version, Reset)
		return
	}

	if command == "-h" || command == "--help" || command == "help" {
		printUsage()
		return
	}

	// Parse command-specific flags
	switch command {
	case "generate", "gen":
		cmdGenerate(os.Args[2:])
	case "build":
		cmdBuild(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	case "test":
		cmdTest(os.Args[2:])
	case "clean":
		cmdClean(os.Args[2:])
	case "new", "init":
		cmdNew(os.Args[2:])
	case "add":
		cmdAdd(os.Args[2:])
	case "remove", "rm":
		cmdRemove(os.Args[2:])
	case "update":
		cmdUpdate(os.Args[2:])
	case "list":
		cmdList(os.Args[2:])
	case "search":
		cmdSearch(os.Args[2:])
	case "info":
		cmdInfo(os.Args[2:])
	case "fmt", "format":
		cmdFmt(os.Args[2:])
	case "lint":
		cmdLint(os.Args[2:])
	case "check":
		cmdCheck(os.Args[2:])
	case "doc":
		cmdDoc(os.Args[2:])
	case "release":
		cmdRelease(os.Args[2:])
	case "upgrade":
		cmdUpgrade(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "%sError:%s Unknown command: %s\n", Red, Reset, command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`%s%sforge%s - C++ Project Generator (like Cargo for Rust)

%sUSAGE:%s
    forge <COMMAND> [OPTIONS]

%sCOMMANDS:%s
    %sgenerate%s    Generate CMake project from forge.yaml (alias: gen)
    %sbuild%s       Compile the project with CMake (-O0/1/2/3/s/fast, --clean)
    %srun%s         Build and run the project
    %stest%s        Build and run tests
    %sclean%s       Remove build artifacts
    %snew%s         Create a new project (in current or new directory)
    %sadd%s         Add a dependency
    %sremove%s      Remove a dependency
    %supdate%s      Update dependencies to latest versions
    %slist%s        List available libraries
    %ssearch%s      Search for libraries
    %sinfo%s        Show detailed library information
    %sfmt%s         Format code with clang-format
    %slint%s        Run clang-tidy static analysis
    %scheck%s       Check code compiles without building
    %sdoc%s         Generate documentation
    %srelease%s     Bump version number
    %supgrade%s     Upgrade forge to the latest version
    %sversion%s     Show version
    %shelp%s        Show this help

EXAMPLES:
    forge new my_project          Create new project in 'my_project/' directory
    forge new my_lib --lib        Create library project
    forge new                     Create project in current directory
    forge new . --lib             Create library in current directory
    forge new -t web-server       Create with template
    forge add spdlog              Add dependency
    forge add --dev catch2        Add dev dependency
    forge generate                Generate CMake project from yaml
    forge build                   Compile with CMake
    forge run                     Build and run
    forge test                    Run tests
    forge fmt                     Format all code
    forge search json             Search for libraries

Run 'forge <COMMAND> --help' for more information on a command.
`, Bold, Cyan, Reset,
		Yellow, Reset,
		Yellow, Reset,
		Green, Reset, // generate
		Green, Reset, // build
		Green, Reset, // run
		Green, Reset, // test
		Green, Reset, // clean
		Green, Reset, // init
		Green, Reset, // new
		Green, Reset, // add
		Green, Reset, // remove
		Green, Reset, // update
		Green, Reset, // list
		Green, Reset, // search
		Green, Reset, // info
		Green, Reset, // fmt
		Green, Reset, // lint
		Green, Reset, // check
		Green, Reset, // doc
		Green, Reset, // release
		Green, Reset, // upgrade
		Green, Reset, // version
		Green, Reset) // help
}

// ============================================================================
// GENERATE COMMAND - Generate CMake project from forge.yaml
// ============================================================================

func cmdGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	configFile := fs.String("config", DefaultCfgFile, "Config file")
	outputDir := fs.String("output", ".", "Output directory")
	features := fs.String("features", "", "Comma-separated features to enable")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.StringVar(configFile, "c", DefaultCfgFile, "Config file (shorthand)")
	fs.StringVar(outputDir, "o", ".", "Output directory (shorthand)")
	fs.Parse(args)

	if err := generateProject(*serverURL, *configFile, *outputDir, *features); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func generateProject(serverURL, configFile, outputDir string, features string) error {
	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file '%s': %w", configFile, err)
	}

	// Parse YAML to get project name
	var config ForgeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	projectName := config.Package.Name
	if projectName == "" {
		projectName = "my_project"
	}

	fmt.Printf("%s📦 Generating project '%s' from %s...%s\n", Cyan, projectName, configFile, Reset)
	fmt.Printf("   Server: %s\n", serverURL)

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(configFile))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		return fmt.Errorf("failed to write form data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Make request to server
	url := fmt.Sprintf("%s/api/forge", serverURL)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w\n\nMake sure the server is running:\n  cd forge-server && uvicorn main:app --port 8000", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	// Read ZIP content
	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Extract ZIP to output directory
	fmt.Printf("%s📦 Extracting to %s...%s\n", Cyan, outputDir, Reset)

	if err := extractZip(zipData, outputDir); err != nil {
		return fmt.Errorf("failed to extract project: %w", err)
	}

	// Generate lock file
	if err := generateLockFile(config, outputDir); err != nil {
		fmt.Printf("%s⚠️  Warning: Could not generate lock file: %v%s\n", Yellow, err, Reset)
	}

	fmt.Printf("%s✅ Project '%s' generated successfully!%s\n\n", Green, projectName, Reset)
	fmt.Printf("Next steps:\n")
	if outputDir != "." {
		fmt.Printf("  cd %s\n", outputDir)
	}
	fmt.Printf("  %sforge build%s      # Compile the project\n", Cyan, Reset)
	fmt.Printf("  %sforge run%s        # Build and run\n", Cyan, Reset)

	return nil
}

// ============================================================================
// BUILD COMMAND - Compile the project with CMake
// ============================================================================

func cmdBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	release := fs.Bool("release", false, "Build in release mode (O2)")
	debug := fs.Bool("debug", false, "Build in debug mode (O0, default)")
	jobs := fs.Int("jobs", 0, "Number of parallel jobs (0 = auto)")
	target := fs.String("target", "", "Specific target to build")
	clean := fs.Bool("clean", false, "Clean build directory before building")
	optLevel := fs.String("opt", "", "Optimization level: 0, 1, 2, 3, s, fast")
	fs.BoolVar(release, "r", false, "Build in release mode (shorthand)")
	fs.IntVar(jobs, "j", 0, "Number of parallel jobs (shorthand)")
	fs.BoolVar(clean, "c", false, "Clean before building (shorthand)")
	fs.StringVar(optLevel, "O", "", "Optimization level (shorthand)")
	fs.Parse(args)

	if err := buildProject(*release, *debug, *jobs, *target, *clean, *optLevel); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func buildProject(release, debug bool, jobs int, target string, clean bool, optLevel string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	projectName := config.Package.Name
	if projectName == "" {
		projectName = "my_project"
	}

	buildDir := "build"

	// Clean if requested
	if clean {
		fmt.Printf("%s🧹 Cleaning build directory...%s\n", Cyan, Reset)
		os.RemoveAll(buildDir)
	}

	// Determine build type and optimization
	buildType := "Debug"
	cxxFlags := ""

	if release {
		buildType = "Release"
	}

	// Handle optimization level
	switch optLevel {
	case "0":
		cxxFlags = "-O0"
		buildType = "Debug"
	case "1":
		cxxFlags = "-O1"
		buildType = "RelWithDebInfo"
	case "2":
		cxxFlags = "-O2"
		buildType = "Release"
	case "3":
		cxxFlags = "-O3"
		buildType = "Release"
	case "s":
		cxxFlags = "-Os"
		buildType = "MinSizeRel"
	case "fast":
		cxxFlags = "-Ofast"
		buildType = "Release"
	}

	optInfo := ""
	if cxxFlags != "" {
		optInfo = fmt.Sprintf(" [%s]", cxxFlags)
	}

	fmt.Printf("%s🔨 Building '%s' (%s%s)...%s\n", Cyan, projectName, buildType, optInfo, Reset)

	// Configure CMake if needed or if clean was done
	needsConfigure := clean
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		needsConfigure = true
	}

	if needsConfigure {
		fmt.Printf("%s⚙️  Configuring CMake...%s\n", Cyan, Reset)
		cmakeArgs := []string{"-B", buildDir, "-DCMAKE_BUILD_TYPE=" + buildType}

		if cxxFlags != "" {
			cmakeArgs = append(cmakeArgs, "-DCMAKE_CXX_FLAGS="+cxxFlags)
		}

		cmd := exec.Command("cmake", cmakeArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	// Build
	fmt.Printf("%s🔧 Compiling...%s\n", Cyan, Reset)
	buildArgs := []string{"--build", buildDir, "--config", buildType}

	if jobs > 0 {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", jobs))
	} else {
		buildArgs = append(buildArgs, "--parallel", fmt.Sprintf("%d", runtime.NumCPU()))
	}

	if target != "" {
		buildArgs = append(buildArgs, "--target", target)
	}

	buildCmd := exec.Command("cmake", buildArgs...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("%s✅ Build complete!%s\n", Green, Reset)
	return nil
}

// ============================================================================
// RUN COMMAND
// ============================================================================

func cmdRun(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	release := fs.Bool("release", false, "Build in release mode")
	target := fs.String("target", "", "Specific target to run")
	fs.Parse(args)

	// Get remaining args to pass to the executable
	execArgs := fs.Args()

	if err := runProject(*release, *target, execArgs); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func runProject(release bool, target string, execArgs []string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	projectName := config.Package.Name
	if projectName == "" {
		projectName = "my_project"
	}

	buildType := "Debug"
	if release {
		buildType = "Release"
	}

	fmt.Printf("%s🔨 Building '%s' (%s)...%s\n", Cyan, projectName, buildType, Reset)

	// Configure CMake if needed
	buildDir := "build"
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		fmt.Printf("%s⚙️  Configuring CMake...%s\n", Cyan, Reset)
		cmd := exec.Command("cmake", "-B", buildDir, "-DCMAKE_BUILD_TYPE="+buildType)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	// Build
	fmt.Printf("%s🔧 Compiling...%s\n", Cyan, Reset)
	buildCmd := exec.Command("cmake", "--build", buildDir, "--config", buildType)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Find and run executable
	execName := projectName
	if runtime.GOOS == "windows" {
		execName += ".exe"
	}

	execPath := filepath.Join(buildDir, execName)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		// Try in build type subdirectory (MSVC)
		execPath = filepath.Join(buildDir, buildType, execName)
	}

	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("executable not found: tried %s", execPath)
	}

	fmt.Printf("\n%s🚀 Running '%s'...%s\n", Green, projectName, Reset)
	fmt.Println(strings.Repeat("─", 50))

	runCmd := exec.Command(execPath, execArgs...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin
	return runCmd.Run()
}

// ============================================================================
// TEST COMMAND
// ============================================================================

func cmdTest(args []string) {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	verbose := fs.Bool("verbose", false, "Show verbose output")
	filter := fs.String("filter", "", "Filter tests by name")
	fs.BoolVar(verbose, "v", false, "Show verbose output (shorthand)")
	fs.Parse(args)

	if err := runTests(*verbose, *filter); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func runTests(verbose bool, filter string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	projectName := config.Package.Name
	fmt.Printf("%s🧪 Running tests for '%s'...%s\n", Cyan, projectName, Reset)

	buildDir := "build"

	// Configure CMake if needed
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		fmt.Printf("%s⚙️  Configuring CMake...%s\n", Cyan, Reset)
		cmd := exec.Command("cmake", "-B", buildDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	// Build tests
	fmt.Printf("%s🔧 Building tests...%s\n", Cyan, Reset)
	buildCmd := exec.Command("cmake", "--build", buildDir)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Run tests with ctest
	fmt.Printf("\n%s🧪 Running tests...%s\n", Green, Reset)
	fmt.Println(strings.Repeat("─", 50))

	ctestArgs := []string{"--test-dir", buildDir, "--output-on-failure"}
	if verbose {
		ctestArgs = append(ctestArgs, "-V")
	}
	if filter != "" {
		ctestArgs = append(ctestArgs, "-R", filter)
	}

	testCmd := exec.Command("ctest", ctestArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	return testCmd.Run()
}

// ============================================================================
// CLEAN COMMAND
// ============================================================================

func cmdClean(args []string) {
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	all := fs.Bool("all", false, "Also remove generated files")
	fs.Parse(args)

	if err := cleanProject(*all); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func cleanProject(all bool) error {
	fmt.Printf("%s🧹 Cleaning build artifacts...%s\n", Cyan, Reset)

	// Remove build directory
	if err := os.RemoveAll("build"); err != nil {
		return fmt.Errorf("failed to remove build directory: %w", err)
	}
	fmt.Println("   ✓ Removed build/")

	// Remove CMake cache
	cacheFiles := []string{
		"CMakeCache.txt",
		"CMakeFiles",
		"cmake_install.cmake",
		"Makefile",
		"compile_commands.json",
	}

	for _, f := range cacheFiles {
		if _, err := os.Stat(f); err == nil {
			os.RemoveAll(f)
			fmt.Printf("   ✓ Removed %s\n", f)
		}
	}

	if all {
		// Remove generated files
		genFiles := []string{LockFile}
		for _, f := range genFiles {
			if _, err := os.Stat(f); err == nil {
				os.Remove(f)
				fmt.Printf("   ✓ Removed %s\n", f)
			}
		}
	}

	fmt.Printf("%s✅ Clean complete!%s\n", Green, Reset)
	return nil
}

// ============================================================================
// NEW COMMAND
// ============================================================================

func cmdNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	templateName := fs.String("template", "", "Use a template")
	isLib := fs.Bool("lib", false, "Create a library project")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.StringVar(templateName, "t", "", "Use a template (shorthand)")
	fs.Parse(args)

	remaining := fs.Args()

	// Default to current directory if no name given
	projectName := "."
	for _, arg := range remaining {
		switch arg {
		case "lib", "library":
			*isLib = true
		case "exe", "bin":
			*isLib = false
		default:
			projectName = arg
		}
	}

	if err := newProject(*serverURL, projectName, *templateName, *isLib); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func newProject(serverURL, projectName, templateName string, isLib bool) error {
	inCurrentDir := projectName == "."

	// If creating in current directory, use folder name as project name
	if inCurrentDir {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectName = filepath.Base(cwd)
	}

	// Validate project name
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`).MatchString(projectName) {
		return fmt.Errorf("invalid project name '%s': must start with letter and contain only letters, numbers, underscores, or hyphens", projectName)
	}

	if inCurrentDir {
		// Check if forge.yaml already exists
		if _, err := os.Stat(DefaultCfgFile); err == nil {
			return fmt.Errorf("forge.yaml already exists in current directory")
		}
		fmt.Printf("%s📁 Initializing project '%s' in current directory...%s\n", Cyan, projectName, Reset)
	} else {
		// Check if directory already exists
		if _, err := os.Stat(projectName); err == nil {
			return fmt.Errorf("directory '%s' already exists", projectName)
		}

		fmt.Printf("%s📁 Creating project '%s'...%s\n", Cyan, projectName, Reset)

		// Create directory
		if err := os.Mkdir(projectName, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Change to the new directory
		if err := os.Chdir(projectName); err != nil {
			return fmt.Errorf("failed to enter directory: %w", err)
		}
	}

	// Create forge.yaml
	var configContent string
	if isLib {
		configContent = fmt.Sprintf(`# forge.yaml - C++ Library Project
package:
  name: %s
  version: "0.1.0"
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

dependencies:
  fmt: {}
`, projectName)
	} else if templateName != "" {
		// Fetch template from server
		url := fmt.Sprintf("%s/api/forge/example/%s", serverURL, templateName)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch template: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("template '%s' not found", templateName)
		}

		data, _ := io.ReadAll(resp.Body)
		// Replace project name in template
		configContent = strings.ReplaceAll(string(data), "my_project", projectName)
		configContent = strings.ReplaceAll(configContent, "hello_world", projectName)
	} else {
		configContent = fmt.Sprintf(`# forge.yaml - C++ Project Dependencies
package:
  name: %s
  version: "0.1.0"
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

dependencies:
  spdlog:
    spdlog_header_only: true
  fmt: {}
`, projectName)
	}

	if err := os.WriteFile(DefaultCfgFile, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("%s✅ Created project '%s'%s\n\n", Green, projectName, Reset)
	fmt.Printf("Next steps:\n")
	if !inCurrentDir {
		fmt.Printf("  cd %s\n", projectName)
	}
	fmt.Printf("  %sforge generate%s   # Generate project files\n", Cyan, Reset)
	fmt.Printf("  %sforge build%s      # Compile the project\n", Cyan, Reset)
	fmt.Printf("  %sforge run%s        # Build and run\n", Cyan, Reset)

	return nil
}

// ============================================================================
// ADD COMMAND
// ============================================================================

func cmdAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	dev := fs.Bool("dev", false, "Add as dev dependency")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) < 1 {
		fmt.Fprintf(os.Stderr, "%sError:%s Library name required\n", Red, Reset)
		fmt.Fprintf(os.Stderr, "Usage: forge add <library> [--dev]\n")
		os.Exit(1)
	}

	libName := remaining[0]
	if err := addDependency(*serverURL, libName, *dev); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func addDependency(serverURL, libName string, dev bool) error {
	// Verify library exists
	lib, err := getLibraryInfo(serverURL, libName)
	if err != nil {
		return fmt.Errorf("library '%s' not found: %w", libName, err)
	}

	// Load current config
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	// Check if already added
	if config.Dependencies == nil {
		config.Dependencies = make(map[string]map[string]interface{})
	}
	if config.DevDependencies == nil {
		config.DevDependencies = make(map[string]map[string]interface{})
	}

	targetDeps := config.Dependencies
	depType := "dependency"
	if dev {
		targetDeps = config.DevDependencies
		depType = "dev-dependency"
	}

	if _, exists := targetDeps[libName]; exists {
		return fmt.Errorf("'%s' is already a %s", libName, depType)
	}

	// Add the dependency
	targetDeps[libName] = make(map[string]interface{})

	fmt.Printf("%s📦 Adding '%s' to %s...%s\n", Cyan, lib.Name, depType, Reset)

	// Save config
	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Printf("%s✅ Added %s (%s)%s\n", Green, lib.Name, lib.Description, Reset)
	fmt.Printf("\nRun %sforge build%s to update your project\n", Cyan, Reset)

	return nil
}

// ============================================================================
// REMOVE COMMAND
// ============================================================================

func cmdRemove(args []string) {
	fs := flag.NewFlagSet("remove", flag.ExitOnError)
	fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) < 1 {
		fmt.Fprintf(os.Stderr, "%sError:%s Library name required\n", Red, Reset)
		fmt.Fprintf(os.Stderr, "Usage: forge remove <library>\n")
		os.Exit(1)
	}

	libName := remaining[0]
	if err := removeDependency(libName); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func removeDependency(libName string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	found := false
	if _, exists := config.Dependencies[libName]; exists {
		delete(config.Dependencies, libName)
		found = true
	}
	if _, exists := config.DevDependencies[libName]; exists {
		delete(config.DevDependencies, libName)
		found = true
	}

	if !found {
		return fmt.Errorf("'%s' is not a dependency", libName)
	}

	fmt.Printf("%s🗑️  Removing '%s'...%s\n", Cyan, libName, Reset)

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Printf("%s✅ Removed %s%s\n", Green, libName, Reset)
	fmt.Printf("\nRun %sforge build%s to update your project\n", Cyan, Reset)

	return nil
}

// ============================================================================
// UPDATE COMMAND
// ============================================================================

func cmdUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.Parse(args)

	remaining := fs.Args()
	var libName string
	if len(remaining) > 0 {
		libName = remaining[0]
	}

	if err := updateDependencies(*serverURL, libName); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func updateDependencies(serverURL, specificLib string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	fmt.Printf("%s🔄 Checking for updates...%s\n", Cyan, Reset)

	// Get all libraries info
	libs, err := getAllLibraries(serverURL)
	if err != nil {
		return err
	}

	libMap := make(map[string]Library)
	for _, lib := range libs {
		libMap[lib.ID] = lib
	}

	updated := 0
	for libName := range config.Dependencies {
		if specificLib != "" && libName != specificLib {
			continue
		}

		if lib, ok := libMap[libName]; ok {
			fmt.Printf("   ✓ %s (up to date)\n", lib.Name)
			updated++
		}
	}

	if updated == 0 {
		fmt.Printf("%s✅ All dependencies are up to date%s\n", Green, Reset)
	} else {
		fmt.Printf("%s✅ Checked %d dependencies%s\n", Green, updated, Reset)
	}

	return nil
}

// ============================================================================
// LIST COMMAND
// ============================================================================

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	category := fs.String("category", "", "Filter by category")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.Parse(args)

	if err := listLibraries(*serverURL, *category); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func listLibraries(serverURL, category string) error {
	libs, err := getAllLibraries(serverURL)
	if err != nil {
		return err
	}

	// Group by category
	categories := make(map[string][]Library)
	for _, lib := range libs {
		if category != "" && lib.Category != category {
			continue
		}
		categories[lib.Category] = append(categories[lib.Category], lib)
	}

	fmt.Printf("%s📚 Available Libraries (%d total)%s\n\n", Bold, len(libs), Reset)

	// Print by category
	categoryOrder := []string{
		"serialization", "logging", "testing", "networking", "cli",
		"configuration", "gui", "formatting", "concurrency", "utility",
		"database", "compression", "math", "cryptography",
	}

	for _, cat := range categoryOrder {
		catLibs, ok := categories[cat]
		if !ok || len(catLibs) == 0 {
			continue
		}

		fmt.Printf("  %s%s:%s\n", Yellow, strings.Title(cat), Reset)
		for _, lib := range catLibs {
			headerOnly := ""
			if lib.HeaderOnly {
				headerOnly = fmt.Sprintf(" %s[header-only]%s", Cyan, Reset)
			}
			fmt.Printf("    • %-20s C++%d%s\n", lib.ID, lib.CppStandard, headerOnly)
		}
		fmt.Println()
	}

	return nil
}

// ============================================================================
// SEARCH COMMAND
// ============================================================================

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) < 1 {
		fmt.Fprintf(os.Stderr, "%sError:%s Search query required\n", Red, Reset)
		fmt.Fprintf(os.Stderr, "Usage: forge search <query>\n")
		os.Exit(1)
	}

	query := strings.Join(remaining, " ")
	if err := searchLibraries(*serverURL, query); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func searchLibraries(serverURL, query string) error {
	libs, err := getAllLibraries(serverURL)
	if err != nil {
		return err
	}

	query = strings.ToLower(query)
	var results []Library

	for _, lib := range libs {
		// Search in id, name, description, tags
		if strings.Contains(strings.ToLower(lib.ID), query) ||
			strings.Contains(strings.ToLower(lib.Name), query) ||
			strings.Contains(strings.ToLower(lib.Description), query) {
			results = append(results, lib)
			continue
		}
		for _, tag := range lib.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, lib)
				break
			}
		}
	}

	if len(results) == 0 {
		fmt.Printf("%s🔍 No libraries found matching '%s'%s\n", Yellow, query, Reset)
		return nil
	}

	fmt.Printf("%s🔍 Found %d libraries matching '%s':%s\n\n", Green, len(results), query, Reset)

	for _, lib := range results {
		fmt.Printf("  %s%s%s (%s)\n", Bold, lib.Name, Reset, lib.ID)
		fmt.Printf("    %s\n", lib.Description)
		if len(lib.Tags) > 0 {
			fmt.Printf("    Tags: %s%s%s\n", Cyan, strings.Join(lib.Tags, ", "), Reset)
		}
		fmt.Println()
	}

	return nil
}

// ============================================================================
// INFO COMMAND
// ============================================================================

func cmdInfo(args []string) {
	fs := flag.NewFlagSet("info", flag.ExitOnError)
	serverURL := fs.String("server", DefaultServer, "Server URL")
	fs.StringVar(serverURL, "s", DefaultServer, "Server URL (shorthand)")
	fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) < 1 {
		fmt.Fprintf(os.Stderr, "%sError:%s Library name required\n", Red, Reset)
		fmt.Fprintf(os.Stderr, "Usage: forge info <library>\n")
		os.Exit(1)
	}

	libName := remaining[0]
	if err := showLibraryInfo(*serverURL, libName); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func showLibraryInfo(serverURL, libName string) error {
	lib, err := getLibraryInfo(serverURL, libName)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s%s%s\n", Bold, lib.Name, Reset)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("ID:          %s\n", lib.ID)
	fmt.Printf("Description: %s\n", lib.Description)
	fmt.Printf("Category:    %s\n", lib.Category)
	fmt.Printf("C++ Standard: C++%d\n", lib.CppStandard)
	fmt.Printf("Header Only: %v\n", lib.HeaderOnly)
	if lib.GithubURL != "" {
		fmt.Printf("GitHub:      %s%s%s\n", Cyan, lib.GithubURL, Reset)
	}
	if len(lib.Tags) > 0 {
		fmt.Printf("Tags:        %s\n", strings.Join(lib.Tags, ", "))
	}

	if len(lib.Options) > 0 {
		fmt.Printf("\n%sOptions:%s\n", Yellow, Reset)
		for _, opt := range lib.Options {
			fmt.Printf("  %s%s%s: %s (default: %v)\n", Cyan, opt.ID, Reset, opt.Description, opt.Default)
		}
	}

	fmt.Printf("\n%sUsage in forge.yaml:%s\n", Yellow, Reset)
	fmt.Printf("  dependencies:\n")
	fmt.Printf("    %s: {}\n", lib.ID)

	return nil
}

// ============================================================================
// FMT COMMAND
// ============================================================================

func cmdFmt(args []string) {
	fs := flag.NewFlagSet("fmt", flag.ExitOnError)
	check := fs.Bool("check", false, "Check formatting without modifying files")
	fs.Parse(args)

	if err := formatCode(*check); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func formatCode(checkOnly bool) error {
	// Check if clang-format is available
	if _, err := exec.LookPath("clang-format"); err != nil {
		return fmt.Errorf("clang-format not found. Please install it first")
	}

	fmt.Printf("%s🎨 Formatting code...%s\n", Cyan, Reset)

	// Find all source files
	var files []string
	extensions := []string{".cpp", ".hpp", ".c", ".h", ".cc", ".cxx", ".hxx"}

	for _, dir := range []string{"src", "include", "tests"} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			for _, ext := range extensions {
				if strings.HasSuffix(path, ext) {
					files = append(files, path)
					break
				}
			}
			return nil
		})
	}

	if len(files) == 0 {
		fmt.Printf("%s✅ No source files found%s\n", Green, Reset)
		return nil
	}

	// Format each file
	formatArgs := []string{"-style=file"}
	if !checkOnly {
		formatArgs = append(formatArgs, "-i")
	} else {
		formatArgs = append(formatArgs, "--dry-run", "--Werror")
	}

	needsFormat := false
	for _, file := range files {
		args := append(formatArgs, file)
		cmd := exec.Command("clang-format", args...)
		output, err := cmd.CombinedOutput()

		if checkOnly && err != nil {
			needsFormat = true
			fmt.Printf("   %s✗ %s needs formatting%s\n", Yellow, file, Reset)
		} else if !checkOnly {
			fmt.Printf("   ✓ %s\n", file)
		}

		if len(output) > 0 && checkOnly {
			fmt.Print(string(output))
		}
	}

	if checkOnly && needsFormat {
		return fmt.Errorf("some files need formatting. Run 'forge fmt' to fix")
	}

	fmt.Printf("%s✅ Formatted %d files%s\n", Green, len(files), Reset)
	return nil
}

// ============================================================================
// LINT COMMAND
// ============================================================================

func cmdLint(args []string) {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	fix := fs.Bool("fix", false, "Automatically fix issues")
	fs.Parse(args)

	if err := lintCode(*fix); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func lintCode(fix bool) error {
	// Check if clang-tidy is available
	if _, err := exec.LookPath("clang-tidy"); err != nil {
		return fmt.Errorf("clang-tidy not found. Please install it first")
	}

	fmt.Printf("%s🔍 Running static analysis...%s\n", Cyan, Reset)

	// Check for compile_commands.json
	compileDb := "build/compile_commands.json"
	if _, err := os.Stat(compileDb); os.IsNotExist(err) {
		fmt.Printf("%s⚙️  Generating compile_commands.json...%s\n", Cyan, Reset)
		cmd := exec.Command("cmake", "-B", "build", "-DCMAKE_EXPORT_COMPILE_COMMANDS=ON")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate compile_commands.json: %w", err)
		}
	}

	// Find source files
	var files []string
	for _, dir := range []string{"src"} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".cpp") || strings.HasSuffix(path, ".cc") {
				files = append(files, path)
			}
			return nil
		})
	}

	if len(files) == 0 {
		fmt.Printf("%s✅ No source files found%s\n", Green, Reset)
		return nil
	}

	// Run clang-tidy
	tidyArgs := []string{"-p", "build"}
	if fix {
		tidyArgs = append(tidyArgs, "-fix")
	}
	tidyArgs = append(tidyArgs, files...)

	cmd := exec.Command("clang-tidy", tidyArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// clang-tidy returns non-zero on warnings
		fmt.Printf("%s⚠️  Analysis complete with warnings%s\n", Yellow, Reset)
		return nil
	}

	fmt.Printf("%s✅ No issues found!%s\n", Green, Reset)
	return nil
}

// ============================================================================
// CHECK COMMAND
// ============================================================================

func cmdCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	fs.Parse(args)

	if err := checkCode(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func checkCode() error {
	fmt.Printf("%s🔎 Checking code...%s\n", Cyan, Reset)

	buildDir := "build"

	// Configure CMake
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); os.IsNotExist(err) {
		fmt.Printf("%s⚙️  Configuring CMake...%s\n", Cyan, Reset)
		cmd := exec.Command("cmake", "-B", buildDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cmake configure failed: %w", err)
		}
	}

	// Build with syntax check only (using -fsyntax-only would be ideal but cmake doesn't support it directly)
	// Instead we do a quick compile
	fmt.Printf("%s🔧 Compiling...%s\n", Cyan, Reset)
	cmd := exec.Command("cmake", "--build", buildDir, "--", "-j", fmt.Sprintf("%d", runtime.NumCPU()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	fmt.Printf("%s✅ Check passed!%s\n", Green, Reset)
	return nil
}

// ============================================================================
// DOC COMMAND
// ============================================================================

func cmdDoc(args []string) {
	fs := flag.NewFlagSet("doc", flag.ExitOnError)
	open := fs.Bool("open", false, "Open documentation in browser")
	fs.Parse(args)

	if err := generateDocs(*open); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func generateDocs(openBrowser bool) error {
	// Check if Doxygen is available
	if _, err := exec.LookPath("doxygen"); err != nil {
		return fmt.Errorf("doxygen not found. Please install it first:\n  macOS: brew install doxygen\n  Ubuntu: sudo apt install doxygen")
	}

	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	fmt.Printf("%s📚 Generating documentation...%s\n", Cyan, Reset)

	// Create Doxyfile if it doesn't exist
	if _, err := os.Stat("Doxyfile"); os.IsNotExist(err) {
		doxyContent := fmt.Sprintf(`PROJECT_NAME           = "%s"
PROJECT_NUMBER         = "%s"
OUTPUT_DIRECTORY       = docs
INPUT                  = src include
RECURSIVE              = YES
EXTRACT_ALL            = YES
GENERATE_HTML          = YES
GENERATE_LATEX         = NO
HTML_OUTPUT            = html
USE_MDFILE_AS_MAINPAGE = README.md
`, config.Package.Name, config.Package.Version)

		if err := os.WriteFile("Doxyfile", []byte(doxyContent), 0644); err != nil {
			return fmt.Errorf("failed to create Doxyfile: %w", err)
		}
		fmt.Printf("   ✓ Created Doxyfile\n")
	}

	// Run Doxygen
	cmd := exec.Command("doxygen")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("doxygen failed: %w", err)
	}

	indexPath := "docs/html/index.html"
	fmt.Printf("%s✅ Documentation generated at %s%s\n", Green, indexPath, Reset)

	if openBrowser {
		var openCmd string
		switch runtime.GOOS {
		case "darwin":
			openCmd = "open"
		case "linux":
			openCmd = "xdg-open"
		case "windows":
			openCmd = "start"
		}

		if openCmd != "" {
			exec.Command(openCmd, indexPath).Start()
		}
	}

	return nil
}

// ============================================================================
// RELEASE COMMAND
// ============================================================================

func cmdRelease(args []string) {
	fs := flag.NewFlagSet("release", flag.ExitOnError)
	fs.Parse(args)

	remaining := fs.Args()
	bumpType := "patch"
	if len(remaining) > 0 {
		bumpType = remaining[0]
	}

	if err := bumpVersion(bumpType); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}
}

func bumpVersion(bumpType string) error {
	config, err := loadConfig(DefaultCfgFile)
	if err != nil {
		return err
	}

	version := config.Package.Version
	if version == "" {
		version = "0.1.0"
	}

	// Parse version
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) < 3 {
		parts = append(parts, make([]string, 3-len(parts))...)
	}

	major, minor, patch := 0, 0, 0
	fmt.Sscanf(parts[0], "%d", &major)
	fmt.Sscanf(parts[1], "%d", &minor)
	fmt.Sscanf(parts[2], "%d", &patch)

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return fmt.Errorf("invalid bump type: %s (use major, minor, or patch)", bumpType)
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	config.Package.Version = newVersion

	fmt.Printf("%s📦 Bumping version: %s → %s%s\n", Cyan, version, newVersion, Reset)

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Printf("%s✅ Version updated to %s%s\n", Green, newVersion, Reset)
	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func loadConfig(path string) (*ForgeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var config ForgeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &config, nil
}

func saveConfig(config *ForgeConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := "# forge.yaml - C++ Project Dependencies\n# Like Cargo.toml for Rust, but for C++!\n\n"
	data = append([]byte(header), data...)

	if err := os.WriteFile(DefaultCfgFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func getAllLibraries(serverURL string) ([]Library, error) {
	url := fmt.Sprintf("%s/api/libraries", serverURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var result struct {
		Libraries []Library `json:"libraries"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Libraries, nil
}

func getLibraryInfo(serverURL, libID string) (*Library, error) {
	libs, err := getAllLibraries(serverURL)
	if err != nil {
		return nil, err
	}

	for _, lib := range libs {
		if lib.ID == libID {
			return &lib, nil
		}
	}

	return nil, fmt.Errorf("library not found")
}

func generateLockFile(config ForgeConfig, outputDir string) error {
	lock := LockConfig{
		Version:      1,
		Dependencies: make(map[string]LockEntry),
	}

	// For now, just record the dependencies without specific commits
	for libID := range config.Dependencies {
		lock.Dependencies[libID] = LockEntry{
			Tag: "latest",
		}
	}

	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}

	header := "# forge.lock - Auto-generated, do not edit\n# This file ensures reproducible builds\n\n"
	data = append([]byte(header), data...)

	return os.WriteFile(filepath.Join(outputDir, LockFile), data, 0644)
}

func extractZip(data []byte, outputDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(outputDir, file.Name)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(absPath, absOutputDir) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0755)

		outFile, err := os.Create(path)
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		fmt.Printf("   📄 %s\n", file.Name)
	}

	return nil
}

// ============================================================================
// UPGRADE COMMAND - Upgrade forge to the latest version
// ============================================================================

func cmdUpgrade(args []string) {
	fmt.Printf("%s🔄 Checking for updates...%s\n", Cyan, Reset)

	// Get latest version from GitHub releases API
	resp, err := http.Get("https://api.github.com/repos/ozacod/forge/releases/latest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to check for updates: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to parse release info: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := Version

	if latestVersion == currentVersion {
		fmt.Printf("%s✓ You're already running the latest version (%s)%s\n", Green, currentVersion, Reset)
		return
	}

	fmt.Printf("%s📦 New version available: %s → %s%s\n", Yellow, currentVersion, latestVersion, Reset)

	// Determine platform and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var binaryName string
	switch goos {
	case "darwin":
		binaryName = fmt.Sprintf("forge-darwin-%s", goarch)
	case "linux":
		binaryName = fmt.Sprintf("forge-linux-%s", goarch)
	case "windows":
		binaryName = fmt.Sprintf("forge-windows-%s.exe", goarch)
	default:
		fmt.Fprintf(os.Stderr, "%sError:%s Unsupported platform: %s\n", Red, Reset, goos)
		os.Exit(1)
	}

	downloadURL := fmt.Sprintf("https://github.com/ozacod/forge/releases/download/%s/%s", release.TagName, binaryName)
	fmt.Printf("%s⬇ Downloading %s...%s\n", Cyan, binaryName, Reset)

	// Download the new binary
	resp, err = http.Get(downloadURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to download: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "%sError:%s Download failed with status %d\n", Red, Reset, resp.StatusCode)
		os.Exit(1)
	}

	binaryData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to read download: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to get executable path: %v\n", Red, Reset, err)
		os.Exit(1)
	}
	execPath, _ = filepath.EvalSymlinks(execPath)

	// Create backup
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to create backup: %v\n", Red, Reset, err)
		fmt.Fprintf(os.Stderr, "Try running with sudo: sudo forge upgrade\n")
		os.Exit(1)
	}

	// Write new binary
	if err := os.WriteFile(execPath, binaryData, 0755); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, execPath)
		fmt.Fprintf(os.Stderr, "%sError:%s Failed to write new binary: %v\n", Red, Reset, err)
		fmt.Fprintf(os.Stderr, "Try running with sudo: sudo forge upgrade\n")
		os.Exit(1)
	}

	// Remove backup
	os.Remove(backupPath)

	fmt.Printf("%s✓ Successfully upgraded to %s!%s\n", Green, latestVersion, Reset)
	fmt.Printf("  Run %sforge version%s to verify.\n", Cyan, Reset)
}

// Unused but kept for potential future use
var _ = bufio.Reader{}
var _ = sort.Strings
