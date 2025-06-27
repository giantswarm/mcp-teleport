package teleport

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Client wraps the tsh CLI for executing Teleport commands
type Client struct {
	dryRun    bool
	debugMode bool
}

// NewClient creates a new Teleport client
func NewClient(dryRun, debugMode bool) *Client {
	return &Client{
		dryRun:    dryRun,
		debugMode: debugMode,
	}
}

// ExecutionResult represents the result of a command execution
type ExecutionResult struct {
	Success      bool   `json:"success"`
	Output       string `json:"output"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	StatusCode   int    `json:"statusCode,omitempty"`
}

// ExecuteCommand executes a tsh command with the given arguments
func (c *Client) ExecuteCommand(command string, args []string) *ExecutionResult {
	// Build the full command
	cmdArgs := []string{command}
	cmdArgs = append(cmdArgs, args...)
	
	fullCommand := fmt.Sprintf("tsh %s", strings.Join(cmdArgs, " "))

	if c.dryRun {
		return &ExecutionResult{
			Success:    true,
			Output:     fmt.Sprintf("DRY RUN: Would execute: %s", fullCommand),
			StatusCode: 0,
		}
	}

	// Create context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute the command
	cmd := exec.CommandContext(ctx, "tsh", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		var statusCode int
		if exitError, ok := err.(*exec.ExitError); ok {
			statusCode = exitError.ExitCode()
		} else {
			statusCode = 1
		}

		// If context was cancelled, it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return &ExecutionResult{
				Success:      false,
				Output:       string(output),
				ErrorMessage: fmt.Sprintf("Command timeout after 30 seconds: %s", err.Error()),
				StatusCode:   statusCode,
			}
		}

		return &ExecutionResult{
			Success:      false,
			Output:       string(output),
			ErrorMessage: err.Error(),
			StatusCode:   statusCode,
		}
	}

	return &ExecutionResult{
		Success:    true,
		Output:     string(output),
		StatusCode: 0,
	}
}

// FormatArgs formats command arguments from parameters
func FormatArgs(params map[string]interface{}) []string {
	var args []string

	for key, value := range params {
		if value == nil {
			continue
		}

		// Map parameter names to tsh flag names
		flagName := mapParameterToFlag(key)
		if flagName == "" {
			continue
		}

		switch v := value.(type) {
		case bool:
			if v {
				args = append(args, fmt.Sprintf("--%s", flagName))
			}
		case string:
			if v != "" {
				if flagName == "login" {
					args = append(args, "-l", v)
				} else {
					args = append(args, fmt.Sprintf("--%s=%s", flagName, v))
				}
			}
		case int, int32, int64:
			args = append(args, fmt.Sprintf("--%s=%v", flagName, v))
		case float32, float64:
			args = append(args, fmt.Sprintf("--%s=%v", flagName, v))
		}
	}

	return args
}

// mapParameterToFlag maps parameter names to tsh flag names
func mapParameterToFlag(param string) string {
	switch param {
	case "loginParam":
		return "login"
	case "proxyParam":
		return "proxy"
	case "userParam":
		return "user"
	case "ttlParam":
		return "ttl"
	case "identityParam":
		return "identity"
	case "insecureParam":
		return "insecure"
	case "debugParam":
		return "debug"
	case "verboseParam":
		return "verbose"
	default:
		// Remove "Param" suffix if present
		if strings.HasSuffix(param, "Param") {
			return strings.ToLower(param[:len(param)-5])
		}
		return strings.ToLower(param)
	}
} 