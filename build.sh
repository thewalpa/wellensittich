#!/bin/bash
go build -o build/wellensittich
GOOS=linux GOARCH=arm64 go build -o build/wellensittich_arm64_linux_0.1