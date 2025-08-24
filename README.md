# GoHTMLX - HTML Components with Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()

GoHTMLX is a modern, powerful tool for building reusable HTML components using Go. It enables developers to write HTML templates with Go-like syntax and compile them into efficient, type-safe Go code for server-side rendering.

## ✨ Features

- 🚀 **Fast & Efficient** - Compiles to native Go code for maximum performance
- 🛡️ **Type Safety** - Full type checking with Go structs and interfaces
- 🔧 **Developer Friendly** - Hot reload, comprehensive CLI, and excellent error messages
- 📦 **Component-Based** - Reusable, composable HTML components
- 🎨 **Modern Syntax** - Clean, intuitive template syntax
- 🔍 **Testing Support** - Built-in testing utilities and helpers
- 📚 **Rich Documentation** - Comprehensive guides and API reference
- 🛠️ **Extensible** - Plugin system and customizable configuration

## 🚀 Quick Start

### Installation

```bash
# Install globally
go install github.com/abdheshnayak/gohtmlx/cmd/gohtmlx@latest

# Or clone and build
git clone https://github.com/abdheshnayak/gohtmlx.git
cd gohtmlx
make install
```

### Try the Example

```bash
git clone https://github.com/abdheshnayak/gohtmlx.git
cd gohtmlx

# Build and run example
make build-example-new && make serve

# Or with development workflow
make dev
```

Visit `http://localhost:3000` to see the example in action!

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

## 📖 Usage

### CLI Commands

GoHTMLX provides a modern CLI with multiple commands:

```bash
# Build components once
gohtmlx build --src=./components --dist=./generated

# Watch for changes and rebuild automatically  
gohtmlx watch --src=./components --dist=./generated --verbose

# Show version
gohtmlx version

# Get help
gohtmlx help
```

### Configuration

Create a `gohtmlx.config.yaml` file in your project root:

```yaml
# Source and output directories
source_dir: "src"
output_dir: "dist" 
package_name: "components"

# Logging
log_level: "info"

# File watching
watch:
  enabled: true
  extensions: [".html", ".htm"]
  debounce_ms: 300

# Code generation
generation:
  format_code: true
  add_comments: true
```

### Project Structure

```
my-project/
├── gohtmlx.config.yaml
├── src/
│   └── components.html
├── dist/
│   └── components/
│       └── components_generated.go
└── main.go
```

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

---

gohtmlx bridges the gap between Go and reusable HTML components, providing developers with an intuitive way to build modern, server-rendered web applications using Go. The examples and usage reflect its ability to simplify server-side HTML generation for projects requiring basic and efficient rendering.
