package ssh

import (
	"context"
	"strings"
	"testing"

	"github.com/giantswarm/mcp-teleport/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

// TestDryRunCommandGeneration tests that our handlers generate the correct tsh commands
func TestDryRunCommandGeneration(t *testing.T) {
	// Create a test server context with dry run enabled
	sc := &server.ServerContext{}
	sc.SetDryRun(true)

	tests := []struct {
		name           string
		handler        string
		params         map[string]interface{}
		expectedInCmd  []string // strings that should appear in the dry run output
	}{
		{
			name:    "enhanced ssh node listing with all parameters",
			handler: "list_ssh_nodes",
			params: map[string]interface{}{
				"search":     "web,db",
				"query":      "labels[\"env\"] == \"prod\"",
				"labels":     "env=prod,team=platform",
				"verbose":    true,
				"all":        true,
				"cluster":    "prod-cluster",
				"proxyParam": "teleport.example.com",
			},
			expectedInCmd: []string{
				"tsh ls",
				"--proxy=teleport.example.com",
				"--format json",
				"--search web,db",
				"--query labels[\"env\"] == \"prod\"",
				"--verbose",
				"--all",
				"--cluster prod-cluster",
				"env=prod,team=platform",
			},
		},
		{
			name:    "advanced ssh connection with port forwarding",
			handler: "ssh",
			params: map[string]interface{}{
				"destination":    "root@web-server",
				"command":        "ls -la",
				"localForward":   "8080:localhost:80",
				"remoteForward":  "9090:localhost:9000",
				"dynamicForward": "1080",
				"openSSHOptions": "StrictHostKeyChecking=no",
				"cluster":        "prod-cluster",
				"tty":            true,
				"proxyParam":     "teleport.example.com",
			},
			expectedInCmd: []string{
				"tsh ssh",
				"--proxy=teleport.example.com",
				"-L 8080:localhost:80",
				"-R 9090:localhost:9000",
				"-D 1080",
				"-o StrictHostKeyChecking=no",
				"--cluster prod-cluster",
				"-t",
				"root@web-server",
				"ls -la",
			},
		},
		{
			name:    "ssh with label selector targeting multiple nodes",
			handler: "ssh",
			params: map[string]interface{}{
				"destination": "root@role=worker,env=prod",
				"command":     "hostname",
				"cluster":     "prod-cluster",
			},
			expectedInCmd: []string{
				"tsh ssh",
				"--cluster prod-cluster",
				"root@role=worker,env=prod",
				"hostname",
			},
		},
		{
			name:    "scp file transfer with options",
			handler: "scp",
			params: map[string]interface{}{
				"source":             "/local/file.txt",
				"destination":        "server:~/file.txt",
				"recursive":          true,
				"preserveAttributes": true,
				"port":               22.0,
				"cluster":            "prod-cluster",
				"quiet":              true,
			},
			expectedInCmd: []string{
				"tsh scp",
				"-r",
				"-p",
				"-P 22",
				"--cluster prod-cluster",
				"-q",
				"/local/file.txt",
				"server:~/file.txt",
			},
		},
		{
			name:    "host resolution",
			handler: "resolve",
			params: map[string]interface{}{
				"host":  "web-server",
				"quiet": true,
			},
			expectedInCmd: []string{
				"tsh resolve",
				"--format json",
				"--quiet",
				"web-server",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with correct structure
			request := createTestRequest(tt.params)

			var result string
			var err error

			// Call the appropriate handler
			switch tt.handler {
			case "list_ssh_nodes":
				mcpResult, mcpErr := handleListSSHNodes(context.Background(), request, sc)
				err = mcpErr
				if mcpResult != nil && len(mcpResult.Content) > 0 {
					result = extractTextFromContent(mcpResult.Content[0])
					if result == "" {
						t.Logf("Failed to extract text from content type: %T, content: %+v", mcpResult.Content[0], mcpResult.Content[0])
					}
				}
			case "ssh":
				mcpResult, mcpErr := handleSSH(context.Background(), request, sc)
				err = mcpErr
				if mcpResult != nil && len(mcpResult.Content) > 0 {
					// Extract text from TextContent
					result = extractTextFromContent(mcpResult.Content[0])
				}
			case "scp":
				mcpResult, mcpErr := handleSCP(context.Background(), request, sc)
				err = mcpErr
				if mcpResult != nil && len(mcpResult.Content) > 0 {
					result = extractTextFromContent(mcpResult.Content[0])
				}
			case "resolve":
				mcpResult, mcpErr := handleResolve(context.Background(), request, sc)
				err = mcpErr
				if mcpResult != nil && len(mcpResult.Content) > 0 {
					result = extractTextFromContent(mcpResult.Content[0])
				}
			default:
				t.Fatalf("Unknown handler: %s", tt.handler)
			}

			// Verify no error
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			// Verify the dry run output contains expected command parts
			if result == "" {
				t.Error("Expected non-empty result from dry run")
				return
			}

			t.Logf("Dry run output: %s", result)

			// Check that expected command parts are present
			for _, expectedPart := range tt.expectedInCmd {
				if !strings.Contains(result, expectedPart) {
					t.Errorf("Expected command part %q not found in dry run output: %s", expectedPart, result)
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

// Helper function to extract text from MCP content
func extractTextFromContent(content interface{}) string {
	// Handle mcp.TextContent specifically
	if textContent, ok := content.(mcp.TextContent); ok {
		return textContent.Text
	}
	
	// Try to extract text from different possible content types
	if textContent, ok := content.(interface{ GetText() string }); ok {
		return textContent.GetText()
	}
	
	// Check if it's a struct with Text field
	if v, ok := content.(struct{ Text string }); ok {
		return v.Text
	}
	
	return ""
}

func TestFormatSSHNodesOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "empty output",
			input:    "",
			expected: "No SSH nodes found",
			hasError: false,
		},
		{
			name:     "empty array",
			input:    "[]",
			expected: "No SSH nodes found",
			hasError: false,
		},
		{
			name:     "single node with all fields",
			input:    `[{"kind":"node","version":"v2","metadata":{"name":"abc123","labels":{"env":"prod","team":"platform"}},"spec":{"hostname":"web-server","addr":"10.0.1.100","cmd_labels":{"role":{"result":"control-plane"}}}}]`,
			expected: "Found 1 SSH node(s):\n\n• web-server (10.0.1.100) [abc123]\n  Labels: env=prod, role=control-plane, team=platform\n\n",
			hasError: false,
		},
		{
			name:     "single node without labels",
			input:    `[{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-server","addr":"10.0.1.100"}}]`,
			expected: "Found 1 SSH node(s):\n\n• web-server (10.0.1.100) [abc123]\n\n",
			hasError: false,
		},
		{
			name:     "multiple nodes",
			input:    `[{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-1","addr":"10.0.1.100"}},{"kind":"node","version":"v2","metadata":{"name":"def456","labels":{"env":"prod"}},"spec":{"hostname":"db-1","addr":"10.0.1.200"}}]`,
			expected: "Found 2 SSH node(s):\n\n• web-1 (10.0.1.100) [abc123]\n\n• db-1 (10.0.1.200) [def456]\n  Labels: env=prod\n\n",
			hasError: false,
		},
		{
			name:     "node with same hostname and addr",
			input:    `[{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-server","addr":"web-server"}}]`,
			expected: "Found 1 SSH node(s):\n\n• web-server [abc123]\n\n",
			hasError: false,
		},
		{
			name:     "invalid JSON",
			input:    "invalid json",
			expected: "",
			hasError: true,
		},
		{
			name:     "malformed JSON object",
			input:    "[{invalid}]",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatSSHNodesOutput(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
				}
			}
		})
	}
}

func TestFormatResolveOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "empty output",
			input:    "",
			expected: "No resolution result",
			hasError: false,
		},
		{
			name:     "basic resolution with all fields",
			input:    `{"kind":"node","version":"v2","metadata":{"name":"abc123","labels":{"env":"prod"}},"spec":{"hostname":"web-server","addr":"10.0.1.100"}}`,
			expected: "Host resolution for: web-server\nAddress: 10.0.1.100\nNode ID: abc123\nLabels: env=prod\n",
			hasError: false,
		},
		{
			name:     "resolution without labels",
			input:    `{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-server","addr":"10.0.1.100"}}`,
			expected: "Host resolution for: web-server\nAddress: 10.0.1.100\nNode ID: abc123\n",
			hasError: false,
		},
		{
			name:     "resolution with only hostname",
			input:    `{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-server","addr":""}}`,
			expected: "Host resolution for: web-server\nNode ID: abc123\n",
			hasError: false,
		},
		{
			name:     "resolution with same hostname and addr",
			input:    `{"kind":"node","version":"v2","metadata":{"name":"abc123"},"spec":{"hostname":"web-server","addr":"web-server"}}`,
			expected: "Host resolution for: web-server\nNode ID: abc123\n",
			hasError: false,
		},
		{
			name:     "invalid JSON",
			input:    "invalid json",
			expected: "",
			hasError: true,
		},
		{
			name:     "malformed JSON object",
			input:    "{invalid}",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatResolveOutput(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
				}
			}
		})
	}
}

// TestSSHRequiresCommand tests that SSH handler requires a command parameter
func TestSSHRequiresCommand(t *testing.T) {
	// Create a test server context with dry run enabled
	sc := &server.ServerContext{}
	sc.SetDryRun(true)

	// Test case: SSH without command should fail
	request := createTestRequest(map[string]interface{}{
		"destination": "root@test-host",
		// No command provided - should trigger error
	})

	result, err := handleSSH(context.Background(), request, sc)

	// Should not return an error from the handler itself
	if err != nil {
		t.Errorf("Expected no error from handler, got: %v", err)
	}

	// Should return MCP error result
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	if len(result.Content) == 0 {
		t.Fatal("Expected content in error result")
	}

	// Extract error message
	errorText := extractTextFromContent(result.Content[0])
	expectedMsg := "Command is required. Interactive shell sessions are not supported"
	
	if !strings.Contains(errorText, expectedMsg) {
		t.Errorf("Expected error message to contain %q, got: %q", expectedMsg, errorText)
	}
}
