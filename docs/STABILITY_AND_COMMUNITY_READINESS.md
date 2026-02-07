# GoHTMLX — Stability, Scale & Community Readiness

This document summarizes actionable steps to make GoHTMLX **stable**, **powerful** for large-scale apps, and **ready for wide community adoption**. It builds on the completed [PLAN_PRODUCTION_GRADE.md](PLAN_PRODUCTION_GRADE.md). **Pick an item, do it, then mark it [x] below.**

---

## Plan checklist (pick and mark done)

Use this list in order (or pick any item). When done with a task, change `[ ]` to `[x]`.

| # | Status | Item |
|---|--------|------|
| 1 | [x] | **Add LICENSE** — e.g. MIT or Apache-2.0; state in README (see §4.1) |
| 2 | [x] | **CONTRIBUTING.md** — clone, build, test, golden update, lint, validator, PR expectations (see §4.2) |
| 3 | [x] | **CODE_OF_CONDUCT.md + SECURITY.md** — link to Go CoC or Covenant; how to report vulnerabilities (see §4.3) |
| 4 | [x] | **Transpiler tests** — add `pkg/transpiler/run_test.go` for Run options and errors (see §2.2) |
| 5 | [x] | **Comparison + v1.0 roadmap** — “Why GoHTMLX” vs templ/go-app/Jet; ROADMAP or README section (see §4.4, §4.5) |
| 6 | [x] | **Optional `--validate-types`** — flag to run go build on generated package and map errors to component/prop (see §2.1) |
| 7 | [x] | **Golden for one-file-per-component** — test with `SingleFile: false` and checked-in goldens (see §2.3) |
| 8 | [x] | **Extra example** — minimal (stdout) or net/http to show framework-agnostic usage (see §4.6) |
| 9 | [x] | **`gohtmlx --version` + release ldflags** — version var + ldflags in release workflow (see §4.7, §6) |
| 10 | [x] | **Error contract docs** — godoc for Run/TranspileError; ensure all failures use TranspileError (see §2.4) |
| 11 | [x] | **Scaling / large apps doc** — short section in README or docs (see §3.2) |
| 12 | [x] | **Troubleshooting / FAQ** — “generated code doesn’t compile”, “wrong line”, “use with my framework” (see §5.2) |
| 13 | [x] | **Go version alignment** — go.mod and CI use same Go (e.g. 1.23); optional: test 1.22 + 1.23 (see §6) |
| 14 | [x] | **Optional incremental transpilation** — `--incremental` flag, re-parse only changed files (see §3.1) |
| 15 | [x] | **Optional: multi-module or dependency note** — separate module for core (no Fiber) or document footprint (see §3.3) |
| 16 | [x] | **Optional: alternative template format** — one component per file / frontmatter; migration path (see §5.1) |

---

## 1. Current Strengths (Already in Place)

- **Deterministic output** and source-aware `TranspileError` with file:line
- **Exit codes** 0/1/2; **one file per component** (default); **per-file imports**
- **Conditionals** (`<if>`/`<elseif>`/`<else>`) and **slots** for layouts
- **Core decoupled** from Fiber (no Fiber imports in `pkg/transpiler`, `pkg/element`, `pkg/gocode`)
- **Tests:** unit (element, gocode, utils, validate), golden, integration build, fuzz (ProcessRaws)
- **CI:** build, test, lint, example build, HTML validator, govulncheck
- **Docs:** README, template reference, production checklist, migration, releasing
- **Releases:** tag-driven workflow with multi-OS binaries

---

## 2. Stability & Robustness

### 2.1 Optional type validation (Phase 4.3) — Checklist item #6

**Goal:** Fail fast in CI on invalid prop types instead of at `go build`.

- [ ] Add a `--validate-types` flag (default off).
- [x] After codegen, run `go build -o /dev/null` on the generated package (or use `go/parser`/type-check) and map compiler errors back to component/prop using existing source mapping.
- [x] Document in production checklist and template reference.

**Impact:** Fewer “mystery” build failures in large codebases.

### 2.2 Transpiler package tests — Checklist item #4

