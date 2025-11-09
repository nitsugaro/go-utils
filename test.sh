#!/usr/bin/env bash
set -e

PKG="./..."
COVER_FILE="cover.out"
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RESET='\033[0m'


log() { echo -e "${CYAN}➤ $1${RESET}"; }


log "Basic testing..."
go test $PKG -v

log "Verifying race conditions..."
go test $PKG -race -v

log "General Report of Global configuration..."
go test -coverpkg=./... -coverprofile=$COVER_FILE ./test -v
go tool cover -func=$COVER_FILE | tail -n 1

log "Benchmarks..."
go test ./test -bench . -benchmem -run ^$


log "Generating HTML report..."
go tool cover -html=$COVER_FILE -o cover.html
echo -e "${GREEN}✅ Archivo generado:${RESET} cover.html"

echo -e "${GREEN}✅ Todo finalizado correctamente.${RESET}"

log "Exeuting benchmark TreeMap..."
go test -benchmem -bench=Benchmark.*Overhead ./test