package catalog

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/stretchr/testify/assert"
)

func TestCreateItem_OK(t *testing.T) {
	resp := `{
		"status": "AVAILABLE",
		"price": {
			"value": 10,
			"originalValue": 0
		},
		"externalCode": "string",
		"sequence": 1,
		"shifts": [{}]
	}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/catalog/v2.0/merchants/merchant_id/categories/category_id/products/product_id", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, "application/json", r.Header["Content-Type"][0])
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
	catalogService := New(adapter, &am)
	assert.NotNil(t, catalogService)
	ci := CategoryItem{
		Name:   "id",
		Status: "AVAILABLE",
		Price:  Price{Value: 10},
		Shifts: []Shift{
			{StartTime: "00:00", EndTime: "23:59", Monday: true},
		},
	}
	respCI, err := catalogService.CreateItem("merchant_id", "category_id", "product_id", ci)
	assert.Nil(t, err)
	assert.Equal(t, float64(10), respCI.Price.Value)
	assert.Equal(t, 1, respCI.Sequence)
}
