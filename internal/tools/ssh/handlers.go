package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

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

	// Build ls command arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Always use JSON format for parsing
	args = append(args, "--format", "json")

	// Add Phase 1 enhanced parameters
	if search, ok := params["search"].(string); ok && search != "" {
		args = append(args, "--search", search)
	}

	if query, ok := params["query"].(string); ok && query != "" {
		args = append(args, "--query", query)
	}

	if verbose, ok := params["verbose"].(bool); ok && verbose {
		args = append(args, "--verbose")
	}

	if all, ok := params["all"].(bool); ok && all {
		args = append(args, "--all")
	}

	if cluster, ok := params["cluster"].(string); ok && cluster != "" {
		args = append(args, "--cluster", cluster)
	}

	// Add labels as positional arguments if provided
	if labels, ok := params["labels"].(string); ok && labels != "" {
		args = append(args, labels)
	}

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

	// Parse JSON output and format for user
	formattedOutput, err := formatSSHNodesOutput(result.Output)
	if err != nil {
		// If JSON parsing fails, return raw output
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: result.Output,
		})
	} else {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: formattedOutput,
		})
	}

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

	// MCP only supports one-time commands, not interactive sessions
	// Validate that a command is provided to prevent interactive shell sessions
	if command == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Command is required. Interactive shell sessions are not supported via MCP - you must provide a specific command to execute.",
				},
			},
			IsError: true,
		}, nil
	}

	// Build SSH arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Handle Phase 1 enhanced parameters
	if localForward, ok := params["localForward"].(string); ok && localForward != "" {
		args = append(args, "-L", localForward)
	}

	if remoteForward, ok := params["remoteForward"].(string); ok && remoteForward != "" {
		args = append(args, "-R", remoteForward)
	}

	if dynamicForward, ok := params["dynamicForward"].(string); ok && dynamicForward != "" {
		args = append(args, "-D", dynamicForward)
	}

	if openSSHOptions, ok := params["openSSHOptions"].(string); ok && openSSHOptions != "" {
		args = append(args, "-o", openSSHOptions)
	}

	if localCommand, ok := params["localCommand"].(string); ok && localCommand != "" {
		args = append(args, "--local", localCommand)
	}

	if noRemoteExec, ok := params["noRemoteExec"].(bool); ok && noRemoteExec {
		args = append(args, "-N")
	}

	if cluster, ok := params["cluster"].(string); ok && cluster != "" {
		args = append(args, "--cluster", cluster)
	}

	if logDir, ok := params["logDir"].(string); ok && logDir != "" {
		args = append(args, "--log-dir", logDir)
	}

	if tty, ok := params["tty"].(bool); ok {
		if tty {
			args = append(args, "-t")
		} else {
			args = append(args, "-T")
		}
	}

	// Add destination
	args = append(args, destination)

	// Add command if provided
	if command != "" {
		args = append(args, command)
	}

	// Handle existing SSH-specific parameters
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

// handleSCP handles the teleport_scp tool
func handleSCP(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Validate required parameters
	source, ok := params["source"].(string)
	if !ok || source == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Source path is required",
				},
			},
			IsError: true,
		}, nil
	}

	destination, ok := params["destination"].(string)
	if !ok || destination == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Destination path is required",
				},
			},
			IsError: true,
		}, nil
	}

	// Build SCP arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Handle SCP-specific parameters
	if recursive, ok := params["recursive"].(bool); ok && recursive {
		args = append(args, "-r")
	}

	if preserveAttributes, ok := params["preserveAttributes"].(bool); ok && preserveAttributes {
		args = append(args, "-p")
	}

	if quiet, ok := params["quiet"].(bool); ok && quiet {
		args = append(args, "-q")
	}

	if port, ok := params["port"].(float64); ok {
		args = append(args, "-P", fmt.Sprintf("%d", int(port)))
	}

	if cluster, ok := params["cluster"].(string); ok && cluster != "" {
		args = append(args, "--cluster", cluster)
	}

	// Add source and destination
	args = append(args, source, destination)

	// Execute SCP command
	result := client.ExecuteCommand("scp", args)

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
		Text: fmt.Sprintf("File transfer completed successfully\n%s", result.Output),
	})

	return &mcp.CallToolResult{
		Content: content,
	}, nil
}

