package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pollAPIResponse = `[{
	"code": "PLACED",
	"correlationId": "1234567890012",
	"createdAt": "2017-05-02T16:01:16.567Z",
	"id": "abc-456-afge-451-n15484"
},
{
	"code": "CANCELLED",
	"correlationId": "9876543210123",
	"createdAt": "2017-05-02T16:01:16.567Z",
	"metadata": {
		"CANCEL_STAGE": "[PRE_CONFIRMED]",
		"CANCEL_CODE": "902",
		"CANCELLATION_OCCURRENCE": {
			"RESTAURANT": {
				"FINANCIAL_OCCURRENCE": "NA",
				"PAYMENT_TYPE": "NA"
			},
			"CONSUMER": {
				"FINANCIAL_OCCURRENCE": "NA",
				"PAYMENT_TYPE": "NA"
			},
			"LOGISTIC": {
				"FINANCIAL_OCCURRENCE": "NA",
				"PAYMENT_TYPE": "NA"
			}
		}
	}
}]`

const pollV2APIResponse = `[{
	"createdAt": "2019-09-19T13:40:11.822Z",
	"fullCode": "PLACED",
	"metadata": {},
	"code": "PLC",
	"orderId": "07110e1b-8191-4670-baed-407219481ffb",
	"id": "cd40582b-0ef2-4d52-bc7c-507fdff12e21"
}]`

func TestPoll_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, pollAPIResponse)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(events))
}

func TestPoll_BadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
}

func TestPoll_Unauthorized(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusUnauthorized)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, ErrUnauthorized, err)
}

func TestPoll_StatusTooManyRequests(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusTooManyRequests)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, ErrReqLimitExceeded, err)
}

func TestPoll_StatusNotFound(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(events))
}

func TestPoll_ValidateErr(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/events:polling", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
}

func TestAcknowledge_ValidateErr(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header["Content-Type"][0], "application/json")
			assert.Equal(t, r.Header["Cache-Control"][0], "no-cache")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, "/v1.0/events/acknowledgment", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events := Events{}
	err := json.Unmarshal([]byte(pollAPIResponse), &events)
	assert.Nil(t, err)
	err = eventsService.Acknowledge(events)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestAcknowledge_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header["Content-Type"][0], "application/json")
			assert.Equal(t, r.Header["Cache-Control"][0], "no-cache")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, "/v1.0/events/acknowledgment", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events := Events{}
	err := json.Unmarshal([]byte(pollAPIResponse), &events)
	assert.Nil(t, err)
	err = eventsService.Acknowledge(events)
	assert.Nil(t, err)
}

func TestAcknowledge_StatusUnauthorized(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header["Content-Type"][0], "application/json")
			assert.Equal(t, r.Header["Cache-Control"][0], "no-cache")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, "/v1.0/events/acknowledgment", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusUnauthorized)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events := Events{}
	err := json.Unmarshal([]byte(pollAPIResponse), &events)
	assert.Nil(t, err)
	err = eventsService.Acknowledge(events)
	assert.NotNil(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestAcknowledge_StatusRequestEntityTooLarge(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Header["Content-Type"][0], "application/json")
			assert.Equal(t, r.Header["Cache-Control"][0], "no-cache")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, "/v1.0/events/acknowledgment", r.URL.Path)
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	var events Events
	err := json.Unmarshal([]byte(pollAPIResponse), &events)
	assert.Nil(t, err)
	err = eventsService.Acknowledge(events)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not get polled")
}

func Test_V2Poll_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/events:polling", r.URL.Path)
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, pollV2APIResponse)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.V2Poll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
}

func Test_V2Poll_NotFound(t *testing.T) {
	resp := `{
	"error": {
		"code": "string",
		"field": "string",
		"details": [null],
		"message": "string"
	}
}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/events:polling", r.URL.Path)
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.V2Poll()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(events))
}

func Test_V2Poll_Forbidden(t *testing.T) {
	resp := `{
	"error": {
		"code": "string",
		"field": "string",
		"details": [null],
		"message": "string"
	}
}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/events:polling", r.URL.Path)
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.V2Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, ErrUnauthorized, err)
}

func Test_V2Poll_StatusRequestEntityTooLarge(t *testing.T) {
	resp := `{
	"error": {
		"code": "string",
		"field": "string",
		"details": [null],
		"message": "too many merchants"
	}
}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/events:polling", r.URL.Path)
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.V2Poll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, "too many merchants", err.Error())
}

func Test_V2Poll_ValidateErr(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/events:polling", r.URL.Path)
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprintf(w, `{}`)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("error"))
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	eventsService := New(adapter, &am, false)
	assert.NotNil(t, eventsService)
	events, err := eventsService.V2Poll()
	require.NotNil(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, "error", err.Error())
}
