# GoHTMLX — Production-Grade Plan

This document outlines the work required to make GoHTMLX production-grade for large-scale applications and to reduce developer complexity. Work is organized in phases with clear deliverables and acceptance criteria.

---

## Phase 1: Foundation (Stability, Errors, Determinism)

**Goal:** Reliable builds, debuggable errors, and reproducible output.

### 1.1 Deterministic code generation ✅

**Problem:** Component order in generated code depends on map iteration; diffs and CI are noisy.

**Deliverables:**
- [x] Sort component names (e.g. alphabetically, or by source file path) before generating.
- [x] Sort struct fields and imports before writing.
- [x] Add a test that runs transpiler twice and asserts byte-identical output.

**Acceptance:** `gohtmlx --src=X --dist=Y` produces the same file contents on repeated runs.

---

### 1.2 Source tracking and error reporting ✅

**Problem:** Transpilation errors don’t point to file:line in source HTML; hard to fix at scale.

**Deliverables:**
- [x] During HTML collection, record `(filepath, content)` and pass file path through the pipeline.
- [x] When parsing sections, associate each component name with its source file (and optional line/offset).
- [x] Define an error type (e.g. `TranspileError`) with: component name, file path, line (or offset), message, optional snippet.
- [x] Use it for all parse/transpile failures (unknown component, invalid props, invalid HTML, etc.).
- [x] Log/print errors with format: `file:line: message` (and snippet when available).
- [x] Exit with non-zero code on any error (already partially there; ensure all paths set exit code).

**Acceptance:** Any parse/transpile error shows which file and (where possible) which line; exit code is non-zero.

---

### 1.3 Exit code and CLI contract ✅

**Problem:** Callers (scripts, CI) need a clear success/failure contract.

**Deliverables:**
- [x] On success: exit 0.
- [x] On any error (flag validation, read, parse, codegen, write): exit 1 (or 2 for usage errors if you want to distinguish).
- [x] Document in README and `--help` (e.g. “Exits 0 on success, 1 on error”).

**Acceptance:** CI can rely on `gohtmlx ... && go build ...`.

---

## Phase 2: Scalable output and incremental builds

**Goal:** Large apps don’t hit a single huge file; builds stay fast and reviewable.

### 2.1 One generated file per component (or per source file) ✅

**Problem:** One `comp_generated.go` for all components doesn’t scale (merge conflicts, compile time).

**Deliverables:**
- [x] Add a mode (e.g. `--one-file-per-component` or make it default and keep `--single-file` for backward compatibility):
  - Emit `{ComponentName}.go` (or `{filename_stem}_gen.go` per source file) under a single output package directory.
- [x] Shared pieces: one small file for shared types/imports (or embed in first file), or a `gohtmlx_runtime.go` that stays minimal.
- [x] Ensure generated package name is configurable (e.g. `--pkg=gohtmlxc` or from path).
- [x] Preserve deterministic ordering (e.g. sort files and contents) so diffs are stable.

**Acceptance:** With 50+ components, output is multiple files; `go build` still works; order is deterministic.

---

### 2.2 Optional incremental / watch mode in CLI ✅

**Problem:** Developers want fast feedback without re-transpiling everything.

**Deliverables:**
- [x] Document that “incremental” can mean: only re-parse changed `.html` files and only regenerate affected components (optional Phase 2.2b); see README “Development and watch”.
- [x] Or: provide a simple `--watch` that re-runs full transpile on file change (no new deps in core; can be a separate cmd or example using fsnotify).
- [x] Keep current Taskfile-based watch as the reference; don’t couple core to fwatcher.

**Acceptance:** There is a documented way to get fast re-transpile during development (full or incremental).

---

## Phase 3: Template language and DX

**Goal:** Richer templates and fewer footguns so developers can stay in HTML-like layer.

### 3.1 Conditional rendering ✅

**Problem:** No first-class `if` in templates; everything is pushed to Go.

**Deliverables:**
- [x] Add `<if condition={expr}>...</if>` (and optional `<else>`, `<elseif condition={expr}>`).
- [x] `condition` must be a boolean expression (e.g. `props.ShowFooter`, `len(props.Items) > 0`).
- [x] Transpile to Go `if ... { ... }` in the generated code.
- [x] Document with examples; add tests (including with `for` and nesting).

