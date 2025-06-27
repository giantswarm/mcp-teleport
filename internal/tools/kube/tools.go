package kube

import (
	"context"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// RegisterKubeTools registers Kubernetes-related tools with the MCP server
func RegisterKubeTools(s *mcpserver.MCPServer, sc *server.ServerContext) error {
	// teleport_kube_list_clusters tool
	listClustersTool := mcp.NewTool("teleport_kube_list_clusters",
		mcp.WithDescription("Get a list of Kubernetes clusters available through Teleport"),
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
		// Enhanced parameters for Kubernetes cluster discovery
		mcp.WithString("search",
			mcp.Description("List of comma separated search keywords or phrases enclosed in quotations (e.g. foo,bar,\"some phrase\")"),
		),
		mcp.WithString("query",
			mcp.Description("Query by predicate language enclosed in single quotes. Supports ==, !=, &&, and || (e.g. 'labels[\"key1\"] == \"value1\" && labels[\"key2\"] != \"value2\"')"),
		),
		mcp.WithString("labels",
			mcp.Description("List of comma separated labels to filter by (e.g. key1=value1,key2=value2)"),
		),
		mcp.WithBoolean("verbose",
			mcp.Description("Show an untruncated list of labels and detailed cluster information"),
		),
		mcp.WithBoolean("all",
			mcp.Description("List Kubernetes clusters from all clusters and proxies"),
		),
		mcp.WithString("cluster",
			mcp.Description("Specify the Teleport cluster to connect"),
		),
		mcp.WithBoolean("quiet",
			mcp.Description("Quiet mode"),
		),
	)

	s.AddTool(listClustersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleKubeListClusters(ctx, request, sc)
	})

	// teleport_kube_login tool
	loginTool := mcp.NewTool("teleport_kube_login",
		mcp.WithDescription("Login to a Kubernetes cluster via Teleport. Updates kubeconfig to enable kubectl access to the specified cluster."),
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
		mcp.WithString("cluster",
			mcp.Description("Specify the Teleport cluster to connect"),
		),
		// Kubernetes cluster login specific parameters
		mcp.WithString("kubeCluster",
			mcp.Description("Name of the Kubernetes cluster to login to. Check 'teleport_kube_list_clusters' for available clusters. If not specified with --all, you'll be prompted to select."),
		),
		mcp.WithString("labels",
			mcp.Description("List of comma separated labels to filter clusters for batch login (e.g. key1=value1,key2=value2). Used with --all."),
		),
		mcp.WithString("query",
			mcp.Description("Query by predicate language for filtering clusters in batch login. Supports ==, !=, &&, and || (e.g. 'labels[\"key1\"] == \"value1\"'). Used with --all."),
		),
		mcp.WithString("asUser",
			mcp.Description("Configure custom Kubernetes user impersonation"),
		),
		mcp.WithString("asGroups",
			mcp.Description("Configure custom Kubernetes group impersonation"),
		),
		mcp.WithString("kubeNamespace",
			mcp.Description("Configure the default Kubernetes namespace"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Generate a kubeconfig with every cluster the user has access to. Mutually exclusive with kubeCluster."),
		),
		mcp.WithString("contextName",
			mcp.Description("Define a custom context name. To use it with --all include \"{{.KubeName}}\""),
		),
		mcp.WithString("requestReason",
			mcp.Description("Reason for requesting access"),
		),
		mcp.WithBoolean("disableAccessRequest",
			mcp.Description("Disable automatic resource access requests"),
		),
	)

	s.AddTool(loginTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleKubeLogin(ctx, request, sc)
	})

	return nil
}
