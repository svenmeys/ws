package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ws "github.com/svenmeys/ws/internal/workspace"
)

var channelSet string

var channelCmd = &cobra.Command{
	Use:   "channel <project>",
	Short: "Get or set Slack channel for a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsPath, data, err := loadWorkspace()
		if err != nil {
			return err
		}
		folder := ws.FindFolder(data, args[0])
		if folder == nil {
			return fmt.Errorf("project not found: %s", args[0])
		}

		if channelSet != "" {
			xpa, ok := folder["x-pa"].(map[string]any)
			if !ok {
				xpa = make(map[string]any)
				folder["x-pa"] = xpa
			}
			xpa["slack_channel"] = channelSet
			if err := ws.WriteWorkspace(wsPath, data); err != nil {
				return err
			}
			fmt.Printf("Set channel for %s: %s\n", ws.GetFolderField(folder, "name"), channelSet)
		} else {
			channel := ws.GetXPAField(folder, "slack_channel")
			if channel != "" {
				fmt.Println(channel)
			} else {
				fmt.Printf("No channel mapped for %s\n", ws.GetFolderField(folder, "name"))
			}
		}
		return nil
	},
}

func init() {
	channelCmd.Flags().StringVarP(&channelSet, "set", "s", "", "set channel ID")
	rootCmd.AddCommand(channelCmd)
}
