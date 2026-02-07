# GoHTMLX Example — Showcase app

This example is a **showcase** of GoHTMLX features: multiple components, props, `<for>` loops, `<if>` conditionals, and a full page layout. It runs as a Fiber app on port 3000.

## What it demonstrates

- **Components:** PageStyle, AppHeader, Hero, Alert, FeatureCard, StatCard, TestimonialCard, Sidebar, AppFooter, **SlotLayout**, **DocSection**, Home.
- **Props and types:** NavLink, Feature, Stat, Testimonial, DocSection (from `src/types`); props passed from Go into the Home component.
- **Loops:** `<for items={props.Features} as="f">` (and similar for stats, testimonials, nav links, doc sections).
- **Conditionals:** `<if condition={props.ShowHero}>` and `<if condition={props.ShowAlert}>` to toggle hero and alert.
- **Slots:** SlotLayout defines `<slot name="header"/>`, `<slot name="body"/>`, `<slot name="footer"/>`; Home fills them with `<SlotLayout><slot name="header">…</slot>…</SlotLayout>`.
- **Template features section:** A dedicated section that explains and labels each feature (components & props, loops, conditionals, slots) with live examples.
- **Coding guide:** Step-by-step teaching section with code snippets: file structure, defining a component, expressions and props, loops, conditionals, slots, and how to transpile and use the generated Go.
- **Layout:** Single-page layout with header, main content (features, stats, testimonials, template features, coding guide), and sidebar.

## Run it

From the **repository root**:

```bash
# Transpile (generates Go under example/dist/gohtmlxc/)
go run . --src=example/src/comps --dist=example/dist

# Run the example app (Fiber server on :3000)
cd example && go run .
```

Then open http://localhost:3000

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
