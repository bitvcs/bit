package securestorage

import (
	"encoding/json/v2"

	"github.com/nipalab/nipa/internal/client/domain"
	"github.com/zalando/go-keyring"
)

const serviceName = "nipa"

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) SaveToken(data *domain.LoginResult) error {
	token := AccessToken{
		Domain:       data.Domain,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
	}
	jsonData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, data.Domain, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) LoadToken(serverDomain string) (*domain.LoginResult, error) {
	data, err := keyring.Get(serviceName, serverDomain)
	if err != nil {
		return nil, err
	}

	var token AccessToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResult{
		Domain:       token.Domain,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
