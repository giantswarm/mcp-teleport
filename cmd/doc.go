// Package cmd provides the command-line interface for mcp-teleport.
//
// The cmd package implements the Cobra-based CLI structure for the MCP Teleport server.
// It provides commands for starting the server, checking version, and self-updating.
//
// # Available Commands
//
// The following commands are available:
//
//   - serve: Start the MCP server with various transport options (stdio, sse, streamable-http)
//   - version: Display version information
//   - selfupdate: Update the binary to the latest version
//
// # Usage
//
// Start the server with default settings (stdio transport):
//
//	mcp-teleport serve
//
// Start with SSE transport:
//
//	mcp-teleport serve --transport=sse --http-addr=:8080
//
// Check version:
//
//	mcp-teleport version
//
// Update to latest version:
//
//	mcp-teleport selfupdate
//
// # Configuration
//
// The server can be configured through command-line flags:
//
//   - Transport type: stdio (default), sse, or streamable-http
//   - HTTP address for web-based transports
//   - Debug mode for verbose logging
//   - Non-destructive mode to prevent destructive operations
//   - Dry-run mode to simulate operations without executing them
//
// # Integration
//
// This MCP server integrates with Teleport CLI (tsh) to provide AI assistants
// with capabilities to interact with Teleport clusters, including:
//
//   - Authentication and session management
//   - SSH access to nodes
//   - Kubernetes cluster operations
//   - Database connectivity
//   - Application access
//
// The server requires 'tsh' to be installed and available in the system PATH.
package cmd