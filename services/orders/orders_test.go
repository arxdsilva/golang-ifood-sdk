package orders

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var orderDetails = `{
    "id": "REFERENCIA",
    "reference": "Id de referencia do pedido",
    "shortReference": "Extranet Id",
    "createdAt": "Timestamp do pedido",
    "type": "Tipo do pedido('DELIVERY' ou 'TOGO')",
    "merchant": {
        "id": "Identificador unico do restaurante⁎",
        "name": "Nome do restaurante",
        "phones": [
            "Telefone do restaurante"
        ],
        "address": {
            "formattedAddress": "Endereço formatado",
            "country": "Pais",
            "state": "Estado",
            "city": "Cidade",
            "neighborhood": "Bairro",
            "streetName": "Endereço (Tipo logradouro + Logradouro)",
            "streetNumber": "Numero",
            "postalCode": "CEP"
        }
    },
    "payments": [
        {
            "name": "Nome da forma de pagamento",
            "code": "Codigo da forma de pagamento⁎⁎⁎",
            "value": "Valor pago na forma",
            "prepaid": "Pedido pago ('true' ou 'false')",
            "issuer": "Bandeira"
        },
        {
            "name": "Nome da forma de pagamento",
            "code": "Codigo da forma de pagamento⁎⁎⁎",
            "value": "Valor pago na forma",
            "prepaid": "Pedido pago ('true' ou 'false')",
            "collector": "Recebedor da forma",
            "issuer": "Bandeira"
        }
    ],
    "customer": {
        "id": "Id do cliente",
        "uuid": "Id Único do cliente",
        "name": "Nome do cliente",
        "taxPayerIdentificationNumber": "CPF/CNPJ do cliente (opcional) ",
        "phone": "0800 + Localizador",
        "ordersCountOnRestaurant":"Qtde de pedidos do cliente nesse restaurante"
    },
    "items": [
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "externalCode": "Código do e-PDV",
            "subItems": [
                {
                    "name": "Nome do item",
                    "quantity": "Quantidade",
                    "price": "Preço",
                    "totalPrice": "Preço total",
                    "discount": "Desconto",
                    "addition": "Adição",
                    "externalCode": "Código do e-PDV"
                }
            ]
        },
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "subItems": [
                {
                    "name": "Nome do item",
                    "quantity": "Quantidade",
                    "price": "Preço",
                    "totalPrice": "Preço total",
                    "discount": "Desconto",
                    "addition": "Adição",
                    "externalCode": "Código e-PDV"
                }
            ]
        },
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "externalCode": "Código do e-PDV",
            "observations": "Observação do item"
        }
    ],
    "subTotal": "Total do pedido(Sem taxa de entrega)",
    "totalPrice": "Total do pedido(Com taxa de entrega)",
    "deliveryFee": "Taxa de entrega",
    "deliveryAddress": {
        "formattedAddress": "Endereço completo de entrega",
        "country": "Pais",
        "state": "Estado",
        "city": "Cidade",
        "coordinates": {
            "latitude": "Latitude do endereço",
            "longitude": "Longitude do endereço"
        },
        "neighborhood": "Bairro",
        "streetName": "Endereço(Tipo logradouro + Logradouro)",
        "streetNumber": "Numero",
        "postalCode": "CEP",
        "reference": "Referencia",
        "complement": "Complemento do endereço"
    },
    "deliveryDateTime": "Timestamp do pedido",
    "preparationTimeInSeconds": "Tempo de preparo do pedido em segundos"
}`

var trackingOK = `{
	"date": 0,
	"deliveryTime": "2020-06-29T15:24:30.405Z",
	"eta": 10,
	"etaToDestination": 0,
	"etaToOrigin": 0,
	"latitude": 0,
	"longitude": 0,
	"orderId": "string",
	"trackDate": "2020-06-29T15:24:30.406Z"
}`

