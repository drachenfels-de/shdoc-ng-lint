#!/bin/sh -eu
# Run shdoc-ng-lint in a container using podman.
# NOTE: Can only lint files below the current working directory.

podman run --rm \
	-v "$(pwd):/data:ro" \
	--workdir /data \
	--network=host \
	ghcr.io/drachenfels-de/shdoc-ng-lint:latest \
	"$@"
