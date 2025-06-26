package ssh

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/giantswarm/mcp-teleport/internal/teleport"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleListSSHNodes handles the teleport_list_ssh_nodes tool
func handleListSSHNodes(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Format arguments
	args := teleport.FormatArgs(params)

	// Execute ls command
	result := client.ExecuteCommand("ls", args)

	// Build MCP response
	var content []mcp.Content
	if !result.Success {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Error: %s\n%s", result.ErrorMessage, result.Output),
		})
		return &mcp.CallToolResult{
			Content: content,
			IsError: true,
		}, nil
	}

	content = append(content, mcp.TextContent{
		Type: "text",
		Text: result.Output,
	})

	return &mcp.CallToolResult{
		Content: content,
	}, nil
}

// handleSSH handles the teleport_ssh tool
func handleSSH(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Validate required destination parameter
	destination, ok := params["destination"].(string)
	if !ok || destination == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Destination host is required",
				},
			},
			IsError: true,
		}, nil
	}

	// Extract command if provided
	command, _ := params["command"].(string)

	// Build SSH arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Add destination
	args = append(args, destination)

	// Add command if provided
	if command != "" {
		args = append(args, command)
	}

	// Handle additional SSH-specific parameters
	if port, ok := params["port"].(float64); ok {
		args = append(args, fmt.Sprintf("--port=%d", int(port)))
	}

	if verbose, ok := params["verbose"].(bool); ok && verbose {
		args = append(args, "--verbose")
	}

	if forwardAgent, ok := params["forwardAgent"].(bool); ok && forwardAgent {
		args = append(args, "--forward-agent")
	}

	// Execute SSH command
	result := client.ExecuteCommand("ssh", args)

	// Build MCP response
	var content []mcp.Content
	if !result.Success {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Error: %s\n%s", result.ErrorMessage, result.Output),
		})
		return &mcp.CallToolResult{
			Content: content,
			IsError: true,
		}, nil
	}

	content = append(content, mcp.TextContent{
		Type: "text",
		Text: result.Output,
	})

	return &mcp.CallToolResult{
		Content: content,
	}, nil
} 