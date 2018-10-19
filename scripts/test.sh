#!/usr/bin/env bash

for t in $(go list ./... | grep -v -E 'lab'); do
  go test -cover ${t}
done