# Contributing to GoHTMLX

Thanks for your interest in contributing. This document explains how to get set up, run checks, and submit changes.

## Prerequisites

- **Go 1.23+** (see [go.mod](../go.mod))
- Clone the repo: `git clone https://github.com/abdheshnayak/gohtmlx.git && cd gohtmlx`

## Build and test

From the repository root:

```bash
# Build everything (including CLI and packages)
go build -v ./...

# Run all tests
go test -v -count=1 ./...
```

Optional: use [Task](https://taskfile.dev/) for `task build` and `task run` (transpile example).

## Updating golden files

The golden test compares transpiler output to checked-in files under `testdata/golden/want/`. If you change the transpiler so that the generated code intentionally differs, update the golden files:

```bash
GOHTMLX_UPDATE_GOLDEN=1 go test -v -count=1 ./...
```

Then commit the updated golden files (`testdata/golden/want/` for single-file mode, `testdata/golden/want_per_component/` for one-file-per-component mode). Do **not** set `GOHTMLX_UPDATE_GOLDEN=1` in CI; use it only locally when you mean to change the expected output.

## Lint and validator

- **golangci-lint** (same as CI):
  ```bash
  golangci-lint run
  ```
  If you don’t have it: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

- **HTML comment validator** (checks `<!-- + define -->` / `<!-- + end -->` and section structure):
  ```bash
  go run scripts/validate.go --src=examples/showcase/src
  ```
  Or with the built CLI: `go build -o gohtmlx . && ./gohtmlx validate --src=examples/showcase/src`

## Example build

Ensure the example app still transpiles and builds:

```bash
go run . --src=examples/showcase/src --dist=examples/showcase/dist
cd example && go build -o /dev/null .
```

## Pull request expectations

1. **Tests** — New behavior should be covered by unit or integration tests. Existing tests must pass.
2. **Docs** — User-facing changes (CLI flags, template syntax, behavior) should be reflected in the README, [docs/TEMPLATE_REFERENCE.md](docs/TEMPLATE_REFERENCE.md), or other docs as appropriate.
3. **Changelog** — Add a note under `[Unreleased]` in [CHANGELOG.md](CHANGELOG.md) for any change that affects users (new features, breaking changes, bug fixes).
4. **CI** — The branch should pass CI (build, test, lint, example build, HTML validator, govulncheck). Run the same commands locally if you can.

## Design and context

- **Template syntax and behavior:** [docs/TEMPLATE_REFERENCE.md](docs/TEMPLATE_REFERENCE.md)
- **Production and stability plan:** [docs/PLAN_PRODUCTION_GRADE.md](docs/PLAN_PRODUCTION_GRADE.md)
- **Stability and community checklist:** [docs/STABILITY_AND_COMMUNITY_READINESS.md](docs/STABILITY_AND_COMMUNITY_READINESS.md)

If you’re proposing a larger change (e.g. new template syntax or CLI flags), opening an issue first to discuss the design is appreciated.
