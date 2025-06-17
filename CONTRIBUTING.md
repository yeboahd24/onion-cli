# Contributing to OnionCLI 🧅

Thank you for your interest in contributing to OnionCLI! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Types of Contributions

We welcome all types of contributions:

- 🐛 **Bug Reports**: Help us identify and fix issues
- ✨ **Feature Requests**: Suggest new features and improvements
- 📝 **Documentation**: Improve docs, examples, and guides
- 🧪 **Testing**: Add tests and improve coverage
- 🎨 **UI/UX**: Enhance the terminal interface
- 🔐 **Security**: Identify and fix security issues
- 🌐 **Tor Integration**: Improve .onion service support

### Getting Started

1. **Fork the Repository**
   ```bash
   git clone https://github.com/yourusername/onioncli.git
   cd onioncli
   ```

2. **Set Up Development Environment**
   ```bash
   make dev  # Install dependencies and build
   ```

3. **Run Tests**
   ```bash
   make test
   make test-coverage
   ```

4. **Try the Application**
   ```bash
   make run
   make run-demo
   ```

## 🛠️ Development Workflow

### 1. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number
```

### 2. Make Changes
- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed
- Run quality checks: `make check`

### 3. Test Your Changes
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Test specific functionality
make run-demo
make run-config-demo
```

### 4. Commit Changes
```bash
git add .
git commit -m "feat: add new feature description"
# or
git commit -m "fix: resolve issue with specific component"
```

### 5. Push and Create PR
```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## 📋 Code Guidelines

### Go Code Style

- Follow standard Go formatting: `make fmt`
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Example Code Structure
```go
// Package comment
package main

import (
    "fmt"
    "log"
    
    "onioncli/pkg/api"
)

// ExampleFunction demonstrates proper Go style
func ExampleFunction(param string) error {
    if param == "" {
        return fmt.Errorf("param cannot be empty")
    }
    
    // Implementation here
    return nil
}
```

### TUI Components

- Use Bubbletea patterns for interactive components
- Implement proper keyboard navigation
- Add help text and user guidance
- Handle window resizing gracefully
- Use consistent styling with Lipgloss

### Testing

- Write unit tests for new functions
- Add integration tests for TUI components
- Test Tor integration scenarios
- Include error case testing
- Aim for >80% test coverage

## 🧪 Testing Guidelines

### Unit Tests
```go
func TestExampleFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected error
    }{
        {"valid input", "test", nil},
        {"empty input", "", fmt.Errorf("param cannot be empty")},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ExampleFunction(tt.input)
            if err != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, err)
            }
        })
    }
}
```

### Integration Tests
- Test complete workflows
- Test Tor connectivity scenarios
- Test configuration management
- Test error handling paths

## 📝 Documentation

### Code Documentation
- Add godoc comments for exported functions
- Include usage examples in comments
- Document complex algorithms and logic
- Keep comments up-to-date with code changes

### User Documentation
- Update README.md for new features
- Add examples to demonstrate usage
- Update keyboard shortcuts documentation
- Include troubleshooting information

## 🐛 Bug Reports

When reporting bugs, please include:

### Bug Report Template
```markdown
**Describe the Bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected Behavior**
What you expected to happen.

**Screenshots/Logs**
If applicable, add screenshots or log output.

**Environment:**
- OS: [e.g. Ubuntu 20.04]
- OnionCLI Version: [e.g. 1.0.0]
- Go Version: [e.g. 1.19]
- Tor Version: [e.g. 0.4.7.8]

**Additional Context**
Any other context about the problem.
```

## ✨ Feature Requests

### Feature Request Template
```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Other solutions or features you've considered.

**Additional context**
Any other context, screenshots, or examples.
```

## 🔐 Security Issues

For security vulnerabilities:

1. **DO NOT** create a public issue
2. Email security@onioncli.dev with details
3. Include steps to reproduce
4. Allow time for investigation and fix
5. We'll coordinate disclosure timing

## 📦 Release Process

### Version Numbering
We follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Release Checklist
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Version bumped in appropriate files
- [ ] Changelog updated
- [ ] Release notes prepared
- [ ] Binaries built for all platforms

## 🎯 Areas for Contribution

### High Priority
- [ ] Import/Export from Postman collections
- [ ] Advanced search and filtering
- [ ] Custom themes and styling
- [ ] Performance optimizations
- [ ] Error handling improvements

### Medium Priority
- [ ] GraphQL support
- [ ] WebSocket connections
- [ ] Plugin system
- [ ] Advanced scripting
- [ ] Team collaboration features

### Low Priority
- [ ] Mobile terminal support
- [ ] Alternative TUI frameworks
- [ ] Integration with other tools
- [ ] Advanced analytics

## 💬 Communication

### Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Email**: security@onioncli.dev for security issues
- **Discord**: [OnionCLI Community](https://discord.gg/onioncli) (coming soon)

### Code of Conduct
- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professional communication
- Respect privacy and security concerns

## 🏆 Recognition

Contributors will be:
- Listed in the project README
- Mentioned in release notes
- Invited to the contributors team
- Given credit in documentation

## 📚 Resources

### Learning Resources
- [Go Documentation](https://golang.org/doc/)
- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea)
- [Tor Developer Documentation](https://2019.www.torproject.org/docs/documentation.html.en)
- [SOCKS5 Protocol](https://tools.ietf.org/html/rfc1928)

### Development Tools
- [golangci-lint](https://golangci-lint.run/) for code linting
- [gotests](https://github.com/cweill/gotests) for test generation
- [delve](https://github.com/go-delve/delve) for debugging
- [pprof](https://golang.org/pkg/net/http/pprof/) for profiling

Thank you for contributing to OnionCLI! 🧅✨
