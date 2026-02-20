package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate workspace file structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, err := resolveWorkspace()
		if err != nil {
			return err
		}

		data, err := ws.ReadWorkspace(wsPath)
		if err != nil {
			return err
		}

		wsDir := filepath.Dir(wsPath)
		folders := ws.GetFolders(data)
		channelIDs := make(map[string]bool)

		var errors, warnings []string

		for _, folder := range folders {
			path := ws.GetFolderField(folder, "path")
			name := ws.GetFolderField(folder, "name")

			fullPath, _ := filepath.Abs(filepath.Join(wsDir, path))
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("Path does not exist: %s", path))
			}

			channel := ws.GetXPAField(folder, "slack_channel")
			if channel != "" {
				if channelIDs[channel] {
					errors = append(errors, fmt.Sprintf("Duplicate channel ID: %s", channel))
				}
				channelIDs[channel] = true
			}

			hasStatus := false
			for _, emoji := range ws.StatusEmojis {
				if strings.Contains(name, emoji) {
					hasStatus = true
					break
				}
			}
			if !hasStatus {
				warnings = append(warnings, fmt.Sprintf("No status emoji in: %s", name))
			}
		}

		if len(errors) > 0 {
			fmt.Println("Errors:")
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
		}

		if len(warnings) > 0 {
			fmt.Println("Warnings:")
			for _, w := range warnings {
				fmt.Printf("  - %s\n", w)
			}
		}

		if len(errors) == 0 && len(warnings) == 0 {
			fmt.Printf("Workspace valid: %d projects\n", len(folders))
		}

		if len(errors) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
