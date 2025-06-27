// Package server provides the core server infrastructure for mcp-teleport.
//
// This package contains the ServerContext type which manages the server's
// lifecycle, configuration, and shared resources. It provides a centralized
// way to manage server state and configuration that can be shared across
// all tools and handlers.
//
// # Key Components
//
// ServerContext: The main server context that holds configuration and shared state.
// It provides methods for accessing configuration like dry-run mode, debug mode,
// and non-destructive mode.
//
// Logger interface: A structured logging interface that can be implemented
// by different logging backends.
//
// # Usage
//
// The ServerContext is created once during server startup and passed to all
// tool handlers. It provides access to configuration and shared resources:
//
//	ctx, err := server.NewServerContext(context.Background(),
//	    server.WithDebugMode(true),
//	    server.WithDryRun(false),
//	)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Shutdown()
//
// Tools can then access configuration through the context:
//
//	if ctx.IsDryRun() {
//	    // Simulate the operation
//	} else {
//	    // Execute the real operation
//	}
package server
