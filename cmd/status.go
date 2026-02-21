package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var statusActive, statusPaused, statusBlocked, statusProgress, statusDormant bool

var statusCmd = &cobra.Command{
	Use:   "status <project>",
	Short: "Update a project's status emoji",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := map[string]bool{
			"active":   statusActive,
			"paused":   statusPaused,
			"blocked":  statusBlocked,
			"progress": statusProgress,
			"dormant":  statusDormant,
		}

		var selected string
		count := 0
		for k, v := range flags {
			if v {
				selected = k
				count++
			}
		}

		if count != 1 {
			return fmt.Errorf("select exactly one status flag")
		}

		wsPath, data, err := loadWorkspace()
		if err != nil {
			return err
		}
		ok, err := ws.UpdateFolderStatus(data, args[0], selected)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("project not found: %s", args[0])
		}

		if err := ws.WriteWorkspace(wsPath, data); err != nil {
			return err
		}

		folder := ws.FindFolder(data, args[0])
		fmt.Printf("Updated: %s\n", ws.GetFolderField(folder, "name"))
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVar(&statusActive, "active", false, "set to active (🟢)")
	statusCmd.Flags().BoolVar(&statusPaused, "paused", false, "set to paused (🟡)")
	statusCmd.Flags().BoolVar(&statusBlocked, "blocked", false, "set to blocked (🔴)")
	statusCmd.Flags().BoolVar(&statusProgress, "progress", false, "set to in progress (🔵)")
	statusCmd.Flags().BoolVar(&statusDormant, "dormant", false, "set to dormant (⚪)")
	rootCmd.AddCommand(statusCmd)
}
