package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/kpango/glg"
)

type (
	Products []Product
	Product  struct {
		ID                  string   `json:"id"`
		Name                string   `json:"name"`
		Description         string   `json:"description"`
		ExternalCode        string   `json:"externalCode"`
		Image               string   `json:"image"`
		Shifts              []Shift  `json:"shifts"`
		Serving             string   `json:"serving"`
		DietaryRestrictions []string `json:"dietaryRestrictions"`
		Ean                 string   `json:"ean"`
	}
	ProductLink struct {
		Status              string   `json:"status"`
		ID                  string   `json:"id"`
		Name                string   `json:"name"`
		Description         string   `json:"description"`
		ExternalCode        string   `json:"externalCode"`
		Image               string   `json:"image"`
		Price               Price    `json:"price"`
		Shifts              []Shift  `json:"shifts"`
		Serving             string   `json:"serving"`
		DietaryRestrictions []string `json:"dietaryRestrictions"`
		Ean                 string   `json:"ean"`
	}

	Pizzas []Pizza
	Pizza  struct {
		ID       string         `json:"id"`
		Sizes    []CategoryItem `json:"sizes"`
		Crusts   []CategoryItem `json:"crusts"`
		Edges    []CategoryItem `json:"edges"`
		Toppings []CategoryItem `json:"toppings"`
		Shifts   []Shift        `json:"shifts"`
	}
)

func (p *Product) verifyFields() (err error) {
	if p.Name == "" {
		return ErrNoProductName
	}
	if len(p.Name) > 65 {
		return errors.New("Name len is higher than 65 characters")
	}
	if len(p.Description) > 2000 {
		return errors.New("Description len is higher than 2000 characters")
	}
	if len(p.Shifts) == 0 {
		return errors.New("Product needs at least 1 shift")
	}
	serving := map[string]string{
		"SERVES_1":       "",
		"SERVES_2":       "",
		"SERVES_3":       "",
		"SERVES_4":       "",
		"NOT_APPLICABLE": "",
	}
	if _, ok := serving[p.Serving]; !ok {
		return errors.New("Serving not valid, verify docs: https://developer.ifood.com.br/reference#productcontroller_createproduct")
	}
	restrictions := map[string]string{
		"VEGETARIAN":      "",
		"VEGAN":           "",
		"ORGANIC":         "",
		"GLUTEN_FREE":     "",
		"SUGAR_FREE":      "",
		"LAC_FREE":        "",
		"ALCOHOLIC_DRINK": "",
		"NATURAL":         "",
	}
	if len(p.DietaryRestrictions) > 0 {
		for _, restriction := range p.DietaryRestrictions {
			if _, ok := restrictions[restriction]; !ok {
				return fmt.Errorf(
					"restriction '%s' does not exist in docs, see: https://developer.ifood.com.br/reference#productcontroller_createproduct",
					restriction)
			}
		}
	}
	return
}

func (p *Pizza) verifyFields() (err error) {
	if len(p.Sizes) == 0 {
		return ErrSizesNotSpecified
	}
	if len(p.Crusts) == 0 {
		return ErrCrustsNotSpecified
	}
	if len(p.Edges) == 0 {
		return ErrEdgesNotSpecified
	}
	if len(p.Toppings) == 0 {
		return ErrToppingsNotSpecified
	}
	if len(p.Shifts) == 0 {
		return ErrShiftsNotSpecified
	}
	for _, size := range p.Sizes {
		if size.Name == "" {
			return ErrSizeNameNotSpecified
		}
		if (size.Status != "AVAILABLE") && (size.Status != "UNAVAILABLE") {
			return ErrInvalidPizzaStatus
		}
		if len(size.AcceptedFractions) == 0 {
			return ErrNoAcceptedFractions
		}
	}
	for _, crust := range p.Crusts {
		if crust.Name == "" {
			return ErrCrustNameNotSpecified
		}
		if (crust.Status != "AVAILABLE") && (crust.Status != "UNAVAILABLE") {
			return ErrInvalidPizzaCrustStatus
		}
	}
	for _, edge := range p.Edges {
		if edge.Name == "" {
			return ErrEdgeNameNotSpecified
		}
		if (edge.Status != "AVAILABLE") && (edge.Status != "UNAVAILABLE") {
			return ErrInvalidPizzaEdgeStatus
		}
	}
	for _, topping := range p.Toppings {
		if topping.Name == "" {
			return ErrToppingNameNotSpecified
		}
		if (topping.Status != "AVAILABLE") && (topping.Status != "UNAVAILABLE") {
			return ErrInvalidPizzaToppingStatus
		}
	}
	for _, shift := range p.Shifts {
		if shift.StartTime == "" {
			return ErrInvalidPizzaStartEndTime
		}
		if shift.EndTime == "" {
			return ErrInvalidPizzaStartEndTime
		}
	}
	return
}

