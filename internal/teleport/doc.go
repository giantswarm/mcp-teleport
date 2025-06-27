// Package teleport provides a wrapper around the Teleport CLI (tsh) for
// executing Teleport commands from the MCP server.
//
// This package abstracts the execution of tsh commands and provides
// structured result handling. It supports dry-run mode for testing
// and debug mode for verbose logging.
//
// # Key Components
//
// Client: The main client for executing tsh commands. It handles command
// construction, execution, and result parsing.
//
// ExecutionResult: A structured representation of command execution results
// including success status, output, and error information.
//
// # Usage
//
// Create a client and execute commands:
//
//	client := teleport.NewClient(false, true) // not dry-run, debug mode
//	result := client.ExecuteCommand("status", []string{})
//	if result.Success {
//	    fmt.Println("Command output:", result.Output)
//	} else {
//	    fmt.Println("Command failed:", result.ErrorMessage)
//	}
//
// Format arguments from parameters:
//
//	params := map[string]interface{}{
//	    "proxy": "teleport.example.com",
//	    "user": "alice",
//	    "debug": true,
//	}
//	args := teleport.FormatArgs(params)
//	result := client.ExecuteCommand("login", args)
package teleport
