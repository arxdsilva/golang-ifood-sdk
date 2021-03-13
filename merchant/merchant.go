package merchant

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
)

const (
	merchantsV1Endpoint = "/v1.0/merchants"
)

var ErrBadRequest = errors.New("Bad request")

type (
	Service interface {
		List() ([]Merchant, error)
	}

	Merchant struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	merchantService struct {
		adapter   adapters.Http
		authToken string
	}
)

func New(adapter adapters.Http, authToken string) *merchantService {
	return &merchantService{adapter, authToken}
}

func (m *merchantService) List() (ml []Merchant, err error) {
	headers := make(map[string]string)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, merchantsV1Endpoint, nil, headers)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = ErrBadRequest
		return
	}
	return ml, json.Unmarshal(resp, &ml)
}
