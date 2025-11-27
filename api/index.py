"""
Vercel serverless function wrapper for Forge API.
"""
import sys
import os
from pathlib import Path

# Get paths - in Vercel, __file__ is in /var/task/api/index.py
api_dir = Path(__file__).parent
project_root = api_dir.parent
forge_server_path = project_root / "forge-server"

# Add forge-server to Python path
sys.path.insert(0, str(forge_server_path))

# Change to forge-server directory to resolve relative imports
# This ensures recipe_loader.py can find the recipes directory
os.chdir(str(forge_server_path))

# Import the FastAPI app
from main import app

# Use Mangum to adapt FastAPI ASGI app for Vercel/Lambda
try:
    from mangum import Mangum
    handler = Mangum(app, lifespan="off")
except ImportError:
    # Fallback: export app directly (might not work but worth trying)
    handler = app

