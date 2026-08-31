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
		Host:         data.Host,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
	}
	jsonData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, data.Host, string(jsonData))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) LoadToken(host string) (*domain.LoginResult, error) {
	data, err := keyring.Get(serviceName, host)
	if err != nil {
		return nil, err
	}

	var token AccessToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResult{
		Host:         token.Host,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
