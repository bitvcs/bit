package usecase

import (
	"context"
)

type Repo struct {
	auth *Auth
}

func NewRepo(auth *Auth) *Repo {
	return &Repo{
		auth: auth,
	}
}

func (r *Repo) Clone(ctx context.Context, host, org, project, branch, path, target string) error {
	err := r.auth.MakeSureLoggedIn(ctx, host)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetDefaultBranch(ctx context.Context, org, project string) (string, error) {
	return "", nil
}

func (r *Repo) GetBranchList(ctx context.Context, org, project string) ([]string, error) {
	return nil, nil
}
