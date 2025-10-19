#!/bin/bash

# SVG to PNG conversion script
# Requires ImageMagick: brew install imagemagick (on macOS)

echo "Converting SVG icons to PNG format..."

# Convert home icons
convert home.svg -resize 81x81 home.png
convert home-active.svg -resize 81x81 home-active.png

# Convert history icons  
convert history.svg -resize 81x81 history.png
convert history-active.svg -resize 81x81 history-active.png

# Convert numbers icons
convert numbers.svg -resize 81x81 numbers.png
convert numbers-active.svg -resize 81x81 numbers-active.png

echo "All icons converted successfully!"
echo "Generated files:"
ls -la *.png