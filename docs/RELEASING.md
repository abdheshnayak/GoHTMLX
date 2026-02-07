# Releasing

How to cut a versioned release and publish binaries.

## Prerequisites

- Write access to the repo; ability to push tags.
- [CHANGELOG.md](../CHANGELOG.md) updated with the release version and notable changes.

## Steps

1. **Update the Changelog**  
   Under `[Unreleased]`, move completed items into a new section, e.g. `## [0.1.0] - 2025-02-07`, and add a link at the bottom for the new version. Document any breaking changes clearly.

2. **Commit and push**  
   Commit the Changelog (and any version bumps), push to the default branch.

3. **Tag and push**  
   Use a tag that matches `v*` (e.g. `v0.1.0`). The release workflow runs on tag push.
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

4. **GitHub Release**  
   The [Release workflow](../.github/workflows/release.yml) runs on tag push: it builds binaries for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64), then creates a GitHub Release with those assets and auto-generated release notes. You can edit the release body on GitHub if needed.

## Installing from a release

Users can install the CLI from a specific version or from the latest tag:

```bash
# Latest (tip of default branch or latest tag, depending on proxy)
go install github.com/abdheshnayak/gohtmlx@latest

# Specific version
go install github.com/abdheshnayak/gohtmlx@v0.1.0
```

Alternatively, they can download a binary from the [Releases](https://github.com/abdheshnayak/gohtmlx/releases) page (e.g. `gohtmlx_linux_amd64`, `gohtmlx_darwin_arm64`, `gohtmlx_windows_amd64.exe`).

## Versioning

- **v0.x** — Pre-production; APIs and behavior may change. Breaking changes are documented in the Changelog.
- **v1.0** — Production-ready; we’ll announce when the project meets the criteria in [PLAN_PRODUCTION_GRADE.md](PLAN_PRODUCTION_GRADE.md).

See [MIGRATION.md](MIGRATION.md) for upgrade guidance when we introduce breaking changes.
