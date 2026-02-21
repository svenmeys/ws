package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var syncDryRun bool

var syncChannelsCmd = &cobra.Command{
	Use:   "sync-channels",
	Short: "Sync channel mappings to capabilities.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, err := resolveWorkspace()
		if err != nil {
			return err
		}

		mappings, err := ws.SyncToCapabilities(wsPath, "", syncDryRun)
		if err != nil {
			return err
		}

		if syncDryRun {
			fmt.Println("Dry run - would sync:")
		} else {
			fmt.Println("Synced to capabilities.yaml:")
		}

		for channel, path := range mappings {
			fmt.Printf("  %s: %s\n", channel, path)
		}

		if len(mappings) == 0 {
			fmt.Println("No channel mappings found in workspace")
		}
		return nil
	},
}

func init() {
	syncChannelsCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "show what would be synced")
	rootCmd.AddCommand(syncChannelsCmd)
}
