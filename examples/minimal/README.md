# Minimal example (stdout)

One component, no HTTP server. Renders HTML to stdout.

## Run

From the **repository root**:

```bash
# Transpile (generates dist/gohtmlxc/)
go run . --src=examples/minimal/src --dist=examples/minimal/dist

# Render to stdout
cd examples/minimal && go run .
```

Output: a single `<div class="greeting">` with "Hello, GoHTMLX!".

## Re-generate

After changing `src/hello.html`:

```bash
# From repo root
go run . --src=examples/minimal/src --dist=examples/minimal/dist
```

The generated package is checked in so `go run .` works without running the CLI first.
