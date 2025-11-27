"""
Vercel serverless function wrapper for Forge API.
"""
import sys
import os
from pathlib import Path

# Get paths - in Vercel, __file__ is in /var/task/api/index.py
# Project root is one level up from api/
project_root = Path(__file__).parent.parent
forge_server_path = project_root / "forge-server"
recipes_path = forge_server_path / "recipes"

# Debug: Print paths (will show in Vercel logs)
print(f"Project root: {project_root}")
print(f"Forge server path: {forge_server_path}")
print(f"Recipes path: {recipes_path}")
print(f"Recipes exists: {recipes_path.exists()}")
print(f"Current working dir: {os.getcwd()}")

# Verify recipes directory exists
if not recipes_path.exists():
    # Try alternative path (in case Vercel structure is different)
    alt_recipes = project_root / "forge-server" / "recipes"
    if not alt_recipes.exists():
        raise RuntimeError(
            f"Recipes directory not found!\n"
            f"  Tried: {recipes_path}\n"
            f"  Tried: {alt_recipes}\n"
            f"  Project root: {project_root}\n"
            f"  Files in project root: {list(project_root.iterdir()) if project_root.exists() else 'NOT FOUND'}"
        )

# Add forge-server to Python path
sys.path.insert(0, str(forge_server_path))

# Change to forge-server directory to resolve relative imports
os.chdir(str(forge_server_path))

# Import the FastAPI app
from main import app

# Vercel Python runtime expects 'handler' to be the ASGI app
handler = app

