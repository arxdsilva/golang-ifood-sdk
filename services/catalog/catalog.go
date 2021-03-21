package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
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
		ListAllCategoriesInCatalog(merchantUUID, catalogID string) (CategoryResponse, error)
		CreateCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template, externalCode string) (CategoryCreateResponse, error)
		GetCategoryInCatalog(merchantUUID, catalogID, categoryID string) (CategoryResponse, error)
		EditCategoryInCatalog(merchantUUID, catalogID, categoryID, name, resourceStatus, externalCode string, sequence int) (CategoryCreateResponse, error)
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
		glg.Error("[SDK] Catalog ListAll adapter.DoRequest: ", err.Error())
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
	return ur, json.Unmarshal(resp, &ur)
}

func (c *catalogService) ListAllCategoriesInCatalog(merchantUUID, catalogID string) (cr CategoryResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog: ", err.Error())
		return
	}
	if catalogID == "" {
		err = errors.New("Catalog ID was not specified")
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories", merchantUUID, catalogID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not list categories in catalog '%s'",
			merchantUUID, catalogID)
		glg.Error("[SDK] Catalog ListAllCategoriesInCatalog err: ", err)
		return
	}
	return cr, json.Unmarshal(resp, &cr)
}

// CreateCategoryInCatalog adds a category in a specified catalog
//
// resource status 	= [AVAILABLE || UNAVAILABLE]
// template 		= [DEFAULT || PIZZA]
//
func (c *catalogService) CreateCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template, externalCode string) (cr CategoryCreateResponse, err error) {
	err = verifyNewCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template)
	if err != nil {
		glg.Error("[SDK] Catalog CreateCategoryInCatalog verifyNewCategoryInCatalog: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog CreateCategoryInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories", merchantUUID, catalogID)
	ci := CategoryItem{Name: name, Status: resourceStatus, Template: template, ExternalCode: externalCode}
	reader, err := httpadapter.NewJsonReader(ci)
	if err != nil {
		glg.Error("[SDK] Catalog CreateCategoryInCatalog NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Catalog CreateCategoryInCatalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusCreated {
		glg.Error("[SDK] Catalog CreateCategoryInCatalog status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not create category in catalog '%s'",
			merchantUUID, catalogID)
		glg.Error("[SDK] Catalog CreateCategoryInCatalog err: ", err)
		return
	}
	return cr, json.Unmarshal(resp, &cr)
}

// GetCategoryInCatalog lists a category in a specified catalog
func (c *catalogService) GetCategoryInCatalog(merchantUUID, catalogID, categoryID string) (cr CategoryResponse, err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		glg.Error("[SDK] Catalog GetCategoryInCatalog: ", err.Error())
		return
	}
	if catalogID == "" {
		err = errors.New("Catalog ID was not specified")
		glg.Error("[SDK] Catalog GetCategoryInCatalog: ", err.Error())
		return
	}
	if categoryID == "" {
		err = errors.New("Category ID was not specified")
		glg.Error("[SDK] Catalog GetCategoryInCatalog: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog GetCategoryInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories/%s", merchantUUID, catalogID, catalogID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog GetCategoryInCatalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog GetCategoryInCatalog status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not get category '%s' in catalog '%s'",
			merchantUUID, categoryID, catalogID)
		glg.Error("[SDK] Catalog GetCategoryInCatalog err: ", err)
		return
	}
	return cr, json.Unmarshal(resp, &cr)
}

// EditCategoryInCatalog changes a category in a specified catalog
//
// resource status = [AVAILABLE || UNAVAILABLE]
func (c *catalogService) EditCategoryInCatalog(merchantUUID, catalogID, categoryID, name, resourceStatus, externalCode string, sequence int) (cr CategoryCreateResponse, err error) {
	err = verifyNewCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, "DEFAULT")
	if err != nil {
		glg.Error("[SDK] Catalog EditCategoryInCatalog verifyNewCategoryInCatalog: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog EditCategoryInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := V2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories/%s", merchantUUID, catalogID, categoryID)
	ci := CategoryItem{
		Status:       resourceStatus,
		ExternalCode: externalCode,
		Name:         name,
		Sequence:     sequence,
	}
	body, err := httpadapter.NewJsonReader(ci)
	if err != nil {
		glg.Error("[SDK] Catalog EditCategoryInCatalog NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPatch, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog EditCategoryInCatalog adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog EditCategoryInCatalog status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not edit category '%s' in catalog '%s'",
			merchantUUID, catalogID, catalogID)
		glg.Error("[SDK] Catalog EditCategoryInCatalog err: ", err)
		return
	}
	return cr, json.Unmarshal(resp, &cr)
}

func verifyNewCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template string) (err error) {
	if merchantUUID == "" {
		err = ErrMerchantNotSpecified
		return
	}
	if catalogID == "" {
		err = errors.New("Catalog ID was not specified")
		return
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
