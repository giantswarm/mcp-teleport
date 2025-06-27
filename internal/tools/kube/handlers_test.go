package kube

import (
	"context"
	"strings"
	"testing"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleKubeListClusters(t *testing.T) {
	// Create test server context with dry run mode
	ctx := context.Background()
	sc, err := server.NewServerContext(ctx,
		server.WithDryRun(true),
		server.WithDebugMode(false),
		server.WithNonDestructiveMode(true),
	)
	if err != nil {
		t.Fatalf("Failed to create server context: %v", err)
	}
	defer sc.Shutdown()

	tests := []struct {
		name   string
		params map[string]interface{}
		wantOk bool
	}{
		{
			name:   "basic list without parameters",
			params: map[string]interface{}{},
			wantOk: true,
		},
		{
			name: "list with search parameter",
			params: map[string]interface{}{
				"search": "test",
			},
			wantOk: true,
		},
		{
			name: "list with query parameter",
			params: map[string]interface{}{
				"query": "labels[\"env\"] == \"prod\"",
			},
			wantOk: true,
		},
		{
			name: "list with labels parameter",
			params: map[string]interface{}{
				"labels": "env=prod,region=us-east",
			},
			wantOk: true,
		},
		{
			name: "list with verbose mode",
			params: map[string]interface{}{
				"verbose": true,
			},
			wantOk: true,
		},
		{
			name: "list all clusters",
			params: map[string]interface{}{
				"all": true,
			},
			wantOk: true,
		},
		{
			name: "list with specific teleport cluster",
			params: map[string]interface{}{
				"cluster": "prod-teleport",
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createTestRequest(tt.params)

			result, err := handleKubeListClusters(ctx, request, sc)
			if err != nil {
				t.Errorf("handleKubeListClusters() error = %v", err)
				return
			}

			if result == nil {
				t.Error("handleKubeListClusters() returned nil result")
				return
			}

			if tt.wantOk && result.IsError {
				t.Errorf("handleKubeListClusters() returned error when success expected")
			}

			if len(result.Content) == 0 {
				t.Error("handleKubeListClusters() returned no content")
			}
		})
	}
}

func TestHandleKubeLogin(t *testing.T) {
	// Create test server context with dry run mode
	ctx := context.Background()
	sc, err := server.NewServerContext(ctx,
		server.WithDryRun(true),
		server.WithDebugMode(false),
		server.WithNonDestructiveMode(true),
	)
	if err != nil {
		t.Fatalf("Failed to create server context: %v", err)
	}
	defer sc.Shutdown()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantOk  bool
		wantErr bool
	}{
		{
			name: "login to specific cluster",
			params: map[string]interface{}{
				"kubeCluster": "test-cluster",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "login to all clusters",
			params: map[string]interface{}{
				"all": true,
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "login with user impersonation",
			params: map[string]interface{}{
				"kubeCluster": "test-cluster",
				"asUser":      "admin",
				"asGroups":    "system:masters",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "login with custom namespace",
			params: map[string]interface{}{
				"kubeCluster":   "test-cluster",
				"kubeNamespace": "production",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "login with custom context name",
			params: map[string]interface{}{
				"kubeCluster": "test-cluster",
				"contextName": "my-context",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:    "login without cluster or all parameter",
			params:  map[string]interface{}{},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "login with both cluster and all parameters",
			params: map[string]interface{}{
				"kubeCluster": "test-cluster",
				"all":         true,
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createTestRequest(tt.params)

			result, err := handleKubeLogin(ctx, request, sc)
			if err != nil {
				t.Errorf("handleKubeLogin() error = %v", err)
				return
			}

			if result == nil {
				t.Error("handleKubeLogin() returned nil result")
				return
			}

			if tt.wantOk && result.IsError {
				t.Errorf("handleKubeLogin() returned error when success expected")
			}

			if !tt.wantOk && !result.IsError {
				t.Errorf("handleKubeLogin() returned success when error expected")
			}

			if len(result.Content) == 0 {
				t.Error("handleKubeLogin() returned no content")
			}
		})
	}
}

func TestFormatKubeClustersOutput(t *testing.T) {
	tests := []struct {
		name       string
		jsonOutput string
		params     map[string]interface{}
		wantErr    bool
		contains   []string
	}{
		{
			name:       "empty output",
			jsonOutput: "",
			params:     map[string]interface{}{},
			wantErr:    false,
			contains:   []string{"No Kubernetes clusters found"},
		},
		{
			name:       "empty array",
			jsonOutput: "[]",
			params:     map[string]interface{}{},
			wantErr:    false,
			contains:   []string{"No Kubernetes clusters found"},
		},
		{
			name: "single cluster without labels",
			jsonOutput: `[
				{
					"kube_cluster_name": "test-cluster",
					"labels": null,
					"selected": false
				}
			]`,
			params:   map[string]interface{}{},
			wantErr:  false,
			contains: []string{"Found 1 Kubernetes cluster", "test-cluster"},
		},
		{
			name: "multiple clusters with labels non-verbose",
			jsonOutput: `[
				{
					"kube_cluster_name": "prod-cluster",
					"labels": {"env": "production", "region": "us-east"},
					"selected": true
				},
				{
					"kube_cluster_name": "dev-cluster",
					"labels": {"env": "development"},
					"selected": false
				}
			]`,
			params:   map[string]interface{}{"verbose": false},
			wantErr:  false,
			contains: []string{"Found 2 Kubernetes cluster", "prod-cluster", "dev-cluster", "(selected)", "verbose=true"},
		},
		{
			name: "multiple clusters with labels verbose",
			jsonOutput: `[
				{
					"kube_cluster_name": "prod-cluster",
					"labels": {"env": "production", "region": "us-east"},
					"selected": true
				},
				{
					"kube_cluster_name": "dev-cluster",
					"labels": {"env": "development"},
					"selected": false
				}
			]`,
			params:   map[string]interface{}{"verbose": true},
			wantErr:  false,
			contains: []string{"Found 2 Kubernetes cluster", "prod-cluster", "dev-cluster", "(selected)", "env=production", "region=us-east", "env=development"},
		},
		{
			name:       "invalid json",
			jsonOutput: `{"invalid": json}`,
			params:     map[string]interface{}{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatKubeClustersOutput(tt.jsonOutput, tt.params)

			if tt.wantErr && err == nil {
				t.Errorf("formatKubeClustersOutput() expected error but got none")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("formatKubeClustersOutput() unexpected error = %v", err)
				return
			}

			if !tt.wantErr {
				for _, contains := range tt.contains {
					if !stringContains(result, contains) {
						t.Errorf("formatKubeClustersOutput() result doesn't contain expected string '%s'. Got: %s", contains, result)
					}
				}
			}
		})
	}
}

// Helper function to create test request
func createTestRequest(params map[string]interface{}) mcp.CallToolRequest {
	var request mcp.CallToolRequest
	request.Params.Arguments = params
	return request
}

// Helper function for case-insensitive string contains check
func stringContains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (haystack == needle ||
		strings.Contains(strings.ToLower(haystack), strings.ToLower(needle)))
}
