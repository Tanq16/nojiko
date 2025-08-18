#!/bin/bash
set -e

STATIC_DIR="$(dirname "$0")/cmd/static"
mkdir -p "$STATIC_DIR/css"
mkdir -p "$STATIC_DIR/js"
mkdir -p "$STATIC_DIR/fonts"

# Download Tailwind CSS
curl -sL "https://cdn.tailwindcss.com" -o "$STATIC_DIR/js/tailwind.js"

# Download Lucide Icons
curl -sL "https://unpkg.com/lucide@latest" -o "$STATIC_DIR/js/lucide.min.js"

# Download Inter font from Google Fonts
curl -sL "https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" \
    -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36" \
    -o "$STATIC_DIR/css/inter.css"

# Parse the downloaded CSS to find and download the actual font files
grep -o "https://fonts.gstatic.com/s/inter/[^)]*" "$STATIC_DIR/css/inter.css" | while read -r url; do
  filename=$(basename "$url")
  curl -sL "$url" -o "$STATIC_DIR/fonts/$filename"
done

# Update the CSS to use the local font files instead of the Google Fonts CDN URLs
sed -i.bak 's|https://fonts.gstatic.com/s/inter/v[0-9]*/|/fonts/|g' "$STATIC_DIR/css/inter.css"
rm "$STATIC_DIR/css/inter.css.bak"
