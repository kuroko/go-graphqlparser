#!/usr/bin/env bash

go run lab/walkergen/main.go -package ast > ast/walker_events.go
go fmt ast/walker_events.go