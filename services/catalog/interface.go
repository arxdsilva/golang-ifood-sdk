package catalog

// Service describes the catalog abstraction
type Service interface {
	ListAllV2(merchantID string) (Catalogs, error)
	ListUnsellableItems(merchantUUID, catalogID string) (UnsellableResponse, error)
	ListAllCategoriesInCatalog(merchantUUID, catalogID string) (CategoryResponse, error)
	CreateCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template, externalCode string) (CategoryCreateResponse, error)
	GetCategoryInCatalog(merchantUUID, catalogID, categoryID string) (CategoryResponse, error)
	EditCategoryInCatalog(merchantUUID, catalogID, categoryID, name, resourceStatus, externalCode string, sequence int) (CategoryCreateResponse, error)
	DeleteCategoryInCatalog(merchantUUID, catalogID, categoryID string) error
	ListProducts(merchantUUID string) (Products, error)
	CreateProduct(merchantUUID string, product Product) (Product, error)
	EditProduct(merchantUUID string, product Product) (Product, error)
	DeleteProduct(merchantUUID, productID string) error
	UpdateProductStatus(merchantUUID, productID, productStatus string) error
	LinkProductToCategory(merchantUUID, categoryID string, product ProductLink) error
	CreatePizza(merchantUUID string, pizza Pizza) (Pizza, error)
	ListPizzas(merchantUUID string) (Pizzas, error)
	UpdatePizza(merchantUUID string, pizza Pizza) error
	UpdatePizzaStatus(merchantUUID, pizzaStatus, pizzaID string) error
	UnlinkPizzaCategory(merchantUUID, pizzaID, categoryID string) error
	LinkPizzaToCategory(merchantUUID, categoryID string, pizza Pizza) error
	UnlinkProductToCategory(merchantUUID, categoryID, productID string) error
}
