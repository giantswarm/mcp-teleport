package kube

import (
	"github.com/giantswarm/mcp-teleport/internal/server"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// RegisterKubeTools registers Kubernetes-related tools with the MCP server
func RegisterKubeTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
	// TODO: Implement Kubernetes tools
	return nil
} 