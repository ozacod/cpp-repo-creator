"""
Vercel serverless function wrapper for Forge API.
"""
import sys
import os
from pathlib import Path

# Add forge-server to Python path
forge_server_path = Path(__file__).parent.parent / "forge-server"
sys.path.insert(0, str(forge_server_path))

# Change to forge-server directory to resolve relative imports
os.chdir(str(forge_server_path))

# Import the FastAPI app
from main import app

# Vercel Python runtime expects 'handler' to be the ASGI app
handler = app

