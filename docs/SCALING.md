# Scaling and large apps

Tips for using GoHTMLX in projects with many components or multiple teams.

## Output layout

- **Prefer one file per component** (the default). Do not use `--single-file` for large codebases; it produces one large file that is harder to review and can slow builds.
- **Write generated code to a dedicated package** — e.g. `--dist=internal/gen` or `--dist=internal/gohtmlx/gen`. Use `--pkg` if you need a specific package name. This keeps generated code separate from hand-written code and makes imports clear.

## Organizing components

- **One `--src` tree:** Put all `.html` component files under a single directory (e.g. `internal/components` or `pkg/ui/templates`) and run a single `gohtmlx --src=... --dist=...`. Component names must be unique across all files; imports are merged and deduplicated.
- **Multiple runs (optional):** If you want to split by domain, run the CLI multiple times with different `--src` and `--dist`/`--pkg` (e.g. `--src=cmd/web/components --dist=internal/gen/web --pkg=web` and `--src=cmd/admin/components --dist=internal/gen/admin --pkg=admin`). Each run produces a separate package.

## CI

- **Validate before transpile:** Run `gohtmlx validate --src=...` before `gohtmlx --src=... --dist=...` so that unclosed or malformed comment blocks fail the build early with a clear file:line message.
- **Validate types (optional):** Use `--validate-types` when running from the module root so that invalid prop types (e.g. typos or missing imports) are reported at transpile time instead of at `go build`. Helps catch mistakes before commit.

## Summary

Use one file per component, a dedicated `--dist` package, validate in CI, and optionally `--validate-types`. For very large repos, consider multiple transpile runs with different `--src`/`--dist`/`--pkg` to split by feature or service.