**Acceptance:** Layouts can show/hide sections and list empty states without extra Go glue.

---

### 3.2 Slots / partials (layout placeholders) ✅

**Problem:** Complex layouts need named placeholders (header, sidebar, content) without prop drilling.

**Deliverables:**
- [x] Design: `<slot name="header"/>` in a layout component and `<Layout><slot name="header">...</slot></Layout>` at call site.
- [x] Transpile slots to named props (e.g. `SlotHeader Element`) in the component struct; layout emits `R(props.SlotHeader)`; callers’ `<slot name="header">...</slot>` become slot props.
- [x] Document and add tests: `SlotNamesFromHTML`, slot placeholder rendering, and caller passing slot content.

**Acceptance:** A layout component can define placeholders that callers fill by name.

---

### 3.3 More robust component syntax (optional but recommended)

**Problem:** Comment-based `<!-- + define "Name" -->` is brittle and not IDE-friendly.

**Deliverables:**
- [x] Keep current syntax supported (backward compatible).
- [x] Option A: Add a strict mode that requires well-formed comments and reports file:line on malformed blocks (validator reports file:line).
- [ ] Option B: Introduce an alternative format (e.g. one component per file with a simple frontmatter or delimiter) and document migration.
- [x] Document “recommended” style and add a lint/validate subcommand that checks for common mistakes (e.g. unclosed define, wrong delimiter).

**Acceptance:** Fewer silent failures; optional validator; path to better tooling (syntax highlighting, etc.).

---

## Phase 4: Imports, modules, and types

**Goal:** Large codebases can split by feature/domain and use proper types.

### 4.1 Per-file (or per-component) imports ✅

**Problem:** Single global “imports” block doesn’t scale; components can’t use different packages.

**Deliverables:**
- [x] Allow imports in each component/file (e.g. repeat `<!-- * define "imports" -->` per file, or a per-component imports block).
- [x] Merge and deduplicate imports when generating (same path → single import with optional alias).
- [x] Emit imports at top of generated file(s); keep deterministic order (e.g. sort by path).

**Acceptance:** Two components can use different external packages without conflict; generated code compiles.

---

### 4.2 Configurable output package and path ✅

**Problem:** Generated code is tied to a single package name and location.

**Deliverables:**
- [x] CLI flags: e.g. `--pkg=name`, `--out=path` (or derive package from `--dist` path).
- [x] Document that apps can point `--dist` to a subpackage of their module (e.g. `internal/gohtmlx/gen`).

**Acceptance:** Generated code can live in a chosen package path and package name.

---

### 4.3 Props type validation (optional)

**Problem:** Props are “name: type” strings; invalid types only surface at `go build`.

**Deliverables:**
- [ ] After generating Go code, optionally run `go/parser` (or `go build -o /dev/null`) on the generated package and map errors back to component/prop (using source mapping from 1.2).
- [ ] Flag: e.g. `--validate-types` (default off for speed; on in CI).
- [ ] Document in “Production checklist”.

**Acceptance:** CI can fail fast on invalid prop types with a clear component/prop reference.

---

## Phase 5: Decouple core and optional integrations ✅

**Goal:** Core is framework-agnostic; Fiber/watcher are optional.

### 5.1 Minimal core ✅

**Deliverables:**
- [x] Move `Run(src, dist string) error` and all parsing/codegen into a package that **does not** import Fiber or fwatcher (e.g. `pkg/transpiler` or `pkg/gohtmlx`).
- [x] `main` (cmd) only: flags, call `transpiler.Run`, exit code. No logging to a specific framework.
- [x] Keep `pkg/element` as the runtime used by generated code; no Fiber there.

**Acceptance:** `go build ./pkg/...` (excluding cmd that use Fiber) builds without Fiber/fwatcher in core.

---

### 5.2 Optional integrations ✅

**Deliverables:**
- [x] Move Fiber logger (and any Fiber-specific helpers) to an optional package, e.g. `pkg/integration/fiber` or `examples/fiber`, and use only from example/main.
- [x] Remove direct dependency on `nxtcoder17/fwatcher` from core. Use it only in example Taskfile or a separate `cmd/gohtmlx-watch` if you add one.
- [x] `go.mod`: core has no Fiber/fwatcher; example (or integration package) has them.

