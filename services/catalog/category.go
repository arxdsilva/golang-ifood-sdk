package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/kpango/glg"
)

// ListAllCategoriesInCatalog gets categories in a catalog
func (c *catalogService) ListAllCategoriesInCatalog(merchantUUID, catalogID string) (cr CategoryResponse, err error) {
	if err = verifyCategoryItems(merchantUUID, catalogID, "category"); err != nil {
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
	endpoint := v2Endpoint + fmt.Sprintf(
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
	glg.Info("[SDK] ListAll Categories success")
	return cr, json.Unmarshal(resp, &cr)
}

// CreateCategoryInCatalog adds a category in a specified catalog
//
// resource status 	= [AVAILABLE ||	UNAVAILABLE]
// template 		= [DEFAULT	 ||	PIZZA]
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
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf(
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
	glg.Info("[SDK] Get Category success")
	return cr, json.Unmarshal(resp, &cr)
}

// GetCategoryInCatalog lists a category in a specified catalog
func (c *catalogService) GetCategoryInCatalog(merchantUUID, catalogID, categoryID string) (cr CategoryResponse, err error) {
	if err = verifyCategoryItems(merchantUUID, catalogID, categoryID); err != nil {
		glg.Error("[SDK] Catalog GetCategoryInCatalog verifyCategoryItems: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog GetCategoryInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories/%s", merchantUUID, catalogID, categoryID)
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
	glg.Info("[SDK] Get Category success")
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
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf(
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
	glg.Info("[SDK] Edit Category success")
	return cr, json.Unmarshal(resp, &cr)
}

// DeleteCategoryInCatalog removes a category in a specified catalog
func (c *catalogService) DeleteCategoryInCatalog(merchantUUID, catalogID, categoryID string) (err error) {
	if err = verifyCategoryItems(merchantUUID, catalogID, categoryID); err != nil {
		glg.Error("[SDK] Catalog DeleteCategoryInCatalog verifyCategoryItems: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog DeleteCategoryInCatalog auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/catalogs/%s/categories/%s", merchantUUID, catalogID, categoryID)
	_, status, err := c.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog DeleteCategoryInCatalog adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog DeleteCategoryInCatalog status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not delete category '%s' in catalog '%s'",
			merchantUUID, catalogID, catalogID)
		glg.Error("[SDK] Catalog DeleteCategoryInCatalog err: ", err)
		return
	}
	glg.Info("[SDK] Delete product success")
	return
}
