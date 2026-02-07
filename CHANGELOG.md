# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) for releases (from v0.x onward).

## [Unreleased]

- CI workflow (build, test, lint, example build, govulncheck)
- Optional validator script (`scripts/validate.go`) for HTML comment structure
- Documentation: template reference, production checklist, example README, API docs (Phase 7–8)

## [0.x] — pre-production

- **Components:** Define with `<!-- + define "Name" -->`, props (YAML), html section.
- **Control flow:** `<for items={...} as="x">`, `<if condition={...}>`, `<elseif>`, `<else>`.
- **Slots:** `<slot name="..."/>` in layouts; `<Layout><slot name="...">...</slot></Layout>` at call sites.
- **Output:** One file per component (default) or `--single-file`; `--pkg`, `--dist`.
- **Errors:** Source-aware `TranspileError` with file/line/snippet; exit codes 0/1/2.
- **Core:** Transpiler in `pkg/transpiler`; runtime in `pkg/element`; no Fiber/fwatcher in core.

Breaking changes (if any) will be noted in release notes and here when releasing v1.0.
