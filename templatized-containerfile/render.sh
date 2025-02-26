#!/usr/bin/env bash

for src in *.yaml; do
  gomplate -d=input="$src" --file=Containerfile.template
  echo ""
done
