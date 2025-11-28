#!/bin/bash

echo "🎵 Spotify Playlist Sorter - Frontend"
echo "====================================="
echo ""

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
    echo ""
fi

# Check if backend is running
echo "🔍 Checking if backend is running..."
if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✅ Backend is running on port 8080"
else
    echo "⚠️  Warning: Backend does not appear to be running on port 8080"
    echo "   Please start the backend server first"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo "🚀 Starting development server..."
echo "   Frontend will be available at: http://localhost:3000"
echo "   API requests will be proxied to: http://localhost:8080/api"
echo ""

npm run dev
