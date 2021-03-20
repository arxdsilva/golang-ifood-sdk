package merchant

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
	merchantsV1Endpoint = "/v1.0/merchants"
)

var ErrBadRequest = errors.New("Bad request")

type (
	Service interface {
		ListAll() ([]Merchant, error)
	}

	Merchant struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Phones  []string `json:"phones"`
		Address Address  `json:"address"`
	}

	Address struct {
		Formattedaddress string `json:"formattedAddress"`
		Country          string `json:"country"`
		State            string `json:"state"`
		City             string `json:"city"`
		Neighborhood     string `json:"neighborhood"`
		Streetname       string `json:"streetName"`
		Streetnumber     string `json:"streetNumber"`
		Postalcode       string `json:"postalCode"`
	}

	merchantService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

func New(adapter adapters.Http, authService auth.Service) *merchantService {
	return &merchantService{adapter, authService}
}

func (m *merchantService) ListAll() (ml []Merchant, err error) {
	err = m.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Merchant auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
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
