package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ozacod/forge/forge-server-go/internal/generator"
	"github.com/ozacod/forge/forge-server-go/internal/recipe"
	"gopkg.in/yaml.v3"
)

const (
	Version    = "1.0.12"
	CLIVersion = "1.0.12"
)

var projectNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

type ProjectConfig struct {
	ProjectName      string             `json:"project_name" binding:"required"`
	CppStandard      int                `json:"cpp_standard"`
	Libraries        []LibrarySelection `json:"libraries"`
	IncludeTests     bool               `json:"include_tests"`
	TestingFramework string             `json:"testing_framework"`
	BuildShared      bool               `json:"build_shared"`
	ClangFormatStyle string             `json:"clang_format_style"`
	ProjectType      string             `json:"project_type"`
}

type LibrarySelection struct {
	LibraryID string         `json:"library_id" binding:"required"`
	Options   map[string]any `json:"options"`
}

type ForgeYAML struct {
	Package struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		CppStandard int    `yaml:"cpp_standard"`
		ProjectType string `yaml:"project_type"`
	} `yaml:"package"`
	Build struct {
		SharedLibs   bool   `yaml:"shared_libs"`
		ClangFormat string `yaml:"clang_format"`
	} `yaml:"build"`
	Testing struct {
		Framework string `yaml:"framework"`
	} `yaml:"testing"`
	Dependencies map[string]any `yaml:"dependencies"`
}

func main() {
	// Initialize recipe loader
	recipesDir := "recipes"
	if envDir := os.Getenv("FORGE_RECIPES_DIR"); envDir != "" {
		recipesDir = envDir
	}
	loader := recipe.NewLoader(recipesDir)

	// Load recipes
	if err := loader.LoadRecipes(); err != nil {
		fmt.Printf("Warning: Failed to load recipes: %v\n", err)
	}

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		api.GET("", apiRoot)
		api.GET("/version", getVersion)
		api.GET("/libraries", getAllLibraries(loader))
		api.GET("/libraries/:id", getLibrary(loader))
		api.GET("/categories", getCategories)
		api.GET("/categories/:id/libraries", getCategoryLibraries(loader))
		api.GET("/search", searchLibraries(loader))
		api.POST("/reload-recipes", reloadRecipes(loader))
		api.POST("/generate", generateProject(loader))
		api.POST("/preview", previewCMake(loader))
		api.GET("/preview", previewCMakeLegacy(loader))
		api.POST("/forge", generateFromForgeYAML(loader))
		api.POST("/forge/dependencies", generateDependenciesOnly(loader))
		api.GET("/forge/template", getForgeTemplate)
		api.GET("/forge/example/:template", getForgeExample)
	}

	// Static file serving
	staticDir := "static"
	hasStatic := false
	if _, err := os.Stat(staticDir); err == nil {
		if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err == nil {
			hasStatic = true
			// Serve static assets
			r.Static("/assets", filepath.Join(staticDir, "assets"))
			r.StaticFile("/forge.svg", filepath.Join(staticDir, "forge.svg"))

			// Serve index.html for root
			r.GET("/", func(c *gin.Context) {
				c.File(filepath.Join(staticDir, "index.html"))
			})

			// Catch-all for SPA routes
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				if strings.HasPrefix(path, "/api") {
					c.JSON(http.StatusNotFound, gin.H{"detail": "Not found"})
					return
				}
				c.File(filepath.Join(staticDir, "index.html"))
			})
		}
	}

	// Fallback root if no static files
	if !hasStatic {
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Forge API - C++ Project Generator",
				"version":     Version,
				"cli_version": CLIVersion,
				"docs":        "/docs",
				"frontend":    "Not built. Run 'make build-frontend-go' to build the UI.",
			})
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Forge server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func apiRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":     "Forge API - C++ Project Generator",
		"version":     Version,
		"cli_version": CLIVersion,
		"docs":        "/docs",
	})
}

func getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     Version,
		"cli_version": CLIVersion,
		"name":        "forge",
		"description": "C++ Project Generator - Like Cargo for Rust, but for C++!",
	})
}

func getAllLibraries(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		libraries, err := loader.GetAllLibraries()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"libraries": libraries})
	}
}

func getLibrary(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		lib, err := loader.GetLibraryByID(id)
		if err != nil || lib == nil {
			c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Library '%s' not found", id)})
			return
		}
		c.JSON(http.StatusOK, lib)
	}
}

func getCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"categories": recipe.Categories})
}

func getCategoryLibraries(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		libraries, err := loader.GetLibrariesByCategory(categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(libraries) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Category '%s' not found or empty", categoryID)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"libraries": libraries})
	}
}

func searchLibraries(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" || len(query) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "Search query must be at least 2 characters"})
			return
		}
		results, err := loader.SearchLibraries(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": results,
			"count":  len(results),
		})
	}
}

func reloadRecipes(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := loader.ReloadRecipes(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		libraries, _ := loader.GetAllLibraries()
		c.JSON(http.StatusOK, gin.H{
			"message": "Recipes reloaded",
			"count":   len(libraries),
		})
	}
}

func generateProject(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		var config ProjectConfig
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}

		// Validate project name
		if !projectNameRegex.MatchString(config.ProjectName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "Project name must start with a letter and contain only letters, numbers, and underscores",
			})
			return
		}

		// Set defaults
		if config.CppStandard == 0 {
			config.CppStandard = 17
		}
		if config.TestingFramework == "" {
			config.TestingFramework = "googletest"
		}
		if config.ClangFormatStyle == "" {
			config.ClangFormatStyle = "Google"
		}
		if config.ProjectType == "" {
			config.ProjectType = "exe"
		}

		// Validate library IDs
		var invalidLibs []string
		var selections []generator.LibrarySelection
		for _, libSel := range config.Libraries {
			lib, err := loader.GetLibraryByID(libSel.LibraryID)
			if err != nil || lib == nil {
				invalidLibs = append(invalidLibs, libSel.LibraryID)
				continue
			}
			options := libSel.Options
			if options == nil {
				options = make(map[string]any)
			}
			selections = append(selections, generator.LibrarySelection{
				LibraryID: libSel.LibraryID,
				Options:   options,
			})
		}

		if len(invalidLibs) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": fmt.Sprintf("Invalid library IDs: %s", strings.Join(invalidLibs, ", ")),
			})
			return
		}

		// Generate ZIP
		zipData, err := generator.CreateProjectZip(
			config.ProjectName,
			config.CppStandard,
			selections,
			config.IncludeTests,
			config.TestingFramework,
			config.BuildShared,
			config.ClangFormatStyle,
			config.ProjectType,
			false, // not flat for web UI
			loader,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": fmt.Sprintf("Failed to generate project: %v", err)})
			return
		}

		c.Data(http.StatusOK, "application/zip", zipData)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", config.ProjectName))
	}
}

func previewCMake(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		var config ProjectConfig
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}

		// Validate project name
		if !projectNameRegex.MatchString(config.ProjectName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "Project name must start with a letter and contain only letters, numbers, and underscores",
			})
			return
		}

		// Set defaults
		if config.CppStandard == 0 {
			config.CppStandard = 17
		}
		if config.TestingFramework == "" {
			config.TestingFramework = "googletest"
		}
		if config.ProjectType == "" {
			config.ProjectType = "exe"
		}

		// Get libraries with their selections
		var librariesWithOptions []generator.LibraryWithOptions
		for _, libSel := range config.Libraries {
			lib, err := loader.GetLibraryByID(libSel.LibraryID)
			if err == nil && lib != nil {
				options := libSel.Options
				if options == nil {
					options = make(map[string]any)
				}
				librariesWithOptions = append(librariesWithOptions, generator.LibraryWithOptions{
					Lib:     lib,
					Options: options,
				})
			}
		}

		cmakeContent, err := generator.GenerateCMakeLists(
			config.ProjectName,
			config.CppStandard,
			librariesWithOptions,
			config.IncludeTests,
			config.TestingFramework,
			config.BuildShared,
			config.ProjectType,
			loader,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cmake_content": cmakeContent})
	}
}

