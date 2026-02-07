# GoHTMLX (HTML Components with Go)

## Overview

gohtmlx enables developers to define and render reusable HTML components using Go. This tool is designed for scenarios where basic HTML rendering is needed or for writing purely server-side components. It simplifies creating dynamic HTML by allowing developers to define components in HTML and use them in Go code. Unlike React or JSX, gohtmlx focuses on server-side rendering and is not intended for building client-side interactive applications.

## Try it now

```bash
git clone https://github.com/abdheshnayak/gohtmlx.git
cd gohtmlx
go mod tidy
go run . --src=example/src --dist=example/dist
cd example
go run .
```

### or use `task` to run it

```bash
git clone https://github.com/abdheshnayak/gohtmlx.git
cd gohtmlx
go mod tidy
cd example
task dev
```

## Concepts

- **Components** — Defined in `.html` with `<!-- + define "Name" -->` and optional `props` (YAML) and required `html` sections. Transpiled to Go structs and `NameComp(props, attrs, children...)` functions.
- **Props** — Declared as `name: type` in the props block; used in HTML as `{props.Name}`. Types can be Go built-ins or custom types (with imports).
- **Control flow** — `<for items={props.Items} as="item">...</for>` for loops; `<if condition={expr}>...</if>`, `<elseif>`, `<else>` for conditionals.
- **Slots** — Layouts use `<slot name="header"/>`; callers pass `<Layout><slot name="header">...</slot></Layout>`.
- **Output** — Default: one generated `.go` file per component under `--dist`; `--single-file` emits one `comp_generated.go`.

## Example Usage

Developers can define reusable components in HTML and use them in their Go applications. Below is an example of defining components and rendering them:

### Defining Components

```html
<!-- + define "Greet" -->
<!-- | define "props" -->
name: string
<!-- | end -->
<!-- | define "html" -->
<div>
  <p>Hello {props.Name}!</p>
</div>
<!-- | end -->
<!-- + end -->

---

<!-- + define "Welcome" -->
<!-- | define "props" -->
projectName: string
<!-- | end -->

<!-- | define "html" -->
<div>
  <p>Welcome to {props.ProjectName}!</p>
</div>
<!-- | end -->
<!-- + end -->

---

<!-- + define "GreetNWelcome" -->
<!-- | define "props" -->
name: string
projectName: string
<!-- | end -->

<!-- | define "html" -->
<div>
  <Greet name={props.Name} ></Greet>
  <Welcome projectName={props.ProjectName} ></Welcome>
</div>
<!-- | end -->
<!-- + end -->
```

### Conditional rendering

Use `<if>`, `<elseif>`, and `<else>` to show or hide blocks based on props:

```html
<if condition={props.ShowFooter}>
  <footer>...</footer>
</if>
<elseif condition={props.ShowAlt}>
  <div>Alternative</div>
</elseif>
<else>
  <div>Default</div>
</else>
```

`condition` is a boolean expression (e.g. `props.ShowFooter`, `len(props.Items) > 0`). It is transpiled to Go `if`/`else` in the generated code.

**Slots:** Layout components can define placeholders with `<slot name="header"/>` (or `name="footer"`, etc.). Callers pass content with `<Layout><slot name="header"><h1>Title</h1></slot></Layout>`. Slots become struct fields (e.g. `SlotHeader Element`) and are rendered as `R(props.SlotHeader)` in the layout.

### Using Components in Go

```go
package main

import (
    gc "github.com/abdheshnayak/gohtmlx/example/dist/gohtmlxc"
)

func main() {
    gc.GreetNWelcomeProps{
		Name:        "Developers",
		ProjectName: "GoHtmlx",
	}.Get().Render(os.Stdout)
}
```

### Rendered HTML

When executed, the rendered HTML will look as follows:

```html
<div>
    <div>
        <p>Hello Developers!</p>
    </div>
    <div>
        <p>Welcome to gohtmlx!</p>
    </div>
</div>
```

