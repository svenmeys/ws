package cmd

import (
	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var workspacePath string

var rootCmd = &cobra.Command{
	Use:           "ws",
	Short:         "Manage VS Code workspace files",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workspacePath, "workspace", "w", "", "workspace file path")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func resolveWorkspace() (string, error) {
	if workspacePath != "" {
		return workspacePath, nil
	}
	return ws.GetDefaultWorkspace()
}