func previewCMakeLegacy(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectName := c.Query("project_name")
		if projectName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"detail": "project_name is required"})
			return
		}

		// Validate project name
		if !projectNameRegex.MatchString(projectName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "Project name must start with a letter and contain only letters, numbers, and underscores",
			})
			return
		}

		cppStandard := 17
		if std := c.Query("cpp_standard"); std != "" {
			fmt.Sscanf(std, "%d", &cppStandard)
		}

		includeTests := c.DefaultQuery("include_tests", "true") == "true"

		// Parse library IDs
		var librariesWithOptions []generator.LibraryWithOptions
		if libraryIDs := c.Query("library_ids"); libraryIDs != "" {
			ids := strings.Split(libraryIDs, ",")
			for _, id := range ids {
				id = strings.TrimSpace(id)
				if id == "" {
					continue
				}
				lib, err := loader.GetLibraryByID(id)
				if err == nil && lib != nil {
					librariesWithOptions = append(librariesWithOptions, generator.LibraryWithOptions{
						Lib:     lib,
						Options: make(map[string]any),
					})
				}
			}
		}

		cmakeContent, err := generator.GenerateCMakeLists(
			projectName,
			cppStandard,
			librariesWithOptions,
			includeTests,
			"googletest",
			false,
			"exe",
			loader,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cmake_content": cmakeContent})
	}
}

func generateFromForgeYAML(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to read file: %v", err)})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to open file: %v", err)})
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to read file: %v", err)})
			return
		}

		var forgeYAML ForgeYAML
		if err := yaml.Unmarshal(data, &forgeYAML); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Invalid YAML format: %v", err)})
			return
		}

		// Extract package info
		projectName := forgeYAML.Package.Name
		if projectName == "" {
			projectName = "my_project"
		}

		// Validate project name
		if !projectNameRegex.MatchString(projectName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": "Project name must start with a letter and contain only letters, numbers, and underscores",
			})
			return
		}

		cppStandard := forgeYAML.Package.CppStandard
		if cppStandard == 0 {
			cppStandard = 17
		}

		projectType := forgeYAML.Package.ProjectType
		if projectType == "" {
			projectType = "exe"
		}
		if projectType != "exe" && projectType != "lib" {
			projectType = "exe"
		}

		// Extract build settings
		buildShared := forgeYAML.Build.SharedLibs
		clangFormatStyle := forgeYAML.Build.ClangFormat
		if clangFormatStyle == "" {
			clangFormatStyle = "Google"
		}

		// Extract testing settings
		testingFramework := forgeYAML.Testing.Framework
		if testingFramework == "" {
			testingFramework = "googletest"
		}
		includeTests := testingFramework != "none"

		// Extract dependencies
		var selections []generator.LibrarySelection
		var invalidLibs []string

		for libID, options := range forgeYAML.Dependencies {
			lib, err := loader.GetLibraryByID(libID)
			if err != nil || lib == nil {
				invalidLibs = append(invalidLibs, libID)
				continue
			}

			opts := make(map[string]any)
			if optionsMap, ok := options.(map[string]any); ok {
				opts = optionsMap
			}

			selections = append(selections, generator.LibrarySelection{
				LibraryID: libID,
				Options:   opts,
			})
		}

		if len(invalidLibs) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail": fmt.Sprintf("Unknown dependencies: %s. Use GET /api/libraries to see available libraries.", strings.Join(invalidLibs, ", ")),
			})
			return
		}

		// Generate ZIP (flat=True for CLI usage)
		zipData, err := generator.CreateProjectZip(
			projectName,
			cppStandard,
			selections,
			includeTests,
			testingFramework,
			buildShared,
			clangFormatStyle,
			projectType,
			true, // flat for CLI
			loader,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": fmt.Sprintf("Failed to generate project: %v", err)})
			return
		}

		c.Data(http.StatusOK, "application/zip", zipData)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", projectName))
	}
}

func generateDependenciesOnly(loader *recipe.Loader) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to read file: %v", err)})
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to open file: %v", err)})
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Failed to read file: %v", err)})
			return
		}

		var forgeYAML ForgeYAML
		if err := yaml.Unmarshal(data, &forgeYAML); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("Invalid YAML format: %v", err)})
			return
		}

		// Extract testing config
		testingFramework := forgeYAML.Testing.Framework
		if testingFramework == "" {
			testingFramework = "none"
		}
		includeTests := testingFramework != "none"

		// Parse dependencies
		var librariesWithOptions []generator.LibraryWithOptions
		for libID, libOptions := range forgeYAML.Dependencies {
			lib, err := loader.GetLibraryByID(libID)
			if err == nil && lib != nil {
				opts := make(map[string]any)
				if optionsMap, ok := libOptions.(map[string]any); ok {
					opts = optionsMap
				}
				librariesWithOptions = append(librariesWithOptions, generator.LibraryWithOptions{
					Lib:     lib,
					Options: opts,
				})
			}
		}

		// Generate dependencies.cmake content
		cmakeContent, err := generator.GenerateDependenciesCMake(
			librariesWithOptions,
			includeTests,
			testingFramework,
			loader,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
			return
		}

		c.String(http.StatusOK, cmakeContent)
	}
}

