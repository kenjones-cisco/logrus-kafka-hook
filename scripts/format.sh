#!/bin/bash

find . \( -path ./vendor \) -prune -o -name "*.go" -exec gofmt -s -w {} \;
find . \( -path ./vendor \) -prune -o -name "*.go" -exec goimports -local "${IMPORT_PATH}" -w {} \;
