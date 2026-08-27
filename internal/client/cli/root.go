package cli

import "github.com/nipalab/nipa/internal/client/usecase"

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
