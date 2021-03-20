package container

import (
	"net/http"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/authentication"
	"github.com/arxdsilva/golang-ifood-sdk/catalog"
	"github.com/arxdsilva/golang-ifood-sdk/events"
	"github.com/arxdsilva/golang-ifood-sdk/merchant"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	"github.com/kpango/glg"
)

type container struct {
	env             int
	timeout         time.Duration
	httpadapter     adapters.Http
	authService     authentication.Service
	merchantService merchant.Service
	catalogService  catalog.Service
	eventsService   events.Service
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

func (c *container) GetAuthenticationService(clientId, clientSecret string) authentication.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetAuthenticationService]: http adapter is nil, please set it with container.GetHttpAdapter")
		return nil
	}
	if c.authService == nil {
		c.authService = authentication.New(c.GetHttpAdapter(), clientId, clientSecret)
	}
	return c.authService
}

func (c *container) GetMerchantService() merchant.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetMerchantService]: http adapter is nil, please set it with container.GetHttpAdapter")
		return nil
	}
	if c.authService == nil {
		glg.Warn("[GetMerchantService]: please set the authentication service")
		return nil
	}
	if c.merchantService == nil {
		c.merchantService = merchant.New(c.GetHttpAdapter(), c.authService)
	}
	return c.merchantService
}

func (c *container) GetCatalogService() catalog.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetCatalogService]: http adapter is nil, please set it with container.GetHttpAdapter")
		return nil
	}
	if c.authService == nil {
		glg.Warn("[GetCatalogService]: please set the authentication service")
		return nil
	}
	if c.catalogService == nil {
		c.catalogService = catalog.New(c.GetHttpAdapter(), c.authService)
	}
	return c.catalogService
}

func (c *container) GetEventsService() events.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetEventsService]: http adapter is nil, please set it with container.GetHttpAdapter")
		return nil
	}
	if c.authService == nil {
		glg.Warn("[GetEventsService]: please set the authentication service")
		return nil
	}
	if c.eventsService == nil {
		c.eventsService = events.New(c.GetHttpAdapter(), c.authService)
	}
	return c.eventsService
}