**Acceptance:** Core has minimal deps; examples show how to plug in Fiber and watch.

---

### 5.3 Logging interface ✅

**Deliverables:**
- [x] Define a small logger interface (e.g. `Info(msg string, kvs ...any)`, `Error(msg string, kvs ...any)`) in core.
- [x] Default: no-op or stdlog; CLI can inject a simple logger (e.g. slog).
- [x] Use this in transpiler for progress and errors (errors still go to structured TranspileError).

**Acceptance:** Core doesn’t force a logging implementation; CLI can log to stdout/slog.

---

## Phase 6: Testing and quality

**Goal:** Safe refactors and regressions caught early.

### 6.1 Unit tests for pipeline ✅

**Deliverables:**
- [x] Tests for: section parsing (comment delimiters), props YAML parsing, HTML→Go code generation for elements (E, R, for, components), struct generation.
- [x] Use table-driven tests and testdata (small .html and expected .go snippets or full output).
- [x] Tests for error cases: malformed comments, invalid YAML, unknown component in HTML, missing required props.

**Acceptance:** Key parsing and codegen paths have unit tests; running tests is part of CI.

---

### 6.2 Golden / integration tests ✅

**Deliverables:**
- [x] Golden test: run transpiler on a fixed `testdata/` tree; compare generated output to checked-in golden files (with deterministic order from Phase 1).
- [x] Integration test: transpile example (or a minimal variant), then `go build` the result and optionally run a single render to stdout.
- [x] Document how to update golden files (`-update` or env var).

**Acceptance:** Changing the transpiler breaks golden or integration test when behavior changes; team updates goldens intentionally.

**How to update golden:** run tests with `GOHTMLX_UPDATE_GOLDEN=1` to overwrite `testdata/golden/want/comp_generated.go` with current output.

---

### 6.3 Fuzz or property tests (optional) ✅

**Deliverables:**
- [x] At least one fuzz target for a critical path: FuzzProcessRaws in pkg/element (processRaws) to catch panics.
- [x] Run in CI with a short run (e.g. -fuzztime=30s) in the build-and-test job.

**Acceptance:** No panics on arbitrary input within a reasonable scope.

---

## Phase 7: Documentation and DX

**Goal:** Onboarding and daily use are straightforward.

### 7.1 User-facing docs ✅

**Deliverables:**
- [x] README: quick start, CLI reference (flags, exit codes), conceptual overview (components, props, for, if, slots).
- [x] Separate doc: “Template reference” ([docs/TEMPLATE_REFERENCE.md](TEMPLATE_REFERENCE.md)) — define, props, html, for, if, slots, attrs.
- [x] “Production checklist” ([docs/PRODUCTION_CHECKLIST.md](PRODUCTION_CHECKLIST.md)): deterministic build, exit codes, one-file-per-component, error handling, CI, security.
- [x] “Migration” or “Upgrading” when you introduce breaking changes (e.g. new CLI flags, new syntax); see [docs/MIGRATION.md](MIGRATION.md).

**Acceptance:** A new developer can run, change a component, and understand errors and flags from docs.

---

### 7.2 Inline and API docs ✅

**Deliverables:**
- [x] Package comments for `pkg/element`, `pkg/gocode`, and the transpiler package.
- [x] Exported functions and types documented (godoc): CompInfo, Html, NewHtml, SlotNamesFromHTML, Attrs, Element, R, E, ConstructStruct, ConstructSource, Run, RunOptions, TranspileError.

**Acceptance:** `go doc` and IDE hover show clear descriptions for public API.

---

### 7.3 Example and optional tooling ✅

**Deliverables:**
- [x] Keep example app working after each phase; add a “complex” example (showcase: layout, conditionals, for, many components) — see [examples/showcase/README.md](../examples/showcase/README.md).
- [x] Optional validator script: `scripts/validate.go` checks comment structure (<!-- + define --> / <!-- + end -->, etc.) and reports file:line; documented in [TEMPLATE_REFERENCE.md](TEMPLATE_REFERENCE.md) and README.

**Acceptance:** Example runs; optional tooling documented.

---

## Phase 8: CI, versioning, and security

