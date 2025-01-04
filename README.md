# Govelte (JSX for Go)

## Overview
govelte introduces a JSX-like syntax for Go, enabling developers to seamlessly write server-side and client-side components in Go. This tool aims to simplify the process of creating dynamic HTML by combining the power of Go with a declarative syntax similar to JSX, commonly used in JavaScript frameworks.

# Try it now

```bash
git clone https://github.com/abdheshnayak/govelte.git
cd govelte
go mod tidy
go run . --src=example/src --dist=example/dist
cd example
go run .
```

### or use `task` to run it

```bash
git clone https://github.com/abdheshnayak/govelte.git
cd govelte
go mod tidy
cd example
task dev
```


## Goals
govelte allows developers to write HTML-like code using Go syntax, which is then transpiled into Go code. The generated Go code can be utilized to render dynamic and reusable components.

## Example Usage
Developers can define reusable components with a JSX-like syntax and use them in their Go applications. Below is an example of defining components and rendering them:

### Defining Components
```jsx
{{- define "Great" }}
<div>
    <p>Hello {name}!</p>
</div>

{{- define "Welcome" }}
<div>
    <p>Welcome to {name}!</p>
    <button onClick={onClick}>Thank You</button>
</div>

{{- define "GreatNWelcome" }}
<div>
    <Great name={name} />
    <Welcome projectName={projectName} />
</div>
```

### Using Components in Go
```go
package main

import (
    . "github.com/abdheshnayak/govelte/pkg/element"
)

func main() {
    GreatNWelcome("Hello Developers", "govelte").Render(os.Stdout)
}

func Great(attrs Attr) Node {
    name := attrs["name"]
    return RenderE("Great", name)
}

func Welcome(attrs Attr) Node {
    projectName := attrs["projectName"]

    onClick := func(e Event) {
        fmt.Println("Thank You")
    }

    return RenderE("Welcome", projectName)
}

func GreatNWelcome(name, projectName string) Node {
    return RenderE("GreatNWelcome", name, projectName)
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
        <p>Welcome to govelte!</p>
        <button onClick="welcomeOnClick">Thank You</button>
    </div>
</div>
<script>
    function welcomeOnClick() {
        console.log("Thank You");
    }
</script>
```

## How It Works
1. **Transpilation:** govelte takes JSX-like syntax written in Go and transpiles it into valid Go code.
2. **Code Replacement:** The transpiler replaces `r.Render` function calls with the generated Go code.
3. **Dynamic Rendering:** The resulting Go code produces dynamic HTML structures, leveraging Go's capabilities for server-side rendering and component-based architecture.

## Benefits
- **Declarative Syntax:** Write HTML-like structures in a readable and reusable manner.
- **Component Reusability:** Define and reuse components efficiently.
- **Seamless Integration:** Combines Go’s performance and JSX’s readability.
- **Dynamic HTML:** Simplifies the creation of dynamic and interactive web content.

## Future Enhancements
- **Improved Error Handling:** Provide detailed errors during transpilation.
- **Enhanced Debugging:** Add tools to visualize the transpilation process.
- **Broader Compatibility:** Extend support for additional libraries and frameworks.

govelte bridges the gap between Go and JSX-like syntax, providing developers with an intuitive way to build modern web applications with Go.


> Current implementation is not as per above mentioned features but it is working in similar way.
