package catalog

type (
	Catalogs []Catalog
	Catalog  struct {
		ID         string   `json:"catalogId"`
		Context    []string `json:"context"`
		Status     string   `json:"status"`
		ModifiedAt float64  `json:"modifiedAt"`
	}

	UnsellableResponse struct {
		Categories []Category `json:"categories"`
	}
	Category struct {
		ID                   string               `json:"id"`
		Status               string               `json:"status"`
		Template             string               `json:"template"`
		Restrictions         []string             `json:"restrictions"`
		UnsellableItems      []UnsellableItem     `json:"unsellableItems"`
		UnsellablePizzaItems UnsellablePizzaItems `json:"unsellablePizzaItems"`
	}
	UnsellableItem struct {
		ID           string   `json:"id"`
		ProductID    string   `json:"productId"`
		Restrictions []string `json:"restrictions"`
	}
	UnsellablePizzaItems struct {
		Toppings []UnsellableItem `json:"toppings"`
		Crusts   []UnsellableItem `json:"crusts"`
		Edges    []UnsellableItem `json:"edges"`
		Sizes    []UnsellableItem `json:"sizes"`
	}
)
