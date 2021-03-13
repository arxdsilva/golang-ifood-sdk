package authentication

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	"github.com/kpango/glg"
)

const (
	authEndpoint   = "/oauth/token"
	valueGrantType = "password"
)

var ErrUnauthorized = errors.New("Unauthorized")

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
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		ExpiresIn   int    `json:"expires_in"`
	}

	authService struct {
		adapter                adapters.Http
		clientId, clientSecret string
	}
)

func New(adapter adapters.Http, clientId, clientSecret string) *authService {
	return &authService{adapter, clientId, clientSecret}
}

func (a *authService) Authenticate(username, password string) (c *Credentials, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("client_id", a.clientId)
	writer.WriteField("client_secret", a.clientSecret)
	writer.WriteField("grant_type", valueGrantType)
	writer.WriteField("username", username)
	writer.WriteField("password", password)
	if err = writer.Close(); err != nil {
		return
	}
	reader := bytes.NewReader(payload.Bytes())
	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Accept"] = "*/*"
	resp, status, err := a.adapter.DoRequest(http.MethodPost, authEndpoint, reader, headers)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		glg.Info("[SDK] Auth: status code ", status)
		err = ErrUnauthorized
		return
	}
	return c, json.Unmarshal(resp, &c)
}
