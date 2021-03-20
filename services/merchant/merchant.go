package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
		Unavailabilities(merchantUUID string) (Unavailabilities, error)
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

	Unavailabilities []Unavailability

	Unavailability struct {
		ID          string    `json:"id"`
		Storeid     string    `json:"storeId"`
		Description string    `json:"description"`
		Authorid    string    `json:"authorId"`
		Start       time.Time `json:"start"`
		End         time.Time `json:"end"`
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
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant ListAll auth.Validate: ", err.Error())
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

func (m *merchantService) Unavailabilities(merchantUUID string) (mu Unavailabilities, err error) {
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant Unavailabilities auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities", merchantsV1Endpoint, merchantUUID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant Unavailabilities adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] Merchant Unavailabilities status code: ", status)
		err = ErrBadRequest
		return
	}
	return mu, json.Unmarshal(resp, &mu)
}
