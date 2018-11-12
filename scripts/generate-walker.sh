#!/usr/bin/env bash

go run lab/walkergen/main.go -path "" > validation/walker.go
go fmt validation/*.go
