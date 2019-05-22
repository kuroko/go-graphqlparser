#!/usr/bin/env bash

typeNames=(
  Error
)

types=
for i in ${!typeNames[@]}; do
  types+=${typeNames[$i]}

  if [ $(expr ${i} + 1) -lt ${#typeNames[@]} ]; then
    types+=,
  fi
done

go run tools/listgen/main.go -package types -types ${types} > graphql/lists.go
go fmt graphql/lists.go
