#!/bin/sh -eu

mkdir -p dist
version="${GITHUB_REF_NAME:-dev}"
for target in \
	"linux amd64" \
	"linux arm64" \
	"darwin amd64" \
	"darwin arm64" \
	"windows amd64"; do
	set -- $target
	ext=""
	if [ "$1" = "windows" ]; then
		ext=".exe"
	fi
	GOOS="$1" GOARCH="$2" go build -ldflags "-s -w -X main.version=$version" -o "dist/shdoc-ng-lint-$1-$2$ext" ./cmd/shdoc-ng-lint
	asset_dir="dist/shdoc-ng-lint-${version}-$1-$2"
	mkdir -p "$asset_dir"
	cp "dist/shdoc-ng-lint-$1-$2$ext" "$asset_dir/"
	cp README.md LICENSE "$asset_dir/" 2>/dev/null || true
	if [ "$1" = "windows" ]; then
		(cd dist && zip -rq "$(basename "$asset_dir").zip" "$(basename "$asset_dir")")
	else
		tar -C dist -czf "${asset_dir}.tar.gz" "$(basename "$asset_dir")"
	fi
done
