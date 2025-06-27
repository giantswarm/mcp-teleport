package server

import (
	"context"
	"fmt"
	"sync"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// ServerContext holds the server configuration and shared resources
type ServerContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.RWMutex

	// Configuration
	nonDestructiveMode bool
	dryRun             bool
	debugMode          bool
	logger             Logger

	// Shared resources would go here (e.g., connection pools, caches)
}

// ServerOption is a functional option for configuring ServerContext
type ServerOption func(*ServerContext)

// WithNonDestructiveMode sets whether destructive operations are allowed
func WithNonDestructiveMode(enabled bool) ServerOption {
	return func(sc *ServerContext) {
		sc.nonDestructiveMode = enabled
	}
}

// WithDryRun sets whether operations should be simulated
func WithDryRun(enabled bool) ServerOption {
	return func(sc *ServerContext) {
		sc.dryRun = enabled
	}
}

// WithDebugMode sets whether debug logging is enabled
func WithDebugMode(enabled bool) ServerOption {
	return func(sc *ServerContext) {
		sc.debugMode = enabled
	}
}

// WithLogger sets the logger for the server context
func WithLogger(logger Logger) ServerOption {
	return func(sc *ServerContext) {
		sc.logger = logger
	}
}

// NewServerContext creates a new server context with the given options
func NewServerContext(ctx context.Context, opts ...ServerOption) (*ServerContext, error) {
	serverCtx, cancel := context.WithCancel(ctx)

	sc := &ServerContext{
		ctx:    serverCtx,
		cancel: cancel,
	}

	// Apply options
	for _, opt := range opts {
		opt(sc)
	}

	// Set default logger if none provided
	if sc.logger == nil {
		sc.logger = &noopLogger{}
	}

	return sc, nil
}

// Context returns the context associated with the server
func (sc *ServerContext) Context() context.Context {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.ctx
}

// IsNonDestructiveMode returns whether destructive operations are disabled
func (sc *ServerContext) IsNonDestructiveMode() bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.nonDestructiveMode
}

// IsDryRun returns whether operations should be simulated
func (sc *ServerContext) IsDryRun() bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.dryRun
}

// IsDebugMode returns whether debug logging is enabled
func (sc *ServerContext) IsDebugMode() bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.debugMode
}

// Logger returns the configured logger
func (sc *ServerContext) Logger() Logger {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.logger
}

// SetDryRun dynamically sets whether operations should be simulated
func (sc *ServerContext) SetDryRun(enabled bool) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.dryRun = enabled
}

// SetDebugMode dynamically sets whether debug logging is enabled
func (sc *ServerContext) SetDebugMode(enabled bool) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.debugMode = enabled
}

// Shutdown gracefully shuts down the server context
func (sc *ServerContext) Shutdown() error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.cancel != nil {
		sc.cancel()
		sc.cancel = nil
	}

	return nil
}

// noopLogger is a logger that does nothing
type noopLogger struct{}

func (l *noopLogger) Debug(msg string, args ...interface{}) {}
func (l *noopLogger) Info(msg string, args ...interface{})  {}
func (l *noopLogger) Warn(msg string, args ...interface{})  {}
func (l *noopLogger) Error(msg string, args ...interface{}) {}

// ExecuteCommand is a helper function for executing shell commands
func (sc *ServerContext) ExecuteCommand(command string, args []string) (string, error) {
	if sc.IsDryRun() {
		sc.Logger().Info("DRY RUN: Would execute command", "command", command, "args", args)
		return fmt.Sprintf("DRY RUN: %s %v", command, args), nil
	}

	// This would contain the actual command execution logic
	// For now, we'll just return a placeholder
	return "", fmt.Errorf("command execution not implemented yet")
}
