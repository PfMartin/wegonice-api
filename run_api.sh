#!/bin/bash

set -e

BINARY_NAME=wegonice-api

ARCH=$(uname -m)
case "$ARCH" in
  aarch64)
    ./bin/$BINARY_NAME-amd64
    ;;
  x86_64)
    ./bin/$BINARY_NAME-amd64
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac