package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/workspace-cli/internal/workspace"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, err := resolveWorkspace()
		if err != nil {
			return err
		}

		data, err := ws.ReadWorkspace(wsPath)
		if err != nil {
			return err
		}

		folders := ws.GetFolders(data)

		if listJSON {
			items := make([]map[string]any, len(folders))
			for i, f := range folders {
				items[i] = map[string]any{
					"name":          ws.GetFolderField(f, "name"),
					"path":          ws.GetFolderField(f, "path"),
					"slack_channel": ws.GetXPAField(f, "slack_channel"),
					"description":   ws.GetXPAField(f, "description"),
				}
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(items)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPATH\tSLACK\tDESCRIPTION")
		for _, f := range folders {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				ws.GetFolderField(f, "name"),
				ws.GetFolderField(f, "path"),
				orDash(ws.GetXPAField(f, "slack_channel")),
				orDash(ws.GetXPAField(f, "description")),
			)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "output as JSON")
	rootCmd.AddCommand(listCmd)
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
