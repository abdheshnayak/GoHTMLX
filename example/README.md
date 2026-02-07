# GoHTMLX Example — Community showcase

A **modern, community-oriented showcase** of GoHTMLX: landing-style layout, hero with CTAs, quick start, features, template syntax, and a full coding guide. Built to attract users and contributors. Runs as a Fiber app on port 3000.

## What it demonstrates

- **Landing layout:** Nav with CTA, hero (badge + title + lead + primary/secondary buttons), quick start block, feature grid, stats, template syntax with slots demo, coding guide, testimonials, footer with GitHub/Docs/Contributing.
- **Components:** PageStyle, AppHeader, Hero (with optional badge and CTAs), Alert, FeatureCard, StatCard, TestimonialCard, Sidebar, AppFooter, SlotLayout, DocSection, CodeBlock, Home.
- **Props and types:** NavLink (with optional `IsCta` for nav styling), Feature, Stat, Testimonial, DocSection, CodeExample; Hero supports badge and dual CTAs.
- **Loops & conditionals:** `<for>`, `<if>`/`<elseif>`/`<else>` used throughout (features, doc sections, code examples, nav).
- **Slots:** SlotLayout with header/body/footer slots; filled at call site in the Template syntax section.
- **Coding guide:** Step-by-step snippets (HTML-escaped so tags display as text): file structure, define component, expressions, loops, conditionals, slots, transpile & use in Go.

## Run it

From the **repository root**:

```bash
# Transpile (generates Go under example/dist/gohtmlxc/)
go run . --src=example/src --dist=example/dist

# Run the example app (Fiber server on :3000)
cd example && go run .
```

Then open [http://localhost:3000](http://localhost:3000).

## Run with Docker

To run the example without installing Go or Node:

```bash
docker run --rm -p 3000:3000 ghcr.io/abdheshnayak/htmlx:example
```

Then open [http://localhost:3000](http://localhost:3000) in your browser.

## Development (watch)

From **example/**:

```bash
task dev
```

This runs transpile watch (root), app watch, and CSS watch in parallel so that edits to `.html` or `.go` trigger re-transpile and app restart.

## Layout

Components are split into **one file per concern** (like frontend frameworks). The transpiler walks `example/src` and merges all `.html` files; imports from any file are merged and deduplicated.

- **example/src/comps/*.html** — One or a few components per file:
  - **styles.html** — PageStyle (global CSS)
  - **header.html** — AppHeader
  - **hero.html** — Hero, Alert
  - **cards.html** — FeatureCard, StatCard, TestimonialCard
  - **sidebar.html** — Sidebar
  - **footer.html** — AppFooter
  - **slots.html** — SlotLayout
  - **doc-section.html** — DocSection
  - **home.html** — Global imports + Home page component
  - **code-block.html** — CodeBlock (for Coding guide snippets)
- **example/src/comps/main.go** — Go code that builds the Home props and calls the generated `gc.Home{...}.Get()`.
- **example/src/types/** — Go types used in props (NavLink, Feature, Stat, Testimonial, DocSection).
- **example/dist/gohtmlxc/** — Generated Go files (one per component by default). Do not edit by hand.

## Build

```bash
cd example && go build -o bin/app .
./bin/app
```

The app serves the showcase page at `/` and uses the generated package to render HTML.
