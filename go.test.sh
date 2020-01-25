#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./pkg/... | grep -v /pkg/app); do
    go test -v -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
