#!/usr/bin/env bash
set -euo pipefail

if ! command -v apt-get >/dev/null 2>&1; then
  echo "This script supports Ubuntu/Debian with apt-get only."
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required. Please install Go 1.23+ first."
  exit 1
fi

WEBKIT_PKG="libwebkit2gtk-4.0-dev"
SOUP_PKG=""
if ! apt-cache show "${WEBKIT_PKG}" >/dev/null 2>&1; then
  WEBKIT_PKG="libwebkit2gtk-4.1-dev"
  SOUP_PKG="libsoup-3.0-dev"
fi

echo "Installing Linux dependencies: ${WEBKIT_PKG}, gtk, pkg-config..."
sudo apt-get update
sudo apt-get install -y \
  build-essential \
  pkg-config \
  libgtk-3-dev \
  libglib2.0-dev \
  libayatana-appindicator3-dev \
  "${WEBKIT_PKG}" \
  ${SOUP_PKG}

echo "Installing Wails CLI..."
go install github.com/wailsapp/wails/v2/cmd/wails@latest

echo "Done."
echo "If ~/go/bin is not in PATH, run:"
echo '  export PATH="$PATH:$(go env GOPATH)/bin"'
echo "Then verify with:"
echo "  wails doctor"
