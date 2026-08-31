package cli

import (
	"context"

	"github.com/nipalab/nipa/internal/client/domain"
	"github.com/spf13/cobra"
)

func (c *Cli) setupCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "clone <url>",
		Short:         "Clone a repository",
		Long:          "Clone a repository from a remote source to your local machine.",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			branch, _ := cmd.Flags().GetString("branch")
			nipaUrl, err := domain.ParseNipaUrl(url)
			if err != nil {
				return err
			}
			return c.useCase.Repo().Clone(context.Background(), nipaUrl.Host, nipaUrl.Org, nipaUrl.Project, branch, nipaUrl.Path)
		},
	}
	cmd.Flags().StringP("branch", "b", "main", "Specify the branch to clone")
	return cmd
}
