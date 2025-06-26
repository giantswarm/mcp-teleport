package teleport

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(true, false)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if !client.dryRun {
		t.Error("Expected dry run to be true")
	}
	if client.debugMode {
		t.Error("Expected debug mode to be false")
	}
}

func TestFormatArgs(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		expected []string
	}{
		{
			name:     "empty params",
			params:   map[string]interface{}{},
			expected: []string{},
		},
		{
			name: "proxy parameter",
			params: map[string]interface{}{
				"proxyParam": "teleport.example.com",
			},
			expected: []string{"--proxy=teleport.example.com"},
		},
		{
			name: "boolean parameter",
			params: map[string]interface{}{
				"debugParam": true,
			},
			expected: []string{"--debug"},
		},
		{
			name: "login parameter",
			params: map[string]interface{}{
				"loginParam": "alice",
			},
			expected: []string{"-l", "alice"},
		},
		{
			name: "multiple parameters",
			params: map[string]interface{}{
				"proxyParam": "teleport.example.com",
				"userParam":  "alice",
				"debugParam": true,
			},
			expected: []string{"--proxy=teleport.example.com", "--user=alice", "--debug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatArgs(tt.params)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d args, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			// Check that all expected args are present (order may vary due to map iteration)
			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected arg %q not found in result: %v", expected, result)
				}
			}
		})
	}
}

func TestExecuteCommandDryRun(t *testing.T) {
	client := NewClient(true, false)
	result := client.ExecuteCommand("status", []string{})

	if !result.Success {
		t.Error("Expected dry run command to succeed")
	}
	if result.StatusCode != 0 {
		t.Errorf("Expected status code 0, got %d", result.StatusCode)
	}
	if result.Output == "" {
		t.Error("Expected output for dry run command")
	}
	if result.ErrorMessage != "" {
		t.Errorf("Expected no error message for dry run, got: %s", result.ErrorMessage)
	}
}

func TestMapParameterToFlag(t *testing.T) {
	tests := []struct {
		param    string
		expected string
	}{
		{"loginParam", "login"},
		{"proxyParam", "proxy"},
		{"userParam", "user"},
		{"debugParam", "debug"},
		{"verboseParam", "verbose"},
		{"customParam", "custom"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		t.Run(tt.param, func(t *testing.T) {
			result := mapParameterToFlag(tt.param)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
} 