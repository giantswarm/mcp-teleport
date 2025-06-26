package database

import (
	"github.com/giantswarm/mcp-teleport/internal/server"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// RegisterDatabaseTools registers database-related tools with the MCP server
func RegisterDatabaseTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
	// TODO: Implement database tools
	return nil
} 