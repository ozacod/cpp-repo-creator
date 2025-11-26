"""
FastAPI backend for C++ Project Creator.
"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import Response
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import re

from recipe_loader import (
    get_all_libraries,
    get_library_by_id,
    get_libraries_by_category,
    search_libraries,
    get_categories,
    reload_recipes,
)
from generator import create_project_zip, generate_cmake_lists


app = FastAPI(
    title="C++ Project Creator API",
    description="API for generating C++ project templates with CMake and FetchContent",
    version="2.0.0",
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
        "message": "C++ Project Creator API",
        "version": "2.0.0",
        "docs": "/docs",
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


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
