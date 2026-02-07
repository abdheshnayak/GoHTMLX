# GoHTMLX Example — Community showcase

A **modern showcase** of GoHTMLX: landing layout, hero, features (with code examples and syntax highlighting), HTMX demo, and feedback form. Runs as a Fiber app on port 3000.

## What it demonstrates

- **Landing layout:** Nav with CTA, hero (badge, title, lead, primary/secondary buttons), alert banner, HTMX demo (server time), Submit with HTMX form, Features section with code blocks, footer.
- **Components:** PageStyle, AppHeader, Hero, Alert, FeatureCard, Sidebar, AppFooter, Home; fragments (ServerTime, FeedbackSuccess, FeedbackErrors); AppScript for copy-code script.
- **Props and types:** NavLink (optional `IsCta`), Feature (with Code and Language for highlighting); Hero supports badge and dual CTAs.
- **Loops & conditionals:** `<for>` for features list; `<if>` for hero, alert, and optional code blocks.
- **HTMX:** Partial updates (Get server time), form submit with server-rendered success/error fragments.

## Run it

From the **repository root**:

```bash
# Transpile (generates Go under examples/showcase/dist/gohtmlxc/)
go run . --src=examples/showcase/src --dist=examples/showcase/dist

# Run the example app (Fiber server on :3000)
cd examples/showcase && go run .
```

Then open [http://localhost:3000](http://localhost:3000).

## Run with Docker

To run the example without installing Go or Node:

```bash
docker run --rm -p 3000:3000 ghcr.io/abdheshnayak/htmlx:example
```

Then open [http://localhost:3000](http://localhost:3000) in your browser.

## Development (watch)

From **examples/showcase/**:

```bash
task dev
```

This runs transpile watch (root), app watch, and CSS watch in parallel so that edits to `.html` or `.go` trigger re-transpile and app restart.

## Layout

Components are split into **one file per concern**. The transpiler walks `examples/showcase/src` and merges all `.html` files; imports are merged and deduplicated.

- **examples/showcase/src/comps/*.html** — PageStyle, AppHeader, Hero, Alert, FeatureCard, Sidebar, AppFooter (styles, header, hero, cards, sidebar, footer); **fragments.html** (ServerTime, FeedbackSuccess, FeedbackErrors); **home.html** (global imports + Home); **slots.html**, **doc-section.html**, **code-block.html** (SlotLayout, DocSection, CodeBlock — available for reuse).
- **examples/showcase/src/comps/main.go** — Builds Home props and calls generated `gc.Home{...}.Get()`; helpers for HTMX fragments.
- **examples/showcase/src/types/** — NavLink, Feature, and other prop types.
- **examples/showcase/dist/gohtmlxc/** — Generated Go (one file per component). Do not edit by hand.

## Build

```bash
cd examples/showcase && go build -o bin/app .
./bin/app
```

The app serves the showcase page at `/` and uses the generated package to render HTML.
