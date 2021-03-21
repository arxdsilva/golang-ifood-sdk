package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/kpango/glg"
)

const (
	V2Endpoint = "/catalog/v2.0"
)

var ErrMerchantNotSpecified = errors.New("merchant not specified")

type (
	Service interface {
		ListAll(merchantID string) (Catalogs, error)
		ListUnsellableItems(merchantUUID, catalogID string) (UnsellableResponse, error)
	}

	catalogService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

func New(adapter adapters.Http, authService auth.Service) *catalogService {
	return &catalogService{adapter, authService}
}

func (c *catalogService) ListAll(merchantUUID string) (ct Catalogs, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Catalog ListAll: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog ListAll auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf("/merchants/%s/catalogs", merchantUUID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog ListAll status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not list catalogs", merchantUUID)
		glg.Error("[SDK] Catalog ListAll err: ", err)
		return
	}
	return ct, json.Unmarshal(resp, &ct)
}

// ListChangelogs not implemented
func (c *catalogService) ListChangelogs(merchantUUID string) (ct Catalogs, err error) {
	return
}

// ListUnsellableItems returns all blocked sellable items and why
func (c *catalogService) ListUnsellableItems(merchantUUID, catalogID string) (ur UnsellableResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Catalog ListUnsellableItems: ", err.Error())
		return
	}
	if catalogID == "" {
		err = errors.New("catalog ID not specified")
		glg.Error("[SDK] Catalog ListUnsellableItems: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog ListUnsellableItems auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/unsellable-items", merchantUUID, catalogID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog ListUnsellableItems status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not list unsellable items, catalog: '%s'",
			merchantUUID, catalogID)
		glg.Error("[SDK] Catalog ListUnsellableItems err: ", err)
		return
	}
	return ur, json.Unmarshal(resp, &ur)
}
