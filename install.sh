#!/usr/bin/env sh
set -e

REPO="https://github.com/ReRixV1/runner.git"
APP_NAME="runner"
TMP_DIR="$(mktemp -d)"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "→ Downloading $APP_NAME..."
git clone "$REPO" "$TMP_DIR"
cd "$TMP_DIR"

echo "→ Building..."
go build -o "$APP_NAME" ./cmd/runner

echo "→ Installing..."
chmod +x "$APP_NAME"
sudo mv "$APP_NAME" $INSTALL_DIR

echo "✓ Installed $APP_NAME"
