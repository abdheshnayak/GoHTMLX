# 🚀 GoHTMLX Dashboard Example

A comprehensive, modern dashboard application built with GoHTMLX showcasing real-world component architecture and best practices.

## ✨ Features

### 🎨 Modern UI Components
- **Responsive Design** - Mobile-first approach with Tailwind CSS
- **Component Library** - Reusable UI components (Button, Card, Table, Modal, etc.)
- **Interactive Elements** - Dropdowns, modals, tooltips, and notifications
- **Data Visualization** - Charts and graphs with Chart.js integration
- **Form Components** - Input fields, validation, and form handling

### 🏗️ Architecture
- **Modular Structure** - Organized component hierarchy
- **Type Safety** - Full TypeScript-like type definitions in Go
- **Server-Side Rendering** - Fast initial page loads
- **Progressive Enhancement** - JavaScript for interactivity
- **RESTful API** - Backend API endpoints for data operations

### 📊 Dashboard Features
- **Statistics Overview** - Key performance metrics
- **User Management** - User listing, editing, and management
- **Activity Feed** - Real-time activity tracking
- **Navigation System** - Multi-level sidebar navigation
- **Notification Center** - System notifications and alerts
- **Search & Filtering** - Advanced data filtering capabilities
- **Pagination** - Efficient data browsing

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or later
- GoHTMLX CLI tool

### Installation

1. **Build the GoHTMLX CLI** (from project root):
   ```bash
   cd ../../
   make build
   ```

2. **Generate Components**:
   ```bash
   cd examples/dashboard
   ../../bin/gohtmlx build --config gohtmlx.config.yaml
   ```

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the Server**:
   ```bash
   go run main.go
   ```

5. **Visit the Dashboard**:
   Open http://localhost:3000 in your browser

### Development Workflow

For development with hot reload:

```bash
# Terminal 1: Watch for component changes
../../bin/gohtmlx watch --config gohtmlx.config.yaml

# Terminal 2: Run server with auto-restart
go run main.go
```

## 📁 Project Structure

```
dashboard/
├── src/                          # Source components
│   ├── components/              # Reusable UI components
│   │   ├── ui.html             # Basic UI components
│   │   ├── table.html          # Table and pagination
│   │   └── navigation.html     # Navigation components
│   ├── layouts/                # Page layouts
│   │   └── main.html          # Main layout template
│   ├── pages/                  # Page components
│   │   └── dashboard.html     # Dashboard page
│   └── types/                  # Type definitions
│       └── main.go            # Go type definitions
├── static/                     # Static assets
│   ├── css/                   # Stylesheets
│   ├── js/                    # JavaScript files
│   │   └── dashboard.js       # Main dashboard JS
│   └── images/                # Image assets
├── dist/                      # Generated Go code
│   └── gohtmlxc/             # Generated components
├── main.go                    # Main server file
├── go.mod                     # Go module definition
├── gohtmlx.config.yaml       # GoHTMLX configuration
└── README.md                  # This file
```

## 🎯 Component Examples

### Basic Button Component
```html
<!-- + define "Button" -->
<!-- | define "props" -->
text: string
variant: t.ButtonVariant
size: t.ButtonSize
disabled: bool
icon: string
onclick: string
<!-- | end -->
<!-- | define "html" -->
<button class="btn btn-{props.variant} btn-{props.size}" {if props.onclick}onclick="{props.onclick}"{end}>
  <if condition="{props.icon != ''}">
    <i class="{props.icon} mr-2"></i>
  </if>
  {props.text}
</button>
<!-- | end -->
<!-- + end -->
```

### Usage in Go
```go
button := gohtmlxc.ButtonComp(gohtmlxc.Button{
    Text:    "Save Changes",
    Variant: types.ButtonPrimary,
    Size:    types.ButtonMedium,
    Icon:    "fas fa-save",
    Onclick: "saveData()",
}, element.Attrs{})
```

### Data Table with Pagination
```go
table := gohtmlxc.TableComp(gohtmlxc.Table{
    Columns:      getUserTableColumns(),
    Data:         formatUsersForTable(users),
    Striped:      true,
    Hoverable:    true,
    EmptyMessage: "No users found",
}, element.Attrs{})
```

## 🎨 Styling

The dashboard uses Tailwind CSS for styling with custom component classes:

- **Button variants**: `btn-primary`, `btn-secondary`, `btn-success`, etc.
- **Badge variants**: `badge-primary`, `badge-success`, `badge-warning`, etc.
- **Avatar sizes**: `avatar-sm`, `avatar-md`, `avatar-lg`
- **Status indicators**: `status-online`, `status-offline`, `status-away`

## 🔧 API Endpoints

The dashboard includes a RESTful API:

- `GET /api/stats` - Dashboard statistics
- `GET /api/users` - User listing with pagination
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `GET /api/notifications` - Get notifications

## 🌟 Best Practices Demonstrated

### 1. **Component Organization**
- Logical grouping of related components
- Consistent naming conventions
- Reusable component patterns

### 2. **Type Safety**
- Comprehensive type definitions
- Proper prop validation
- Type-safe data handling

### 3. **Performance**
- Server-side rendering for fast initial loads
- Efficient component generation
- Optimized asset loading

### 4. **Accessibility**
- Semantic HTML structure
- ARIA labels and roles
- Keyboard navigation support

### 5. **Responsive Design**
- Mobile-first approach
- Flexible grid layouts
- Adaptive component behavior

## 🚀 Advanced Features

### Real-time Updates with HTMX
```html
<div hx-get="/api/stats" hx-trigger="every 30s" hx-swap="innerHTML">
  <!-- Stats will auto-refresh -->
</div>
```

### Interactive Charts
```javascript
// Chart initialization in dashboard.js
function initializeCharts(data) {
    new Chart(ctx, {
        type: 'line',
        data: data.revenueData,
        options: chartOptions
    });
}
```

### Modal Dialogs
```go
modal := gohtmlxc.ModalComp(gohtmlxc.Modal{
    Title:    "Add New User",
    Size:     "lg",
    Show:     true,
    Closable: true,
}, element.Attrs{}, 
    // Modal content here
)
```

## 🔮 Extending the Dashboard

### Adding New Components
1. Create component in `src/components/`
2. Define props and HTML structure
3. Rebuild with GoHTMLX
4. Use in your Go code

### Adding New Pages
1. Create page component in `src/pages/`
2. Add route handler in `main.go`
3. Update navigation if needed

### Custom Styling
1. Add CSS classes in the layout
2. Follow Tailwind conventions
3. Use component-specific modifiers

## 📚 Learning Resources

This example demonstrates:
- **Component-based architecture**
- **Server-side rendering with Go**
- **Modern web development practices**
- **Type-safe template systems**
- **Progressive enhancement patterns**

## 🤝 Contributing

This example serves as:
- **Reference implementation** for GoHTMLX best practices
- **Starting point** for new dashboard projects
- **Learning resource** for component architecture
- **Testing ground** for new GoHTMLX features

## 📄 License

This example is part of the GoHTMLX project and follows the same license terms.

---

**Built with ❤️ using GoHTMLX** - The modern way to build server-side rendered web applications with Go.
