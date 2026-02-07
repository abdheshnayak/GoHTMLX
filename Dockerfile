# Build stage: transpile GoHTMLX, build example CSS, build example binary
FROM golang:1.23-bookworm AS builder
RUN apt-get update && apt-get install -y nodejs npm && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY . .

# Build gohtmlx CLI and transpile example
RUN go build -o /gohtmlx . && /gohtmlx --src=examples/showcase/src --dist=examples/showcase/dist

# Build example static assets (CSS + app.js) — Tailwind v4 needs tailwindcss in node_modules for @import
RUN cd examples/showcase && npm install && mkdir -p dist/static && \
    npx @tailwindcss/cli -i ./src/input.css -o ./dist/static/main.css

# Build example app
RUN go build -o /example-app ./examples/showcase

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /example-app /app/example-app
COPY --from=builder /build/examples/showcase/dist/static /app/dist/static
EXPOSE 3000
CMD ["/app/example-app"]
