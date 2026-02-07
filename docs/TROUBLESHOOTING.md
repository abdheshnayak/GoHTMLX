# Troubleshooting and FAQ

## Generated code doesn’t compile

**Cause:** Usually invalid prop types or missing imports in your template.

- **Prop types:** In the `<!-- | define "props" -->` block, use valid Go type names. For slices or external types use a string, e.g. `items: "[]pkg.Item"`. Ensure the package is listed in `<!-- * define "imports" -->`.
- **Catch it earlier:** Run `gohtmlx --src=... --dist=... --validate-types` from your module root. The CLI will run `go build` on the generated package and report file:line for type errors. Use this in CI to fail before commit.

See [Template reference](TEMPLATE_REFERENCE.md) for props and imports.

## Error points to the wrong line

Transpilation errors often report the **line where the component starts** (the `<!-- + define "Name" -->` line), not the exact line of the mistake inside the component.

- Use the **snippet** in the error message: it shows a few lines of context around that component so you can find the real issue (e.g. a typo in props YAML or in the `html` block).
- If the error says “component X”, open that component in the file and check props, HTML, and comment delimiters (`<!-- | end -->`, etc.).

## How do I use GoHTMLX with my framework?

Generated components implement `element.Element` and have a `Render(io.Writer) (int, error)` method. Any HTTP stack that can write to an `io.Writer` works.

- **net/http:** Pass `http.ResponseWriter` (it implements `io.Writer`). See [examples/nethttp](../examples/nethttp/README.md).
- **Fiber:** Pass the Fiber context (it implements `io.Writer`). See [examples/showcase](../examples/showcase/README.md) (showcase app).
- **Echo, Chi, etc.:** Get the response writer from the request context and call `el.Render(w)`.

No framework-specific code is required in the generated package; it only depends on `pkg/element`.
