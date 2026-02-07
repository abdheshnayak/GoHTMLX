# GoHTMLX — Roadmap and v1.0

## v1.0 criteria (production-grade)

We consider the project ready for **v1.0** when it meets the following (from [PLAN_PRODUCTION_GRADE.md](PLAN_PRODUCTION_GRADE.md)):

| Criterion | Status |
|-----------|--------|
| **Stability** — Deterministic builds, non-zero exit on error, no panics on valid input | ✅ Done |
| **Debuggability** — Every transpile error has file (and line when possible) and clear message | ✅ Done |
| **Scalability** — One generated file per component, per-file imports, configurable package/path | ✅ Done |
| **Expressiveness** — Conditionals and slots for layouts | ✅ Done |
| **Maintainability** — Core decoupled from Fiber/fwatcher; tests (unit + golden + integration) in CI | ✅ Done |
| **Documentation** — README, template reference, production checklist, API docs | ✅ Done |
| **Releases** — Versioned releases, changelog, binary distribution; CI and dependency checks | ✅ Done |

The core production plan (phases 1–8) is **complete**. v1.0 will be declared when we’ve closed the remaining items we care about for a stable, community-ready release (see below).

## What’s left before v1.0 (optional but recommended)

Tracked in [STABILITY_AND_COMMUNITY_READINESS.md](STABILITY_AND_COMMUNITY_READINESS.md). Summary:

- **Community:** LICENSE ✅, CONTRIBUTING ✅, CODE_OF_CONDUCT ✅, SECURITY ✅, comparison + roadmap ✅
- **Stability:** Transpiler tests ✅; optional `--validate-types`; golden for one-file-per-component; error contract docs
- **DX:** Optional `gohtmlx --version` + release ldflags; troubleshooting/FAQ; scaling doc; extra example (e.g. net/http)
- **Later:** Optional incremental transpilation, alternative template format, multi-module layout

When the maintainers are satisfied with the above (and any remaining checklist items), we’ll release **v1.0** and note it in the [CHANGELOG](../CHANGELOG.md).

## Versioning

- **v0.x** — Pre-production; APIs and behavior may change. Breaking changes are documented in the Changelog.
- **v1.0** — Production-ready; we’ll avoid breaking changes without a major version bump and will announce deprecations before removal.

See [MIGRATION.md](MIGRATION.md) for upgrade guidance when we introduce breaking changes.
