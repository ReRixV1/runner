#!/bin/bash
set -e

echo "building..."
go build -o runner ./cmd/runner
echo "installing..."
sudo mv runner /usr/local/bin/

echo "installed" 
