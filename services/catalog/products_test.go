package catalog

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
)

func TestProduct_verifyFields_noname(t *testing.T) {
	p := Product{}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestProduct_verifyFields_long_name(t *testing.T) {
	p := Product{
		Name: "produtoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoproduto",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "len is higher")
}

func TestProduct_verifyFields_long_description(t *testing.T) {
	p := Product{
		Name:        "nome",
		Description: "produtoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoprodutoproduto",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Description")
}

func TestProduct_verifyFields_no_shift(t *testing.T) {
	p := Product{
		Name: "nome",
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "shift")
}

func TestProduct_verifyFields_no_serving(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Serving")
}

func TestProduct_verifyFields_OK(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving: "SERVES_1",
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}

func TestProduct_verifyFields_invalid_restriction(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"JAPONES"},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "restriction")
}

func TestProduct_verifyFields_OK_restriction(t *testing.T) {
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}

func TestPizza_verifyFields_ErrSizesNotSpecified(t *testing.T) {
	p := Pizza{}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrSizesNotSpecified, err)
}

func TestPizza_verifyFields_ErrCrustsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrCrustsNotSpecified, err)
}

func TestPizza_verifyFields_ErrEdgesNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrEdgesNotSpecified, err)
}

func TestPizza_verifyFields_ErrToppingsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrToppingsNotSpecified, err)
}

func TestPizza_verifyFields_ErrShiftsNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "item"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrShiftsNotSpecified, err)
}

func TestPizza_verifyFields_ErrSizeNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{ID: "id"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrSizeNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStatus, err)
}

func TestPizza_verifyFields_ErrNoAcceptedFractions(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE"},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoAcceptedFractions, err)
}

func TestPizza_verifyFields_ErrCrustNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{ID: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrCrustNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaCrustStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaCrustStatus, err)
}

func TestPizza_verifyFields_ErrEdgeNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{ID: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrEdgeNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaEdgeStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaEdgeStatus, err)
}

func TestPizza_verifyFields_ErrToppingNameNotSpecified(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{ID: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrToppingNameNotSpecified, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaToppingStatus(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping"},
		},
		Shifts: []Shift{
			{StartTime: "edge", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaToppingStatus, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaStartEndTime(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "", Monday: true},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStartEndTime, err)
}

func TestPizza_verifyFields_ErrInvalidPizzaEndTime(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "start", EndTime: ""},
		},
	}
	err := p.verifyFields()
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPizzaStartEndTime, err)
}

func TestPizza_verifyFields_OK(t *testing.T) {
	p := Pizza{
		Sizes: []CategoryItem{
			{Name: "id", Status: "AVAILABLE", AcceptedFractions: []float64{1}},
		},
		Crusts: []CategoryItem{
			{Name: "crust", Status: "AVAILABLE"},
		},
		Edges: []CategoryItem{
			{Name: "edge", Status: "AVAILABLE"},
		},
		Toppings: []CategoryItem{
			{Name: "topping", Status: "AVAILABLE"},
		},
		Shifts: []Shift{
			{StartTime: "start", EndTime: "end"},
		},
	}
	err := p.verifyFields()
	assert.Nil(t, err)
}

func TestListProducts_OK(t *testing.T) {
	categories := `[{
		"id": "string",
		"name": "pizza",
		"description": "string",
		"externalCode": "string",
		"image": "string",
		"shifts": [],
		"serving": "NOT_APPLICABLE",
		"dietaryRestrictions": [],
		"ean": "string"
	},
	{
		"id": "string",
		"name": "string",
		"description": "string",
		"externalCode": "string",
		"image": "string",
		"shifts": [],
		"serving": "NOT_APPLICABLE",
		"dietaryRestrictions": [],
		"ean": "string"
	}]`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/products", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, categories)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	list, err := catalogService.ListProducts("merchant_id")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "pizza", list[0].Name)
}

func TestListProducts_NoMerchantID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	_, err := catalogService.ListProducts("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}

func TestListProducts_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	_, err := catalogService.ListProducts("merchant_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestListProducts_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/products", r.URL.Path)
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
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	_, err := catalogService.ListProducts("merchant_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not get all products")
}

func TestListProducts_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	_, err := catalogService.ListProducts("merchant_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestCreateProduct_OK(t *testing.T) {
	product := `{
		"id": "12134",
		"name": "string",
		"description": "string",
		"externalCode": "string",
		"image": "string",
		"shifts": [{}],
		"serving": "NOT_APPLICABLE",
		"dietaryRestrictions": [
			"VEGETARIAN"
		],
		"ean": "string"
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/products", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, product)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	productResp, err := catalogService.CreateProduct("merchant_id", p)
	assert.Nil(t, err)
	assert.Equal(t, "12134", productResp.ID)
}

func TestCreateProduct_verifyProductErr(t *testing.T) {
	am := auth.AuthMock{}
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	p := Product{}
	_, err := catalogService.CreateProduct("merchant_id", p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoProductName, err)
}

func TestCreateProduct_NoMerchantID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	_, err := catalogService.CreateProduct("", Product{})
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}

func TestCreateProduct_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	_, err := catalogService.CreateProduct("merchant_id", p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestCreateProduct_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/products", r.URL.Path)
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
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	_, err := catalogService.CreateProduct("merchant_id", p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not create product")
}

func TestCreateProduct_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	p := Product{
		Name: "nome",
		Shifts: []Shift{
			{
				StartTime: "00:00",
				EndTime:   "23:59",
				Monday:    true,
			},
		},
		Serving:             "SERVES_1",
		DietaryRestrictions: []string{"SUGAR_FREE"},
	}
	_, err := catalogService.CreateProduct("merchant_id", p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}
