package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var addName, addSlack, addDesc string

var addCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add a new project to the workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, data, err := loadWorkspace()
		if err != nil {
			return err
		}
		if ws.FindFolder(data, args[0]) != nil {
			return fmt.Errorf("project with path '%s' already exists", args[0])
		}

		folder := ws.AddFolder(data, args[0], addName, addSlack, addDesc)
		if err := ws.WriteWorkspace(wsPath, data); err != nil {
			return err
		}

		fmt.Printf("Added: %s\n", folder["name"])
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "display name with emoji (required)")
	_ = addCmd.MarkFlagRequired("name")
	addCmd.Flags().StringVarP(&addSlack, "slack", "s", "", "Slack channel ID")
	addCmd.Flags().StringVarP(&addDesc, "desc", "d", "", "project description")
	rootCmd.AddCommand(addCmd)
}
