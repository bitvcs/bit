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
