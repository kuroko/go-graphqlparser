#!/usr/bin/env bash

go run tools/walkergen/cmd/walkergen/main.go --ast-path "./ast" \
  --package "validation" \
  > validation/walker.go

go fmt validation/*.go
