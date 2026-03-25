#!/usr/bin/env bash
set -e

APP_NAME="nrg2iso"

platforms=(
    "linux amd64"
    "linux arm64"
    "windows amd64"
    "windows arm64"
    "darwin amd64"
    "darwin arm64"
)

for p in "${platforms[@]}"; do
    read -r GOOS GOARCH <<<"$p"

    outfile="${APP_NAME}"
    [[ "$GOOS" == "windows" ]] && outfile="${outfile}_win" || outfile="${outfile}_${GOOS}"
    [[ "$GOARCH" == "amd64" ]] && outfile="${outfile}_x64" || outfile="${outfile}_${GOARCH}"
    [[ "$GOOS" == "windows" ]] && outfile="${outfile}.exe"

    echo "→ Building $outfile"
    GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -ldflags="-w -s" -o "$outfile" .
done

echo "✔ Done"
