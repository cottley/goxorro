#!/bin/bash

# Build script for goxorro - creates a self-contained binary

echo "Building goxorro binary..."

# Build for current platform
go build -ldflags="-s -w" -o goxorro .

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: goxorro"
    echo "File size: $(du -h goxorro | cut -f1)"
else
    echo "Build failed!"
    exit 1
fi