package cli

import (
	"github.com/nipalab/nipa/internal/client/usecase"
	"github.com/spf13/cobra"
)

type usecaseContainer interface {
	Auth() *usecase.Auth
}

type Cli struct {
	useCase usecaseContainer
}

func NewCli(useCase usecaseContainer) *Cli {
	return &Cli{
		useCase: useCase,
	}
}

func (c *Cli) Run() error {
	var rootCmd = &cobra.Command{
		Use:   "nipa",
		Short: "nipa is centralized version control system for your project",
	}
	rootCmd.AddCommand(c.setupCloneCmd())
	return rootCmd.Execute()
}
