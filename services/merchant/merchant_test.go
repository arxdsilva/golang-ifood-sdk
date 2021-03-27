package merchant

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var merchants = `[
    {
        "id": "0fd7a60b-930c-49f5-a8d9-b721bb86f7c0",
        "name": "Test"
    }
]`

var unavailabilities = `[
    {
        "id": "84ab1175-5360-4b03-8598-3d16faaa560d",
        "storeId": "3d1b6527-99f2-498b-a6ad-23b4d2bf9999",
        "description": "Teste de Fechamento",
        "authorId": "Id",
        "start": "2020-04-04T14:30:00",
        "end": "2020-04-04T18:10:00"
    }
]`

var unavNowResponse = `{
    "id": "d0fd503f-7a2f-4bbb-8a5b-cee335ee4233",
    "storeId": "3d1b6527-99f2-498b-a6ad-23b4d2bfc999",
    "description": "Teste de Pausa Programada | Client id: username",
    "authorId": "9999999",
    "start": "2020-10-19T11:20:41.640899",
    "end": "2020-10-19T11:35:41.640899"
}`

var available = `
[
	{
		"context": "delivery",
		"available": true,
		"state": "OK",
		"reopenable": {
			"identifier": null,
			"type": null,
			"reopenable": false
			},
		"validations": [
			{
			"id": "opening-hours",
			"code": "during.opening-hours.config",
			"state": "OK",
			"message": {
				"title": "Dentro do horário de funcionamento",
				"subtitle": "quarta-feira, das 00:00 às 23:59",
				"description": "",
				"priority": 27
				}
			},
			{
			"id": "is-connected",
			"code": "is.connected.config",
			"state": "OK",
			"message": {
				"title": "Loja conectada à rede do iFood",
				"subtitle": "",
				"description": "",
				"priority": 999
				}
			}
		],
		"message": {
			"title": "Loja aberta",
			"subtitle": "",
			"description": "",
			"priority": 999
		}
	}
]`

func TestListAll_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/merchants", r.URL.Path)
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, merchants)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	events, err := merchantService.ListAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
}

func TestListAll_StatusNotFound(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/merchants", r.URL.Path)
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	events, err := merchantService.ListAll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
}

func TestListAll_ValidateErr(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/merchants", r.URL.Path)
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	events, err := merchantService.ListAll()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(events))
}

func TestUnavailabilities_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, unavailabilities)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	unavs, err := merchantService.Unavailabilities(id.String())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(unavs))
	assert.NotNil(t, unavs[0].ID)
}

func TestUnavailabilities_ErrMerchantNotSpecified(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, unavailabilities)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	unavs, err := merchantService.Unavailabilities("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
	assert.Equal(t, 0, len(unavs))
}

func TestUnavailabilities_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	unavs, err := merchantService.Unavailabilities(id.String())
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(unavs))
}

func TestCreateUnavailabilyNow_ErrMerchantNotSpecified(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities:now")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, unavNowResponse)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	unav, err := merchantService.CreateUnavailabilyNow("", "", 10)
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
	assert.Equal(t, UnavailabilityResponse{}, unav)
}

func TestCreateUnavailabilyNow_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities:now")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, unavNowResponse)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	unav, err := merchantService.CreateUnavailabilyNow(id.String(), "", 10)
	assert.Nil(t, err)
	assert.NotNil(t, unav.ID)
}

func TestCreateUnavailabilyNow_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities:now")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	unav, err := merchantService.CreateUnavailabilyNow(id.String(), "", 10)
	assert.NotNil(t, err)
	assert.Equal(t, UnavailabilityResponse{}, unav)
	assert.Contains(t, err.Error(), "not create")
}

func TestDeleteUnavailabily_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.Contains(t, r.URL.Path, "v1.0")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodDelete)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	err = merchantService.DeleteUnavailabily(id.String(), id.String())
	assert.Nil(t, err)
}

func TestDeleteUnavailabily_NoMerchant(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.Contains(t, r.URL.Path, "v1.0")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodDelete)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	err := merchantService.DeleteUnavailabily("", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantORUnavailabilityIDNotSpecified, err)
}

func TestDeleteUnavailabily_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "unavailabilities")
			assert.Contains(t, r.URL.Path, "v1.0")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodDelete)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	err = merchantService.DeleteUnavailabily(id.String(), id.String())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not delete")
}

func TestDeleteAvailabily_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/availabilities")
			assert.Contains(t, r.URL.Path, "/merchant/v2.0/merchants/")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, available)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	ar, err := merchantService.Availabily(id.String())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ar))
	assert.Equal(t, true, ar[0].Available)
}

func TestDeleteAvailabily_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/availabilities")
			assert.Contains(t, r.URL.Path, "/merchant/v2.0/merchants/")
			assert.NotNil(t, r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	merchantService := New(adapter, &am)
	assert.NotNil(t, merchantService)
	id, err := uuid.NewV1()
	assert.Nil(t, err)
	ar, err := merchantService.Availabily(id.String())
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(ar))
}