// ListProducts from a merchant
func (c *catalogService) ListProducts(merchantUUID string) (ps Products, err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog ListProducts verifyCategoryItems: ", err.Error())
		return
	}
	err = c.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Catalog ListProducts auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/products", merchantUUID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog ListProducts adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog ListProducts status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not get all products", merchantUUID)
		glg.Error("[SDK] Catalog ListProducts err: ", err)
		return
	}
	glg.Infof("[SDK] List products success, merchant '%s'", merchantUUID)
	return ps, json.Unmarshal(resp, &ps)
}

// CreateProduct in a merchant
func (c *catalogService) CreateProduct(merchantUUID string, product Product) (cp Product, err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog CreateProduct verifyCategoryItems: ", err.Error())
		return
	}
	if err = product.verifyFields(); err != nil {
		glg.Error("[SDK] Catalog CreateProduct verifyFields: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog CreateProduct auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/products", merchantUUID)
	body, err := httpadapter.NewJsonReader(product)
	if err != nil {
		glg.Error("[SDK] Catalog CreateProduct NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPost, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog CreateProduct adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusCreated {
		glg.Error("[SDK] Catalog CreateProduct status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create product", merchantUUID)
		glg.Error("[SDK] Catalog CreateProduct err: ", err)
		return
	}
	glg.Infof("[SDK] Create product id '%s' success, merchant '%s'", cp.ID, merchantUUID)
	return cp, json.Unmarshal(resp, &cp)
}

// EditProduct in a merchant
func (c *catalogService) EditProduct(merchantUUID string, product Product) (cp Product, err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog EditProduct verifyCategoryItems: ", err.Error())
		return
	}
	if product.ID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog EditProduct err: ", err.Error())
		return
	}
	if err = product.verifyFields(); err != nil {
		glg.Error("[SDK] Catalog EditProduct verifyFields: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog EditProduct auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/products/%s",
		merchantUUID, product.ID)
	body, err := httpadapter.NewJsonReader(product)
	if err != nil {
		glg.Error("[SDK] Catalog EditProduct NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPut, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog EditProduct adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog EditProduct status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not edit product id '%s'", merchantUUID, product.ID)
		glg.Error("[SDK] Catalog EditProduct err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog EditProduct id '%s' success, merchant '%s'", product.ID, merchantUUID)
	return cp, json.Unmarshal(resp, &cp)
}

// DeleteProduct in a merchant
func (c *catalogService) DeleteProduct(merchantUUID, productID string) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog DeleteProduct verifyCategoryItems: ", err.Error())
		return
	}
	if productID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog DeleteProduct err: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog DeleteProduct auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/products/%s",
		merchantUUID, productID)
	resp, status, err := c.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog DeleteProduct adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		badResp := &apiError{}
		err = json.Unmarshal(resp, badResp)
		if err != nil {
			glg.Error("[SDK] Catalog DeleteProduct Unmarshal: ", err)
			return
		}
		glg.Error("[SDK] Catalog DeleteProduct status code: ", status, " merchant: ", merchantUUID, " detail: ", badResp.Details.Code)
		err = fmt.Errorf("Merchant '%s' could not delete product id '%s', code: '%s'", merchantUUID, productID, badResp.Details.Code)
		glg.Error("[SDK] Catalog DeleteProduct err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog DeleteProduct id '%s' success, merchant '%s'", productID, merchantUUID)
	return
}

