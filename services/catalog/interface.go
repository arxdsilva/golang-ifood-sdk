package catalog

type Service interface {
	ListAll(merchantID string) (Catalogs, error)
	ListUnsellableItems(merchantUUID, catalogID string) (UnsellableResponse, error)
	ListAllCategoriesInCatalog(merchantUUID, catalogID string) (CategoryResponse, error)
	CreateCategoryInCatalog(merchantUUID, catalogID, name, resourceStatus, template, externalCode string) (CategoryCreateResponse, error)
	GetCategoryInCatalog(merchantUUID, catalogID, categoryID string) (CategoryResponse, error)
	EditCategoryInCatalog(merchantUUID, catalogID, categoryID, name, resourceStatus, externalCode string, sequence int) (CategoryCreateResponse, error)
	DeleteCategoryInCatalog(merchantUUID, catalogID, categoryID string) error
	ListProducts(merchantUUID string) (Products, error)
	CreateProduct(merchantUUID string, product Product) (Product, error)
	EditProduct(merchantUUID, productID string, product Product) (Product, error)
	DeleteProduct(merchantUUID, productID string) error
	UpdateProductStatus(merchantUUID, productID, productStatus string) error
}
