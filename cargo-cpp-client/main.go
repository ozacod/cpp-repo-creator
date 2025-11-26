package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Version        = "1.0.0"
	DefaultServer  = "http://localhost:8000"
	DefaultCfgFile = "cpp-cargo.yaml"
)

// CargoConfig represents the cpp-cargo.yaml structure
type CargoConfig struct {
	Package struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		CppStandard int    `yaml:"cpp_standard"`
	} `yaml:"package"`
	Build struct {
		SharedLibs  bool   `yaml:"shared_libs"`
		ClangFormat string `yaml:"clang_format"`
	} `yaml:"build"`
	Testing struct {
		Framework string `yaml:"framework"`
	} `yaml:"testing"`
	Dependencies map[string]map[string]interface{} `yaml:"dependencies"`
}

func main() {
	// Define flags
	var (
		serverURL    string
		configFile   string
		outputDir    string
		showVersion  bool
		initProject  bool
		listLibs     bool
		templateName string
	)

	flag.StringVar(&serverURL, "server", DefaultServer, "Server URL")
	flag.StringVar(&serverURL, "s", DefaultServer, "Server URL (shorthand)")
	flag.StringVar(&configFile, "config", DefaultCfgFile, "Config file path")
	flag.StringVar(&configFile, "c", DefaultCfgFile, "Config file path (shorthand)")
	flag.StringVar(&outputDir, "output", ".", "Output directory")
	flag.StringVar(&outputDir, "o", ".", "Output directory (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version (shorthand)")
	flag.BoolVar(&initProject, "init", false, "Initialize a new cpp-cargo.yaml")
	flag.BoolVar(&listLibs, "list", false, "List available libraries")
	flag.StringVar(&templateName, "template", "", "Use a template (minimal, web-server, game, cli-tool, networking, data-processing)")
	flag.StringVar(&templateName, "t", "", "Use a template (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `cargo-cpp - C++ Project Generator (like Cargo for Rust)

Usage:
  cargo-cpp [flags]
  cargo-cpp build      Generate project from cpp-cargo.yaml
  cargo-cpp init       Create a new cpp-cargo.yaml
  cargo-cpp list       List available libraries

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  cargo-cpp build                    # Build from cpp-cargo.yaml in current dir
  cargo-cpp build -c myconfig.yaml   # Build from specific config
  cargo-cpp build -o ./myproject     # Output to specific directory
  cargo-cpp init                     # Create cpp-cargo.yaml template
  cargo-cpp init -t web-server       # Create from template
  cargo-cpp list                     # Show available libraries
  cargo-cpp -s http://myserver:8000  # Use custom server

`)
	}

	flag.Parse()

	// Handle version flag
	if showVersion {
		fmt.Printf("cargo-cpp version %s\n", Version)
		return
	}

	// Get positional argument (command)
	args := flag.Args()
	command := "build"
	if len(args) > 0 {
		command = args[0]
	}

	// Override with flags
	if initProject {
		command = "init"
	}
	if listLibs {
		command = "list"
	}

	switch command {
	case "build":
		if err := buildProject(serverURL, configFile, outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "init":
		if err := initConfig(serverURL, templateName, configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "list":
		if err := listLibraries(serverURL); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}
}

func buildProject(serverURL, configFile, outputDir string) error {
	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file '%s': %w", configFile, err)
	}

	// Parse YAML to get project name
	var config CargoConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	projectName := config.Package.Name
	if projectName == "" {
		projectName = "my_project"
	}

	fmt.Printf("🔨 Building project '%s'...\n", projectName)
	fmt.Printf("   Server: %s\n", serverURL)
	fmt.Printf("   Config: %s\n", configFile)

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
	url := fmt.Sprintf("%s/api/cargo", serverURL)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w\n\nMake sure the server is running:\n  cd cargo-cpp-server && uvicorn main:app --port 8000", err)
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

	// Extract ZIP to output directory (files go directly into outputDir)
	fmt.Printf("📦 Extracting to %s...\n", outputDir)

	if err := extractZip(zipData, outputDir); err != nil {
		return fmt.Errorf("failed to extract project: %w", err)
	}

	fmt.Printf("✅ Project '%s' created successfully!\n\n", projectName)
	fmt.Printf("Next steps:\n")
	if outputDir != "." {
		fmt.Printf("  cd %s\n", outputDir)
	}
	fmt.Printf("  cmake -B build\n")
	fmt.Printf("  cmake --build build\n")

	return nil
}

func extractZip(data []byte, outputDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	// Get absolute path of output directory
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

		// Prevent path traversal
		if !strings.HasPrefix(absPath, absOutputDir) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// Extract file
		outFile, err := os.Create(path)
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}

		fmt.Printf("   📄 %s\n", file.Name)
	}

	return nil
}

func initConfig(serverURL, templateName, outputFile string) error {
	var url string
	if templateName != "" {
		url = fmt.Sprintf("%s/api/cargo/example/%s", serverURL, templateName)
		fmt.Printf("📋 Fetching '%s' template...\n", templateName)
	} else {
		url = fmt.Sprintf("%s/api/cargo/template", serverURL)
		fmt.Printf("📋 Fetching default template...\n")
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w\n\nMake sure the server is running:\n  cd cargo-cpp-server && uvicorn main:app --port 8000", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file '%s' already exists. Use a different name or delete it first", outputFile)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Created %s\n\n", outputFile)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  1. Edit %s to customize your project\n", outputFile)
	fmt.Printf("  2. Run: cargo-cpp build\n")

	return nil
}

func listLibraries(serverURL string) error {
	url := fmt.Sprintf("%s/api/libraries", serverURL)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w\n\nMake sure the server is running:\n  cd cargo-cpp-server && uvicorn main:app --port 8000", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Libraries []struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Category    string   `json:"category"`
			HeaderOnly  bool     `json:"header_only"`
			CppStandard int      `json:"cpp_standard"`
			Tags        []string `json:"tags"`
		} `json:"libraries"`
	}

	if err := parseJSON(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Group by category
	categories := make(map[string][]struct {
		ID          string
		Name        string
		Description string
		HeaderOnly  bool
		CppStandard int
	})

	for _, lib := range result.Libraries {
		categories[lib.Category] = append(categories[lib.Category], struct {
			ID          string
			Name        string
			Description string
			HeaderOnly  bool
			CppStandard int
		}{lib.ID, lib.Name, lib.Description, lib.HeaderOnly, lib.CppStandard})
	}

	fmt.Printf("📚 Available Libraries (%d total)\n\n", len(result.Libraries))

	// Print by category
	categoryOrder := []string{
		"serialization", "logging", "testing", "networking", "cli",
		"configuration", "gui", "formatting", "concurrency", "utility",
		"database", "compression", "math", "cryptography",
	}

	for _, cat := range categoryOrder {
		libs, ok := categories[cat]
		if !ok || len(libs) == 0 {
			continue
		}

		fmt.Printf("  %s:\n", strings.Title(cat))
		for _, lib := range libs {
			headerOnly := ""
			if lib.HeaderOnly {
				headerOnly = " [header-only]"
			}
			fmt.Printf("    • %-20s C++%d%s\n", lib.ID, lib.CppStandard, headerOnly)
		}
		fmt.Println()
	}

	fmt.Printf("Use in cpp-cargo.yaml:\n")
	fmt.Printf("  dependencies:\n")
	fmt.Printf("    spdlog:\n")
	fmt.Printf("      spdlog_header_only: true\n")
	fmt.Printf("    nlohmann_json: {}\n")

	return nil
}

func parseJSON(r io.Reader, v interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Simple JSON parsing without encoding/json import
	// We'll use a basic approach
	return yaml.Unmarshal(data, v) // YAML is a superset of JSON
}