// handleResolve handles the teleport_resolve tool
func handleResolve(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Validate required host parameter
	host, ok := params["host"].(string)
	if !ok || host == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Host is required",
				},
			},
			IsError: true,
		}, nil
	}

	// Build resolve arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Always use JSON format for parsing
	args = append(args, "--format", "json")

	// Handle quiet mode
	if quiet, ok := params["quiet"].(bool); ok && quiet {
		args = append(args, "--quiet")
	}

	// Add host
	args = append(args, host)

	// Execute resolve command
	result := client.ExecuteCommand("resolve", args)

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

	// Parse JSON output and format for user
	formattedOutput, err := formatResolveOutput(result.Output)
	if err != nil {
		// If JSON parsing fails, return raw output
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: result.Output,
		})
	} else {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: formattedOutput,
		})
	}

	return &mcp.CallToolResult{
		Content: content,
	}, nil
}

// formatSSHNodesOutput formats JSON output from tsh ls command
func formatSSHNodesOutput(jsonOutput string) (string, error) {
	if strings.TrimSpace(jsonOutput) == "" {
		return "No SSH nodes found", nil
	}

	var nodes []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &nodes); err != nil {
		return "", fmt.Errorf("failed to parse JSON output: %w", err)
	}

	if len(nodes) == 0 {
		return "No SSH nodes found", nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d SSH node(s):\n\n", len(nodes)))

	for _, node := range nodes {
		// Extract from the real Teleport node structure
		metadata, _ := node["metadata"].(map[string]interface{})
		spec, _ := node["spec"].(map[string]interface{})

		if metadata == nil || spec == nil {
			continue
		}

		// Get hostname from spec.hostname
		hostname, _ := spec["hostname"].(string)
		// Get addr from spec.addr (often empty)
		addr, _ := spec["addr"].(string)
		// Get UUID from metadata.name
		uuid, _ := metadata["name"].(string)
		// Get static labels from metadata.labels
		labels, _ := metadata["labels"].(map[string]interface{})

		result.WriteString(fmt.Sprintf("• %s", hostname))
		if addr != "" && addr != hostname {
			result.WriteString(fmt.Sprintf(" (%s)", addr))
		}
		if uuid != "" {
			result.WriteString(fmt.Sprintf(" [%s]", uuid))
		}
		result.WriteString("\n")

		// Combine static labels with dynamic labels if present
		allLabels := make(map[string]interface{})
		if labels != nil {
			for k, v := range labels {
				allLabels[k] = v
			}
		}

		// Add dynamic labels from spec.cmd_labels
		if cmdLabels, ok := spec["cmd_labels"].(map[string]interface{}); ok {
			for k, v := range cmdLabels {
				if labelData, ok := v.(map[string]interface{}); ok {
					if result, ok := labelData["result"].(string); ok {
						allLabels[k] = result
					}
				}
			}
		}

		if len(allLabels) > 0 {
			result.WriteString("  Labels: ")
			// Sort labels for consistent output
			var labelKeys []string
			for k := range allLabels {
				labelKeys = append(labelKeys, k)
			}
			sort.Strings(labelKeys)

			var labelPairs []string
			for _, k := range labelKeys {
				labelPairs = append(labelPairs, fmt.Sprintf("%s=%v", k, allLabels[k]))
			}
			result.WriteString(strings.Join(labelPairs, ", "))
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

// formatResolveOutput formats JSON output from tsh resolve command
func formatResolveOutput(jsonOutput string) (string, error) {
	if strings.TrimSpace(jsonOutput) == "" {
		return "No resolution result", nil
	}

	var resolveData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &resolveData); err != nil {
		return "", fmt.Errorf("failed to parse JSON output: %w", err)
	}

	// Extract from the real Teleport node structure
	metadata, _ := resolveData["metadata"].(map[string]interface{})
	spec, _ := resolveData["spec"].(map[string]interface{})

	if metadata == nil || spec == nil {
		return "Invalid resolution result structure", nil
	}

	// Get hostname from spec.hostname
	hostname, _ := spec["hostname"].(string)
	// Get addr from spec.addr (often empty)
	addr, _ := spec["addr"].(string)
	// Get UUID from metadata.name
	uuid, _ := metadata["name"].(string)

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Host resolution for: %s\n", hostname))
	if addr != "" && addr != hostname {
		result.WriteString(fmt.Sprintf("Address: %s\n", addr))
	}
	if uuid != "" {
		result.WriteString(fmt.Sprintf("Node ID: %s\n", uuid))
	}

	// Show labels for additional context
	if labels, ok := metadata["labels"].(map[string]interface{}); ok && len(labels) > 0 {
		result.WriteString("Labels: ")
		// Sort labels for consistent output
		var labelKeys []string
		for k := range labels {
			labelKeys = append(labelKeys, k)
		}
		sort.Strings(labelKeys)

		var labelPairs []string
		for _, k := range labelKeys {
			labelPairs = append(labelPairs, fmt.Sprintf("%s=%v", k, labels[k]))
		}
		result.WriteString(strings.Join(labelPairs, ", "))
		result.WriteString("\n")
	}

	return result.String(), nil
}
