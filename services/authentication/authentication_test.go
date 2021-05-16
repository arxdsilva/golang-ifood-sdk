package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	"github.com/stretchr/testify/assert"
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
