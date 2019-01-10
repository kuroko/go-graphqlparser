#!/usr/bin/env bash

go run tools/walkergen/cmd/walkergen/main.go --ast-path "./ast" \
  --package "ast" \
  > ast/walker.go

go fmt ast/walker.go
