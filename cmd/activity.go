package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var activityProject string

var activityCmd = &cobra.Command{
	Use:   "activity <state>",
	Short: "Set activity indicator on workspace title",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state := args[0]
		if _, ok := ws.Indicators[state]; !ok {
			keys := make([]string, 0, len(ws.Indicators))
			for k := range ws.Indicators {
				keys = append(keys, k)
			}
			return fmt.Errorf("invalid state: %s (use: %s)", state, strings.Join(keys, ", "))
		}
		wsPath, data, err := loadWorkspace()
		if err != nil {
			return err
		}
		var folder map[string]any
		if activityProject != "" {
			folder = ws.FindFolder(data, activityProject)
		} else {
			cwd, _ := os.Getwd()
			_, folder = ws.GetProjectFromCwd(cwd)
			if folder != nil {
				folder = ws.FindFolder(data, ws.GetFolderField(folder, "path"))
			}
		}

		if folder == nil {
			return fmt.Errorf("could not find project")
		}

		originalName := ws.GetFolderField(folder, "name")
		ws.UpdateActivityIndicator(folder, state)

		if ws.GetFolderField(folder, "name") != originalName {
			if err := ws.WriteWorkspace(wsPath, data); err != nil {
				return err
			}
			fmt.Printf("Updated: %s\n", ws.GetFolderField(folder, "name"))
		} else {
			fmt.Printf("No change: %s\n", ws.GetFolderField(folder, "name"))
		}
		return nil
	},
}

func init() {
	activityCmd.Flags().StringVarP(&activityProject, "project", "p", "", "project name/path (auto-detects from cwd)")
	rootCmd.AddCommand(activityCmd)
}
