package ssh

import (
	"context"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// RegisterSSHTools registers SSH-related tools with the MCP server
func RegisterSSHTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
	// teleport_list_ssh_nodes tool
	listSSHNodesTool := mcp.NewTool("teleport_list_ssh_nodes",
		mcp.WithDescription("List SSH nodes available through Teleport"),
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

	s.AddTool(listSSHNodesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListSSHNodes(ctx, request, sc)
	})

	// teleport_ssh tool
	sshTool := mcp.NewTool("teleport_ssh",
		mcp.WithDescription("Run shell or execute a command on a remote SSH node"),
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
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("Remote host to connect to"),
		),
		mcp.WithString("command",
			mcp.Description("Command to execute on the remote host"),
		),
		mcp.WithNumber("port",
			mcp.Description("SSH port on the remote host"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Verbose output"),
		),
		mcp.WithBoolean("forwardAgent",
			mcp.Description("Forward SSH agent to the remote host"),
		),
	)

	s.AddTool(sshTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleSSH(ctx, request, sc)
	})

	return nil
} 