## Usage

### Installation

To install gohtmlx, you can use the following command:

```bash
go install github.com/abdheshnayak/gohtmlx@latest
```

### Transpilation

To use gohtmlx, you can run the following command:

```bash
gohtmlx --src=path/to/src --dist=path/to/dist
```

This command will transpile HTML components from the `src` directory and generate Go code in the `dist` directory.

### CLI reference

| Flag | Required | Description |
|------|----------|-------------|
| `--src` | Yes | Source directory containing `.html` component files (walked recursively). |
| `--dist` | Yes | Destination directory for generated Go code (e.g. `dist/gohtmlxc/`). |
| `--single-file` | No | Emit one `comp_generated.go` (legacy). Default: one `.go` file per component. |
| `--pkg` | No | Generated package name (default `gohtmlxc`). |
| `--validate-types` | No | After codegen, run `go build` on the generated package and fail with file/line on error. Run from module root. |
| `--incremental` | No | Skip transpilation if no `.html` under `--src` is newer than generated `.go` files; useful in watch scripts. |
| `--version` | No | Print version and exit (set at build time via ldflags in releases). |

Example: `gohtmlx --src=example/src --dist=example/dist --pkg=gohtmlxc`. Use `--validate-types` in CI to catch invalid prop types before commit.

**Validate (optional):** Check comment structure (unclosed `define`/`end`, mismatched delimiters) without transpiling:

```bash
gohtmlx validate --src=path/to/html/dir
```

Exit 0 if all `.html` files pass; 1 if any error (messages include file:line). Useful in CI or before committing.

Imports from each file are merged and deduplicated by path; order is deterministic. See **[docs/TEMPLATE_REFERENCE.md](docs/TEMPLATE_REFERENCE.md)** for define, props, for, if, slots, and attrs.

### Exit codes

- **0** — Success. Generated code was written to `--dist`.
- **1** — Transpilation failed (parse error, codegen error, or write error). Errors include file and line when available.
- **2** — Invalid arguments or missing required flags (e.g. missing `--src` or `--dist`).

Scripts and CI can rely on `gohtmlx --src=... --dist=... && go build ...`.

### Development and watch

For fast re-transpile during development, use the **Taskfile-based watch** (no extra dependencies in the core CLI):

- **From repo root:** `task dev` — watches `go` and `html` under the repo, re-runs full transpile then exits (restart to run again), or use `nodemon -e go,html -i example/dist --exec "go run . --src=example/src --dist=example/dist"` to loop.
- **From example:** `task dev` — runs transpile watch (root), app watch, and CSS watch in parallel so that changing `.html` triggers a re-transpile and app restart.

Use **`--incremental`** in watch scripts: the CLI skips work when no `.html` file is newer than the generated `.go` files. Each run without `--incremental` is a full transpile.

### Scaling / large apps

For many components: use the default **one file per component**, point `--dist` at a dedicated package (e.g. `internal/gen`), run **`gohtmlx validate --src=...`** in CI before transpile, and optionally **`--validate-types`** to catch invalid prop types early. See **[docs/SCALING.md](docs/SCALING.md)** for layout, organizing by domain, and CI tips.

## How It Works

1. **Transpilation:** gohtmlx takes HTML components defined with placeholders and transpiles them into valid Go code.
3. **Dynamic Rendering:** The resulting Go code produces dynamic HTML structures, leveraging Go's capabilities for server-side rendering and component-based architecture.

## Benefits

- **Declarative Syntax:** Write HTML-like structures in a readable and reusable manner.
- **Component Reusability:** Define and reuse server-side components efficiently.
- **Seamless Integration:** Combines Go’s performance and HTML's clarity.
- **Dynamic HTML:** Simplifies the creation of dynamic server-side web content.

## Why GoHTMLX? (comparison)

