# Migration and upgrading

When we introduce **breaking changes** (new CLI flags, syntax changes, or behavior changes), they will be documented here and in [CHANGELOG.md](../CHANGELOG.md) under the relevant release.

## How to upgrade

1. **Check release notes**  
   Read the [Changelog](../CHANGELOG.md) for the version you’re upgrading to. Look for “Breaking” or “Changed” items.

2. **Validate before/after**  
   Run the validator to ensure comment structure is valid (helps avoid silent issues after re-transpiling):
   ```bash
   gohtmlx validate --src=path/to/your/html
   ```

3. **Re-transpile**  
   Run the CLI with your usual flags; the generated package may change (imports, struct fields, or call signatures). Fix any compile errors in the code that *uses* the generated components (e.g. prop renames, new required props).

4. **Tests and CI**  
   Run your tests and any integration builds (e.g. `go build ./...`) to confirm nothing is broken.

## Pre–v1.0

Until v1.0, we may change behavior or APIs without a major version bump. We will still list breaking or notable changes in the Changelog. When we approach v1.0, we’ll summarize any migration steps in a single “Upgrading to v1.0” section.

## After v1.0

Breaking changes will be reserved for major versions (e.g. v2.0). Minor and patch releases will stay backward compatible where possible; deprecations will be announced in the Changelog before removal.