// UpdateProductStatus in a merchant
func (c *catalogService) UpdateProductStatus(merchantUUID, productID, productStatus string) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog UpdateProductStatus verifyCategoryItems: ", err.Error())
		return
	}
	if productID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog UpdateProductStatus err: ", err.Error())
		return
	}
	if (productStatus != "AVAILABLE") && (productStatus != "UNAVAILABLE") {
		err = fmt.Errorf("product status '%s' should be 'AVAILABLE' or 'UNAVAILABLE'", productStatus)
		glg.Error("[SDK] Catalog UpdateProductStatus err: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog UpdateProductStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/products/%s/status",
		merchantUUID, productID)
	bodyStatus := struct {
		Status string `json:"status"`
	}{Status: productStatus}
	body, err := httpadapter.NewJsonReader(bodyStatus)
	if err != nil {
		glg.Error("[SDK] Catalog UpdateProductStatus NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := c.adapter.DoRequest(http.MethodPatch, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog UpdateProductStatus adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog UpdateProductStatus status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not update product id '%s'", merchantUUID, productID)
		glg.Error("[SDK] Catalog UpdateProductStatus err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog UpdateProductStatus id '%s' success, merchant '%s'", productID, merchantUUID)
	return
}

// LinkProductToCategory in a merchant
func (c *catalogService) LinkProductToCategory(merchantUUID, categoryID string, product ProductLink) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", categoryID); err != nil {
		glg.Error("[SDK] Catalog LinkProductToCategory verifyCategoryItems: ", err.Error())
		return
	}
	if product.ID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog LinkProductToCategory err: ", err.Error())
		return
	}
	if (product.Status != "AVAILABLE") && (product.Status != "UNAVAILABLE") {
		err = fmt.Errorf("product status '%s' should be 'AVAILABLE' or 'UNAVAILABLE'", product.Status)
		glg.Error("[SDK] Catalog LinkProductToCategory err: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog LinkProductToCategory auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/categories/%s/products/%s",
		merchantUUID, categoryID, product.ID)
	body, err := httpadapter.NewJsonReader(product)
	if err != nil {
		glg.Error("[SDK] Catalog LinkProductToCategory NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := c.adapter.DoRequest(http.MethodPatch, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog LinkProductToCategory adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusCreated {
		glg.Error("[SDK] Catalog LinkProductToCategory status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not link product id '%s' to category '%s'",
			merchantUUID, product.ID, categoryID)
		glg.Error("[SDK] Catalog LinkProductToCategory err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog LinkProductToCategory id '%s' success, merchant '%s'", product.ID, merchantUUID)
	return
}

// UnlinkProductToCategory in a merchant
func (c *catalogService) UnlinkProductToCategory(merchantUUID, categoryID, productID string) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", categoryID); err != nil {
		glg.Error("[SDK] Catalog UnlinkProductToCategory verifyCategoryItems: ", err.Error())
		return
	}
	if productID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog UnlinkProductToCategory err: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog UnlinkProductToCategory auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/categories/%s/products/%s",
		merchantUUID, categoryID, productID)
	_, status, err := c.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog UnlinkProductToCategory adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog UnlinkProductToCategory status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could not unlink product id '%s' to category '%s'",
			merchantUUID, productID, categoryID)
		glg.Error("[SDK] Catalog UnlinkProductToCategory err: ", err)
		return
	}
	glg.Infof("[SDK] Catalog UnlinkProductToCategory id '%s' success, merchant '%s'", productID, merchantUUID)
	return
}

// CreatePizza in a merchant
func (c *catalogService) CreatePizza(merchantUUID string, pizza Pizza) (cp Pizza, err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog CreatePizza verifyCategoryItems: ", err.Error())
		return
	}
	if err = pizza.verifyFields(); err != nil {
		glg.Error("[SDK] Catalog CreatePizza verifyFields: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog CreatePizza auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/pizzas", merchantUUID)
	body, err := httpadapter.NewJsonReader(pizza)
	if err != nil {
		glg.Error("[SDK] Catalog CreatePizza NewJsonReader error: ", err.Error())
		return
	}
	resp, status, err := c.adapter.DoRequest(http.MethodPost, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog CreatePizza adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusCreated {
		glg.Error("[SDK] Catalog CreatePizza status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create pizza", merchantUUID)
		glg.Error("[SDK] Catalog CreatePizza err: ", err)
		return
	}
	glg.Infof("[SDK] Create pizza id '%s' success, merchant '%s'", cp.ID, merchantUUID)
	return cp, json.Unmarshal(resp, &cp)
}

// ListPizzas in a merchant
func (c *catalogService) ListPizzas(merchantUUID string) (pz Pizzas, err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog ListPizzas verifyCategoryItems: ", err.Error())
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog ListPizzas auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/pizzas", merchantUUID)
	resp, status, err := c.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog ListPizzas adapter.DoRequest: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Catalog ListPizzas status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not list pizzas", merchantUUID)
		glg.Error("[SDK] Catalog ListPizzas err: ", err)
		return
	}
	glg.Infof("[SDK] List pizzas merchant '%s' success", merchantUUID)
	return pz, json.Unmarshal(resp, &pz)
}

// UpdatePizza in a merchant
func (c *catalogService) UpdatePizza(merchantUUID string, pizza Pizza) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog UpdatePizza verifyCategoryItems: ", err.Error())
		return
	}
	if pizza.ID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog UpdatePizza verifyFields: ", err.Error(), " merchant ", merchantUUID)
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog UpdatePizza auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/pizzas/%s", merchantUUID, pizza.ID)
	body, err := httpadapter.NewJsonReader(pizza)
	if err != nil {
		glg.Error("[SDK] Catalog UpdatePizza NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := c.adapter.DoRequest(http.MethodPut, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog UpdatePizza adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog UpdatePizza status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not create pizza", merchantUUID)
		glg.Error("[SDK] Catalog UpdatePizza err: ", err)
		return
	}
	glg.Infof("[SDK] Update pizza id '%s' success, merchant '%s'", pizza.ID, merchantUUID)
	return
}

