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
	V1Endpoint = "/v1.0/merchants"
	V2Endpoint = "/v2.0/merchants"
)

var (
	ErrMerchantNotSpecified                   = errors.New("merchant not specified")
	ErrMerchantORUnavailabilityIDNotSpecified = errors.New("merchant or unavailability not specified")
)

type (
	Service interface {
		ListAll() ([]Merchant, error)
		Unavailabilities(merchantUUID string) (Unavailabilities, error)
		CreateUnavailabilyNow(merchantUUID, description string, pauseMinutes int32) (UnavailabilityResponse, error)
		DeleteUnavailabily(merchantUUID, unavailabilityID string) error
		Availabily(merchantUUID string) (AvailabilityResponse, error)
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
	endpoint := fmt.Sprintf("%s/%s/unavailabilities", V1Endpoint, merchantUUID)
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

// CreateUnavailabilyNow cadastra indisponibilidade no merchant
func (m *merchantService) CreateUnavailabilyNow(merchantUUID, description string, pauseMinutes int32) (ur UnavailabilityResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Merchant CreateUnavailabilyNow: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilyNow auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities:now", V1Endpoint, merchantUUID)
	unv := unavailability{Description: description, Minutes: pauseMinutes}
	reader, err := httpadapter.NewJsonReader(unv)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilyNow NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := m.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Merchant CreateUnavailabilyNow adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant CreateUnavailabilyNow status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create 'unavailability'", merchantUUID)
		glg.Error("[SDK] Merchant CreateUnavailabilyNow err: ", err)
		return
	}
	return ur, json.Unmarshal(resp, &ur)
}

// DeleteUnavailabily remove indisponibilidade no merchant
func (m *merchantService) DeleteUnavailabily(merchantUUID, unavailabilityID string) (err error) {
	if (merchantUUID == "") || (unavailabilityID == "") {
		err = ErrMerchantORUnavailabilityIDNotSpecified
		glg.Error("[SDK] Merchant DeleteUnavailabily: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant DeleteUnavailabily auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/unavailabilities/%s", V1Endpoint, merchantUUID, unavailabilityID)
	_, status, err := m.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant DeleteUnavailabily adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant DeleteUnavailabily status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not delete unavailability id '%s' ", merchantUUID, unavailabilityID)
		glg.Error("[SDK] Merchant DeleteUnavailabily err: ", err)
		return
	}
	return
}

// Availabily recebe o status de disponibilidade de um merchant
func (m *merchantService) Availabily(merchantUUID string) (ar AvailabilityResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Merchant Availabily: ", err.Error())
		return
	}
	if err = m.auth.Validate(); err != nil {
		glg.Error("[SDK] Merchant Availabily auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", m.auth.GetToken())
	endpoint := fmt.Sprintf("/merchant%s/%s/availabilities", V2Endpoint, merchantUUID)
	resp, status, err := m.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Merchant Availabily adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Merchant Availabily status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not get availability", merchantUUID)
		glg.Error("[SDK] Merchant Availabily err: ", err)
		return
	}
	return ar, json.Unmarshal(resp, &ar)
}
