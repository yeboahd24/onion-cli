# OnionCLI 🧅

**A powerful terminal-based API client specifically designed for .onion services and Tor networks**

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/yourusername/onioncli)

OnionCLI is a Postman alternative built specifically for testing and interacting with .onion APIs through the Tor network. It provides a beautiful terminal user interface with comprehensive features for API development, testing, and debugging in the dark web ecosystem.

## 📸 Screenshots

### Main Interface
```
╭─ OnionCLI - Request Builder ─────────────────────────────────────╮
│ URL: http://3g2upl4pq6kufc4m.onion/api/search                    │
│ Method: [GET] POST PUT DELETE PATCH HEAD OPTIONS                 │
│                                                                  │
│ Headers:                                                         │
│ ┃ User-Agent: OnionCLI/1.0                                       │
│ ┃ Accept: application/json                                       │
│                                                                  │
│ Request Body:                                                    │
│ ┃ {                                                              │
│ ┃   "query": "privacy tools",                                    │
│ ┃   "limit": 10                                                  │
│ ┃ }                                                              │
│                                                                  │
│ [ Send Request ]                                                 │
╰──────────────────────────────────────────────────────────────────╯
```

### Response Viewer
```
╭─ Response - 200 OK (2.3s via Tor) ──────────────────────────────╮
│ {                                                                │
│   "results": [                                                  │
│     {                                                            │
│       "title": "Privacy Tools",                                 │
│       "url": "http://example.onion/tools",                      │
│       "description": "Essential privacy tools..."               │
│     }                                                            │
│   ],                                                             │
│   "total": 42                                                   │
│ }                                                                │
╰──────────────────────────────────────────────────────────────────╯
```

## 🎯 Problem Statement

Traditional API clients like Postman, Insomnia, or curl don't provide seamless integration with Tor networks and .onion services. Developers working with:

- **Dark web APIs** and .onion services
- **Privacy-focused applications** requiring Tor routing
- **Decentralized services** on hidden networks
- **Security research** and penetration testing

...face challenges with:
- ❌ Complex Tor proxy configuration
- ❌ Poor error handling for Tor-specific issues
- ❌ No built-in .onion URL validation
- ❌ Lack of Tor network diagnostics
- ❌ No understanding of Tor latency patterns

## ✨ Solution: OnionCLI

OnionCLI solves these problems by providing:

- 🧅 **Native Tor Integration**: Automatic SOCKS5 proxy configuration
- 🔍 **Smart .onion Detection**: Automatic routing for .onion URLs
- 🎨 **Beautiful TUI**: Terminal interface built with Bubbletea/Lipgloss
- 🚀 **Performance Optimized**: Designed for Tor's higher latency
- 🔐 **Security First**: Built with privacy and security in mind

## 🚀 Features

### 🌐 Core Functionality
- **Tor Network Integration**: Seamless SOCKS5 proxy support for .onion services
- **HTTP Methods**: Support for GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **Request Builder**: Interactive form-based request construction
- **Response Viewer**: Pretty-printed JSON, XML, and text responses
- **Real-time Feedback**: Loading spinners and status indicators

### 🔐 Authentication & Security
- **Multiple Auth Methods**: API Keys, Bearer Tokens, Basic Auth, Custom Headers
- **Secure Storage**: Encrypted credential management
- **Session Management**: Persistent authentication across requests
- **Custom Headers**: Full control over request headers

### 📚 Organization & Workflow
- **Request Collections**: Organize related requests into collections
- **Environment Management**: Multiple environments (dev, staging, prod)
- **Variable Substitution**: Use `{{variables}}` in URLs and headers
- **Request History**: Persistent history with search and replay
- **Save & Load**: Save frequently used requests

### 🎯 Tor-Specific Features
- **Automatic .onion Detection**: Smart routing for hidden services
- **Tor Connection Testing**: Built-in connectivity diagnostics
- **Error Analysis**: Tor-specific error messages and suggestions
- **Latency Optimization**: UI optimized for Tor's network characteristics
- **Circuit Information**: Display Tor circuit details (when available)

### 🎨 User Experience
- **Interactive TUI**: Beautiful terminal interface with keyboard shortcuts
- **Syntax Highlighting**: JSON/XML response highlighting
- **Progress Indicators**: Visual feedback for long-running requests
- **Error Handling**: Comprehensive error analysis with actionable suggestions
- **Keyboard Shortcuts**: Efficient navigation and quick actions

### ⚙️ Configuration & Customization
- **Flexible Configuration**: YAML-based configuration management
- **Proxy Settings**: Customizable Tor proxy configuration
- **Themes**: Dark/light theme support
- **Timeouts**: Configurable request timeouts for Tor networks
- **Export/Import**: Backup and share configurations