**Goal:** Reliable releases and safe dependencies.

### 8.1 CI pipeline ✅

**Deliverables:**
- [x] CI (GitHub Actions): on push/PR run `go build ./...`, `go test ./...`, HTML validator (scripts/validate.go), golangci-lint, example transpile+build, govulncheck.
- [x] Linting: golangci-lint (via golangci-lint-action).
- [x] Example build in CI (transpile then `go build` in examples/showcase/).

**Acceptance:** PRs that break build or tests are rejected; main stays green.

---

### 8.2 Versioning and releases ✅

**Deliverables:**
- [x] Semantic versioning: v0.x for pre-production; v1.0 when “production-grade” criteria are met.
- [x] Changelog (CHANGELOG.md) for notable changes; document breaking changes when releasing.
- [x] Optional: tag releases and publish binaries (e.g. GitHub Releases) for `go install github.com/...@latest`; see [docs/RELEASING.md](RELEASING.md) and `.github/workflows/release.yml`.

**Acceptance:** Users can depend on a specific version and read what changed.

---

### 8.3 Dependencies and security ✅

**Deliverables:**
- [x] Keep `go.mod` minimal; prefer standard library and well-maintained deps (`golang.org/x/net`, `sigs.k8s.io/yaml`). Core has no Fiber/fwatcher.
- [x] Run govulncheck in CI (vulncheck job); fix high/critical issues. Use `go mod tidy` before commits.
- [x] No secrets or credentials in repo; document that generated code may import user packages (see PRODUCTION_CHECKLIST and TEMPLATE_REFERENCE).

**Acceptance:** No known high/critical vulnerabilities in direct deps; CI enforces checks.

---

## Implementation order (suggested)

| Order | Phase / item | Rationale |
|-------|----------------|-----------|
| 1 | 1.1 Deterministic output | Unblocks stable golden tests and clean diffs |
| 2 | 1.2 Source tracking & errors | Foundation for all other error reporting |
| 3 | 1.3 Exit code & CLI contract | Required for CI and scripts |
| 4 | 6.1 Unit tests + 6.2 Golden | Lock behavior before refactors |
| 5 | 5.1–5.3 Decouple core | Enables clean dependency graph and optional integrations |
| 6 | 2.1 One file per component | Scalability for large apps |
| 7 | 3.1 Conditionals | High impact on template expressiveness |
| 8 | 4.1 Per-file imports, 4.2 Configurable pkg/out | Needed for multi-package apps |
| 9 | 3.2 Slots | Improves layout reuse |
| 10 | 2.2 Watch / incremental | DX after core is stable |
| 11 | 3.3 Robust syntax / validator | Reduces footguns |
| 12 | 4.3 Type validation | Optional but valuable in CI |
| 13 | Phase 7 (docs), Phase 8 (CI, versioning, security) | Ongoing; ramp up before v1 |

---

## Definition of “production-grade” for v1.0

- **Stability:** Deterministic builds, non-zero exit on error, no panics on valid input.
- **Debuggability:** Every transpile error has file (and line when possible) and clear message.
- **Scalability:** One generated file per component (or equivalent), per-file imports, configurable package/path.
- **Expressiveness:** Conditionals and (at least) a simple slot/partial mechanism.
- **Maintainability:** Core decoupled from Fiber/fwatcher; tests (unit + golden/integration) in CI.
- **Documentation:** README, template reference, production checklist, and basic API docs.
- **Releases:** Versioned releases, changelog, and optional binary distribution; CI and dependency checks in place.

---

## Summary checklist (high level)

- [x] Phase 1: Determinism, source-aware errors, exit codes
- [x] Phase 2: One file per component, optional watch/incremental
- [x] Phase 3: Conditionals, slots, more robust syntax/validator
- [x] Phase 4: Per-file imports, configurable pkg/out, optional type validation
- [x] Phase 5: Core decoupled; optional Fiber/fwatcher; logging interface
- [x] Phase 6: Unit + golden + integration tests; optional fuzz
- [x] Phase 7: User docs, API docs, example, optional IDE/tooling
- [x] Phase 8: CI, versioning, changelog, dependency/security checks

This plan, executed in the suggested order, should make GoHTMLX production-grade and suitable for large-scale applications while reducing developer complexity.
