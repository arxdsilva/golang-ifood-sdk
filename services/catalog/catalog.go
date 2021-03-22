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
	v2Endpoint = "/catalog/v2.0"
)

type catalogService struct {
	adapter adapters.Http
	auth    auth.Service
}

// New returns an implementation of the catalog service
func New(adapter adapters.Http, authService auth.Service) *catalogService {
	return &catalogService{adapter, authService}
}

// ListAll catalogs from a Merchant
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
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/catalogs", merchantUUID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog ListAll adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog ListAll status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not list catalogs", merchantUUID)
		glg.Error("[SDK] Catalog ListAll err: ", err)
		return
	}
	glg.Info("[SDK] ListAll catalogs success")
	return ct, json.Unmarshal(resp, &ct)
}

// ListChangelogs not implemented
func (c *catalogService) ListChangelogs(merchantUUID string) (ct Catalogs, err error) {
	return
}

// ListUnsellableItems returns all blocked sellable items and why
func (c *catalogService) ListUnsellableItems(merchantUUID, catalogID string) (ur UnsellableResponse, err error) {
	if err = verifyCategoryItems(merchantUUID, catalogID, "category"); err != nil {
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
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/unsellable-items", merchantUUID, catalogID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog ListUnsellableItems adapter.DoRequest: ", err.Error())
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
	glg.Info("[SDK] List Unsellable Items success")
	return ur, json.Unmarshal(resp, &ur)
}

func verifyNewCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template string) (err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		return
	}
	if catalogID == "" {
		return ErrCatalogNotSpecified
	}
	if len(name) > 100 {
		err = errors.New("Category name needs to have less than 100 characters")
		return
	}
	if name == "" {
		err = fmt.Errorf("Category name on catalog '%s' was not specified", catalogID)
		return
	}
	if (resourceStatus != "AVAILABLE") && (resourceStatus != "UNAVAILABLE") {
		err = fmt.Errorf(
			"Category status on catalog '%s' should be 'AVAILABLE' or 'UNAVAILABLE' and is '%s'",
			catalogID, resourceStatus)
		return
	}
	if (template != "DEFAULT") && (template != "PIZZA") {
		err = fmt.Errorf(
			"Category template on catalog '%s' should be 'DEFAULT' or 'PIZZA' and is '%s'",
			catalogID, template)
		return
	}
	return
}

func verifyCategoryItems(merchantID, catalogID, categoryID string) (err error) {
	if merchantID == "" {
		return ErrMerchantNotSpecified
	}
	if catalogID == "" {
		return ErrCatalogNotSpecified
	}
	if categoryID == "" {
		return ErrCategoryNotSpecified
	}
	return
}