## 📦 Installation

### Pre-built Binaries

Download the latest release for your platform:

```bash
# Linux
wget https://github.com/yourusername/onioncli/releases/latest/download/onioncli_linux_amd64.tar.gz
tar -xzf onioncli_linux_amd64.tar.gz
sudo mv onioncli /usr/local/bin/

# macOS
wget https://github.com/yourusername/onioncli/releases/latest/download/onioncli_darwin_amd64.tar.gz
tar -xzf onioncli_darwin_amd64.tar.gz
sudo mv onioncli /usr/local/bin/

# Windows
# Download onioncli_windows_amd64.zip and extract
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/onioncli.git
cd onioncli

# Build using Makefile
make build

# Or build manually
go build -o onioncli ./cmd/onioncli
```

### Using Go Install

```bash
go install github.com/yourusername/onioncli/cmd/onioncli@latest
```

### Package Managers

```bash
# Homebrew (macOS/Linux)
brew install onioncli

# Snap (Linux)
sudo snap install onioncli

# AUR (Arch Linux)
yay -S onioncli

# Chocolatey (Windows)
choco install onioncli
```

### Docker

```bash
# Run with Docker
docker run -it --rm onioncli/onioncli

# Build Docker image
docker build -t onioncli .
```

## 🚀 Quick Start

### 1. Start OnionCLI
```bash
onioncli
```

### 2. Configure Tor (if needed)
OnionCLI automatically detects and uses Tor proxy at `127.0.0.1:9050`. If you need custom configuration:
- Press `c` for settings
- Configure proxy address and port
- Test connection

### 3. Make Your First Request
1. Enter a .onion URL (e.g., `http://example.onion/api/users`)
2. Select HTTP method (GET, POST, etc.)
3. Add headers if needed
4. Add request body for POST/PUT requests
5. Press Enter to send

### 4. Explore Features
- Press `h` to view request history
- Press `c` to browse collections
- Press `v` to manage environments
- Press `a` to configure authentication
- Press `?` for keyboard shortcuts

## 🎮 Usage Examples

### Basic GET Request
```
URL: http://3g2upl4pq6kufc4m.onion/api/search?q=privacy
Method: GET
Headers: 
  User-Agent: OnionCLI/1.0
  Accept: application/json
```

### POST with Authentication
```
URL: http://example.onion/api/users
Method: POST
Headers:
  Authorization: Bearer {{api_token}}
  Content-Type: application/json
Body:
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

### Using Environment Variables
```
# Development Environment
base_url = http://dev-api.example.onion:8080
api_token = dev-token-123

# Request
URL: {{base_url}}/api/users
Headers:
  Authorization: Bearer {{api_token}}
```

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Navigate between fields |
| `Enter` | Send request / Select item |
| `Esc` | Go back / Cancel |
| `h` | View request history |
| `c` | Browse collections |
| `v` | Manage environments |
| `a` | Configure authentication |
| `s` | Save current request |
| `r` | Retry last request |
| `e` | View error details |
| `?` | Toggle help |
| `q` / `Ctrl+C` | Quit application |

## 🔧 Configuration

OnionCLI stores configuration in `~/.onioncli/`:

```
~/.onioncli/
├── config.yaml          # Main configuration
├── environments.json    # Environment variables
├── collections/         # Request collections
│   ├── collection1.json
│   └── collection2.json
└── history.json         # Request history
```

### Sample Configuration

```yaml
tor:
  enabled: true
  proxy_addr: "127.0.0.1"
  proxy_port: 9050
  timeout: 30
  auto_detect: true

http:
  timeout: 30
  follow_redirects: true
  max_redirects: 10
  verify_ssl: true
  user_agent: "OnionCLI/1.0"

ui:
  theme: "dark"
  show_line_numbers: true
  auto_save: true
  confirm_exit: false

history:
  enabled: true
  max_entries: 100
  auto_save: true

default_headers:
  User-Agent: "OnionCLI/1.0"
  Accept: "application/json, text/plain, */*"
```

## 🛠️ Development

### Prerequisites
- Go 1.19 or later
- Tor daemon running (for testing .onion requests)
- Make (optional, for using Makefile)

### Building
```bash
# Using Makefile (recommended)
make dev          # Full development setup
make build        # Build for current platform
make build-all    # Build for all platforms

