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

func TestListAllCategoriesInCatalog_OK(t *testing.T) {
	categories := `{
		"id": "string",
		"sequence": 10,
		"name": "string",
		"externalCode": "string",
		"status": "string",
		"items": [{}],
		"template": "string",
		"pizza": {
			"id": "string",
			"sizes": []
		}
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog_id/categories", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, categories)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	catalog, err := ordersService.ListAllCategoriesInCatalog("merchant_id", "catalog_id")
	assert.Nil(t, err)
	assert.Equal(t, 10, catalog.Sequence)
}

func TestListAllCategoriesInCatalog_NoMerchantID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllCategoriesInCatalog("", "catalog_id")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}
func TestListAllCategoriesInCatalog_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllCategoriesInCatalog("merchant_id", "catalog_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestListAllCategoriesInCatalog_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog_id/categories", r.URL.Path)
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
	_, err := ordersService.ListAllCategoriesInCatalog("merchant_id", "catalog_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not list categories in catalog")
}

func TestListAllCategoriesInCatalog_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllCategoriesInCatalog("reference_id", "catalog_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestCreateCategoryInCatalog_OK(t *testing.T) {
	resp := `{
		"id":"string",
		"name":"string",
		"externalCode":"string",
		"status":"AVAILABLE",
		"sequence":10,
		"template":"DEFAULT"
		}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog_id/categories", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusCreated)
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
	category, err := ordersService.CreateCategoryInCatalog("merchant_id", "catalog_id", "name", "AVAILABLE", "DEFAULT", "")
	assert.Nil(t, err)
	assert.Equal(t, 10, category.Sequence)
}

func TestCreateCategoryInCatalog_NoMerchantID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.CreateCategoryInCatalog("", "catalog_id", "name", "AVAILABLE", "DEFAULT", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}
func TestCreateCategoryInCatalog_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.CreateCategoryInCatalog("merchant_id", "catalog_id", "name", "AVAILABLE", "DEFAULT", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestCreateCategoryInCatalog_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog_id/categories", r.URL.Path)
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
	_, err := ordersService.CreateCategoryInCatalog("merchant_id", "catalog_id", "name", "AVAILABLE", "DEFAULT", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not create category")
}

func TestCreateCategoryInCatalog_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.CreateCategoryInCatalog("reference_id", "catalog_id", "name", "AVAILABLE", "DEFAULT", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}
