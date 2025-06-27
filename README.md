# MCP Teleport

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/giantswarm/mcp-teleport.svg)](https://github.com/giantswarm/mcp-teleport/releases)

A Model Context Protocol (MCP) server that enables AI assistants like Claude to interact with [Teleport](https://goteleport.com/) infrastructure through the `tsh` CLI.

## Overview

MCP Teleport bridges the gap between AI assistants and Teleport infrastructure by providing a standardized interface for:

- **Authentication**: Login and session management
- **SSH Access**: Secure shell access to remote nodes
- **Kubernetes**: Cluster discovery & authentication
- **Application Access**: Web application tunneling

## Features

### 🔐 **Authentication Tools**
- `teleport_login` - Login to Teleport clusters
- `teleport_status` - Check current authentication status
- `teleport_list_clusters` - List available clusters

### 🖥️ **SSH Tools**
- `teleport_list_ssh_nodes` - List available SSH nodes
- `teleport_ssh` - Execute commands on remote SSH nodes

### ☸️ **Kubernetes Tools**
- `teleport_kube_list_clusters` - List available Kubernetes clusters
- `teleport_kube_login` - Login to Kubernetes clusters and update kubeconfig

### 🌐 **Application Tools** *(Coming Soon)*
- `teleport_apps` - List available applications
- `teleport_app_login` - Access web applications

### 🛠️ **Operational Features**
- **Multiple Transports**: stdio, SSE, streamable HTTP
- **Dry Run Mode**: Test operations safely
- **Debug Logging**: Comprehensive troubleshooting
- **Command Timeouts**: Prevent hanging operations
- **Structured Responses**: Consistent error handling

## Prerequisites

- **Go 1.24+** (for building from source)
- **Teleport CLI (`tsh`)** installed and configured
- **Active Teleport cluster** access

### Install Teleport CLI

```bash
# macOS
brew install teleport

# Linux (Ubuntu/Debian)
curl -O https://get.gravitational.com/teleport-v15.1.0-linux-amd64-bin.tar.gz
tar -xzf teleport-v15.1.0-linux-amd64-bin.tar.gz
sudo mv teleport/tsh /usr/local/bin/

# Verify installation
tsh version
```

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/giantswarm/mcp-teleport/releases).

```bash
# Linux/macOS
wget https://github.com/giantswarm/mcp-teleport/releases/latest/download/mcp-teleport_Linux_x86_64.tar.gz
tar -xzf mcp-teleport_Linux_x86_64.tar.gz
sudo mv mcp-teleport /usr/local/bin/

# Verify installation
mcp-teleport version
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/giantswarm/mcp-teleport.git
cd mcp-teleport

# Build
go build -o mcp-teleport

# Optional: Install to PATH
sudo mv mcp-teleport /usr/local/bin/
```

### Go Install

```bash
go install github.com/giantswarm/mcp-teleport@latest
```

## Usage

### Standalone Server

#### Start with Default Settings (stdio)
```bash
mcp-teleport
```

#### Web-based Deployment (SSE)
```bash
mcp-teleport serve --transport=sse --http-addr=:8080
```

#### Debug Mode with Dry Run
```bash
mcp-teleport serve --debug --dry-run
```

### Integration with Claude Desktop

1. **Install mcp-teleport** (see installation section above)

2. **Configure Claude Desktop**:
   
   Open Claude Desktop settings and add to your MCP servers configuration:

   ```json
   {
     "mcpServers": {
       "teleport": {
         "command": "mcp-teleport",
         "args": ["serve"]
       }
     }
   }
   ```

   **Alternative with full path**:
   ```json
   {
     "mcpServers": {
       "teleport": {
         "command": "/usr/local/bin/mcp-teleport",
         "args": ["serve", "--debug"]
       }
     }
   }
   ```

3. **Restart Claude Desktop**

4. **Start Interacting**:
   - "Check my Teleport login status"
   - "List all SSH nodes available through Teleport"
   - "SSH to server web-01 and run 'uptime'"
   - "Login to Teleport cluster teleport.example.com as user alice"
   - "List all Kubernetes clusters I can access"
   - "Login to the production Kubernetes cluster"

### Command Reference

```bash
# Server commands
mcp-teleport                          # Start with stdio transport
mcp-teleport serve --transport=sse    # Start with SSE transport
mcp-teleport serve --debug            # Enable debug logging
mcp-teleport serve --dry-run          # Simulate operations

# Utility commands
mcp-teleport version                  # Show version
mcp-teleport selfupdate               # Update to latest version
mcp-teleport --help                   # Show help
```

## Configuration

### Server Options

| Flag | Description | Default |
|------|-------------|---------|
| `--transport` | Transport type: stdio, sse, streamable-http | `stdio` |
| `--http-addr` | HTTP server address | `:8080` |
| `--debug` | Enable debug logging | `false` |
| `--dry-run` | Simulate operations | `false` |
| `--non-destructive` | Prevent destructive operations | `true` |

### Transport Types

#### STDIO (Default)
Perfect for Claude Desktop integration:
```bash
mcp-teleport serve
```

#### Server-Sent Events (SSE)  
For web applications requiring real-time updates:
```bash
mcp-teleport serve --transport=sse --http-addr=:8080
```
Access at: `http://localhost:8080/sse`

#### Streamable HTTP
For modern web applications:
```bash
mcp-teleport serve --transport=streamable-http --http-addr=:8080
```
Access at: `http://localhost:8080/mcp`

## Example Interactions

### Authentication
```
User: "Check my current Teleport login status"
AI: Uses teleport_status tool
Response: Current login information and certificate details

User: "Login to teleport.example.com as user alice"  
AI: Uses teleport_login tool with proxy and user parameters
Response: Login success confirmation
```

### SSH Operations
```
User: "List all SSH servers I can access"
AI: Uses teleport_list_ssh_nodes tool
Response: Available SSH nodes with details

User: "SSH to web-server-01 and check disk usage"
AI: Uses teleport_ssh tool with destination and command
Response: Command output from remote server
```

### Kubernetes Operations
```
User: "List all Kubernetes clusters I can access"
AI: Uses teleport_kube_list_clusters tool
Response: Available K8s clusters with labels and metadata

User: "Login to the production Kubernetes cluster"
AI: Uses teleport_kube_login tool with kubeCluster parameter
Response: Kubeconfig updated, kubectl access configured

User: "Show me all clusters with environment=production label"
AI: Uses teleport_kube_list_clusters with query parameter
Response: Filtered list of production clusters
```

## Development

### Project Structure

```
mcp-teleport/
├── cmd/                    # CLI commands and entry points
│   ├── root.go            # Root command setup
│   ├── serve.go           # Server command with transport options
│   ├── version.go         # Version command
│   └── selfupdate.go      # Self-update functionality
├── internal/
│   ├── server/            # Server context and configuration
│   │   ├── context.go     # Server context management
│   │   └── doc.go         # Package documentation
│   ├── teleport/          # Teleport CLI wrapper
│   │   ├── client.go      # tsh command execution
│   │   ├── client_test.go # Unit tests
│   │   └── doc.go         # Package documentation
│   └── tools/             # MCP tool implementations
│       ├── auth/          # Authentication tools
│       ├── ssh/           # SSH tools
│       ├── kube/          # Kubernetes tools
│       ├── database/      # Database tools (stubs)
│       └── apps/          # Application tools (stubs)
├── .goreleaser.yaml       # Release configuration
├── go.mod                 # Go module definition
├── LICENSE                # MIT license
├── Makefile              # Build automation
└── README.md             # This file
```

### Building

```bash
# Install dependencies
go mod download

# Run tests
go test -v ./...

# Build binary
go build -o mcp-teleport

# Build for all platforms
goreleaser build --snapshot --clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Test server in dry-run mode
./mcp-teleport serve --dry-run --debug
```

### Adding New Tools

1. **Create tool package** in `internal/tools/`
2. **Implement tool registration** function
3. **Add handler functions** for each tool
4. **Register in serve.go** 
5. **Add tests** for new functionality

Example:
```go
// internal/tools/example/tools.go
func RegisterExampleTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
    tool := mcp.NewTool("teleport_example", ...)
    s.AddTool(tool, handleExample)
    return nil
}
```

## Security Considerations

- **Principle of Least Privilege**: Run with minimal required permissions
- **Network Security**: Use HTTPS for web transports in production
- **Teleport RBAC**: Ensure proper Teleport role-based access controls
- **Command Validation**: All tsh commands are validated before execution
- **Timeout Protection**: Commands timeout after 30 seconds to prevent hanging

### Production Deployment

```bash
# Use non-destructive mode in production
mcp-teleport serve --non-destructive

# Enable logging for audit trail
mcp-teleport serve --debug

# Consider running as dedicated user
sudo -u teleport-mcp mcp-teleport serve
```

## Troubleshooting

### Common Issues

**❌ "tsh command not found"**
```bash
# Verify tsh installation
which tsh
tsh version

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin
```

**❌ "Command timeout after 30 seconds"**
- Check network connectivity to Teleport cluster
- Verify proxy address is correct
- Try with `--debug` flag for detailed logging

**❌ "Permission denied"**
- Ensure valid Teleport login session
- Check Teleport user permissions
- Verify RBAC policies allow requested operations

**❌ "Failed to start MCP server"**
- Check if port is already in use (for web transports)
- Verify Go binary has execute permissions
- Try with `--debug` flag for detailed error information

### Debug Mode

Enable comprehensive logging:
```bash
mcp-teleport serve --debug --dry-run
```

### Logs Analysis

```bash
# Check recent logs (if using systemd)
journalctl -u mcp-teleport -f

# View server output
mcp-teleport serve --debug 2>&1 | tee mcp-teleport.log
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Make** your changes and add tests
4. **Test** your changes: `make test`
5. **Commit** your changes: `git commit -m 'Add amazing feature'`
6. **Push** to the branch: `git push origin feature/amazing-feature`
7. **Open** a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/mcp-teleport.git
cd mcp-teleport

# Install dependencies
go mod download

# Run tests
make test

# Build and test locally
make build
./mcp-teleport serve --dry-run --debug
```

## Roadmap

- [x] **v1.1**: Complete Kubernetes tools implementation ✅
- [ ] **v1.2**: Database tools implementation  
- [ ] **v1.3**: Application tools implementation
- [ ] **v1.4**: Resource management tools
- [ ] **v2.0**: Advanced workflow automation
- [ ] **v2.1**: Teleport Connect integration
- [ ] **v2.2**: Multi-cluster management

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Teleport](https://goteleport.com/) - Modern infrastructure access platform
- [Model Context Protocol](https://modelcontextprotocol.io/) - Standardized AI integration
- [mcp-kubernetes](https://github.com/giantswarm/mcp-kubernetes) - Structural inspiration
- [Claude](https://claude.ai/) - AI assistant integration target

## Support

- **Documentation**: [GitHub Wiki](https://github.com/giantswarm/mcp-teleport/wiki)
- **Issues**: [GitHub Issues](https://github.com/giantswarm/mcp-teleport/issues)
- **Discussions**: [GitHub Discussions](https://github.com/giantswarm/mcp-teleport/discussions)
- **Email**: [support@giantswarm.io](mailto:support@giantswarm.io)

---

**Made with ❤️ by [Giant Swarm](https://giantswarm.io)** 

## SSH Command Execution

The `teleport_ssh` tool supports two destination formats:

### 1. Direct Hostname Targeting
```bash
# Target a specific node
"destination": "root@hostname"
```

### 2. Label Selector Targeting (Multi-Node)
```bash
# Target all nodes matching label criteria
"destination": "root@role=worker,env=prod"
"destination": "root@cluster=wallaby,role=control-plane"
```

**Important:** Only one-time commands are supported. Interactive shell sessions are not supported via MCP.

### Expected Output Formats

**Single Node Output:**
```
total 408
drwxr-xr-x. 8 root root  4096 Mar 26 12:18 kubernetes
drwxr-xr-x. 1 root root  4096 Mar 26 12:18 sysctl.d
```

**Multi-Node Output (Label Selector):**
```
WARNING: Multiple nodes matched label selector, running command on all.
Running command on wallaby-9wldd:
Running command on wallaby-rd565:
[wallaby-9wldd]  07:15:15 up 92 days, 18:59,  0 user,  load average: 1.71, 1.55, 1.47
[wallaby-rd565]  07:15:15 up 92 days, 19:05,  0 user,  load average: 1.10, 0.92, 0.90

[wallaby-9wldd] success
[wallaby-rd565] success

2 host(s) succeeded; 0 host(s) failed
```

## SSH Node Listing

The `teleport_list_ssh_nodes` tool returns JSON formatted node information:

```json
[
  {
    "kind": "node",
    "version": "v2",
    "metadata": {
      "name": "41c3ee63-af98-44b1-9ec6-14cb19ba7e6b",
      "labels": {
        "azure/environment": "prod",
        "cluster": "wallaby",
        "role": "control-plane"
      }
    },
    "spec": {
      "hostname": "wallaby-9wldd",
      "addr": "",
      "cmd_labels": {
        "arch": {"result": "x86_64"},
        "role": {"result": "control-plane"}
      }
    }
  }
]
```

This is parsed and presented in a user-friendly format showing hostname, labels, and dynamic command labels.

## Kubernetes Cluster Management

### Kubernetes Cluster Discovery

The `teleport_kube_list_clusters` tool provides comprehensive cluster discovery with filtering capabilities:

#### Basic Usage
```bash 
# List all accessible clusters
"No additional parameters needed"
```

#### Advanced Filtering
```bash
# Search by name or description
"search": "production,staging"

# Filter by labels 
"labels": "env=prod,region=us-east"

# Query with predicate language
"query": "labels[\"environment\"] == \"production\""

# Show detailed information
"verbose": true

# List from all Teleport clusters
"all": true
```

#### Expected Output Format
```
Found 2 Kubernetes cluster(s):

• prod-east-k8s (selected)
  Labels: 2 available (use verbose=true to see details)

• staging-west-k8s
  Labels: env=staging, region=us-west-1

Tip: Use verbose=true to see detailed label information for each cluster.
```

### Kubernetes Cluster Authentication

The `teleport_kube_login` tool handles cluster authentication and kubeconfig management:

#### Single Cluster Login
```bash
# Login to specific cluster
"kubeCluster": "prod-east-k8s"

# Login with custom context name
"kubeCluster": "prod-east-k8s",
"contextName": "my-prod-context"

# Login with user/group impersonation
"kubeCluster": "prod-east-k8s",
"asUser": "admin",
"asGroups": "system:masters"

# Login with default namespace
"kubeCluster": "prod-east-k8s", 
"kubeNamespace": "production"
```

#### Batch Login (All Clusters)
```bash
# Login to all accessible clusters
"all": true

# Login to filtered clusters
"all": true,
"labels": "env=prod"

# Login with custom context template
"all": true,
"contextName": "teleport-{{.KubeName}}"
```

#### Expected Behavior
- Updates your local kubeconfig automatically
- Provides clear success/failure feedback
- Handles MFA challenges gracefully
- Supports both single and batch operations
- Only uses `tsh kube login` - no manual kubeconfig manipulation

**Important**: Either `kubeCluster` or `all=true` must be specified, but not both. 