# Manual build
go build -o onioncli ./cmd/onioncli
```

### Testing
```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run demos
make run-demo
make run-config-demo
make run-performance-demo
```

### Project Structure
```
OnionCLI/
├── cmd/onioncli/         # Main application entry point
├── pkg/
│   ├── api/              # HTTP client and authentication
│   ├── collections/      # Collections and environments
│   ├── config/           # Configuration management
│   ├── history/          # Request history
│   └── tui/              # Terminal UI components
├── examples/             # Demo applications
├── Makefile             # Build automation
└── README.md            # This file
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Areas for Contribution
- 🐛 Bug fixes and improvements
- ✨ New features and enhancements
- 📚 Documentation improvements
- 🧪 Test coverage expansion
- 🎨 UI/UX improvements
- 🔐 Security enhancements

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubbletea and Lipgloss libraries
- [Tor Project](https://www.torproject.org/) for privacy and anonymity tools
- The Go community for excellent tooling and libraries

## 🔒 Security Considerations

### Tor Network Safety
- OnionCLI routes .onion requests through Tor automatically
- Regular HTTP requests can optionally use Tor proxy
- No request data is logged or transmitted outside Tor network
- Credentials are stored locally with encryption

### Privacy Features
- No telemetry or analytics collection
- Local-only configuration and history storage
- Optional request history (can be disabled)
- Secure credential storage with encryption

### Best Practices
- Always verify .onion URLs before making requests
- Use environment variables for sensitive data
- Regularly update OnionCLI for security patches
- Review saved requests before sharing collections

## 🚨 Troubleshooting

### Common Issues

**Tor Connection Failed**
```bash
# Check if Tor is running
sudo systemctl status tor

# Start Tor daemon
sudo systemctl start tor

# Test Tor connectivity
curl --socks5 127.0.0.1:9050 http://check.torproject.org
```

**Permission Denied**
```bash
# Make binary executable
chmod +x onioncli

# Check file permissions
ls -la onioncli
```

**Module Download Issues**
```bash
# Clean and rebuild
make clean-all
make deps
make build
```

**Configuration Issues**
```bash
# Reset configuration
rm -rf ~/.onioncli
onioncli  # Will recreate default config
```

### Getting Help
1. Check this README and [MAKEFILE_USAGE.md](MAKEFILE_USAGE.md)
2. Run `onioncli --help` for command-line options
3. Press `?` in the TUI for keyboard shortcuts
4. Check [Issues](https://github.com/yourusername/onioncli/issues) for known problems

## 🗺️ Roadmap

### Version 2.0 (Planned)
- [ ] GraphQL support for .onion APIs
- [ ] WebSocket connections through Tor
- [ ] Plugin system for custom authentication
- [ ] Advanced scripting and automation
- [ ] Team collaboration features
- [ ] API documentation generation

### Version 1.5 (In Progress)
- [x] Request collections and environments
- [x] Variable substitution
- [x] Performance optimizations
- [x] Enhanced error handling
- [ ] Import/export from Postman
- [ ] Advanced filtering and search
- [ ] Custom themes and styling

### Version 1.0 (Current)
- [x] Core Tor integration
- [x] Interactive TUI
- [x] Authentication support
- [x] Request history
- [x] Configuration management
- [x] Error diagnostics

## 📊 Performance

OnionCLI is optimized for Tor network characteristics:

- **Latency Handling**: UI designed for higher Tor latencies
- **Connection Reuse**: Efficient SOCKS5 connection management
- **Memory Usage**: ~10-15MB RAM usage
- **Binary Size**: ~13MB (statically linked)
- **Startup Time**: <100ms cold start
- **Request Throughput**: Limited by Tor network, not OnionCLI

## 🌟 Why OnionCLI?

### vs. Postman
- ✅ Native Tor support vs. ❌ Complex proxy setup
- ✅ .onion URL validation vs. ❌ No dark web support
- ✅ Terminal-based vs. ❌ GUI-only
- ✅ Privacy-focused vs. ❌ Cloud-based features

### vs. curl
- ✅ Interactive TUI vs. ❌ Command-line only
- ✅ Request history vs. ❌ No persistence
- ✅ Collections vs. ❌ No organization
- ✅ Error analysis vs. ❌ Basic error messages

### vs. HTTPie
- ✅ Tor integration vs. ❌ Manual proxy setup
- ✅ Interactive mode vs. ❌ Command-line only
- ✅ Authentication management vs. ❌ Per-request auth
- ✅ .onion specialization vs. ❌ General purpose

## 📞 Support

- 📖 Documentation: [Wiki](https://github.com/yourusername/onioncli/wiki)
- 🐛 Bug Reports: [Issues](https://github.com/yourusername/onioncli/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/yourusername/onioncli/discussions)
- 📧 Email: support@onioncli.dev

---

**Made with ❤️ for the privacy-conscious developer community**

*OnionCLI: Because privacy shouldn't be complicated* 🧅
