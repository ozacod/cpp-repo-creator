package recipe

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const Version = "1.0.12"
const CLIVersion = "1.0.12"

type LibraryOption struct {
	ID                       string   `yaml:"id" json:"id"`
	Name                     string   `yaml:"name" json:"name"`
	Description              string   `yaml:"description" json:"description"`
	Type                     string   `yaml:"type" json:"type"` // boolean, string, choice, integer
	Default                  any      `yaml:"default" json:"default,omitempty"`
	Choices                  []string `yaml:"choices" json:"choices,omitempty"`
	CMakeVar                 string   `yaml:"cmake_var" json:"cmake_var,omitempty"`
	CMakeDefine              string   `yaml:"cmake_define" json:"cmake_define,omitempty"`
	AffectsLink              bool     `yaml:"affects_link" json:"affects_link,omitempty"`
	LinkLibrariesWhenEnabled []string `yaml:"link_libraries_when_enabled" json:"link_libraries_when_enabled,omitempty"`
}

type FetchContent struct {
	Repository   string `yaml:"repository" json:"repository"`
	Tag          string `yaml:"tag" json:"tag"`
	SourceSubdir string `yaml:"source_subdir" json:"source_subdir,omitempty"`
}

type Library struct {
	ID              string          `yaml:"id" json:"id"`
	Name            string          `yaml:"name" json:"name"`
	Description     string          `yaml:"description" json:"description"`
	Category        string          `yaml:"category" json:"category"`
	GitHubURL       string          `yaml:"github_url" json:"github_url"`
	CppStandard     int             `yaml:"cpp_standard" json:"cpp_standard"`
	HeaderOnly      bool            `yaml:"header_only" json:"header_only"`
	Tags            []string        `yaml:"tags" json:"tags"`
	Alternatives    []string        `yaml:"alternatives" json:"alternatives"`
	FetchContent    *FetchContent   `yaml:"fetch_content" json:"fetch_content,omitempty"`
	LinkLibraries   []string        `yaml:"link_libraries" json:"link_libraries"`
	Options         []LibraryOption `yaml:"options" json:"options"`
	CMakePre        string          `yaml:"cmake_pre" json:"cmake_pre,omitempty"`
	CMakePost       string          `yaml:"cmake_post" json:"cmake_post,omitempty"`
	SystemPackage   bool            `yaml:"system_package" json:"system_package,omitempty"`
	FindPackageName string          `yaml:"find_package_name" json:"find_package_name,omitempty"`
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

var Categories = []Category{
	{ID: "serialization", Name: "Serialization", Icon: "📦", Description: "JSON, XML, Binary serialization"},
	{ID: "logging", Name: "Logging", Icon: "📝", Description: "Logging and diagnostics"},
	{ID: "testing", Name: "Testing", Icon: "🧪", Description: "Unit testing and mocking frameworks"},
	{ID: "networking", Name: "Networking", Icon: "🌐", Description: "HTTP, TCP/UDP, async I/O"},
	{ID: "cli", Name: "CLI", Icon: "💻", Description: "Command line argument parsing"},
	{ID: "configuration", Name: "Configuration", Icon: "⚙️", Description: "Config file parsing (YAML, TOML)"},
	{ID: "gui", Name: "GUI", Icon: "🖼️", Description: "Graphical user interfaces"},
	{ID: "formatting", Name: "Formatting", Icon: "✨", Description: "String formatting and text processing"},
	{ID: "concurrency", Name: "Concurrency", Icon: "⚡", Description: "Threading, async, lock-free structures"},
	{ID: "utility", Name: "Utility", Icon: "🔧", Description: "General utilities and helpers"},
	{ID: "database", Name: "Database", Icon: "💾", Description: "Database clients and ORMs"},
	{ID: "compression", Name: "Compression", Icon: "🗜️", Description: "Data compression libraries"},
	{ID: "math", Name: "Math", Icon: "📐", Description: "Mathematics and linear algebra"},
	{ID: "cryptography", Name: "Cryptography", Icon: "🔐", Description: "Encryption and cryptographic functions"},
}

type Loader struct {
	recipesDir string
	fs         fs.FS
	libraries  map[string]*Library
	loaded     bool
}

func NewLoader(recipesDir string) *Loader {
	if recipesDir == "" {
		recipesDir = "recipes"
	}
	return &Loader{
		recipesDir: recipesDir,
		fs:         nil,
		libraries:  make(map[string]*Library),
		loaded:     false,
	}
}

func NewLoaderWithFS(recipesFS fs.FS, recipesDir string) *Loader {
	return &Loader{
		recipesDir: recipesDir,
		fs:         recipesFS,
		libraries:  make(map[string]*Library),
		loaded:     false,
	}
}

func (l *Loader) LoadRecipes() error {
	if l.loaded {
		return nil
	}

	var entries []fs.DirEntry
	var err error

	if l.fs != nil {
		entries, err = fs.ReadDir(l.fs, l.recipesDir)
		if err != nil {
			return fmt.Errorf("failed to read embedded recipes directory: %w", err)
		}
	} else {
		if _, err := os.Stat(l.recipesDir); os.IsNotExist(err) {
			return fmt.Errorf("recipes directory not found: %s", l.recipesDir)
		}
		entries, err = os.ReadDir(l.recipesDir)
		if err != nil {
			return fmt.Errorf("failed to read recipes directory: %w", err)
		}
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		if strings.HasPrefix(entry.Name(), "_") {
			continue
		}

		filepath := filepath.Join(l.recipesDir, entry.Name())
		lib, err := l.loadRecipeFile(filepath)
		if err != nil {
			fmt.Printf("Warning: Failed to load recipe %s: %v\n", filepath, err)
			continue
		}
		if lib != nil {
			l.libraries[lib.ID] = lib
		}
	}

	l.loaded = true
	return nil
}

func (l *Loader) loadRecipeFile(filepath string) (*Library, error) {
	var data []byte
	var err error

	if l.fs != nil {
		data, err = fs.ReadFile(l.fs, filepath)
	} else {
		data, err = os.ReadFile(filepath)
	}

	if err != nil {
		return nil, err
	}

	var lib Library
	if err := yaml.Unmarshal(data, &lib); err != nil {
		return nil, err
	}

	if lib.ID == "" {
		return nil, fmt.Errorf("missing id field")
	}

	// Set defaults
	if lib.Name == "" {
		lib.Name = lib.ID
	}
	if lib.Category == "" {
		lib.Category = "utility"
	}
	if lib.CppStandard == 0 {
		lib.CppStandard = 11
	}
	if lib.LinkLibraries == nil {
		lib.LinkLibraries = []string{}
	}
	if lib.Options == nil {
		lib.Options = []LibraryOption{}
	}
	if lib.Tags == nil {
		lib.Tags = []string{}
	}
	if lib.Alternatives == nil {
		lib.Alternatives = []string{}
	}

	return &lib, nil
}

func (l *Loader) GetAllLibraries() ([]*Library, error) {
	if err := l.LoadRecipes(); err != nil {
		return nil, err
	}
	libraries := make([]*Library, 0, len(l.libraries))
	for _, lib := range l.libraries {
		libraries = append(libraries, lib)
	}
	return libraries, nil
}

func (l *Loader) GetLibraryByID(id string) (*Library, error) {
	if err := l.LoadRecipes(); err != nil {
		return nil, err
	}
	return l.libraries[id], nil
}

func (l *Loader) GetLibrariesByCategory(category string) ([]*Library, error) {
	if err := l.LoadRecipes(); err != nil {
		return nil, err
	}
	var result []*Library
	for _, lib := range l.libraries {
		if lib.Category == category {
			result = append(result, lib)
		}
	}
	return result, nil
}

func (l *Loader) SearchLibraries(query string) ([]*Library, error) {
	if err := l.LoadRecipes(); err != nil {
		return nil, err
	}
	query = strings.ToLower(query)
	var result []*Library
	for _, lib := range l.libraries {
		if strings.Contains(strings.ToLower(lib.Name), query) ||
			strings.Contains(strings.ToLower(lib.Description), query) {
			result = append(result, lib)
			continue
		}
		for _, tag := range lib.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				result = append(result, lib)
				break
			}
		}
	}
	return result, nil
}

func (l *Loader) ReloadRecipes() error {
	l.libraries = make(map[string]*Library)
	l.loaded = false
	return l.LoadRecipes()
}