GoHTMLX is **best for server-rendered HTML with a component model and minimal dependencies**. How it compares to other options:

| | GoHTMLX | [templ](https://templ.host) | [go-app](https://go-app.dev) | [Jet](https://github.com/CloudyKit/jet) |
|--|--------|-----------------------------|------------------------------|----------------------------------------|
| **Template format** | HTML + comment blocks | Go-like `.templ` syntax | Go structs + components | HTML with Jet expressions |
| **Output** | Generated Go (one file per component or single file) | Generated Go | Go + optional WASM | Runtime template execution |
| **Runtime** | Pure Go (no JS, no WASM) | Pure Go | Can target WASM for interactivity | Pure Go |
| **Framework** | Any (net/http, Fiber, Echo, etc.) | Any | go-app runtime | Any |
| **Focus** | Server-side components only; HTML as source of truth | Type-safe server components + optional HTMX | Full-stack (server + client) | Server-side templates |

Use **GoHTMLX** when you want to author components in **HTML**, keep the **server-only** model (no client JS framework), and get **generated Go** that fits into any HTTP stack. Use templ for type-safe Go-native component syntax; go-app for full-stack with client interactivity; Jet for runtime template rendering with a different expression language.

## Documentation

- **[Template reference](docs/TEMPLATE_REFERENCE.md)** — define, props, html, for, if, slots, attrs.
- **[Production checklist](docs/PRODUCTION_CHECKLIST.md)** — deterministic build, exit codes, one-file-per-component, CI, security.
- **[Scaling / large apps](docs/SCALING.md)** — output layout, organizing components, and CI for large codebases.
- **[Troubleshooting / FAQ](docs/TROUBLESHOOTING.md)** — generated code won’t compile, error line numbers, using with your framework.
- **[Migration / upgrading](docs/MIGRATION.md)** — how to upgrade when we introduce breaking changes.
- **[Releasing](docs/RELEASING.md)** — how to tag versions and publish binaries (maintainers).
- **[Production-grade plan](docs/PLAN_PRODUCTION_GRADE.md)** — full roadmap (phases 1–8: determinism, errors, scaling, template language, testing, docs, CI).
- **[Roadmap & v1.0](docs/ROADMAP.md)** — v1.0 criteria, what’s done, and what’s left before production release.
- **[Stability & community readiness](docs/STABILITY_AND_COMMUNITY_READINESS.md)** — steps for large-scale apps and wide community adoption.
- **[Example README](example/README.md)** — showcase app (components, for, if, layout) and how to run it.
- **[examples/minimal](examples/minimal/README.md)** — one component, no framework; renders HTML to stdout.
- **[examples/nethttp](examples/nethttp/README.md)** — same component with `net/http`; framework-agnostic server.

**Optional:** Run `go run scripts/validate.go --src=DIR` to check .html files for unclosed or mismatched comment blocks (see template reference).

**Releases:** See [CHANGELOG.md](CHANGELOG.md) for notable changes and versioning (v0.x pre-production).

### Dependencies

The **core** (CLI, `pkg/transpiler`, `pkg/element`, `pkg/gocode`) is framework-agnostic and does not import Fiber or any HTTP stack. The repository’s `go.mod` includes Fiber because the **example** app and **`pkg/integration/fiber`** use it. If you only use the transpiler or generated code with `net/http` (or another framework), you do not need Fiber at runtime; the core remains minimal.

## License

GoHTMLX is open source under the [MIT License](LICENSE). Contributions are welcome—see [CONTRIBUTING.md](CONTRIBUTING.md). We follow the [Go Community Code of Conduct](https://go.dev/conduct/) ([CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)). For security issues, see [SECURITY.md](SECURITY.md).

---

gohtmlx bridges the gap between Go and reusable HTML components, providing developers with an intuitive way to build modern, server-rendered web applications using Go. The examples and usage reflect its ability to simplify server-side HTML generation for projects requiring basic and efficient rendering.