var v2OrderDetails = `{
	"benefits": [{"targetId": "string","sponsorshipValues": [{"name": "string","value": 0}],"value": 0,"target": "string"}],
	"orderType": "DELIVERY",
	"payments": {
	  "methods": [
		{
		  "wallet": {"name": "string"},
		  "method": "string",
		  "prepaid": true,
		  "currency": "string",
		  "type": "ONLINE",
		  "value": 0,
		  "cash": {
			"changeFor": 0
		  },
		  "card": {"brand": "string"}
		}
	  ],
	  "pending": 0,
	  "prepaid": 0
	},
	"merchant": {"name": "string","id": "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
	"salesChannel": "string",
	"orderTiming": "IMMEDIATE",
	"createdAt": "2021-05-23T14:57:03.193Z",
	"total": {
	  "benefits": 0,
	  "deliveryFee": 0,
	  "orderAmount": 0,
	  "subTotal": 0
	},
	"preparationStartDateTime": "2021-05-23T14:57:03.193Z",
	"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
	"displayId": "string",
	"items": [
	  {
		"unitPrice": 0,
		"quantity": 0,
		"externalCode": "string",
		"totalPrice": 0,
		"index": 0,
		"unit": "string",
		"ean": "string",
		"price": 0,
		"observations": "string",
		"name": "string",
		"options": [
		  {
			"unitPrice": 0,
			"unit": "string",
			"ean": "string",
			"quantity": 0,
			"externalCode": "string",
			"price": 0,
			"name": "string",
			"index": 0,
			"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
		  }
		],
		"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"optionsPrice": 0
	  }
	],
	"customer": {
	  "phone": {
		"number": "string",
		"localizer": "string",
		"localizerExpiration": "2021-05-23T14:57:03.193Z"
	  },
	  "documentNumber": "string",
	  "name": "string",
	  "ordersCountOnMerchant": 0,
	  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	},
	"extraInfo": "string",
	"delivery": {
	  "mode": "DEFAULT",
	  "deliveredBy": "IFOOD",
	  "deliveryAddress": {
		"reference": "string",
		"country": "string",
		"streetName": "string",
		"formattedAddress": "string",
		"streetNumber": "string",
		"city": "string",
		"postalCode": "string",
		"coordinates": {"latitude": 0,"longitude": 0},
		"neighborhood": "string",
		"state": "string",
		"complement": "string"
	  },
	  "deliveryDateTime": "2021-05-23T14:57:03.193Z"
	},
	"schedule": {
	  "deliveryDateTimeStart": "2021-05-23T14:57:03.193Z",
	  "deliveryDateTimeEnd": "2021-05-23T14:57:03.193Z"
	},
	"indoor": {
	  "mode": "DEFAULT",
	  "deliveryDateTime": "2021-05-23T14:57:03.193Z",
	  "table": "string"
	},
	"takeout": {"mode": "DEFAULT","takeoutDateTime": "2021-05-23T14:57:03.193Z"}
}`

func TestGetDetails_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, orderDetails)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	od, err := ordersService.GetDetails("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, "REFERENCIA", od.ID)
}

func TestGetDetails_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestGetDetails_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestGetDetails_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	od, err := ordersService.GetDetails("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, OrderDetails{}, od)
	assert.Contains(t, err.Error(), "could not retrieve details")
}

func TestGetDetails_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetIntegrateStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/integration", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetIntegrateStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetIntegrateStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetIntegrateStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/integration", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be integrated")
}

func TestSetIntegrateStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetConfirmStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/confirmation", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetConfirmStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetConfirmStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetConfirmStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/confirmation", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be confirmed")
}

func TestSetConfirmStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetDispatchStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/dispatch", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetDispatchStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetDispatchStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetDispatchStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/dispatch", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be dispatched")
}

func TestSetDispatchStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetReadyToDeliverStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/readyToDeliver", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetReadyToDeliverStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetReadyToDeliverStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetReadyToDeliverStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/readyToDeliver", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not be set as 'ready to deliver'")
}

func TestSetReadyToDeliverStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetCancelStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id/statuses/cancellationRequested", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.Nil(t, err)
}

func TestSetCancelStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("", "501")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetCancelStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetCancelStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id/statuses/cancellationRequested", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not be set as 'cancelled' code")
}

func TestSetCancelStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestClientCancellationStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationAccepted", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.Nil(t, err)
}

func TestClientCancellationStatus_OK_NotAccepted(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationDenied", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", false)
	assert.Nil(t, err)
}

func TestClientCancellationStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("", true)
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestClientCancellationStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestClientCancellationStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationAccepted", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not set 'client cancellation' status")
}

func TestClientCancellationStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestTracking_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/tracking", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, trackingOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	tr, err := ordersService.Tracking("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, 10, tr.Eta)
}

func TestTracking_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestTracking_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestTracking_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/tracking", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not get tracking information")
}

func TestTracking_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestDeliveryInformation_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/delivery-information", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, trackingOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	tr, err := ordersService.DeliveryInformation("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, 10, tr.Eta)
}

func TestDeliveryInformation_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestDeliveryInformation_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestDeliveryInformation_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/delivery-information", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could get delivery information")
}

func TestDeliveryInformation_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func Test_verifyCancel_OK(t *testing.T) {
	err := verifyCancel("reference", "501")
	assert.Nil(t, err)
}

func Test_verifyCancel_NoReferenceID(t *testing.T) {
	err := verifyCancel("", "501")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_verifyCancel_NoCode(t *testing.T) {
	err := verifyCancel("reference", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrCancelCodeNotSpecified, err)
}

func Test_verifyCancel_InvalidCode(t *testing.T) {
	err := verifyCancel("reference", "12344")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is invalid, verify docs")
}

func Test_V2GetDetails_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, v2OrderDetails)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	od, err := ordersService.V2GetDetails("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, "3fa85f64-5717-4562-b3fc-2c963f66afa6", od.ID)
}

func Test_V2GetDetails_BadRequest(t *testing.T) {
	resp := `{
		"error": {
			"code": "string",
			"field": "string",
			"details": [null],
			"message": "bad request"
		}
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.V2GetDetails("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, "bad request", err.Error())
}

func Test_V2GetDetails_NoOrderUUID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.V2GetDetails("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_V2GetDetails_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("error"))
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.V2GetDetails("98989898989")
	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Error())
}

func Test_V2SetConfirmStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id/confirm", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetConfirmStatus("reference_id")
	assert.Nil(t, err)
}

func Test_V2SetConfirmStatus_BadRequest(t *testing.T) {
	resp := `{
		"error": {
			"code": "string",
			"field": "string",
			"details": [null],
			"message": "bad request"
		}
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id/confirm", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetConfirmStatus("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, "bad request", err.Error())
}

func Test_V2SetConfirmStatus_NoRID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetConfirmStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_V2SetConfirmStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("validate err"))
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetConfirmStatus("123123123")
	assert.NotNil(t, err)
	assert.Equal(t, "validate err", err.Error())
}

func Test_V2SetDispatchStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/orders/reference_id/dispatch", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetDispatchStatus("reference_id")
	assert.Nil(t, err)
}

func Test_V2SetDispatchStatus_BadRequest(t *testing.T) {
	resp := `{
		"error": {
			"code": "string",
			"field": "string",
			"details": [null],
			"message": "bad request"
		}
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id/dispatch", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetDispatchStatus("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, "bad request", err.Error())
}

func Test_V2SetDispatchStatus_NoRID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetDispatchStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_V2SetDispatchStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("validate err"))
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetDispatchStatus("123123123")
	assert.NotNil(t, err)
	assert.Equal(t, "validate err", err.Error())
}

func Test_V2SetReadyToPickupStatus_NoRID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetReadyToPickupStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_V2SetReadyToPickupStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("validate err"))
	adapter := httpadapter.New(http.DefaultClient, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetReadyToPickupStatus("123123123")
	assert.NotNil(t, err)
	assert.Equal(t, "validate err", err.Error())
}

func Test_V2SetReadyToPickupStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/order/v1.0/orders/reference_id/readyToPickup", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetReadyToPickupStatus("reference_id")
	assert.Nil(t, err)
}

func Test_V2SetReadyToPickupStatus_BadRequest(t *testing.T) {
	resp := `{
		"error": {
			"code": "string",
			"field": "string",
			"details": [null],
			"message": "bad request"
		}
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, newV2Endpoint+"reference_id/readyToPickup", r.URL.Path)
			require.Equal(t, "Bearer token", r.Header["Authorization"][0])
			require.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, resp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.V2SetReadyToPickupStatus("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, "bad request", err.Error())
}