func getForgeTemplate(c *gin.Context) {
	projectType := c.DefaultQuery("project_type", "exe")

	if projectType == "lib" {
		template := `# forge.yaml - C++ Library Project
# Like Cargo.toml for Rust, but for C++!

package:
  name: my_library
  version: "0.1.0"
  cpp_standard: 17  # 11, 14, 17, 20, or 23
  project_type: lib  # lib = library only (no executable)

build:
  shared_libs: false
  clang_format: Google  # Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU

testing:
  framework: googletest  # googletest, catch2, doctest, or none

dependencies:
  fmt: {}
`
		c.String(http.StatusOK, template)
		return
	}

	template := `# forge.yaml - C++ Project Dependencies
# Like Cargo.toml for Rust, but for C++!

package:
  name: my_awesome_project
  version: "1.0.0"
  cpp_standard: 17  # 11, 14, 17, 20, or 23
  project_type: exe  # exe = executable, lib = library only

build:
  shared_libs: false
  clang_format: Google  # Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU

testing:
  framework: googletest  # googletest, catch2, doctest, or none

# Dependencies and their options
# Use: curl http://localhost:8000/api/libraries to see all available libraries
dependencies:
  # JSON library
  nlohmann_json:
    json_diagnostics: false
    json_install: false
  
  # Logging library with options
  spdlog:
    spdlog_header_only: true
    spdlog_fmt_external: false
  
  # Fast formatting library
  fmt:
    fmt_install: false

# Example: Minimal config
# ---
# package:
#   name: hello_world
# dependencies:
#   fmt: {}
`
	c.String(http.StatusOK, template)
}

func getForgeExample(c *gin.Context) {
	templateName := c.Param("template")
	projectType := c.DefaultQuery("project_type", "exe")

	templates := map[string]string{
		"minimal": fmt.Sprintf(`# Minimal C++ project
package:
  name: hello_cpp
  cpp_standard: 17
  project_type: %s

dependencies:
  fmt: {}
`, projectType),
		"web-server": fmt.Sprintf(`# Web server project
package:
  name: my_web_server
  cpp_standard: 17
  project_type: %s

build:
  clang_format: Google

testing:
  framework: catch2

dependencies:
  crow:
    crow_enable_ssl: false
  nlohmann_json: {}
  spdlog:
    spdlog_header_only: true
`, projectType),
		"game": fmt.Sprintf(`# Game development project
package:
  name: my_game
  cpp_standard: 17
  project_type: %s

build:
  clang_format: Google

testing:
  framework: none

dependencies:
  raylib:
    raylib_build_examples: false
  glm: {}
  entt: {}
  spdlog:
    spdlog_header_only: true
`, projectType),
		"cli-tool": fmt.Sprintf(`# Command-line tool project
package:
  name: my_cli_tool
  cpp_standard: 17
  project_type: %s

build:
  clang_format: Google

testing:
  framework: doctest

dependencies:
  cli11: {}
  fmt: {}
  spdlog:
    spdlog_header_only: true
  indicators: {}
  tabulate: {}
`, projectType),
		"networking": fmt.Sprintf(`# Networking project
package:
  name: my_network_app
  cpp_standard: 17
  project_type: %s

build:
  clang_format: Google

testing:
  framework: googletest

dependencies:
  asio: {}
  nlohmann_json: {}
  spdlog:
    spdlog_header_only: true
  xxhash: {}
`, projectType),
		"data-processing": fmt.Sprintf(`# Data processing project
package:
  name: data_processor
  cpp_standard: 20
  project_type: %s

build:
  clang_format: LLVM

testing:
  framework: catch2

dependencies:
  simdjson: {}
  range_v3: {}
  taskflow: {}
  fmt: {}
  spdlog:
    spdlog_header_only: true
`, projectType),
	}

	template, ok := templates[templateName]
	if !ok {
		keys := make([]string, 0, len(templates))
		for k := range templates {
			keys = append(keys, k)
		}
		c.JSON(http.StatusNotFound, gin.H{
			"detail": fmt.Sprintf("Template '%s' not found. Available: %s", templateName, strings.Join(keys, ", ")),
		})
		return
	}

	c.String(http.StatusOK, template)
}

