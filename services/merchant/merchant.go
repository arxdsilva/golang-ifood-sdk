package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/kpango/glg"
)

const (
	V1Endpoint = "/v1.0/merchants"
)

type (
	Service interface {
		ListAll() ([]Merchant, error)
		Unavailabilities(merchantUUID string) (Unavailabilities, error)
		CreateUnavailabily(merchantUUID, description string, pauseMinutes int32) (UnavailabilityResponse, error)
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

	unavailability struct {
		Description string `json:"description"`
		Minutes     int32  `json:"minutes"`
	}

	UnavailabilityResponse struct {
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

// ListAll lista merchants cuja autenticacao tem permissao
func (m *merchantService) ListAll() (ml []Merchant, err error) {
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant ListAll auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	resp, status, err := m.adapter.DoRequest(http.MethodGet, V1Endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant ListAll adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant ListAll status code: ", status)
		err = errors.New("Could not list merchants")
		glg.Error("[SDK] Merchant ListAll err: ", err)
		return
	}
	return ml, json.Unmarshal(resp, &ml)
}

// Unavailabilities lista indisponibilidades do merchant
func (m *merchantService) Unavailabilities(merchantUUID string) (mu Unavailabilities, err error) {
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant Unavailabilities auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities", V1Endpoint, merchantUUID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant Unavailabilities adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant Unavailabilities status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create 'unavailability'", merchantUUID)
		glg.Error("[SDK] Merchant Unavailabilities err: ", err)
		return
	}
	return mu, json.Unmarshal(resp, &mu)
}

// CreateUnavailabily cadastra indisponibilidade no merchant
func (m *merchantService) CreateUnavailabily(merchantUUID, description string, pauseMinutes int32) (ur UnavailabilityResponse, err error) {
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabily auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities:now", V1Endpoint, merchantUUID)
	unv := unavailability{Description: description, Minutes: pauseMinutes}
	reader, err := httpadapter.NewJsonReader(unv)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabily NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := m.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabily adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant CreateUnavailabily status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create 'unavailability'", merchantUUID)
		glg.Error("[SDK] Merchant CreateUnavailabily err: ", err)
		return
	}
	return ur, json.Unmarshal(resp, &ur)
}
