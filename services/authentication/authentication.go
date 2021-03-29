package authentication

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	"github.com/kpango/glg"
)

const (
	authEndpoint   = "/oauth/token"
	valueGrantType = "password"
)

var ErrUnauthorized = errors.New("Unauthorized")

type (
	// Service describes the auth service abstraction
	Service interface {
		Authenticate(username, password string) (*Credentials, error)
		Validate() error
		GetToken() string
	}

	// Credentials describes the API credential type
	Credentials struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		ExpiresIn   int    `json:"expires_in"`
	}

	authService struct {
		adapter                adapters.Http
		clientId, clientSecret string
		username, password     string
		currentExpiration      time.Time
		Token                  string
	}
)

// New returns an auth service implementation
func New(adapter adapters.Http, clientId, clientSecret string) *authService {
	return &authService{adapter: adapter, clientId: clientId, clientSecret: clientSecret}
}

// Authenticate queries the iFood API for a credential
func (a *authService) Authenticate(username, password string) (c *Credentials, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("client_id", a.clientId)
	writer.WriteField("client_secret", a.clientSecret)
	writer.WriteField("grant_type", valueGrantType)
	writer.WriteField("username", username)
	writer.WriteField("password", password)
	if err = writer.Close(); err != nil {
		glg.Error("[SDK] Auth writer.Close: ", err.Error())
		return
	}
	reader := bytes.NewReader(payload.Bytes())
	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Accept"] = "*/*"
	resp, status, err := a.adapter.DoRequest(http.MethodPost, authEndpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Auth adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] Auth: status code ", status)
		err = ErrUnauthorized
		return
	}
	if err = json.Unmarshal(resp, &c); err != nil {
		glg.Error("[SDK] Unmarshal: ", err)
		return
	}
	glg.Info("[SDK] Authenticate success")
	a.currentExpiration = time.Now().Add(time.Hour)
	a.username = username
	a.password = password
	a.Token = c.AccessToken
	return
}

// Validate validates or renews a token auth
func (a *authService) Validate() (err error) {
	if !time.Now().After(a.currentExpiration) {
		return
	}
	glg.Info("[SDK] Renew Auth")
	_, err = a.Authenticate(a.username, a.password)
	return
}

// GetToken returns the last valid token
func (a *authService) GetToken() (token string) {
	return a.Token
}
