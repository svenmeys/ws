package cmd

import (
	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
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

func Execute() error {
	return rootCmd.Execute()
}

func resolveWorkspace() (string, error) {
	if workspacePath != "" {
		return workspacePath, nil
	}
	return ws.GetDefaultWorkspace()
}

func loadWorkspace() (string, map[string]any, error) {
	wsPath, err := resolveWorkspace()
	if err != nil {
		return "", nil, err
	}
	data, err := ws.ReadWorkspace(wsPath)
	if err != nil {
		return "", nil, err
	}
	return wsPath, data, nil
}
