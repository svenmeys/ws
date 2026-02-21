package cmd

import (
	"os"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Claude Code hook handler (reads JSON from stdin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, indicatorType, ok := ws.HandleHook(os.Stdin)
		if !ok {
			return nil
		}

		wsPath, folder := ws.GetProjectFromCwd(cwd)
		if wsPath == "" || folder == nil {
			return nil
		}

		data, err := ws.ReadWorkspace(wsPath)
		if err != nil {
			return nil
		}

		folder = ws.FindFolder(data, ws.GetFolderField(folder, "path"))
		if folder == nil {
			return nil
		}

		originalName := ws.GetFolderField(folder, "name")
		ws.UpdateActivityIndicator(folder, indicatorType)

		if ws.GetFolderField(folder, "name") != originalName {
			_ = ws.WriteWorkspace(wsPath, data)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
