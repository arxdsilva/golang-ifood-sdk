package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
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
		adapter   adapters.Http
		authToken string
	}

	Catalogs []Catalog
	Catalog  struct {
		ID         string   `json:"catalogId"`
		Context    []string `json:"context"`
		Status     string   `json:"status"`
		ModifiedAt float64  `json:"modifiedAt"`
	}
)

func New(adapter adapters.Http, authToken string) *catalogService {
	return &catalogService{adapter, authToken}
}

func (m *catalogService) ListAll(merchantID string) (ct Catalogs, err error) {
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.authToken)
	endpoint := catalogV2Endpoint + fmt.Sprintf(listAllEndpoint, merchantID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
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
