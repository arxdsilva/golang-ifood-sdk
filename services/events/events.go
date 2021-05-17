package events

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
	v3Endpoint    = "/v3.0/events"
	v1Endpoint    = "/v1.0/events"
	v2APIEndpoint = "/order/v1.0"
)

// ErrUnauthorized api error
var ErrUnauthorized = errors.New("Unauthorized request")

// ErrReqLimitExceeded API query limit exceeded
var ErrReqLimitExceeded = errors.New("EVENTS POLL REQUEST LIMIT EXCEEDED")

type (
	// Service describes the event abstraction
	Service interface {
		Poll() ([]Event, error)
		Acknowledge([]Event) (err error)
	}

	eventACK struct {
		ID string `json:"id"`
	}

	// Events is a group of Event
	Events []Event

	// Event returned by the API
	Event struct {
		Code          string                 `json:"code"`
		CorrelationID string                 `json:"correlationId"`
		CreatedAt     time.Time              `json:"createdAt"`
		ID            string                 `json:"id,omitempty"`
		Metadata      map[string]interface{} `json:"metadata,omitempty"`
	}

	V2Event struct {
		Createdat time.Time              `json:"createdAt"`
		Fullcode  string                 `json:"fullCode"`
		Metadata  map[string]interface{} `json:"metadata"`
		Code      string                 `json:"code"`
		Orderid   string                 `json:"orderId"`
		ID        string                 `json:"id"`
	}

	eventService struct {
		adapter adapters.Http
		auth    auth.Service
		v2      bool
	}
)

// New returns the event service implementation
func New(adapter adapters.Http, authService auth.Service, v2 bool) *eventService {
	return &eventService{adapter, authService, v2}
}

// Poll queries the iFood API for new events
func (ev *eventService) Poll() (ml []Event, err error) {
	err = ev.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Event auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", ev.auth.GetToken())
	endpoint := v3Endpoint + ":polling"
	resp, status, err := ev.adapter.DoRequest(
		http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Event adapter.DoRequest: ", err.Error())
		return
	}
	if status == http.StatusNotFound {
		glg.Info("[SDK] Event adapter.DoRequest No events to poll")
		return
	}
	if status == http.StatusTooManyRequests {
		err = ErrReqLimitExceeded
		glg.Warn("[SDK] Event adapter.DoRequest REQUEST LIMIT EXCEEDED")
		return
	}
	if status == http.StatusUnauthorized {
		err = ErrUnauthorized
		glg.Warn("[SDK] Event adapter.DoRequest no auth")
		return
	}
	if status != http.StatusOK {
		err = errors.New("Events could not get polled")
		glg.Errorf("[SDK] Event adapter.DoRequest status '%d' err: %s", status, err.Error())
		return
	}
	glg.Info("[SDK] Poll was successfull")
	return ml, json.Unmarshal(resp, &ml)
}

// V2Poll queries the iFood API for new events
func (ev *eventService) V2Poll() (ml []V2Event, err error) {
	err = ev.auth.Validate()
	if err != nil {
		glg.Error("[SDK] (Event V2Poll) auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", ev.auth.GetToken())
	endpoint := v2APIEndpoint + "/events:polling"
	resp, status, err := ev.adapter.DoRequest(
		http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] (Event V2Poll) adapter.DoRequest: ", err.Error())
		return
	}
	if status == http.StatusNotFound {
		glg.Info("[SDK] (Event V2Poll) adapter.DoRequest No events to poll")
		return
	}
	if status == http.StatusTooManyRequests {
		err = ErrReqLimitExceeded
		glg.Warn("[SDK] (Event V2Poll) adapter.DoRequest REQUEST LIMIT EXCEEDED")
		return
	}
	if status == http.StatusUnauthorized {
		err = ErrUnauthorized
		glg.Warn("[SDK] (Event V2Poll) adapter.DoRequest no auth")
		return
	}
	if status != http.StatusOK {
		err = errors.New("Events could not get polled")
		glg.Errorf("[SDK] (Event V2Poll) adapter.DoRequest status '%d' err: %s", status, err.Error())
		return
	}
	glg.Info("[SDK] (V2Poll) was successfull")
	return ml, json.Unmarshal(resp, &ml)
}

// Acknowledge queries the iFood API to set events as 'polled'
func (ev *eventService) Acknowledge(events []Event) (err error) {
	err = ev.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Event auth.Validate: ", err.Error())
		return
	}
	eACK := []eventACK{}
	for _, e := range events {
		eACK = append(eACK, eventACK{e.ID})
	}
	reader, err := httpadapter.NewJsonReader(eACK)
	if err != nil {
		glg.Error("[SDK] Event NewJsonReader: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Cache-Control"] = "no-cache"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", ev.auth.GetToken())
	endpoint := v1Endpoint + "/acknowledgment"
	_, status, err := ev.adapter.DoRequest(
		http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Event adapter.DoRequest: ", err.Error())
		return
	}
	if status == http.StatusUnauthorized {
		glg.Warn("[SDK] Event AUTH error status code: ", status)
		err = ErrUnauthorized
		return
	}
	if status != http.StatusOK {
		err = errors.New("Events could not get polled")
		glg.Errorf("[SDK] Event Acknowledge status '%d' err: %s", status, err.Error())
		return
	}
	glg.Info("[SDK] Acknowledge was successfull")
	return
}
