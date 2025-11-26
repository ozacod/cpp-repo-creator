"""
Recipe loader module for loading library recipes from YAML files.
"""

import os
import yaml
from pathlib import Path
from typing import Dict, List, Optional, Any, TypedDict


class LibraryOption(TypedDict, total=False):
    """Library build option definition."""
    id: str
    name: str
    description: str
    type: str  # boolean, string, choice, integer
    default: Any
    choices: List[str]
    cmake_var: str
    cmake_define: str
    affects_link: bool
    link_libraries_when_enabled: List[str]


class FetchContent(TypedDict, total=False):
    """FetchContent configuration."""
    repository: str
    tag: str
    source_subdir: str


class Library(TypedDict, total=False):
    """Library definition loaded from recipe."""
    id: str
    name: str
    description: str
    category: str
    github_url: str
    cpp_standard: int
    header_only: bool
    tags: List[str]
    alternatives: List[str]
    fetch_content: FetchContent
    link_libraries: List[str]
    options: List[LibraryOption]
    cmake_pre: str
    cmake_post: str
    system_package: bool
    find_package_name: str


class Category(TypedDict):
    """Category definition."""
    id: str
    name: str
    icon: str
    description: str


# Category definitions
CATEGORIES: List[Category] = [
    {"id": "serialization", "name": "Serialization", "icon": "📦", "description": "JSON, XML, Binary serialization"},
    {"id": "logging", "name": "Logging", "icon": "📝", "description": "Logging and diagnostics"},
    {"id": "testing", "name": "Testing", "icon": "🧪", "description": "Unit testing and mocking frameworks"},
    {"id": "networking", "name": "Networking", "icon": "🌐", "description": "HTTP, TCP/UDP, async I/O"},
    {"id": "cli", "name": "CLI", "icon": "💻", "description": "Command line argument parsing"},
    {"id": "configuration", "name": "Configuration", "icon": "⚙️", "description": "Config file parsing (YAML, TOML)"},
    {"id": "gui", "name": "GUI", "icon": "🖼️", "description": "Graphical user interfaces"},
    {"id": "formatting", "name": "Formatting", "icon": "✨", "description": "String formatting and text processing"},
    {"id": "concurrency", "name": "Concurrency", "icon": "⚡", "description": "Threading, async, lock-free structures"},
    {"id": "utility", "name": "Utility", "icon": "🔧", "description": "General utilities and helpers"},
    {"id": "database", "name": "Database", "icon": "💾", "description": "Database clients and ORMs"},
    {"id": "compression", "name": "Compression", "icon": "🗜️", "description": "Data compression libraries"},
    {"id": "math", "name": "Math", "icon": "📐", "description": "Mathematics and linear algebra"},
    {"id": "cryptography", "name": "Cryptography", "icon": "🔐", "description": "Encryption and cryptographic functions"},
]


class RecipeLoader:
    """Loads and manages library recipes from YAML files."""

    def __init__(self, recipes_dir: Optional[str] = None):
        """Initialize the recipe loader.
        
        Args:
            recipes_dir: Path to the recipes directory. Defaults to ./recipes.
        """
        if recipes_dir is None:
            recipes_dir = Path(__file__).parent / "recipes"
        self.recipes_dir = Path(recipes_dir)
        self._libraries: Dict[str, Library] = {}
        self._loaded = False

    def load_recipes(self) -> None:
        """Load all recipe files from the recipes directory."""
        if self._loaded:
            return

        if not self.recipes_dir.exists():
            raise RuntimeError(f"Recipes directory not found: {self.recipes_dir}")

        for recipe_file in self.recipes_dir.glob("*.yaml"):
            # Skip schema file
            if recipe_file.name.startswith("_"):
                continue

            try:
                library = self._load_recipe_file(recipe_file)
                if library:
                    self._libraries[library["id"]] = library
            except Exception as e:
                print(f"Warning: Failed to load recipe {recipe_file}: {e}")

        self._loaded = True

    def _load_recipe_file(self, filepath: Path) -> Optional[Library]:
        """Load a single recipe file.
        
        Args:
            filepath: Path to the YAML recipe file.
            
        Returns:
            Library dictionary or None if loading failed.
        """
        with open(filepath, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)

        if not data or "id" not in data:
            return None

        # Ensure required fields have defaults
        library: Library = {
            "id": data["id"],
            "name": data.get("name", data["id"]),
            "description": data.get("description", ""),
            "category": data.get("category", "utility"),
            "github_url": data.get("github_url", ""),
            "cpp_standard": data.get("cpp_standard", 11),
            "header_only": data.get("header_only", False),
            "tags": data.get("tags", []),
            "alternatives": data.get("alternatives", []),
            "link_libraries": data.get("link_libraries", []),
            "options": data.get("options", []),
            "system_package": data.get("system_package", False),
        }

        # Optional fields
        if "fetch_content" in data:
            library["fetch_content"] = data["fetch_content"]
        if "cmake_pre" in data:
            library["cmake_pre"] = data["cmake_pre"]
        if "cmake_post" in data:
            library["cmake_post"] = data["cmake_post"]
        if "find_package_name" in data:
            library["find_package_name"] = data["find_package_name"]

        return library

    def get_all_libraries(self) -> List[Library]:
        """Get all loaded libraries."""
        self.load_recipes()
        return list(self._libraries.values())

    def get_library_by_id(self, library_id: str) -> Optional[Library]:
        """Get a library by its ID."""
        self.load_recipes()
        return self._libraries.get(library_id)

    def get_libraries_by_category(self, category: str) -> List[Library]:
        """Get all libraries in a specific category."""
        self.load_recipes()
        return [lib for lib in self._libraries.values() if lib["category"] == category]

    def search_libraries(self, query: str) -> List[Library]:
        """Search libraries by name, description, or tags."""
        self.load_recipes()
        query = query.lower()
        results = []
        for lib in self._libraries.values():
            if (
                query in lib["name"].lower()
                or query in lib["description"].lower()
                or any(query in tag.lower() for tag in lib["tags"])
            ):
                results.append(lib)
        return results

    def reload_recipes(self) -> None:
        """Force reload all recipes."""
        self._libraries = {}
        self._loaded = False
        self.load_recipes()


# Global recipe loader instance
_recipe_loader: Optional[RecipeLoader] = None


def get_recipe_loader() -> RecipeLoader:
    """Get the global recipe loader instance."""
    global _recipe_loader
    if _recipe_loader is None:
        _recipe_loader = RecipeLoader()
    return _recipe_loader


def get_all_libraries() -> List[Library]:
    """Get all available libraries."""
    return get_recipe_loader().get_all_libraries()


def get_library_by_id(library_id: str) -> Optional[Library]:
    """Get a library by its ID."""
    return get_recipe_loader().get_library_by_id(library_id)


def get_libraries_by_category(category: str) -> List[Library]:
    """Get all libraries in a specific category."""
    return get_recipe_loader().get_libraries_by_category(category)


def search_libraries(query: str) -> List[Library]:
    """Search libraries by name, description, or tags."""
    return get_recipe_loader().search_libraries(query)


def get_categories() -> List[Category]:
    """Get all library categories."""
    return CATEGORIES


def reload_recipes() -> None:
    """Force reload all recipes."""
    get_recipe_loader().reload_recipes()

