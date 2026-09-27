#!/usr/bin/env bash

CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o target/hardhatkv-linux-amd64 ./
CGO_ENABLED=0  GOOS=linux GOARCH=arm64 go build -o target/hardhatkv-linux-arm64 ./
CGO_ENABLED=0  GOOS=darwin GOARCH=amd64  go build -o target/hardhatkv-darwin-amd64 ./
CGO_ENABLED=0  GOOS=darwin GOARCH=arm64  go build -o target/hardhatkv-darwin-arm64 ./
CGO_ENABLED=0  GOOS=windows GOARCH=amd64  go build -o target/hardhatkv-windows-amd64.exe ./
