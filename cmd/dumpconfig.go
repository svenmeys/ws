package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var dumpSection string

var dumpConfigCmd = &cobra.Command{
	Use:   "dump-config",
	Short: "Dump workspace config as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, data, err := loadWorkspace()
		if err != nil {
			return err
		}
		folders := ws.GetFolders(data)
		wsDir := filepath.Dir(wsPath)

		projects := make([]map[string]any, len(folders))
		for i, f := range folders {
			relPath := ws.GetFolderField(f, "path")
			absPath, _ := filepath.Abs(filepath.Join(wsDir, relPath))
			projects[i] = map[string]any{
				"name":          ws.GetFolderField(f, "name"),
				"path":          absPath,
				"slack_channel": ws.GetXPAField(f, "slack_channel"),
				"description":   ws.GetXPAField(f, "description"),
			}
		}

		config := map[string]any{
			"workspace_path":   wsPath,
			"projects":         projects,
			"channel_mappings": ws.GetChannelMappings(data),
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		switch dumpSection {
		case "channels":
			return enc.Encode(config["channel_mappings"])
		case "projects":
			return enc.Encode(config["projects"])
		case "":
			return enc.Encode(config)
		default:
			return fmt.Errorf("unknown section: %s (use: channels, projects)", dumpSection)
		}
	},
}

func init() {
	dumpConfigCmd.Flags().StringVar(&dumpSection, "section", "", "only dump specific section (channels, projects)")
	rootCmd.AddCommand(dumpConfigCmd)
}
