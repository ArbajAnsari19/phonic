#!/bin/bash

# Development status checker for Phonic AI Calling Agent

echo "🎵 Phonic AI Calling Agent - Development Status"
echo "=============================================="
echo

echo "📁 Project Structure:"
find . -name "README.md" | sort | sed 's|./||' | sed 's|/README.md|| ' | head -10

echo
echo "🔧 Go Environment:"
echo "Go Version: $(go version)"
echo "Module Path: $(grep 'module ' go.mod)"

echo
echo "🛠️ Tools Available:"
command -v protoc >/dev/null && echo "✅ protoc: $(protoc --version | head -1)" || echo "❌ protoc: not found"
command -v protoc-gen-go >/dev/null && echo "✅ protoc-gen-go: available" || echo "❌ protoc-gen-go: not found"  
command -v protoc-gen-go-grpc >/dev/null && echo "✅ protoc-gen-go-grpc: available" || echo "❌ protoc-gen-go-grpc: not found"
command -v docker >/dev/null && echo "✅ docker: $(docker --version)" || echo "❌ docker: not found"
command -v make >/dev/null && echo "✅ make: available" || echo "❌ make: not found"

echo
echo "🎯 Next Steps:"
echo "1. Initialize Git repository (Step 4)"
echo "2. Create Makefile (Step 5)" 
echo "3. Set up Docker environment (Step 6)"
echo
