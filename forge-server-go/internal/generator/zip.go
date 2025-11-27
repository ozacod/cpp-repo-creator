package generator

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/ozacod/forge/forge-server-go/internal/recipe"
)

func CreateProjectZip(
	projectName string,
	cppStandard int,
	librarySelections []LibrarySelection,
	includeTests bool,
	testingFramework string,
	buildShared bool,
	clangFormatStyle string,
	projectType string,
	flat bool,
	loader *recipe.Loader,
) ([]byte, error) {
	// Get library objects with their options
	var librariesWithOptions []LibraryWithOptions
	var allLibraries []*recipe.Library

	for _, selection := range librarySelections {
		lib, err := loader.GetLibraryByID(selection.LibraryID)
		if err != nil {
			continue
		}
		if lib != nil {
			options := selection.Options
			if options == nil {
				options = make(map[string]any)
			}
			librariesWithOptions = append(librariesWithOptions, LibraryWithOptions{
				Lib:     lib,
				Options: options,
			})
			allLibraries = append(allLibraries, lib)
		}
	}

	// Separate test libraries from main libraries
	var testLibraries, mainLibraries []LibraryWithOptions
	for _, lwo := range librariesWithOptions {
		if lwo.Lib.Category == "testing" {
			testLibraries = append(testLibraries, lwo)
		} else {
			mainLibraries = append(mainLibraries, lwo)
		}
	}

	// Add selected testing framework if not already present
	if includeTests && testingFramework != "" && testingFramework != "none" {
		existingTestIDs := make(map[string]bool)
		for _, lwo := range testLibraries {
			existingTestIDs[lwo.Lib.ID] = true
		}
		if !existingTestIDs[testingFramework] {
			testLib, err := loader.GetLibraryByID(testingFramework)
			if err == nil && testLib != nil {
				testLibraries = append([]LibraryWithOptions{{Lib: testLib, Options: map[string]any{}}}, testLibraries...)
			}
		}
	}

	testLibsOnly := make([]*recipe.Library, 0, len(testLibraries))
	for _, lwo := range testLibraries {
		testLibsOnly = append(testLibsOnly, lwo.Lib)
	}

	// Create in-memory ZIP file
	var zipBuffer bytes.Buffer
	zw := zip.NewWriter(&zipBuffer)

	// Use empty prefix for flat mode (CLI), project_name for wrapped mode (web UI)
	prefix := ""
	if !flat {
		prefix = projectName + "/"
	}

	// .cmake/forge/dependencies.cmake
	depsCMake, err := GenerateDependenciesCMake(librariesWithOptions, includeTests, testingFramework, loader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dependencies.cmake: %w", err)
	}
	if err := writeZipFile(zw, prefix+".cmake/forge/dependencies.cmake", depsCMake); err != nil {
		return nil, err
	}

	// CMakeLists.txt
	cmakeLists, err := GenerateCMakeLists(projectName, cppStandard, librariesWithOptions, includeTests, testingFramework, buildShared, projectType, loader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CMakeLists.txt: %w", err)
	}
	if err := writeZipFile(zw, prefix+"CMakeLists.txt", cmakeLists); err != nil {
		return nil, err
	}

	// README.md
	readme := GenerateReadme(projectName, allLibraries, cppStandard, projectType)
	if err := writeZipFile(zw, prefix+"README.md", readme); err != nil {
		return nil, err
	}

	// .gitignore
	gitignore := GenerateGitignore()
	if err := writeZipFile(zw, prefix+".gitignore", gitignore); err != nil {
		return nil, err
	}

	// .clang-format
	clangFormat := GenerateClangFormat(clangFormatStyle)
	if err := writeZipFile(zw, prefix+".clang-format", clangFormat); err != nil {
		return nil, err
	}

	// Include directory
	header := GenerateLibHeader(projectName)
	if err := writeZipFile(zw, prefix+fmt.Sprintf("include/%s/%s.hpp", projectName, projectName), header); err != nil {
		return nil, err
	}

	// Source directory - only include main.cpp for executable projects
	if projectType == "exe" {
		mainCpp := GenerateMainCpp(projectName, allLibraries)
		if err := writeZipFile(zw, prefix+fmt.Sprintf("src/main.cpp", projectName), mainCpp); err != nil {
			return nil, err
		}
	}
	libSource := GenerateLibSource(projectName, allLibraries)
	if err := writeZipFile(zw, prefix+fmt.Sprintf("src/%s.cpp", projectName), libSource); err != nil {
		return nil, err
	}

	// Tests directory
	if includeTests {
		testCMake := GenerateTestCMake(projectName, testLibraries, mainLibraries, projectType)
		if err := writeZipFile(zw, prefix+"tests/CMakeLists.txt", testCMake); err != nil {
			return nil, err
		}
		testMain := GenerateTestMain(projectName, testLibsOnly)
		if err := writeZipFile(zw, prefix+"tests/test_main.cpp", testMain); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return zipBuffer.Bytes(), nil
}

func writeZipFile(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create zip entry %s: %w", name, err)
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write zip entry %s: %w", name, err)
	}
	return nil
}

