package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command for the mcp-teleport application.
// It is the entry point when the application is called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "mcp-teleport",
	Short: "MCP server for Teleport operations",
	Long: `mcp-teleport is a Model Context Protocol (MCP) server that provides
tools for interacting with Teleport clusters through the tsh CLI. It offers various
capabilities including authentication, SSH access, Kubernetes operations, database
connectivity, and application access.

When run without subcommands, it starts the MCP server (equivalent to 'mcp-teleport serve').`,
	// SilenceUsage prevents Cobra from printing the usage message on errors that are handled by the application.
	// This is useful for providing cleaner error output to the user.
	SilenceUsage: true,
}

// SetVersion sets the version for the root command.
// This function is typically called from the main package to inject the application version at build time.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute is the main entry point for the CLI application.
// It initializes and executes the root command, which in turn handles subcommands and flags.
// This function is called by main.main().
func Execute() {
	// SetVersionTemplate defines a custom template for displaying the version.
	// This is used when the --version flag is invoked.
	rootCmd.SetVersionTemplate(`{{printf "mcp-teleport version %s\n" .Version}}`)

	// If no subcommand is provided, run the serve command by default
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "serve")
	}

	err := rootCmd.Execute()
	if err != nil {
		// Cobra itself usually prints the error. Exiting with a non-zero status code
		// indicates that an error occurred during execution.
		os.Exit(1)
	}
}

// init is a special Go function that is executed when the package is initialized.
// It is used here to add subcommands to the root command.
func init() {
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newSelfUpdateCmd())
	rootCmd.AddCommand(newServeCmd())
} 