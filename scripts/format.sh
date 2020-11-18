#!/bin/bash

find . \( -path ./vendor -o -path ./.mod \) -prune -o -name "*.go" -exec gofmt -s -w {} \;
find . \( -path ./vendor -o -path ./.mod \) -prune -o -name "*.go" -exec goimports -local "${IMPORT_PATH}" -w {} \;
