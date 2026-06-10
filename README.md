# shdoc-ng-lint

`shdoc-ng-lint` checks shell script documentation parsed by
[`shdoc-ng`](https://github.com/jdevera/shdoc-ng).

It uses `shdoc-ng` v0.9.0 to parse shell script documentation and then applies
lint rules to the resulting documents.

## Installation

Install the latest version with Go:

```bash
go install github.com/drachenfels-de/shdoc-ng-lint/cmd/shdoc-ng-lint@latest
```

Install a specific version:

```bash
go install github.com/drachenfels-de/shdoc-ng-lint/cmd/shdoc-ng-lint@v0.1.0
```

Prebuilt release archives are published on the
[GitHub Releases](https://github.com/drachenfels-de/shdoc-ng-lint/releases) page.

Docker images are published to GitHub Container Registry:

```bash
docker pull ghcr.io/drachenfels-de/shdoc-ng-lint:latest
```

## Rules

- every function must have a non-empty `@description`
- every function must have a non-empty `@example`
- every function must have either `@noargs` or at least one non-zero `@arg`
- every function must have at least one `@exitcode`
- every function must start with a configured prefix when `--function-prefix` is set

## Usage

```bash
shdoc-ng-lint --function-prefix doc_ script.sh
cat script.sh | shdoc-ng-lint --include-undocumented
```
