package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func (c *Cli) setupCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <url>",
		Short: "Clone a repository",
		Long:  "Clone a repository from a remote source to your local machine.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			branch, _ := cmd.Flags().GetString("branch")
			slog.Debug("Cloning repository", "url", url, "branch", branch)
		},
	}
	cmd.Flags().StringP("branch", "b", "main", "Specify the branch to clone")
	return cmd
}
