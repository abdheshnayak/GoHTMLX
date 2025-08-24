# Contributing to GoHTMLX

Thank you for your interest in contributing to GoHTMLX! This guide will help you get started.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Release Process](#release-process)

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your changes
5. Make your changes
6. Test your changes
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile)
- Task (optional, for using Taskfile)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gohtmlx.git
cd gohtmlx

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run the example
make example
```

## Making Changes

### Project Structure

```
gohtmlx/
├── cmd/gohtmlx/          # CLI entry point
├── internal/             # Internal packages
│   ├── cli/             # CLI implementation
│   ├── compiler/        # Compilation logic
│   ├── config/          # Configuration handling
│   ├── generator/       # Code generation
│   ├── parser/          # Template parsing
│   └── watcher/         # File watching
├── pkg/                 # Public packages
│   ├── element/         # Element rendering
│   └── logger/          # Logging utilities
├── example/             # Example project
└── docs/                # Documentation
```

### Branch Naming

Use descriptive branch names:
- `feature/add-new-component-syntax`
- `fix/parser-error-handling`
- `docs/update-readme`
- `refactor/improve-performance`

### Commit Messages

Follow conventional commits:
```
type(scope): description

body (optional)

footer (optional)
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Examples:
```
feat(parser): add support for nested components
fix(cli): handle missing config file gracefully
docs(readme): update installation instructions
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/parser

# Run tests with verbose output
go test -v ./...
```

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Include edge cases and error conditions
- Aim for high test coverage (>80%)

Example test structure:
```go
func TestParseComponent(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected Component
        wantErr  bool
    }{
        {
            name:  "valid component",
            input: "<!-- + define ... -->",
            expected: Component{...},
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Code Style

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Linting

Run the linter before submitting:
```bash
make lint
```

Fix any linting issues:
```bash
golangci-lint run --fix
```

### Documentation

- Document all exported functions and types
- Use examples in documentation comments
- Update README.md for user-facing changes
- Add inline comments for complex logic

## Submitting Changes

### Pull Request Process

1. **Update your branch** with the latest main:
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-feature-branch
   git rebase main
   ```

2. **Run quality checks**:
   ```bash
   make check  # Runs fmt, lint, and test
   ```

3. **Create a pull request**:
   - Use a descriptive title
   - Fill out the PR template
   - Link related issues
   - Add screenshots for UI changes

4. **Address feedback**:
   - Respond to review comments
   - Make requested changes
   - Push updates to your branch

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Added new tests
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## Release Process

Releases are handled by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push a git tag
4. GitHub Actions will build and publish the release

## Getting Help

- 💬 Discussions: Use GitHub Discussions for questions
- 🐛 Issues: Use GitHub Issues for bugs and feature requests
- 📧 Email: Contact maintainers directly for sensitive issues

## Recognition

Contributors will be recognized in:
- CHANGELOG.md
- GitHub releases
- README.md contributors section

Thank you for contributing to GoHTMLX! 🎉
