package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/kpango/glg"
)

// CreateItem product-category association
//
// POST
//
// 201 created
// 400 bad req
// 404 not found
// 409 conflict
//
// Response OK:
// {
// 		"status":"AVAILABLE",
// 		"price":{
// 			"value":0,
// 			"originalValue":0
// 		},
// 		"externalCode":"string",
// 		"sequence":0,
// 		"shifts":[{...}]
// }
func (c *catalogService) CreateItem(merchantID, categoryID, productID string, ci CategoryItem) (cp ProductLink, err error) {
	err = verifyCategoryItems(merchantID, categoryID, productID)
	if err != nil {
		glg.Error("[SDK] Catalog CreateItem verifyCategoryItems: ", err.Error())
		return
	}
	if err = ci.verify(); err != nil {
		glg.Error("[SDK] Catalog CreateItem verify: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog CreateItem auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/categories/%s/products/%s", merchantID, categoryID, productID)
	reader, err := httpadapter.NewJsonReader(ci)
	if err != nil {
		glg.Error("[SDK] Catalog CreateItem NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Catalog CreateItem adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusCreated {
		badResp := &apiError{}
		err = json.Unmarshal(resp, badResp)
		if err != nil {
			glg.Error("[SDK] Catalog DeleteProduct Unmarshal: ", err)
			return
		}
		glg.Error("[SDK] Catalog CreateItem status code: ", status, " merchant: ", merchantID)
		err = fmt.Errorf(
			"Merchant '%s' could not create item category '%s', code: '%s'",
			merchantID, categoryID, badResp.Details.Code)
		glg.Error("[SDK] Catalog CreateItem err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog CreateItem success product id '%s', merchant '%s'", productID, merchantID)
	return cp, json.Unmarshal(resp, &cp)
}

// EditItem product-catalog association
//
// PATCH
//
// 200 OK
// 400 bad req
//
// Response:
// {
// 		"status":"AVAILABLE",
// 		"price":{
// 			"value":0,
// 			"originalValue":0
// 		},
// 		"externalCode":"string",
// 		"sequence":0,
// 		"shifts":[{...}]
// }
func (c *catalogService) EditItem(merchantID, categoryID, productID string, ci CategoryItem) (cp ProductLink, err error) {
	err = verifyCategoryItems(merchantID, categoryID, productID)
	if err != nil {
		glg.Error("[SDK] Catalog EditItem verifyCategoryItems: ", err.Error())
		return
	}
	if err = ci.verify(); err != nil {
		glg.Error("[SDK] Catalog EditItem verify: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog EditItem auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/categories/%s/products/%s", merchantID, categoryID, productID)
	reader, err := httpadapter.NewJsonReader(ci)
	if err != nil {
		glg.Error("[SDK] Catalog EditItem NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPatch, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Catalog EditItem adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		badResp := &apiError{}
		err = json.Unmarshal(resp, badResp)
		if err != nil {
			glg.Error("[SDK] Catalog DeleteProduct Unmarshal: ", err)
			return
		}
		glg.Error("[SDK] Catalog EditItem status code: ", status, " merchant: ", merchantID)
		err = fmt.Errorf(
			"Merchant '%s' could not create item category '%s', code: '%s'",
			merchantID, categoryID, badResp.Details.Code)
		glg.Error("[SDK] Catalog EditItem err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog EditItem success product id '%s', merchant '%s'", productID, merchantID)
	return cp, json.Unmarshal(resp, &cp)
}

// DeleteItem product-catalog association
//
// PATCH
//
// 200 OK
// 400 bad req
// 404 not found
//
func (c *catalogService) DeleteItem(merchantID, categoryID, productID string) (err error) {
	err = verifyCategoryItems(merchantID, categoryID, productID)
	if err != nil {
		glg.Error("[SDK] Catalog DeleteItem verifyCategoryItems: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog DeleteItem auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/categories/%s/products/%s", merchantID, categoryID, productID)
	resp, status, err := c.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog DeleteItem adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		badResp := &apiError{}
		err = json.Unmarshal(resp, badResp)
		if err != nil {
			glg.Error("[SDK] Catalog DeleteProduct Unmarshal: ", err)
			return
		}
		glg.Error("[SDK] Catalog DeleteItem status code: ", status, " merchant: ", merchantID)
		err = fmt.Errorf(
			"Merchant '%s' could not create item category '%s', code: '%s'",
			merchantID, categoryID, badResp.Details.Code)
		glg.Error("[SDK] Catalog DeleteItem err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog DeleteItem success product id '%s', merchant '%s'", productID, merchantID)
	return
}
