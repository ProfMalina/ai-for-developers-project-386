#!/bin/bash

# Script to start Prism mock server for API development
# This emulates the backend API based on the OpenAPI specification

OPENAPI_SPEC="../typespec/tsp-output/schema/openapi.yaml"

if [ ! -f "$OPENAPI_SPEC" ]; then
  echo "Error: OpenAPI specification not found at $OPENAPI_SPEC"
  echo "Please run 'make openapi' first to generate the specification"
  exit 1
fi

echo "Starting Prism mock server..."
echo "API will be available at: http://localhost:4010"
echo "Press Ctrl+C to stop"

npx @stoplight/prism-cli mock -h 0.0.0.0 -p 4010 "$OPENAPI_SPEC"
