package container

import (
	"net/http"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	"github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/arxdsilva/golang-ifood-sdk/services/catalog"
	"github.com/arxdsilva/golang-ifood-sdk/services/events"
	"github.com/arxdsilva/golang-ifood-sdk/services/merchant"
	"github.com/arxdsilva/golang-ifood-sdk/services/orders"
	"github.com/kpango/glg"
)

type Container struct {
	env             int
	timeout         time.Duration
	httpadapter     adapters.Http
	AuthService     authentication.Service
	MerchantService merchant.Service
	CatalogService  catalog.Service
	EventsService   events.Service
	OrdersService   orders.Service
}

func New(env int, timeout time.Duration) *Container {
	return &Container{env: env, timeout: timeout}
}

func (c *Container) GetHttpAdapter() adapters.Http {
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

// Do a start method to instantiate all services instead of each separated

func (c *Container) GetAuthenticationService(clientId, clientSecret string) authentication.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetAuthenticationService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		c.AuthService = authentication.New(c.GetHttpAdapter(), clientId, clientSecret)
	}
	return c.AuthService
}

func (c *Container) GetMerchantService() merchant.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetMerchantService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		glg.Warn("[GetMerchantService]: please set the authentication service")
		return nil
	}
	if c.MerchantService == nil {
		c.MerchantService = merchant.New(c.GetHttpAdapter(), c.AuthService)
	}
	return c.MerchantService
}

func (c *Container) GetCatalogService() catalog.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetCatalogService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		glg.Warn("[GetCatalogService]: please set the authentication service")
		return nil
	}
	if c.CatalogService == nil {
		c.CatalogService = catalog.New(c.GetHttpAdapter(), c.AuthService)
	}
	return c.CatalogService
}

func (c *Container) GetEventsService() events.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetEventsService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		glg.Warn("[GetEventsService]: please set the authentication service")
		return nil
	}
	if c.EventsService == nil {
		c.EventsService = events.New(c.GetHttpAdapter(), c.AuthService)
	}
	return c.EventsService
}

func (c *Container) GetOrdersService() orders.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetOrdersService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		glg.Warn("[GetOrdersService]: please set the authentication service")
		return nil
	}
	if c.OrdersService == nil {
		c.OrdersService = orders.New(c.GetHttpAdapter(), c.AuthService)
	}
	return c.OrdersService
}
