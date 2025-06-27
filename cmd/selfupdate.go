package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

// newSelfUpdateCmd creates the selfupdate command
func newSelfUpdateCmd() *cobra.Command {
	var (
		checkOnly bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "selfupdate",
		Short: "Update mcp-teleport to the latest version",
		Long: `Update mcp-teleport to the latest version from GitHub releases.
		
By default, this command will check for updates and prompt before installing.
Use --force to skip the confirmation prompt.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelfUpdate(checkOnly, force)
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates, don't install")
	cmd.Flags().BoolVar(&force, "force", false, "Force update without confirmation")

	return cmd
}

// runSelfUpdate handles the self-update logic
func runSelfUpdate(checkOnly, force bool) error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("giantswarm/mcp-teleport"))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}

	currentVersion := rootCmd.Version
	if currentVersion == "dev" {
		log.Println("Development version detected, checking for updates...")
	}

	if !found {
		fmt.Printf("No release found for repository giantswarm/mcp-teleport\n")
		return nil
	}

	if latest.LessOrEqual(currentVersion) && currentVersion != "dev" {
		fmt.Printf("Current version %s is the latest\n", currentVersion)
		return nil
	}

	fmt.Printf("Found new version %s (current: %s)\n", latest.Version(), currentVersion)

	if checkOnly {
		return nil
	}

	if !force {
		fmt.Print("Do you want to update? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			fmt.Println("Update cancelled")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %w", err)
	}

	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	fmt.Printf("Successfully updated to version %s\n", latest.Version())
	return nil
}
