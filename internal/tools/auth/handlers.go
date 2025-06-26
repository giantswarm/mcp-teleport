package auth

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/giantswarm/mcp-teleport/internal/teleport"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleLogin handles the teleport_login tool
func handleLogin(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
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

	// Execute login command
	result := client.ExecuteCommand("login", args)

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

// handleStatus handles the teleport_status tool
func handleStatus(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
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

	// Execute status command
	result := client.ExecuteCommand("status", args)

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

// handleListClusters handles the teleport_list_clusters tool
func handleListClusters(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
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

	// Execute clusters command
	result := client.ExecuteCommand("clusters", args)

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