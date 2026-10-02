#!/usr/bin/env bash
# zkompiluje emulator do web/main.wasm a zkopiruje odpovidajici wasm_exec.js
set -euo pipefail
cd "$(dirname "$0")"

(cd emul && GOOS=js GOARCH=wasm go build -o ../web/main.wasm .)
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

echo "OK -> web/main.wasm"
