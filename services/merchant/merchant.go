package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/kpango/glg"
)

const (
	v1Endpoint = "/v1.0/merchants"
	v2Endpoint = "/v2.0/merchants"
)

var (
	// ErrMerchantNotSpecified no merchant
	ErrMerchantNotSpecified = errors.New("merchant not specified")
	// ErrMerchantORUnavailabilityIDNotSpecified no merchant or unavailability
	ErrMerchantORUnavailabilityIDNotSpecified = errors.New("merchant or unavailability not specified")
)

type (
	// Service describes the merchant API abstraction
	Service interface {
		ListAll() ([]Merchant, error)
		Unavailabilities(merchantUUID string) (Unavailabilities, error)
		CreateUnavailabilityNow(merchantUUID, description string, pauseMinutes int32) (UnavailabilityResponse, error)
		DeleteUnavailability(merchantUUID, unavailabilityID string) error
		Availability(merchantUUID string) (AvailabilityResponse, error)
	}

	merchantService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

// New returns a new merchant service
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
	resp, status, err := m.adapter.DoRequest(http.MethodGet,
		v1Endpoint, nil, headers)
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
	glg.Info("[SDK] Merchant ListAll success")
	return ml, json.Unmarshal(resp, &ml)
}

// Unavailabilities lista indisponibilidades do merchant
func (m *merchantService) Unavailabilities(merchantUUID string) (mu Unavailabilities, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Merchant Unavailabilities: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant Unavailabilities auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities", v1Endpoint, merchantUUID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant Unavailabilities adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant Unavailabilities status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not get 'unavailabilities'", merchantUUID)
		glg.Error("[SDK] Merchant Unavailabilities err: ", err)
		return
	}
	glg.Info("[SDK] Merchant Unavailabilities success")
	return mu, json.Unmarshal(resp, &mu)
}

// CreateUnavailabilityNow cadastra indisponibilidade no merchant
func (m *merchantService) CreateUnavailabilityNow(merchantUUID, description string, pauseMinutes int32) (ur UnavailabilityResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Merchant CreateUnavailabilityNow: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilityNow auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities:now", v1Endpoint, merchantUUID)
	unv := unavailability{Description: description, Minutes: pauseMinutes}
	reader, err := httpadapter.NewJsonReader(unv)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilityNow NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := m.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilityNow adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant CreateUnavailabilityNow status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create 'unavailability'", merchantUUID)
		glg.Error("[SDK] Merchant CreateUnavailabilityNow err: ", err)
		return
	}
	return ur, json.Unmarshal(resp, &ur)
}

// DeleteUnavailability remove indisponibilidade no merchant
func (m *merchantService) DeleteUnavailability(merchantUUID, unavailabilityID string) (err error) {
	if (merchantUUID == "") || (unavailabilityID == "") {
		err = ErrMerchantORUnavailabilityIDNotSpecified
		glg.Error("[SDK] Merchant DeleteUnavailability: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant DeleteUnavailability auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities/%s", v1Endpoint, merchantUUID, unavailabilityID)
	_, status, err := m.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant DeleteUnavailability adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant DeleteUnavailability status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not delete unavailability id '%s' ", merchantUUID, unavailabilityID)
		glg.Error("[SDK] Merchant DeleteUnavailability err: ", err)
		return
	}
	return
}

// Availability recebe o status de disponibilidade de um merchant
func (m *merchantService) Availability(merchantUUID string) (ar AvailabilityResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Merchant Availability: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant Availability auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("/merchant%s/%s/availabilities", v2Endpoint, merchantUUID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant Availability adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant Availability status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not get availability", merchantUUID)
		glg.Error("[SDK] Merchant Availability err: ", err)
		return
	}
	return ar, json.Unmarshal(resp, &ar)
}
