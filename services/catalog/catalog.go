package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/kpango/glg"
)

const (
	catalogV2Endpoint = "/catalog/v2.0"
	listAllEndpoint   = "/merchants/%s/catalogs"
)

var ErrBadRequest = errors.New("Bad request")

type (
	Service interface {
		ListAll(merchantID string) (Catalogs, error)
	}

	catalogService struct {
		adapter adapters.Http
		auth    auth.Service
	}

	Catalogs []Catalog
	Catalog  struct {
		ID         string   `json:"catalogId"`
		Context    []string `json:"context"`
		Status     string   `json:"status"`
		ModifiedAt float64  `json:"modifiedAt"`
	}
)

func New(adapter adapters.Http, authService auth.Service) *catalogService {
	return &catalogService{adapter, authService}
}

func (c *catalogService) ListAll(merchantID string) (ct Catalogs, err error) {
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := catalogV2Endpoint + fmt.Sprintf(listAllEndpoint, merchantID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] Catalog ListAll status code: ", status)
		err = ErrBadRequest
		return
	}
	return ct, json.Unmarshal(resp, &ct)
}
