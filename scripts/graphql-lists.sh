#!/usr/bin/env bash

typeNames=(
  Error
  Location
  PathNode
)

types=
for i in ${!typeNames[@]}; do
  types+=${typeNames[$i]}

  if [ $(expr ${i} + 1) -lt ${#typeNames[@]} ]; then
    types+=,
  fi
done

go run lab/generics/main.go -package graphql -types ${types} > graphql/lists.go
go fmt graphql/lists.go
