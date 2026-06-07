#!/usr/bin/env bash
#
# generate.sh - Generate ogen API client from OpenAPI spec
#
# This script generates the internal API client from the OpenAPI specification.
# The generated code is placed in internal/api/ and should not be edited manually.
#
# Usage:
#   ./generate.sh
#
# Prerequisites:
#   - Go 1.24+
#   - ogen (will be installed if missing)
#
set -euo pipefail

echo "=== Aha-Go Code Generation ==="

# Check ogen is installed
if ! command -v ogen &> /dev/null; then
    echo "Installing ogen..."
    go install github.com/ogen-go/ogen/cmd/ogen@latest
fi

# Check OpenAPI spec exists
SPEC_FILE="openapi/aha.yaml"
if [[ ! -f "$SPEC_FILE" ]]; then
    echo "Error: $SPEC_FILE not found"
    exit 1
fi

# Create target directory if needed
mkdir -p internal/api

# Generate code
echo "Generating API client from $SPEC_FILE..."
ogen --package api --target internal/api --clean "$SPEC_FILE"

# Tidy dependencies
echo "Running go mod tidy..."
go mod tidy

# Verify build
echo "Verifying build..."
go build ./...

echo ""
echo "=== Generation complete ==="
echo ""
echo "Generated files are in internal/api/"
echo "Do not edit generated files directly."
echo ""
echo "Next steps:"
echo "  1. Update wrapper code if API changed"
echo "  2. Run tests: go test ./..."
echo "  3. Commit changes"
