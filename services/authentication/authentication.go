package authentication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	"github.com/kpango/glg"
)

const (
	authRoot         = "/authentication/v1.0"
	authEndpoint     = "/oauth/token"
	userCodeEndpoint = "/oauth/userCode"
	valueGrantType   = "password"
)

// ErrUnauthorized API no auth error
var ErrUnauthorized = errors.New("Unauthorized")
var ErrGrantType = errors.New("Grant type is invalid, should be 'client_credentials', 'authorization_code' or 'refresh_token'")
var ErrNoRefreshToken = errors.New("Grant type 'refresh_token', should have a refresh token provided")
var ErrNoAuthCodeOrVerifier = errors.New("Grant type 'authorization_code', should have both 'authorizationCode' and 'authorizationCodeVerifier' provided")

type (
	// Service describes the auth service abstraction
	Service interface {
		V2GetUserCode() (*UserCode, error)
		Authenticate(username, password string) (*Credentials, error)
		V2Authenticate(authType, authCode, authCodeVerifier, refreshToken string) (c *V2Credentials, err error)
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

	// V2Credentials describes the API credential type
	V2Credentials struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		Type         string `json:"type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	authError struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"error"`
		} `json:"error"`
	}

	UserCode struct {
		Usercode                  string `json:"userCode"`
		AuthorizationCodeVerifier string `json:"authorizationCodeVerifier"`
		VerificationURL           string `json:"verificationUrl"`
		VerificationURLComplete   string `json:"verificationUrlComplete"`
		ExpiresIn                 int    `json:"expiresIn"`
	}

	authService struct {
		adapter                adapters.Http
		clientId, clientSecret string
		username, password     string
		currentExpiration      time.Time
		Token                  string
		refreshToken           string
		v2                     bool
	}
)

// New returns an auth service implementation
func New(adapter adapters.Http, clientId, clientSecret string, v2 bool) Service {
	return &authService{adapter: adapter, clientId: clientId, clientSecret: clientSecret, v2: v2}
}

func (a *authService) V2GetUserCode() (uc *UserCode, err error) {
	params := url.Values{}
	params.Add("clientId", a.clientId)
	body := strings.NewReader(params.Encode())
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, status, err := a.adapter.DoRequest(http.MethodPost, authRoot+userCodeEndpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] (V2GetUserCode::DoRequest) error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] (V2GetUserCode::status.check): status code ", status)
		err = ErrUnauthorized
		return
	}
	if err = json.Unmarshal(resp, &uc); err != nil {
		glg.Error("[SDK] (V2GetUserCode::Unmarshal) error: ", err)
		return
	}
	glg.Info("[SDK] V2GetUserCode success")
	return
}

// V2Authenticate queries the iFood API for a credential
func (a *authService) V2Authenticate(authType, authCode, authCodeVerifier, refreshToken string) (c *V2Credentials, err error) {
	if err = verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken); err != nil {
		glg.Error("[SDK] (V2Authenticate::verifyV2Inputs) error: ", err.Error())
		return
	}
	data := url.Values{}
	data.Set("client_id", a.clientId)
	data.Set("client_secret", a.clientSecret)
	data.Set("grant_type", valueGrantType)
	data.Set("authorizationCode", authCode)
	data.Set("authorizationCodeVerifier", authCodeVerifier)
	data.Set("refreshToken", refreshToken)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := strings.NewReader(data.Encode())
	resp, status, err := a.adapter.DoRequest(http.MethodPost, authRoot+authEndpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] (V2Authenticate::DoRequest) error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] V2Authenticate: status code ", status)
		authErr := authError{}
		json.Unmarshal(resp, &authErr)
		warn := fmt.Sprintf("[SDK] (V2GetUserCode::status.code): code '%s' message '%s'",
			authErr.Error.Code, authErr.Error.Message)
		glg.Warn(warn)
		err = ErrUnauthorized
		return
	}
	if err = json.Unmarshal(resp, &c); err != nil {
		glg.Error("[SDK] (V2Authenticate::Unmarshal) error: ", err)
		return
	}
	glg.Info("[SDK] (V2Authenticate) success")
	a.currentExpiration = time.Now().Add(time.Hour * 6)
	a.Token = c.AccessToken
	a.refreshToken = c.RefreshToken
	return
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
		glg.Error("[SDK] (Auth::writer.Close) error: ", err.Error())
		return
	}
	reader := bytes.NewReader(payload.Bytes())
	headers := make(map[string]string)
	headers["Content-Type"] = writer.FormDataContentType()
	headers["Accept"] = "*/*"
	resp, status, err := a.adapter.DoRequest(http.MethodPost, authEndpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] (Auth::DoRequest) error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] (Auth::status.code): status code ", status)
		err = ErrUnauthorized
		return
	}
	if err = json.Unmarshal(resp, &c); err != nil {
		glg.Error("[SDK] (Auth::Unmarshal) error: ", err)
		return
	}
	glg.Info("[SDK] (Authenticate) success")
	a.currentExpiration = time.Now().Add(time.Hour)
	a.username = username
	a.password = password
	a.Token = c.AccessToken
	return
}

// Validate validates or renews a token auth
func (a *authService) Validate() (err error) {
	if !time.Now().After(a.currentExpiration) {
		glg.Debug("[SDK] (auth::Validate) not time")
		return
	}
	glg.Info("[SDK] (auth::Validate) Renewing Auth")
	if a.v2 {
		_, err = a.V2Authenticate("refresh_token", "", "", a.refreshToken)
		return
	}
	_, err = a.Authenticate(a.username, a.password)
	return
}

// GetToken returns the last valid token
func (a *authService) GetToken() (token string) {
	return a.Token
}

func verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken string) (err error) {
	if (authType != "client_credentials") && (authType != "authorization_code") && (authType != "refresh_token") {
		return ErrGrantType
	}
	switch authType {
	case "authorization_code":
		if authCode == "" || authCodeVerifier == "" {
			return ErrNoAuthCodeOrVerifier
		}
	case "refresh_token":
		if refreshToken == "" {
			return ErrNoRefreshToken
		}
	}
	return
}
