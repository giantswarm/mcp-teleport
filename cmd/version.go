package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newVersionCmd creates the version command
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number of mcp-teleport",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mcp-teleport version %s\n", rootCmd.Version)
		},
	}
} 