**Goal:** Direct tests for `transpiler.Run` and options (e.g. one-file vs one-per-component, `--pkg`).

- [x] Add `pkg/transpiler/run_test.go` (or equivalent) with table-driven tests:
  - Run with `SingleFile: true` and `SingleFile: false`, assert expected files exist and content shape.
  - Run with custom `Pkg`, assert package name in generated code.
  - Run on a small testdata tree with intentional errors, assert `TranspileError` and file path.
- Keeps refactors safe and documents expected behavior.

### 2.3 Golden test for one-file-per-component mode — Checklist item #7

**Goal:** Avoid regressions when the default mode (one file per component) changes.

- [x] Add a second golden test (or extend testdata) that runs with `SingleFile: false` and compares generated files (e.g. sorted list of files + content) to checked-in goldens.
- Optional: same testdata as current golden, different output layout.

### 2.4 Error contract for programmatic use — Checklist item #10

- [x] Ensure all code paths that can fail return `*TranspileError` (or wrap with `errors.As`), so callers can distinguish validation vs I/O and show file:line in UIs or tooling.
- [x] Document in godoc that `Run` returns `*TranspileError` with `FilePath`, `Line`, `Message`, `Snippet`.

---

## 3. Large-Scale Readiness

### 3.1 Optional incremental transpilation (Phase 2.2b) — Checklist item #14

**Goal:** Only re-parse changed `.html` files and regenerate affected components.

- [x] Add `--incremental` flag: skip transpilation when no `.html` under src is newer than generated `.go` files (best-effort; reduces work in watch scripts).
- Optional: per-file mtimes/hashes and regenerate only affected components (not implemented).

### 3.2 Document patterns for very large apps — Checklist item #11

- [x] Add a short “Scaling” or “Large apps” section (in README or docs):
  - Prefer one file per component; use `--dist` into a dedicated package (e.g. `internal/gen`).
  - Split components by domain/feature (multiple `--src` dirs or multiple runs with different `--dist`/`--pkg` if you support it).
  - Run `gohtmlx validate --src=...` in CI before transpile.
  - Optional: run with `--validate-types` in CI.

### 3.3 Minimal dependency footprint for library users (optional) — Checklist item #15

- Today the root `go.mod` includes Fiber (used by example and `pkg/integration/fiber`). Core code does not import Fiber, but `go get github.com/abdheshnayak/gohtmlx/pkg/transpiler` still pulls the whole module.
- **Option A:** Move the CLI and optional integrations (e.g. Fiber) into a separate Go module (e.g. `cmd/gohtmlx` and `example` as a separate module) and keep a root module that contains only `pkg/...` with minimal deps. Then library-only users get no Fiber.
- **Option B:** Document that “core” is framework-agnostic and that Fiber is only required for the example and `pkg/integration/fiber`; accept that the single-module layout keeps the repo simpler.

---

## 4. Community Adoption

### 4.1 License and project metadata — Checklist item #1

- [x] **Add a LICENSE file** (e.g. MIT or Apache-2.0). Without it, many organizations will not adopt the project.
- [x] Ensure README clearly states license and that the project is open for contributions (if that’s the intent).

### 4.2 CONTRIBUTING.md — Checklist item #2

- [x] How to clone, build, run tests, and update golden files (`GOHTMLX_UPDATE_GOLDEN=1`).
- [x] How to run lint and validator.
- [x] PR expectations: tests and docs for new behavior; changelog note for user-facing changes.
- [x] Point to template reference and production plan for design context.

### 4.3 CODE_OF_CONDUCT and SECURITY — Checklist item #3

