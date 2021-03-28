package catalog

// CreateItem product-catalog association
//
// POST
//
// 201 created
// 400 bad req
// 404 not found
// 409 conflict
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
func (c *catalogService) CreateItem(merchantID, categoryID, productID string) (cp Product, err error) {

	return
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
func (c *catalogService) EditItem(merchantID, categoryID, productID string) (cp Product, err error) {

	return
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

	return
}
