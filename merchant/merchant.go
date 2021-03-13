package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	"github.com/kpango/glg"
)

const (
	merchantsV1Endpoint = "/v1.0/merchants"
)

var ErrBadRequest = errors.New("Bad request")

type (
	Service interface {
		ListAll() ([]Merchant, error)
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

func (m *merchantService) ListAll() (ml []Merchant, err error) {
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.authToken)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, merchantsV1Endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant ListAll adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] Merchant ListAll status code: ", status)
		err = ErrBadRequest
		return
	}
	return ml, json.Unmarshal(resp, &ml)
}