- [x] Add **CODE_OF_CONDUCT.md** (e.g. link to [Go Community Code of Conduct](https://go.dev/conduct/) or Contributor Covenant).
- [x] Add **SECURITY.md** with how to report vulnerabilities (e.g. private email or GitHub Security Advisories).

### 4.4 Comparison and positioning — Checklist item #5

- [x] Add a short “Comparison” or “Why GoHTMLX” section (README or docs):
  - Compare with **templ**, **go-app**, **Jet**, or other server-side Go HTML/template solutions (focus on: server-only components, no JS runtime, HTML-in-Go generation, optional framework integration).
  - Clarify: “Best for server-rendered HTML with a component model and minimal dependencies.”

### 4.5 Clear v1.0 criteria and roadmap — Checklist item #5

- [x] In README or a dedicated **ROADMAP.md**, list the v1.0 criteria (e.g. from PLAN_PRODUCTION_GRADE: stability, determinism, errors, scalability, tests, docs, CI, releases).
- [x] Mark what’s done and what’s left (e.g. optional type validation, incremental mode). This sets expectations and builds trust.

### 4.6 Additional examples — Checklist item #8

- [x] **Minimal example:** One component, no framework; just `gohtmlx` + `go run` printing HTML to stdout (already partly in README; could be a `examples/minimal` directory).
- [x] **net/http example:** Same components rendered via `net/http` instead of Fiber, to show framework-agnostic usage.
- Optional: **Echo** or **Chi** example in a separate dir or repo to show integration pattern.

### 4.7 Discoverability and release hygiene — Checklist item #9

- [x] **Releases:** Keep tagging and CHANGELOG up to date; consider adding version to the binary (e.g. `gohtmlx --version` with `-ldflags="-X main.Version=v0.1.0"` in release workflow).
- [ ] **Go doc:** Ensure `go doc` for `pkg/transpiler`, `pkg/element`, and `pkg/gocode` is clear; pkg.go.dev will surface this.
- Optional: **Go version** — support at least the last two minor versions (e.g. 1.22 and 1.23) in CI to signal compatibility.

---

## 5. Documentation and DX

### 5.1 Alternative template format (Phase 3.3 Option B) — Checklist item #16

- **Goal:** Improve IDE support and reduce comment-syntax footguns.
- [x] Consider an optional format: e.g. one component per file with a simple frontmatter or delimiter (e.g. `---` + YAML props + HTML block). Document migration path from current comment-based format.
- Documented in [TEMPLATE_REFERENCE.md](TEMPLATE_REFERENCE.md) “Future: alternative format”; not implemented in the CLI yet.

### 5.2 Troubleshooting / FAQ — Checklist item #12

- [x] Add a short **Troubleshooting** or **FAQ** (in docs or README):
  - “Generated code doesn’t compile” → check prop types and imports; mention future `--validate-types`.
  - “Error points to wrong line” → explain that line is component-start; suggest snippet.
  - “How to use with my framework?” → point to Fiber example and net/http pattern.

### 5.3 Performance (optional)

- If users ask about “huge HTML trees” or “many components,” add a note (or benchmark): transpilation is one-time; runtime is just Go function calls and `element` tree building. No need to over-engineer before demand appears.

---

## 6. CI and Release Tweaks — Checklist items #9, #13

- [x] **Go version:** Align `go.mod` and CI (e.g. use `1.23` in both, or `1.23.x` in go.mod and `1.23` in CI). Consider testing on 1.22 and 1.23 for compatibility.
- [x] **Release binary:** Add `-ldflags="-s -w -X main.Version=$TAG"` (or similar) so `gohtmlx --version` (if you add it) shows the release tag.
- **Changelog:** Keep a strict “Unreleased” vs versioned sections and mention breaking changes in release notes.

---

## Summary

- **Stability:** Add transpiler tests (#4), optional type validation (#6), and a golden for default one-file-per-component mode (#7). Keep error contract clear (#10).
- **Large-scale:** Document scaling patterns (#11); optionally add incremental transpilation (#14) and dependency footprint note or multi-module (#15).
- **Community:** Add LICENSE (#1), CONTRIBUTING (#2), CODE_OF_CONDUCT + SECURITY (#3); comparison and v1.0 roadmap (#5); extra example (#8); `--version` and ldflags (#9).
- **CI/Release:** Align Go versions (#13), add version ldflag (#9), keep changelog and release notes clear.

Pick an item from the **Plan checklist** at the top, complete it, then mark it `[x]`.
