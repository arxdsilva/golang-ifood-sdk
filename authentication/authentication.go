package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/raulinoneto/golang-ifood-sdk/container"
	"github.com/raulinoneto/golang-ifood-sdk/httpadapter"
)

const (
	authEndpoint   = "/oauth/token"
	valueGrantType = "password"
)

var ErrUnauthorized = errors.New("")

type (
	Service interface {
		Authenticate(username, password string) (*Credentials, error)
	}

	Authentication struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		GrantType    string `json:"grant_type"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}

	Credentials struct {
		Key            string
		ExpirationDate time.Time
	}

	authService struct {
		adapter                container.HttpAdapter
		clientId, clientSecret string
	}
)

func New(adapter container.HttpAdapter, clientId, clientSecret string) *authService {
	return &authService{adapter, clientId, clientSecret}
}

func (a *authService) Authenticate(username, password string) (*Credentials, error) {
	auth := Authentication{
		ClientId:     a.clientId,
		ClientSecret: a.clientSecret,
		GrantType:    valueGrantType,
		Username:     username,
		Password:     password,
	}

	reader, boundary, err := httpadapter.NewMultipartReader(auth)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	headers["Content-Type"] = fmt.Sprintf("multipart/related; boundary=%s", boundary)
	headers["Accept"] = "*/*"

	_, status, err := a.adapter.DoRequest(http.MethodPost, authEndpoint, reader, headers)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, ErrUnauthorized
	}
	return nil, nil
}
