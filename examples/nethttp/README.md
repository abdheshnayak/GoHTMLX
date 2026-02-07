# net/http example

Same component as [minimal](../minimal), served with the standard library only. No Fiber or other framework.

## Run

From the **repository root**:

```bash
# Transpile
go run . --src=examples/nethttp/src --dist=examples/nethttp/dist

# Start server (http://localhost:8080)
cd examples/nethttp && go run .
```

Open http://localhost:8080 for "Hello, GoHTMLX!". Use `?name=World` for "Hello, World!".

## Re-generate

After changing `src/hello.html`:

```bash
# From repo root
go run . --src=examples/nethttp/src --dist=examples/nethttp/dist
```
