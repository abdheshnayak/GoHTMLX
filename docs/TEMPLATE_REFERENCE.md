# GoHTMLX — Template reference

This document describes the template syntax: how to define components, props, HTML, and use control flow (for, if, slots).

---

## File and section structure

- **Source:** One or more `.html` files under the directory passed to `--src`. Each file can define multiple components.
- **Global imports:** Optional block `<!-- * define "imports" --> ... <!-- * end -->` at the top of a file. Imports from all files are merged and deduplicated.
- **Components:** Each component is wrapped in `<!-- + define "ComponentName" --> ... <!-- + end -->`. Use `---` between components for readability.

---

## Defining a component

```html
<!-- + define "MyCard" -->
<!-- | define "props" -->
title: string
count: int
<!-- | end -->
<!-- | define "html" -->
<div class="card">
  <h2>{props.Title}</h2>
  <p>Count: {props.Count}</p>
</div>
<!-- | end -->
<!-- + end -->
```

- **`<!-- + define "Name" -->`** — Starts a component named `Name`. The name must be unique across all files.
- **`<!-- | define "props" -->`** — Optional. YAML list of prop names and Go types (e.g. `title: string`, `items: "[]mypkg.Item"`). Props become struct fields (e.g. `Title`, `Items`) and are available as `props.Title`, `props.Items` in the HTML.
- **`<!-- | define "html" -->`** — Required. The component’s HTML template. Standard HTML tags and custom component tags are allowed; see below.

Delimiters: `<!-- * ... -->` for global imports, `<!-- + ... -->` for component boundaries, `<!-- | ... -->` for section blocks inside a component.

---

## Props and expressions

- **In HTML:** Use `{props.PropName}` for a single expression (e.g. `{props.Title}`). The first letter of the prop name is capitalized in the generated struct.
- **Multiple expressions in one text:** `{props.Author} — {props.Role}` is supported; each `{...}` is emitted as a separate expression (comma-separated in generated code).
- **In attributes:** `attr={props.Value}` or `class={props.ClassName}`. The value is a Go expression.
- **Types:** Use Go type names in the props block. For slice or external types use a string, e.g. `items: "[]pkg.Item"` or `item: "mypkg.Type"`. The generated struct will reference those types; ensure the package is imported via the global imports block.

---

## HTML: standard elements and components

- **Standard HTML elements** (e.g. `div`, `span`, `a`, `form`) are transpiled to `E(\`tag\`, Attrs{...}, children...)`. Attributes become `Attrs{ \`key\`: value }`.
- **Custom components** are tags whose name matches a defined component (case-insensitive in the parser). Use `<ComponentName prop={value}>` or `<ComponentName></ComponentName>`. Children are passed as the trailing arguments to the component function.
- **Slots:** See “Slots” below.

---

## Loops: `<for>`

```html
<for items={props.Items} as="item">
  <li>{item.Name}</li>
</for>
```

- **`items`** — A Go expression (e.g. `props.Items`, `props.Links`) that is rangeable (slice, map, etc.).
- **`as`** — Loop variable name (e.g. `item`, `link`). Use this name inside the loop body.
- The body is emitted inside a `for _, as := range items { ... }` in the generated code.

---

## Conditionals: `<if>`, `<elseif>`, `<else>`

```html
<if condition={props.Show}>
  <p>Visible</p>
</if>
<elseif condition={props.Alt}>
  <p>Alternative</p>
</elseif>
<else>
  <p>Default</p>
</else>
```

- **`condition`** — A boolean Go expression (e.g. `props.Show`, `len(props.Items) > 0`). Required on `<if>` and `<elseif>`.
- **`<else>`** — Optional; no attribute. If no branch matches and there is no `<else>`, nothing is rendered for that block.
- Adjacent `<elseif>` and `<else>` must immediately follow the opening `<if>` (no other elements in between).

---

## Slots (layout placeholders)

**In a layout component:** Define a placeholder with `<slot name="..."/>` (or `<slot name="..."></slot>`).

```html
<!-- + define "PageLayout" -->
<!-- | define "html" -->
<div class="layout">
  <header><slot name="header"/></header>
  <main><slot name="main"/></main>
  <footer><slot name="footer"/></footer>
</div>
<!-- | end -->
<!-- + end -->
```

The transpiler discovers slot names from the HTML and adds fields like `SlotHeader`, `SlotMain`, `SlotFooter` (type `Element`) to the component struct. The layout emits `R(props.SlotHeader)` etc.

**At the call site:** Pass content into slots with `<slot name="...">...</slot>` as direct children of the layout component.

```html
<PageLayout>
  <slot name="header"><AppHeader/></slot>
  <slot name="main"><p>Main content</p></slot>
  <slot name="footer"><AppFooter/></slot>
</PageLayout>
```

Each slot’s content is rendered and passed as the corresponding slot prop; any other children are passed as the component’s default children.

---

## Attributes and special props

- **`Attrs`:** Every component struct includes an `Attrs Attrs` field. The runtime can use it for extra attributes. In HTML you can pass attributes on the component tag; if they are not listed in the component’s props, they go into `Attrs` (e.g. `id`, `class` when not declared as props).
- **Literal attributes:** `class="foo"` → `\`class\`: \`foo\``. **Expression attributes:** `class={props.Class}` → `\`class\`: props.Class`.

---

## Recommended style

- One component per logical block; separate components with `---`.
- Keep prop names in `props` YAML in camelCase; they are capitalized in the struct (e.g. `heroTitle` → `HeroTitle`).
- Put global imports in one place per file; use the same import path only once (merged at generation time).
- For complex layouts, use slots so callers can inject header/main/footer without prop drilling.

For a full example, see the **example** app in the repository (`example/src/comps/main.html` and [example/README.md](../example/README.md)).

---

## Optional: Validator script

A small checker for comment structure is provided so you can catch unclosed blocks before running the transpiler:

```bash
go run scripts/validate.go --src=path/to/your/html/dir
```

It reports file:line for:

- Unclosed `<!-- + define -->` (missing `<!-- + end -->`), and similarly for `<!-- | define -->` / `<!-- | end -->`, `<!-- * define -->` / `<!-- * end -->`.
- Mismatched closes (e.g. `<!-- | end -->` when the last open block was `<!-- + define -->`).

Exit 0 if all files pass; 1 if any error. Optional in normal workflow; useful in CI or before committing.
