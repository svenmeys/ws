package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var resolveJSON bool

var resolveChannelCmd = &cobra.Command{
	Use:   "resolve-channel <channel_id>",
	Short: "Resolve a Slack channel ID to its project path",
	Long:  "Scans all workspaces to find the project mapped to the channel. Returns the absolute path.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID := args[0]
		wsDir := ws.GetWorkspacesDir()

		matches, err := filepath.Glob(filepath.Join(wsDir, "*.code-workspace"))
		if err != nil || len(matches) == 0 {
			if resolveJSON {
				json.NewEncoder(os.Stdout).Encode(map[string]string{
					"error":   "workspaces directory not found",
					"channel": channelID,
				})
			}
			return fmt.Errorf("workspaces directory not found: %s", wsDir)
		}
		sort.Strings(matches)

		for _, wsPath := range matches {
			data, err := ws.ReadWorkspace(wsPath)
			if err != nil {
				continue
			}

			mappings := ws.GetChannelMappings(data)
			if relPath, ok := mappings[channelID]; ok {
				absPath, _ := filepath.Abs(filepath.Join(filepath.Dir(wsPath), relPath))
				if resolveJSON {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "  ")
					return enc.Encode(map[string]string{
						"channel":   channelID,
						"path":      absPath,
						"workspace": wsPath,
					})
				}
				fmt.Println(absPath)
				return nil
			}
		}

		if resolveJSON {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error":   "channel not found",
				"channel": channelID,
			})
		}
		return fmt.Errorf("channel not found: %s", channelID)
	},
}

func init() {
	resolveChannelCmd.Flags().BoolVarP(&resolveJSON, "json", "j", false, "output as JSON")
	rootCmd.AddCommand(resolveChannelCmd)
}
