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

var catalogsv2 = `[
	{
		"catalogId":"string",
		"status":"AVAILABLE",
		"context":[
			"DEFAULT",
			"INDOOR"
		],
		"modifiedAt":"2021-03-28T13:16:56.574Z"
	},
	{
		"catalogId":"string",
		"status":"AVAILABLE",
		"context":[
			"DEFAULT",
			"INDOOR"
		],
		"modifiedAt":"2021-03-28T13:16:56.574Z"
	}
]`

var unsellableResp = `{
	"categories": [{
		"id": "string",
		"status": "string",
		"template": "string",
		"restrictions": [
			"HAS_VIOLATION",
			"CATEGORY_PAUSED"
		],
		"unsellableItems": [{
			"id": "string",
			"productId": "string"
		}]
	}]
}`

func Test_verifyNewCategoryInCatalog_no_merchant(t *testing.T) {
	err := verifyNewCategoryInCatalog("", "", "", "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrMerchantNotSpecified)
}

func Test_verifyNewCategoryInCatalog_no_category(t *testing.T) {
	err := verifyNewCategoryInCatalog("merchant", "", "", "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrCatalogNotSpecified)
}

func Test_verifyNewCategoryInCatalog_name_too_big(t *testing.T) {
	name := "namenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamenamename"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "less than 100 characters")
}

func Test_verifyNewCategoryInCatalog_no_name(t *testing.T) {
	name := ""
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name on catalog")
}

func Test_verifyNewCategoryInCatalog_no_resource(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "should be 'AVAILABLE' or 'UNAVAILABLE'")
}

func Test_verifyNewCategoryInCatalog_no_template(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "AVAILABLE", "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "'DEFAULT' or 'PIZZA' and")
}

func Test_verifyNewCategoryInCatalog_OK(t *testing.T) {
	name := "name"
	err := verifyNewCategoryInCatalog("merchant", "catalog", name, "AVAILABLE", "DEFAULT")
	assert.Nil(t, err)
}

func TestListAllV2_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, catalogsv2)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	list, err := ordersService.ListAllV2("merchant_id")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
}

func TestListAllV2_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllV2("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}

func TestListAllV2_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllV2("merchant_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestListAllV2_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs", r.URL.Path)
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
	_, err := ordersService.ListAllV2("merchant_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not list catalogs")
}

func TestListAllV2_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListAllV2("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestListUnsellableItems_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog/unsellable-items", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, unsellableResp)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	list, err := ordersService.ListUnsellableItems("merchant_id", "catalog")
	assert.Nil(t, err)
	assert.Equal(t, "string", list.Categories[0].ID)
}

func TestListUnsellableItems_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListUnsellableItems("", "catalog")
	assert.NotNil(t, err)
	assert.Equal(t, ErrMerchantNotSpecified, err)
}

func TestListUnsellableItems_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListUnsellableItems("merchant_id", "catalog")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestListUnsellableItems_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/catalogs/catalog/unsellable-items", r.URL.Path)
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
	_, err := ordersService.ListUnsellableItems("merchant_id", "catalog")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not list unsellable items")
}

func TestListUnsellableItems_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.ListUnsellableItems("reference_id", "catalog")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func Test_verifyCategoryItems_NoCatalog(t *testing.T) {
	err := verifyCategoryItems("m_id", "", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrCatalogNotSpecified, err)
}

func Test_verifyCategoryItems_NoCategory(t *testing.T) {
	err := verifyCategoryItems("m_id", "catalog", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrCategoryNotSpecified, err)
}
