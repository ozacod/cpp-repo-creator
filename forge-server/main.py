"""
FastAPI backend for C++ Project Creator.
"""

from fastapi import FastAPI, HTTPException, File, UploadFile
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import Response, PlainTextResponse
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import re
import yaml

from recipe_loader import (
    get_all_libraries,
    get_library_by_id,
    get_libraries_by_category,
    search_libraries,
    get_categories,
    reload_recipes,
)
from generator import create_project_zip, generate_cmake_lists

# Version info
VERSION = "1.0.11"
CLI_VERSION = "1.0.11"

app = FastAPI(
    title="Forge API",
    description="API for generating C++ project templates with CMake and FetchContent",
    version=VERSION,
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class LibrarySelection(BaseModel):
    """Library selection with options."""
    library_id: str
    options: Dict[str, Any] = Field(default_factory=dict)


class ProjectConfig(BaseModel):
    """Project configuration for generation."""
    project_name: str = Field(..., min_length=1, max_length=50, description="Project name")
    cpp_standard: int = Field(default=17, ge=11, le=23, description="C++ standard version")
    libraries: List[LibrarySelection] = Field(default=[], description="List of library selections with options")
    include_tests: bool = Field(default=True, description="Include test configuration")
    testing_framework: str = Field(default="googletest", description="Testing framework (none, googletest, catch2, doctest)")
    build_shared: bool = Field(default=False, description="Build as shared libraries")
    clang_format_style: str = Field(default="Google", description="Clang-format style (Google, LLVM, Chromium, Mozilla, WebKit, Microsoft, GNU)")
    project_type: str = Field(default="exe", description="Project type (exe for executable, lib for library)")

    class Config:
        json_schema_extra = {
            "example": {
                "project_name": "my_project",
                "cpp_standard": 17,
                "libraries": [
                    {"library_id": "spdlog", "options": {"spdlog_header_only": True}},
                    {"library_id": "nlohmann_json", "options": {}},
                ],
                "include_tests": True,
                "testing_framework": "googletest",
                "build_shared": False,
                "clang_format_style": "Google",
            }
        }


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "message": "Forge API - C++ Project Generator",
        "version": VERSION,
        "cli_version": CLI_VERSION,
        "docs": "/docs",
    }


@app.get("/api/version")
async def get_version():
    """Get version information."""
    return {
        "version": VERSION,
        "cli_version": CLI_VERSION,
        "name": "forge",
        "description": "C++ Project Generator - Like Cargo for Rust, but for C++!",
    }


@app.get("/api/libraries")
async def get_all_libraries_endpoint():
    """Get all available libraries with their options."""
    return {"libraries": get_all_libraries()}


@app.get("/api/libraries/{library_id}")
async def get_library(library_id: str):
    """Get a specific library by ID."""
    library = get_library_by_id(library_id)
    if not library:
        raise HTTPException(status_code=404, detail=f"Library '{library_id}' not found")
    return library


@app.get("/api/categories")
async def get_all_categories():
    """Get all library categories."""
    return {"categories": get_categories()}


@app.get("/api/categories/{category_id}/libraries")
async def get_category_libraries(category_id: str):
    """Get all libraries in a specific category."""
    libraries = get_libraries_by_category(category_id)
    if not libraries:
        raise HTTPException(status_code=404, detail=f"Category '{category_id}' not found or empty")
    return {"libraries": libraries}


@app.get("/api/search")
async def search(q: str):
    """Search libraries by name, description, or tags."""
    if not q or len(q) < 2:
        raise HTTPException(status_code=400, detail="Search query must be at least 2 characters")
    results = search_libraries(q)
    return {"query": q, "results": results, "count": len(results)}


@app.post("/api/reload-recipes")
async def reload_recipes_endpoint():
    """Reload all recipes from files."""
    reload_recipes()
    return {"message": "Recipes reloaded", "count": len(get_all_libraries())}


