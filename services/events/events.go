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
	V3Endpoint = "/v3.0/events"
)

var ErrUnauthorized = errors.New("Unauthorized request")

type (
	Service interface {
		Poll() ([]Event, error)
		Acknowledge([]Event) (err error)
	}

	eventACK struct {
		ID string `json:"id"`
	}

	Event struct {
		Code          string            `json:"code"`
		CorrelationID string            `json:"correlationId"`
		CreatedAt     time.Time         `json:"createdAt"`
		ID            string            `json:"id,omitempty"`
		Metadata      map[string]string `json:"metadata,omitempty"`
	}

	eventService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

func New(adapter adapters.Http, authService auth.Service) *eventService {
	return &eventService{adapter, authService}
}

func (ev *eventService) Poll() (ml []Event, err error) {
	err = ev.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Event auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", ev.auth.GetToken())
	endpoint := V3Endpoint + ":polling"
	resp, status, err := ev.adapter.DoRequest(
		http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Event adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Info("[SDK] Event Poll status code: ", status)
		return
	}
	return ml, json.Unmarshal(resp, &ml)
}

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
	endpoint := V3Endpoint + ":polling"
	_, status, err := ev.adapter.DoRequest(
		http.MethodGet, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Event adapter.DoRequest: ", err.Error())
		return
	}
	if status == http.StatusUnauthorized {
		glg.Info("[SDK] Event AUTH error status code: ", status)
		err = ErrUnauthorized
		return
	}
	if status != http.StatusOK {
		glg.Info("[SDK] Event ACK status code: ", status)
		return
	}
	return
}
