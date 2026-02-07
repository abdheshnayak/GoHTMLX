# GoHTMLX — Production checklist

Use this checklist when deploying GoHTMLX in production or in CI.

---

## Build and transpilation

- [ ] **Deterministic builds:** Run `gohtmlx --src=... --dist=...` twice and confirm generated files are identical (no map-iteration noise in diffs). The CLI sorts component names, struct fields, and imports.
- [ ] **Exit codes:** Scripts/CI should use `gohtmlx ... && go build ...`. Exit `0` = success; `1` = transpilation failed; `2` = invalid or missing flags.
- [ ] **One file per component (default):** Prefer the default mode (no `--single-file`) so each component lives in its own `.go` file. This reduces merge conflicts and keeps builds scalable.
- [ ] **Output location:** Use `--dist` to write into a subpackage of your module (e.g. `internal/gen` or `example/dist`) and `--pkg` if you need a specific package name (default `gohtmlxc`).

---

## Error handling and debugging

- [ ] **Source-aware errors:** All transpilation errors use `TranspileError` with component name, file path, and (when available) line and snippet. Ensure logs or CLI output show `file:line: message`.
- [ ] **No silent failures:** The CLI exits with non-zero code on any parse, codegen, or write error. Do not ignore exit codes in CI.

---

## Types and validation

- [ ] **Prop types:** Props are declared as “name: type” in the template. Invalid types surface at `go build` of the generated package. Use `--validate-types` in CI to fail fast with file/line (run from module root); see [PLAN_PRODUCTION_GRADE.md](PLAN_PRODUCTION_GRADE.md) Phase 4.3.
- [ ] **Imports:** Use the global `<!-- * define "imports" -->` block for packages required by prop types (e.g. `t "yourmod/types"`). Imports are merged and deduplicated across files.

---

## Development workflow

- [ ] **Watch mode:** Use the Taskfile-based watch (e.g. `task dev` from the example or root) for re-transpile on file changes. No extra deps in the core CLI.
- [ ] **Golden / integration tests:** If you add golden tests, run them in CI. Update goldens only when you intend to change output (`GOHTMLX_UPDATE_GOLDEN=1` for the repo’s golden test).

---

## Security and dependencies

- [ ] **Minimal deps:** Core has no Fiber/fwatcher; optional integrations live in `pkg/integration` or the example. Keep `go.mod` minimal. See README “Dependencies” for framework-agnostic core and when Fiber is used.
- [ ] **No secrets:** Do not put credentials or secrets in templates or generated code. Generated code may import your packages; ensure those packages do not expose secrets.
- [ ] **Vulnerability checks:** Run `go mod tidy` and (if available) `govulncheck` or Dependabot in CI.

---

## CI and releases

- [ ] **CI:** The repo includes `.github/workflows/ci.yml`: on push/PR it runs `go build ./...`, `go test ./...`, HTML validator, golangci-lint, example build, and govulncheck. Keep main green.
- [ ] **Changelog:** See [CHANGELOG.md](../CHANGELOG.md) for notable changes. Use semantic versioning (v0.x pre-production; v1.0 when production-grade).
- [ ] **Releases:** Tag releases and document breaking changes in the changelog or GitHub Releases.

## Documentation and onboarding

- [ ] **README:** New developers should use the README for quick start, CLI flags, and exit codes.
- [ ] **Template reference:** Use [docs/TEMPLATE_REFERENCE.md](TEMPLATE_REFERENCE.md) for define, props, html, for, if, slots, and attrs.
- [ ] **Production plan:** See [docs/PLAN_PRODUCTION_GRADE.md](PLAN_PRODUCTION_GRADE.md) for the full roadmap (CI, versioning, type validation, etc.).

---

## Summary

- Use **deterministic** output and **non-zero exit on error**.
- Prefer **one file per component** and **source-aware errors**.
- In CI: **run `gohtmlx` then `go build`** and fail on non-zero exit; optionally run tests and golden updates under control.
- Keep **dependencies minimal** and **documentation** (README, template reference, this checklist) up to date.
