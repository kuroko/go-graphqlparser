#!/usr/bin/env bash

go run lab/walkergen/main.go -package validation > validation/walker_events.go
go fmt validation/walker_events.go