// UpdatePizzaStatus in a merchant
//
// pizzaStatus = [AVAILABLE || UNAVAILABLE]
func (c *catalogService) UpdatePizzaStatus(merchantUUID, pizzaStatus, pizzaID string) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", "categoryID"); err != nil {
		glg.Error("[SDK] Catalog UpdatePizzaStatus verifyCategoryItems: ", err.Error())
		return
	}
	if pizzaID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog UpdatePizzaStatus verifyFields: ", err.Error(), " merchant ", merchantUUID)
		return
	}
	if (pizzaStatus != "AVAILABLE") && (pizzaStatus != "UNAVAILABLE") {
		err = fmt.Errorf("pizza status '%s' should be 'AVAILABLE' or 'UNAVAILABLE'", pizzaStatus)
		glg.Error("[SDK] Catalog UpdatePizzaStatus verifyFields: ", err.Error(), " merchant ", merchantUUID)
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog UpdatePizzaStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf("/merchants/%s/pizzas/%s", merchantUUID, pizzaID)
	updateBody := struct {
		Status string `json:"status"`
	}{pizzaStatus}
	body, err := httpadapter.NewJsonReader(updateBody)
	if err != nil {
		glg.Error("[SDK] Catalog UpdatePizzaStatus NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := c.adapter.DoRequest(http.MethodPatch, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog UpdatePizzaStatus adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog UpdatePizzaStatus status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not update pizza id '%s' status", merchantUUID, pizzaID)
		glg.Error("[SDK] Catalog UpdatePizzaStatus err: ", err)
		return
	}
	glg.Infof("[SDK] Update pizza id '%s' success, merchant '%s'", pizzaID, merchantUUID)
	return
}

// LinkPizzaToCategory in a merchant
func (c *catalogService) LinkPizzaToCategory(merchantUUID, categoryID string, pizza Pizza) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", categoryID); err != nil {
		glg.Error("[SDK] Catalog LinkPizzaToCategory verifyCategoryItems: ", err.Error())
		return
	}
	if pizza.ID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog LinkPizzaToCategory verifyFields: ", err.Error(), " merchant ", merchantUUID)
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog LinkPizzaToCategory auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	headers["Content-Type"] = "application/json"
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/pizzas/%s/categories/%s", merchantUUID, pizza.ID, categoryID)
	body, err := httpadapter.NewJsonReader(pizza)
	if err != nil {
		glg.Error("[SDK] Catalog LinkPizzaToCategory NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := c.adapter.DoRequest(http.MethodPost, endpoint, body, headers)
	if err != nil {
		glg.Error("[SDK] Catalog LinkPizzaToCategory adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog LinkPizzaToCategory status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf("Merchant '%s' could not link pizza id '%s' to category '%s'",
			merchantUUID, pizza.ID, categoryID)
		glg.Error("[SDK] Catalog LinkPizzaToCategory err: ", err)
		return
	}
	glg.Infof("[SDK] Update pizza id '%s' success, merchant '%s'", pizza.ID, merchantUUID)
	return
}

// UnlinkPizzaCategory in a merchant category
func (c *catalogService) UnlinkPizzaCategory(merchantUUID, pizzaID, categoryID string) (err error) {
	if err = verifyCategoryItems(merchantUUID, "catalogID", categoryID); err != nil {
		glg.Error("[SDK] Catalog UnlinkPizzaCategory verifyCategoryItems: ", err.Error())
		return
	}
	if pizzaID == "" {
		err = ErrNoProductID
		glg.Error("[SDK] Catalog UnlinkPizzaCategory verifyFields: ", err.Error(), " merchant ", merchantUUID)
		return
	}
	if err = c.auth.Validate(); err != nil {
		glg.Error("[SDK] Catalog UnlinkPizzaCategory auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.auth.GetToken())
	endpoint := v2Endpoint + fmt.Sprintf(
		"/merchants/%s/pizzas/%s/categories/%s", merchantUUID, pizzaID, categoryID)
	_, status, err := c.adapter.DoRequest(http.MethodDelete, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Catalog UnlinkPizzaCategory adapter.DoRequest: ", err.Error())
		return
	}
	if status >= http.StatusBadRequest {
		glg.Error("[SDK] Catalog UnlinkPizzaCategory status code: ", status, " merchant: ", merchantUUID)
		err = fmt.Errorf(
			"Merchant '%s' could unlink pizza id '%s' from category '%s'",
			merchantUUID, pizzaID, categoryID)
		glg.Error("[SDK] Catalog UnlinkPizzaCategory err: ", err)
		return
	}
	glg.Infof("[SDK] Update pizza id '%s' success, merchant '%s'", pizzaID, merchantUUID)
	return
}
