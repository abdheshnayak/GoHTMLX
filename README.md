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

## Goals

gohtmlx allows developers to write reusable HTML components, which are then transpiled into Go code. The generated Go code can be utilized to render dynamic and reusable server-side components efficiently. The focus is on providing a simple way to create server-rendered HTML with a declarative and reusable approach.

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

### Options

- `--src`: Directory containing the source `.html` component files.
- `--dist`: Directory where generated Go code is written (e.g. `dist/gohtmlxc/`).
- `--single-file`: Emit one `comp_generated.go` (legacy). Default is one file per component plus `imports.go`.
- `--pkg`: Generated package name (default `gohtmlxc`). You can point `--dist` to a subpackage of your module (e.g. `internal/gen` or `example/dist`) so the generated package lives where you want.

Imports from each file are merged and deduplicated by path (same path in multiple `<!-- * define "imports" -->` blocks → single import); order is deterministic.

### Exit codes

- **0** — Success. Generated code was written to `--dist`.
- **1** — Transpilation failed (parse error, codegen error, or write error). Errors include file and line when available.
- **2** — Invalid arguments or missing required flags (e.g. missing `--src` or `--dist`).

Scripts and CI can rely on `gohtmlx --src=... --dist=... && go build ...`.

### Development and watch

For fast re-transpile during development, use the **Taskfile-based watch** (no extra dependencies in the core CLI):

- **From repo root:** `task dev` — watches `go` and `html` under the repo, re-runs full transpile then exits (restart to run again), or use `nodemon -e go,html -i example/dist --exec "go run . --src=example/src --dist=example/dist"` to loop.
- **From example:** `task dev` — runs transpile watch (root), app watch, and CSS watch in parallel so that changing `.html` triggers a re-transpile and app restart.

Each run is a **full transpile** (all `.html` files under `--src`). A future enhancement could support incremental mode (only re-parse changed files and regenerate affected components).

## How It Works

1. **Transpilation:** gohtmlx takes HTML components defined with placeholders and transpiles them into valid Go code.
3. **Dynamic Rendering:** The resulting Go code produces dynamic HTML structures, leveraging Go's capabilities for server-side rendering and component-based architecture.

## Benefits

- **Declarative Syntax:** Write HTML-like structures in a readable and reusable manner.
- **Component Reusability:** Define and reuse server-side components efficiently.
- **Seamless Integration:** Combines Go’s performance and HTML's clarity.
- **Dynamic HTML:** Simplifies the creation of dynamic server-side web content.

## Future Enhancements

- **Improved Error Handling:** Provide detailed errors during transpilation.
- **Enhanced Debugging:** Add tools to visualize the transpilation process.
- **Broader Compatibility:** Extend support for additional libraries and frameworks.

## Production-grade roadmap

A detailed plan to make GoHTMLX production-ready for large-scale apps (deterministic builds, source-aware errors, one-file-per-component, conditionals, slots, testing, docs, CI) is in **[docs/PLAN_PRODUCTION_GRADE.md](docs/PLAN_PRODUCTION_GRADE.md)**.

---

gohtmlx bridges the gap between Go and reusable HTML components, providing developers with an intuitive way to build modern, server-rendered web applications using Go. The examples and usage reflect its ability to simplify server-side HTML generation for projects requiring basic and efficient rendering.
