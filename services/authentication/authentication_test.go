package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	mock := new(mocks.HttpClientMock)
	adapter := httpadapter.New(mock, "")
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
}

func TestAuth_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/oauth/token", r.URL.Path)
			assert.Contains(t, r.Header["Content-Type"][0], "multipart")
			fmt.Fprintf(w, `{"access_token":"token","expires_in":3600}`)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
	c, err := as.Authenticate("user", "pass")
	assert.Nil(t, err)
	assert.Equal(t, "token", c.AccessToken)
}

func TestAuth_NotOK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/oauth/token", r.URL.Path)
			assert.Contains(t, r.Header["Content-Type"][0], "multipart")
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
	c, err := as.Authenticate("user", "pass")
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrUnauthorized)
}

func TestAuth_BadResp(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/oauth/token", r.URL.Path)
			assert.Contains(t, r.Header["Content-Type"][0], "multipart")
			fmt.Fprintf(w, `{`)
		}),
	)
	defer ts.Close()
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
	c, err := as.Authenticate("user", "pass")
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "JSON")
}

func Test_verifyV2Inputs_ErrGrantType(t *testing.T) {
	authType := ""
	authCode := ""
	authCodeVerifier := ""
	refreshToken := ""
	err := verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken)
	assert.NotNil(t, err)
	assert.Equal(t, ErrGrantType, err)
}

func Test_verifyV2Inputs_ErrNoAuthCodeOrVerifier(t *testing.T) {
	authType := "authorization_code"
	authCode := ""
	authCodeVerifier := ""
	refreshToken := ""
	err := verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoAuthCodeOrVerifier, err)
}

func Test_verifyV2Inputs_ErrNoRefreshToken(t *testing.T) {
	authType := "refresh_token"
	authCode := "testCode"
	authCodeVerifier := "testToken"
	refreshToken := ""
	err := verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoRefreshToken, err)
}

func Test_verifyV2Inputs_OK_authorization_code(t *testing.T) {
	authType := "authorization_code"
	authCode := "testCode"
	authCodeVerifier := "testToken"
	refreshToken := ""
	err := verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken)
	assert.Nil(t, err)
}

func Test_verifyV2Inputs_OK_refresh_token(t *testing.T) {
	authType := "refresh_token"
	authCode := ""
	authCodeVerifier := ""
	refreshToken := "TOKEN"
	err := verifyV2Inputs(authType, authCode, authCodeVerifier, refreshToken)
	assert.Nil(t, err)
}

func TestV2Auth_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/oauth/token", r.URL.Path)
			require.Contains(t, r.Header["Content-Type"][0], "application/x-www-form-urlencoded")
			fmt.Fprintf(w, `{"accessToken":"token","expiresIn":3600}`)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
	c, err := as.V2Authenticate("authorization_code", "testCode", "verifier", "refresh")
	assert.Nil(t, err)
	assert.Equal(t, "token", c.AccessToken)
}

func TestV2Auth_RespNotOK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/oauth/token", r.URL.Path)
			require.Contains(t, r.Header["Content-Type"][0], "application/x-www-form-urlencoded")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": {"code": "BadRequest","message": "Invalid grant type"}}`)
		}),
	)
	defer ts.Close()
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	as := New(adapter, "client", "secret", false)
	assert.NotNil(t, as)
	c, err := as.V2Authenticate("authorization_code", "testCode", "verifier", "refresh")
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}
