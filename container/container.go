package container

import (
	"net/http"
	"os"
	"strconv"
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

// Container is the SDK abstractions holder to facilitate the API manipulation
type Container struct {
	env             int
	v2              bool
	timeout         time.Duration
	httpadapter     adapters.Http
	AuthService     authentication.Service
	MerchantService merchant.Service
	CatalogService  catalog.Service
	EventsService   events.Service
	OrdersService   orders.Service
}

// New returns a new container
func New(env int, timeout time.Duration, v2 bool) *Container {
	return &Container{env: env, timeout: timeout, v2: v2}
}

func Create(clientId, clientSecret string, env int, v2 bool) (c *Container) {
	c = New(env, time.Minute, v2)
	c.GetHttpAdapter()
	c.GetAuthenticationService(clientId, clientSecret)
	c.GetMerchantService()
	c.GetCatalogService()
	c.GetEventsService()
	c.GetOrdersService()
	return
}

// CreateFromEnvs creates a new instance of the container struct
// from envs
// 		"IFOOD_CLIENT_ID"
// 		"IFOOD_CLIENT_SECRET"
//		"IFOOD_ENV" (default to Production)
//		always uses the api v2
func CreateFromEnvs() (c *Container) {
	clientID := os.Getenv("IFOOD_CLIENT_ID")
	clientSecret := os.Getenv("IFOOD_CLIENT_SECRET")
	env, _ := strconv.Atoi(os.Getenv("IFOOD_ENV"))
	return Create(clientID, clientSecret, env, true)
}

// GetHttpAdapter returns new HTTP adapter according to the env
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
		if c.v2 {
			c.httpadapter = httpadapter.New(client, v2urlProduction)
			return c.httpadapter
		}
		c.httpadapter = httpadapter.New(client, urlProduction)
	case EnvSandBox:
		if c.v2 {
			c.httpadapter = httpadapter.New(client, v2urlProduction)
			return c.httpadapter
		}
		c.httpadapter = httpadapter.New(client, urlSandbox)
	}
	return c.httpadapter
}

// Do a start method to instantiate all services instead of each separated

// GetAuthenticationService instantiates an auth service, also adds it to the container
func (c *Container) GetAuthenticationService(clientId, clientSecret string) authentication.Service {
	if c.httpadapter == nil {
		glg.Warn("[GetAuthenticationService]: http adapter is nil, please set it with Container.GetHttpAdapter")
		return nil
	}
	if c.AuthService == nil {
		c.AuthService = authentication.New(c.GetHttpAdapter(), clientId, clientSecret, c.v2)
	}
	return c.AuthService
}

// GetMerchantService instantiates an merchant service, also adds it to the container
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

// GetCatalogService instantiates an catalog service, also adds it to the container
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

// GetEventsService instantiates an events service, also adds it to the container
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
		c.EventsService = events.New(c.GetHttpAdapter(), c.AuthService, true)
	}
	return c.EventsService
}

// GetOrdersService instantiates an orders service, also adds it to the container
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