@app.post("/api/generate")
async def generate_project(config: ProjectConfig):
    """Generate a C++ project and return it as a ZIP file."""
    
    # Validate project name (alphanumeric and underscores only)
    if not re.match(r'^[a-zA-Z][a-zA-Z0-9_]*$', config.project_name):
        raise HTTPException(
            status_code=400,
            detail="Project name must start with a letter and contain only letters, numbers, and underscores"
        )
    
    # Validate library IDs
    invalid_libs = []
    for lib_selection in config.libraries:
        if not get_library_by_id(lib_selection.library_id):
            invalid_libs.append(lib_selection.library_id)
    
    if invalid_libs:
        raise HTTPException(
            status_code=400,
            detail=f"Invalid library IDs: {', '.join(invalid_libs)}"
        )
    
    # Generate the ZIP file
    try:
        zip_content = create_project_zip(
            project_name=config.project_name,
            cpp_standard=config.cpp_standard,
            library_selections=config.libraries,
            include_tests=config.include_tests,
            testing_framework=config.testing_framework,
            build_shared=config.build_shared,
            clang_format_style=config.clang_format_style,
            project_type=config.project_type,
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to generate project: {str(e)}")
    
    return Response(
        content=zip_content,
        media_type="application/zip",
        headers={
            "Content-Disposition": f"attachment; filename={config.project_name}.zip"
        }
    )


@app.post("/api/preview")
async def preview_cmake(config: ProjectConfig):
    """Preview the generated CMakeLists.txt without downloading."""
    
    # Validate project name
    if not re.match(r'^[a-zA-Z][a-zA-Z0-9_]*$', config.project_name):
        raise HTTPException(
            status_code=400,
            detail="Project name must start with a letter and contain only letters, numbers, and underscores"
        )
    
    # Get libraries with their selections
    libraries_with_options = []
    for lib_selection in config.libraries:
        lib = get_library_by_id(lib_selection.library_id)
        if lib:
            libraries_with_options.append((lib, lib_selection.options))
    
    cmake_content = generate_cmake_lists(
        config.project_name,
        config.cpp_standard,
        libraries_with_options,
        config.include_tests,
        config.testing_framework,
        config.build_shared,
    )
    
    return {"cmake_content": cmake_content}


# Legacy endpoint for backwards compatibility
@app.get("/api/preview")
async def preview_cmake_legacy(
    project_name: str,
    cpp_standard: int = 17,
    library_ids: Optional[str] = None,
    include_tests: bool = True,
):
    """Preview the generated CMakeLists.txt (legacy endpoint)."""
    
    # Validate project name
    if not re.match(r'^[a-zA-Z][a-zA-Z0-9_]*$', project_name):
        raise HTTPException(
            status_code=400,
            detail="Project name must start with a letter and contain only letters, numbers, and underscores"
        )
    
    # Parse library IDs and create selections with default options
    libraries_with_options = []
    if library_ids:
        lib_id_list = [lid.strip() for lid in library_ids.split(",") if lid.strip()]
        for lib_id in lib_id_list:
            lib = get_library_by_id(lib_id)
            if lib:
                libraries_with_options.append((lib, {}))
    
    cmake_content = generate_cmake_lists(
        project_name,
        cpp_standard,
        libraries_with_options,
        include_tests,
        False,  # build_shared
    )
    
    return {"cmake_content": cmake_content}


@app.post("/api/forge")
async def generate_from_cargo(file: UploadFile = File(...)):
    """
    Generate a C++ project from a forge.yaml file.
    
    The YAML file format:
    ```yaml
    package:
      name: my_project
      version: "1.0.0"
      cpp_standard: 17
    
    build:
      shared_libs: false
      clang_format: Google
    
    testing:
      framework: googletest  # googletest, catch2, doctest, or none
    
    dependencies:
      spdlog:
        header_only: true
      nlohmann_json: {}
      fmt:
        header_only: false
    ```
    
    Usage:
        curl -X POST -F "file=@forge.yaml" http://localhost:8000/api/forge -o project.zip
    """
    
    # Read and parse the YAML file
    try:
        content = await file.read()
        cargo_config = yaml.safe_load(content.decode('utf-8'))
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML format: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Failed to read file: {str(e)}")
    
    # Extract package info
    package = cargo_config.get("package", {})
    project_name = package.get("name", "my_project")
    cpp_standard = package.get("cpp_standard", 17)
    project_type = package.get("project_type", "exe")  # exe or lib
    
    # Validate project name
    if not re.match(r'^[a-zA-Z][a-zA-Z0-9_]*$', project_name):
        raise HTTPException(
            status_code=400,
            detail="Project name must start with a letter and contain only letters, numbers, and underscores"
        )
    
    # Validate project type
    if project_type not in ("exe", "lib"):
        project_type = "exe"
    
    # Extract build settings
    build = cargo_config.get("build", {})
    build_shared = build.get("shared_libs", False)
    clang_format_style = build.get("clang_format", "Google")
    
    # Extract testing settings
    testing = cargo_config.get("testing", {})
    testing_framework = testing.get("framework", "googletest")
    include_tests = testing_framework != "none"
    
    # Extract dependencies
    dependencies = cargo_config.get("dependencies", {})
    library_selections = []
    invalid_libs = []
    
    for lib_id, options in dependencies.items():
        if not get_library_by_id(lib_id):
            invalid_libs.append(lib_id)
        else:
            # Options can be a dict or empty/null
            opts = options if isinstance(options, dict) else {}
            library_selections.append(LibrarySelection(library_id=lib_id, options=opts))
    
    if invalid_libs:
        raise HTTPException(
            status_code=400,
            detail=f"Unknown dependencies: {', '.join(invalid_libs)}. Use GET /api/libraries to see available libraries."
        )
    
    # Generate the ZIP file (flat=True means no wrapper folder, files go directly in current dir)
    try:
        zip_content = create_project_zip(
            project_name=project_name,
            cpp_standard=cpp_standard,
            library_selections=library_selections,
            include_tests=include_tests,
            testing_framework=testing_framework,
            build_shared=build_shared,
            clang_format_style=clang_format_style,
            project_type=project_type,
            flat=True,  # CLI usage: extract directly to current directory
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to generate project: {str(e)}")
    
    return Response(
        content=zip_content,
        media_type="application/zip",
        headers={
            "Content-Disposition": f"attachment; filename={project_name}.zip"
        }
    )


@app.post("/api/forge/dependencies")
async def generate_dependencies_only(file: UploadFile = File(...)):
    """
    Generate only the dependencies.cmake file from a forge.yaml file.
    Used by 'forge add' and 'forge remove' commands.
    
    Returns plain text content of dependencies.cmake.
    """
    from generator import generate_dependencies_cmake
    from pydantic import BaseModel
    
    # Read and parse the YAML file
    try:
        content = await file.read()
        cargo_config = yaml.safe_load(content.decode('utf-8'))
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML format: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Failed to read file: {str(e)}")
    
    # Extract testing config
    testing_config = cargo_config.get("testing", {})
    include_tests = testing_config.get("framework", "none") != "none"
    testing_framework = testing_config.get("framework", "none")
    
    # Parse dependencies
    dependencies = cargo_config.get("dependencies", {})
    
    class LibSelection:
        def __init__(self, library_id: str, options: dict):
            self.library_id = library_id
            self.options = options
    
    library_selections = []
    for lib_id, lib_options in dependencies.items():
        if lib_options is None:
            lib_options = {}
        library_selections.append(LibSelection(lib_id, lib_options))
    
    # Get library objects with their options
    libraries_with_options = []
    for selection in library_selections:
        lib = get_library_by_id(selection.library_id)
        if lib:
            libraries_with_options.append((lib, selection.options))
    
    # Generate dependencies.cmake content
    cmake_content = generate_dependencies_cmake(
        libraries_with_options,
        include_tests,
        testing_framework
    )
    
    return PlainTextResponse(content=cmake_content, media_type="text/plain")


@app.get("/api/forge/template")
async def get_cargo_template(project_type: str = "exe"):
    """Get a sample forge.yaml template.
    
    Args:
        project_type: "exe" for executable project, "lib" for library project
    """
    if project_type == "lib":
        template = """# forge.yaml - C++ Library Project
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
"""
    else:
        template = """# forge.yaml - C++ Project Dependencies
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
"""
    return PlainTextResponse(content=template, media_type="text/yaml")


@app.get("/api/forge/example/{template_name}")
async def get_cargo_example(template_name: str, project_type: str = "exe"):
    """Get example forge.yaml templates for common use cases.
    
    Args:
        template_name: Template name (minimal, web-server, game, cli-tool, networking, data-processing)
        project_type: "exe" for executable project, "lib" for library project
    """
    
    templates = {
        "minimal": """# Minimal C++ project
package:
  name: hello_cpp
  cpp_standard: 17
  project_type: {project_type}

dependencies:
  fmt: {{}}
""",
        "web-server": """# Web server project
package:
  name: my_web_server
  cpp_standard: 17
  project_type: {project_type}

build:
  clang_format: Google

testing:
  framework: catch2

dependencies:
  crow:
    crow_enable_ssl: false
  nlohmann_json: {{}}
  spdlog:
    spdlog_header_only: true
""",
        "game": """# Game development project
package:
  name: my_game
  cpp_standard: 17
  project_type: {project_type}

build:
  clang_format: Google

testing:
  framework: none

dependencies:
  raylib:
    raylib_build_examples: false
  glm: {{}}
  entt: {{}}
  spdlog:
    spdlog_header_only: true
""",
        "cli-tool": """# Command-line tool project
package:
  name: my_cli_tool
  cpp_standard: 17
  project_type: {project_type}

build:
  clang_format: Google

testing:
  framework: doctest

dependencies:
  cli11: {{}}
  fmt: {{}}
  spdlog:
    spdlog_header_only: true
  indicators: {{}}
  tabulate: {{}}
""",
        "networking": """# Networking project
package:
  name: my_network_app
  cpp_standard: 17
  project_type: {project_type}

build:
  clang_format: Google

testing:
  framework: googletest

dependencies:
  asio: {{}}
  nlohmann_json: {{}}
  spdlog:
    spdlog_header_only: true
  xxhash: {{}}
""",
        "data-processing": """# Data processing project
package:
  name: data_processor
  cpp_standard: 20
  project_type: {project_type}

build:
  clang_format: LLVM

testing:
  framework: catch2

dependencies:
  simdjson: {{}}
  range_v3: {{}}
  taskflow: {{}}
  fmt: {{}}
  spdlog:
    spdlog_header_only: true
""",
    }
    
    if template_name not in templates:
        raise HTTPException(
            status_code=404,
            detail=f"Template '{template_name}' not found. Available: {', '.join(templates.keys())}"
        )
    
    # Format template with project_type
    content = templates[template_name].format(project_type=project_type)
    return PlainTextResponse(content=content, media_type="text/yaml")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
