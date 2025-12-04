#!/usr/bin/env python3
import http.server
import socketserver
from pathlib import Path
import os
import pathlib

PORT = 8080
OUTPUT_DIR = Path("/output")

class Handler(http.server.SimpleHTTPRequestHandler):
    def translate_path(self, path):
        # Serve everything relative to /output
        if path == "/":
            path = "/index.html"
        return str(OUTPUT_DIR / path.lstrip("/"))

with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print(f"Serving /output on port {PORT}")
    httpd.serve_forever()

