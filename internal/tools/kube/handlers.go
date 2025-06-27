package kube

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

// KubeCluster represents a Kubernetes cluster from tsh kube ls JSON output
type KubeCluster struct {
	KubeClusterName string                 `json:"kube_cluster_name"`
	Labels          map[string]interface{} `json:"labels"`
	Selected        bool                   `json:"selected"`
}

// handleKubeListClusters handles the teleport_kube_list_clusters tool
func handleKubeListClusters(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Build kube ls command arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Always use JSON format for parsing
	args = append(args, "--format", "json")

	// Add enhanced parameters
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

	if quiet, ok := params["quiet"].(bool); ok && quiet {
		args = append(args, "--quiet")
	}

	// Add labels as positional arguments if provided
	if labels, ok := params["labels"].(string); ok && labels != "" {
		args = append(args, labels)
	}

	// Execute kube ls command
	result := client.ExecuteCommand("kube ls", args)

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
	formattedOutput, err := formatKubeClustersOutput(result.Output, params)
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

// handleKubeLogin handles the teleport_kube_login tool
func handleKubeLogin(ctx context.Context, request mcp.CallToolRequest, sc *server.ServerContext) (*mcp.CallToolResult, error) {
	// Create teleport client
	client := teleport.NewClient(sc.IsDryRun(), sc.IsDebugMode())

	// Extract parameters
	params := make(map[string]interface{})
	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
			params = argsMap
		}
	}

	// Build kube login command arguments
	var args []string

	// Add common parameters (proxy, user, etc.)
	commonArgs := teleport.FormatArgs(params)
	args = append(args, commonArgs...)

	// Add cluster parameter if specified
	if cluster, ok := params["cluster"].(string); ok && cluster != "" {
		args = append(args, "--cluster", cluster)
	}

	// Add Kubernetes cluster login specific parameters
	if labels, ok := params["labels"].(string); ok && labels != "" {
		args = append(args, "--labels", labels)
	}

	if query, ok := params["query"].(string); ok && query != "" {
		args = append(args, "--query", query)
	}

	if asUser, ok := params["asUser"].(string); ok && asUser != "" {
		args = append(args, "--as", asUser)
	}

	if asGroups, ok := params["asGroups"].(string); ok && asGroups != "" {
		args = append(args, "--as-groups", asGroups)
	}

	if kubeNamespace, ok := params["kubeNamespace"].(string); ok && kubeNamespace != "" {
		args = append(args, "--kube-namespace", kubeNamespace)
	}

	if all, ok := params["all"].(bool); ok && all {
		args = append(args, "--all")
	}

	if contextName, ok := params["contextName"].(string); ok && contextName != "" {
		args = append(args, "--set-context-name", contextName)
	}

	if requestReason, ok := params["requestReason"].(string); ok && requestReason != "" {
		args = append(args, "--request-reason", requestReason)
	}

	if disableAccessRequest, ok := params["disableAccessRequest"].(bool); ok && disableAccessRequest {
		args = append(args, "--disable-access-request")
	}

	// Add the Kubernetes cluster name as the final argument if specified
	if kubeCluster, ok := params["kubeCluster"].(string); ok && kubeCluster != "" {
		args = append(args, kubeCluster)
	}

	// Validate that either kubeCluster or --all is specified
	kubeCluster, hasKubeCluster := params["kubeCluster"].(string)
	all, hasAll := params["all"].(bool)

	if !hasKubeCluster && (!hasAll || !all) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: Either 'kubeCluster' must be specified for single cluster login, or 'all' must be true for batch login to all accessible clusters.",
				},
			},
			IsError: true,
		}, nil
	}

	if hasKubeCluster && kubeCluster != "" && hasAll && all {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: 'kubeCluster' and 'all' are mutually exclusive. Specify either a specific cluster name or use --all for batch login.",
				},
			},
			IsError: true,
		}, nil
	}

	// Execute kube login command
	result := client.ExecuteCommand("kube login", args)

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

	// Format success message
	var successMessage strings.Builder
	if hasAll && all {
		successMessage.WriteString("Successfully logged in to all accessible Kubernetes clusters.\n")
	} else {
		successMessage.WriteString(fmt.Sprintf("Successfully logged in to Kubernetes cluster: %s\n", kubeCluster))
	}

	successMessage.WriteString("Your kubeconfig has been updated. You can now use kubectl to interact with the cluster(s).\n\n")

	if result.Output != "" {
		successMessage.WriteString("Command output:\n")
		successMessage.WriteString(result.Output)
	}

	content = append(content, mcp.TextContent{
		Type: "text",
		Text: successMessage.String(),
	})

	return &mcp.CallToolResult{
		Content: content,
	}, nil
}

// formatKubeClustersOutput formats JSON output from tsh kube ls command
func formatKubeClustersOutput(jsonOutput string, params map[string]interface{}) (string, error) {
	if strings.TrimSpace(jsonOutput) == "" {
		return "No Kubernetes clusters found", nil
	}

	var clusters []KubeCluster
	if err := json.Unmarshal([]byte(jsonOutput), &clusters); err != nil {
		return "", fmt.Errorf("failed to parse JSON output: %w", err)
	}

	if len(clusters) == 0 {
		return "No Kubernetes clusters found", nil
	}

	// Check if verbose mode is requested
	verbose, _ := params["verbose"].(bool)

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d Kubernetes cluster(s):\n\n", len(clusters)))

	// Sort clusters by name for consistent output
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].KubeClusterName < clusters[j].KubeClusterName
	})

	for _, cluster := range clusters {
		result.WriteString(fmt.Sprintf("• %s", cluster.KubeClusterName))

		if cluster.Selected {
			result.WriteString(" (selected)")
		}
		result.WriteString("\n")

		// Show labels if available and verbose mode is enabled or if labels exist
		if cluster.Labels != nil && len(cluster.Labels) > 0 {
			if verbose {
				result.WriteString("  Labels: ")
				// Sort labels for consistent output
				var labelKeys []string
				for k := range cluster.Labels {
					labelKeys = append(labelKeys, k)
				}
				sort.Strings(labelKeys)

				var labelPairs []string
				for _, k := range labelKeys {
					labelPairs = append(labelPairs, fmt.Sprintf("%s=%v", k, cluster.Labels[k]))
				}
				result.WriteString(strings.Join(labelPairs, ", "))
				result.WriteString("\n")
			} else {
				// In non-verbose mode, just show count of labels
				result.WriteString(fmt.Sprintf("  Labels: %d available (use verbose=true to see details)\n", len(cluster.Labels)))
			}
		}
		result.WriteString("\n")
	}

	if !verbose && len(clusters) > 0 {
		result.WriteString("Tip: Use verbose=true to see detailed label information for each cluster.\n")
	}

	return result.String(), nil
}
