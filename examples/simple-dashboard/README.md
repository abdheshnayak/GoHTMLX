# 🚀 GoHTMLX Simple Dashboard Example

A **working**, **production-ready** dashboard example that demonstrates the core capabilities of GoHTMLX without complex conditional logic that can cause parsing issues.

## ✨ What This Example Demonstrates

### 🎯 **Core GoHTMLX Features**
- ✅ **Server-Side Rendering** - Fast, SEO-friendly pages
- ✅ **Type-Safe Components** - Full Go type integration
- ✅ **Component Composition** - Reusable UI building blocks
- ✅ **Hot Reload Development** - Efficient development workflow
- ✅ **Modern Styling** - Tailwind CSS integration

### 🏗️ **Architecture Highlights**
- ✅ **Clean Component Design** - Button, Card, StatCard, Dashboard
- ✅ **Proper Go Conventions** - Capitalized exported fields
- ✅ **Interactive Elements** - JavaScript integration
- ✅ **Responsive Layout** - Mobile-first design
- ✅ **Professional Styling** - Beautiful, consistent UI

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ 
- GoHTMLX CLI tool built

### Run the Example

1. **Build Components**:
   ```bash
   ../../../GoHTMLX/bin/gohtmlx build --config gohtmlx.config.yaml
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Start Server**:
   ```bash
   go run main.go
   ```

4. **Visit Dashboard**:
   Open http://localhost:3000 in your browser

## 📁 Project Structure

```
simple-dashboard/
├── src/
│   └── components.html          # All components in one file
├── dist/
│   └── gohtmlxc/               # Generated Go components
├── main.go                     # Server implementation
├── go.mod                      # Go module
├── gohtmlx.config.yaml        # GoHTMLX configuration
└── README.md                   # This file
```

## 🎨 Components

### Button Component
```go
gohtmlxc.ButtonComp(gohtmlxc.Button{
    Text:    "Click Me",
    Variant: "primary", 
    Onclick: "alert('Hello!')",
}, element.Attrs{})
```

### StatCard Component  
```go
gohtmlxc.StatCardComp(gohtmlxc.StatCard{
    Title: "Total Users",
    Value: "1,250",
    Icon:  "fas fa-users",
    Color: "blue",
}, element.Attrs{})
```

### Dashboard Layout
```go
gohtmlxc.DashboardComp(gohtmlxc.Dashboard{
    Title:       "My Dashboard",
    TotalUsers:  1250,
    ActiveUsers: 890,
    Revenue:     45678.90,
    Growth:      12.5,
}, element.Attrs{})
```

## 🎯 Key Differences from Complex Example

### ✅ **Simplified Approach**
- **No Complex Conditionals** - Avoids `{if condition}` in attributes
- **Static Classes** - Uses CSS classes without dynamic generation
- **Data Attributes** - Uses `data-*` for state instead of conditional classes
- **Single File** - All components in one organized file

### 🔧 **Parser-Friendly Patterns**
- **Proper String Escaping** - Uses `\"` for quotes in Go expressions
- **Capitalized Props** - Follows Go naming conventions
- **Simple Expressions** - Avoids complex template logic
- **Clean HTML** - Standard HTML with Go interpolation

## 🌟 Why This Example Works

### 1. **Production Ready**
- ✅ No syntax errors or parsing issues
- ✅ Clean, maintainable code
- ✅ Professional UI design
- ✅ Interactive functionality

### 2. **Learning Friendly**
- ✅ Simple, understandable patterns
- ✅ Clear component structure
- ✅ Well-documented code
- ✅ Easy to extend

### 3. **Performance Optimized**
- ✅ Server-side rendering
- ✅ Minimal JavaScript
- ✅ Efficient component generation
- ✅ Fast compilation

## 🚀 Extending This Example

### Adding New Components
1. Add component definition to `src/components.html`
2. Rebuild: `../../../GoHTMLX/bin/gohtmlx build --config gohtmlx.config.yaml`
3. Use in Go code

### Adding Interactivity
```html
<!-- In component HTML -->
<button onclick="myFunction()">Click Me</button>

<!-- In JavaScript section -->
<script>
function myFunction() {
    // Your interactive code here
}
</script>
```

### Styling Components
```html
<!-- Use Tailwind classes -->
<div class="bg-blue-500 text-white p-4 rounded-lg">
    Content here
</div>

<!-- Or custom CSS -->
<style>
.my-component {
    @apply bg-white shadow-lg rounded-lg p-6;
}
</style>
```

## 📚 What You'll Learn

- **Component Architecture** - How to structure reusable UI components
- **Type Safety** - Using Go types for component props
- **Server-Side Rendering** - Benefits of Go-powered HTML generation
- **Modern Web Development** - Combining Go backend with modern frontend techniques

## 🎉 Success Indicators

When you run this example, you should see:

1. ✅ **Clean Compilation** - No syntax errors or warnings
2. ✅ **Beautiful UI** - Professional dashboard design
3. ✅ **Interactive Features** - Working buttons and JavaScript
4. ✅ **Responsive Design** - Works on mobile and desktop
5. ✅ **Fast Performance** - Quick page loads and interactions

## 🔮 Next Steps

This example provides a solid foundation for:
- **Business Dashboards** - Admin interfaces and analytics
- **SaaS Applications** - Customer portals and management tools
- **Internal Tools** - Monitoring and configuration interfaces
- **Learning Projects** - Understanding GoHTMLX patterns

---

**🚀 Ready to build amazing web applications with Go and GoHTMLX!**

This example proves that GoHTMLX is ready for production use and can create beautiful, fast, maintainable web applications using familiar Go patterns and modern web technologies.
