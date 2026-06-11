`shdoc-ng-lint` is a Go CLI module at `github.com/drachenfels-de/shdoc-ng-lint`.

Purpose:
- lint shell script documentation parsed with `shdoc-ng` v0.9.0
- use `shdoc.ParseDocumentWithOptions` and validate the returned documents

Lint rules:
- every function must have a non-empty `@description`
- every function must have a non-empty `@example`
- every function must define `@noargs` or at least one non-`$0` `@arg`
- every function must define at least one `@exitcode`
- when `--function-prefix` is set, every function name must start with that prefix

CLI behavior:
- reads files or stdin
- supports `--include-undocumented`, `--include-internal`, `--include-all`, and `--werror`

Project expectations:
- Go code should stay `gofumpt`-formatted
- changes should pass `golangci-lint`, `go test ./...`, and `go build ./...`

Workflows present:
- `.github/workflows/ci.yml` runs formatting checks, `golangci-lint`, tests, and build
- `.github/workflows/release.yml` builds and publishes release archives on version tags
- `.github/workflows/docker.yml` builds and publishes multi-arch Docker images to GHCR
- `.gitea/workflows/ci.yml` runs formatting checks, `golangci-lint`, tests, and build on Gitea Actions
- `.gitea/workflows/release.yml` builds archives and publishes release assets on version tags in Gitea
- `.gitea/workflows/docker.yml` builds and publishes multi-arch Docker images to the Gitea container registry

Commit note requirement:
- when creating a commit or commit message/comment for this repo, always add a Codex co-author note in the commit message body
- include the current session's model, reasoning mode, and plugin version details instead of hard-coding values
- if the exact model, reasoning mode, or plugin version is not visible in the current session, ask the user before writing the commit note
