#!/usr/bin/env bash

typeNames=(
  Argument
  Definition
  Directive
  EnumValueDefinition
  FieldDefinition
  InputValueDefinition
  Location
  OperationTypeDefinition
  PathNode
  RootOperationTypeDefinition
  Selection
  Type
  VariableDefinition
)

types=
for i in ${!typeNames[@]}; do
  types+=${typeNames[$i]}

  if [ $(expr ${i} + 1) -lt ${#typeNames[@]} ]; then
    types+=,
  fi
done

go run tools/listgen/main.go -package ast -types ${types} > ast/lists.go
go fmt ast/lists.go
