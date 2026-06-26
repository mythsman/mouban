#!/usr/bin/env bash
set -euo pipefail

if [ -d "web" ]; then
  pushd web >/dev/null
  if [ -f "package-lock.json" ]; then
    npm ci --no-audit --no-fund
  else
    npm install --no-audit --no-fund
  fi
  npm run build
  popd >/dev/null
fi

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main .
