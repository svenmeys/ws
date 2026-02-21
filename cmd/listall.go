package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var listAllJSON bool

var listAllCmd = &cobra.Command{
	Use:   "list-all",
	Short: "List all workspaces and their projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		wsDir := ws.GetWorkspacesDir()
		matches, err := filepath.Glob(filepath.Join(wsDir, "*.code-workspace"))
		if err != nil || len(matches) == 0 {
			return fmt.Errorf("no workspaces found in %s", wsDir)
		}
		sort.Strings(matches)

		var all []map[string]any
		for _, wsPath := range matches {
			data, err := ws.ReadWorkspace(wsPath)
			if err != nil {
				continue
			}
			folders := ws.GetFolders(data)
			projects := make([]map[string]any, len(folders))
			for i, f := range folders {
				projects[i] = map[string]any{
					"name":          ws.GetFolderField(f, "name"),
					"path":          ws.GetFolderField(f, "path"),
					"slack_channel": ws.GetXPAField(f, "slack_channel"),
				}
			}

			stem := strings.TrimSuffix(filepath.Base(wsPath), ".code-workspace")

			all = append(all, map[string]any{
				"workspace":     stem,
				"path":          wsPath,
				"project_count": len(folders),
				"projects":      projects,
			})
		}

		if listAllJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(all)
		}

		for _, info := range all {
			fmt.Printf("\n%s (%v projects)\n", info["workspace"], info["project_count"])
			fmt.Printf("  %s\n", info["path"])
			if projects, ok := info["projects"].([]map[string]any); ok {
				for _, p := range projects {
					channel := ""
					if ch, ok := p["slack_channel"].(string); ok && ch != "" {
						channel = " #" + ch
					}
					fmt.Printf("  • %s%s\n", p["name"], channel)
				}
			}
		}
		return nil
	},
}

func init() {
	listAllCmd.Flags().BoolVarP(&listAllJSON, "json", "j", false, "output as JSON")
	rootCmd.AddCommand(listAllCmd)
}
