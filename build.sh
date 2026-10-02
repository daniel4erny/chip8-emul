#!/usr/bin/env bash
# zkompiluje emulator do web/main.wasm, zkopiruje odpovidajici wasm_exec.js
# a vygeneruje seznam rom pro web (web/roms/list.json)
set -euo pipefail
cd "$(dirname "$0")"

(cd emul && GOOS=js GOARCH=wasm go build -o ../web/main.wasm .)
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

{
	printf '['
	sep=''
	for f in web/roms/*.ch8; do
		[ -e "$f" ] || continue
		name=$(basename "$f")
		name=${name//\\/\\\\}
		name=${name//\"/\\\"}
		printf '%s"%s"' "$sep" "$name"
		sep=','
	done
	printf ']\n'
} > web/roms/list.json

echo "OK -> web/main.wasm, web/roms/list.json"
