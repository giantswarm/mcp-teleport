package auth

import (
	"context"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// RegisterAuthTools registers authentication-related tools with the MCP server
func RegisterAuthTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
	// teleport_login tool
	loginTool := mcp.NewTool("teleport_login",
		mcp.WithDescription("Login to a Teleport cluster"),
		mcp.WithString("loginParam",
			mcp.Description("Remote host login"),
		),
		mcp.WithString("proxyParam",
			mcp.Description("Teleport proxy address"),
		),
		mcp.WithString("userParam",
			mcp.Description("Teleport user, defaults to current local user"),
		),
		mcp.WithString("ttlParam",
			mcp.Description("Minutes to live for a session"),
		),
		mcp.WithString("identityParam",
			mcp.Description("Identity file"),
		),
		mcp.WithBoolean("insecureParam",
			mcp.Description("Do not verify server's certificate and host name. Use only in test environments"),
		),
		mcp.WithBoolean("debugParam",
			mcp.Description("Verbose logging to stdout"),
		),
	)

	s.AddTool(loginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleLogin(ctx, request, sc)
	})

	// teleport_status tool
	statusTool := mcp.NewTool("teleport_status",
		mcp.WithDescription("Display the list of proxy servers and retrieved certificates"),
		mcp.WithString("loginParam",
			mcp.Description("Remote host login"),
		),
		mcp.WithString("proxyParam",
			mcp.Description("Teleport proxy address"),
		),
		mcp.WithString("userParam",
			mcp.Description("Teleport user, defaults to current local user"),
		),
		mcp.WithString("ttlParam",
			mcp.Description("Minutes to live for a session"),
		),
		mcp.WithString("identityParam",
			mcp.Description("Identity file"),
		),
		mcp.WithBoolean("insecureParam",
			mcp.Description("Do not verify server's certificate and host name. Use only in test environments"),
		),
		mcp.WithBoolean("debugParam",
			mcp.Description("Verbose logging to stdout"),
		),
	)

	s.AddTool(statusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleStatus(ctx, request, sc)
	})

	// teleport_list_clusters tool
	listClustersTool := mcp.NewTool("teleport_list_clusters",
		mcp.WithDescription("List available Teleport clusters"),
		mcp.WithString("loginParam",
			mcp.Description("Remote host login"),
		),
		mcp.WithString("proxyParam",
			mcp.Description("Teleport proxy address"),
		),
		mcp.WithString("userParam",
			mcp.Description("Teleport user, defaults to current local user"),
		),
		mcp.WithString("ttlParam",
			mcp.Description("Minutes to live for a session"),
		),
		mcp.WithString("identityParam",
			mcp.Description("Identity file"),
		),
		mcp.WithBoolean("insecureParam",
			mcp.Description("Do not verify server's certificate and host name. Use only in test environments"),
		),
		mcp.WithBoolean("debugParam",
			mcp.Description("Verbose logging to stdout"),
		),
	)

	s.AddTool(listClustersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListClusters(ctx, request, sc)
	})

	return nil
} 