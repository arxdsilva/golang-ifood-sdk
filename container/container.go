package container

import (
	"net/http"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	"github.com/arxdsilva/golang-ifood-sdk/authentication"
	"github.com/arxdsilva/golang-ifood-sdk/httpadapter"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
)

type container struct {
	env         int
	timeout     time.Duration
	httpadapter adapters.Http
	authService authentication.Service
}

func New(env int, timeout time.Duration) *container {
	return &container{env: env, timeout: timeout}
}

func (c *container) GetHttpAdapter() adapters.Http {
	if c.httpadapter != nil {
		return c.httpadapter
	}
	client := &http.Client{
		Timeout: c.timeout,
	}
	switch c.env {
	case EnvDevelopment:
		c.httpadapter = httpadapter.New(new(mocks.HttpClientMock), "")
	case EnvProduction:
		c.httpadapter = httpadapter.New(client, urlProduction)
	case EnvSandBox:
		c.httpadapter = httpadapter.New(client, urlSandbox)
	}
	return c.httpadapter
}

func (c container) GetAuthenticationService(clientId, clientSecret string) authentication.Service {
	if c.authService == nil {
		c.authService = authentication.New(c.GetHttpAdapter(), clientId, clientSecret)
	}
	return c.authService
}
