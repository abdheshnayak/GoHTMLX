# GoHTMLX Example — Showcase app

This example is a **showcase** of GoHTMLX features: multiple components, props, `<for>` loops, `<if>` conditionals, and a full page layout. It runs as a Fiber app on port 3000.

## What it demonstrates

- **Components:** PageStyle, AppHeader, Hero, Alert, FeatureCard, StatCard, TestimonialCard, Sidebar, AppFooter, Home.
- **Props and types:** NavLink, Feature, Stat, Testimonial (from `src/types`); props passed from Go into the Home component.
- **Loops:** `<for items={props.Features} as="f">` (and similar for stats, testimonials, nav links).
- **Conditionals:** `<if condition={props.ShowHero}>` and `<if condition={props.ShowAlert}>` to toggle hero and alert.
- **Layout:** Single-page layout with header, main content (features, stats, testimonials), and sidebar.

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

- **example/src/comps/main.html** — All component definitions (comment-based `<!-- + define "Name" -->` and `<!-- | define "props" -->` / `<!-- | define "html" -->`).
- **example/src/comps/main.go** — Go code that builds the Home props (nav links, features, stats, testimonials, hero/alert flags) and calls the generated `gc.Home{...}.Get()`.
- **example/src/types/** — Go types used in props (NavLink, Feature, Stat, Testimonial).
- **example/dist/gohtmlxc/** — Generated Go files (one per component by default). Do not edit by hand.

## Build

```bash
cd example && go build -o bin/app .
./bin/app
```

The app serves the showcase page at `/` and uses the generated package to render HTML.
