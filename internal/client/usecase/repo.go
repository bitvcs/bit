package usecase

import (
	"context"

	"github.com/nipalab/nipa/internal/client/domain"
)

type Repo struct {
	auth *Auth
}

func NewRepo(auth *Auth) *Repo {
	return &Repo{
		auth: auth,
	}
}

func (r *Repo) Clone(ctx context.Context, host, org, project, branch, path string) error {
	loggedIn, err := r.auth.isLoggedIn(host)
	if err != nil {
		return err
	}
	if !loggedIn {
		return domain.NewUserError("user is not logged in")
	}
	return nil
}

func (r *Repo) GetDefaultBranch(ctx context.Context, org, project string) (string, error) {
	return "", nil
}

func (r *Repo) GetBranchList(ctx context.Context, org, project string) ([]string, error) {
	return nil, nil
}
