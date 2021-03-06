package catalog

type (
	// Catalogs group of Catalog
	Catalogs []Catalog

	// Catalog API response
	Catalog struct {
		ID         string   `json:"catalogId"`
		Context    []string `json:"context"`
		Status     string   `json:"status"`
		ModifiedAt string   `json:"modifiedAt"`
	}

	// UnsellableResponse API response
	UnsellableResponse struct {
		Categories []Category `json:"categories"`
	}

	// Category struct
	Category struct {
		ID                   string           `json:"id"`
		Status               string           `json:"status"`
		Template             string           `json:"template"`
		Restrictions         []string         `json:"restrictions"`
		UnsellableItems      []UnsellableItem `json:"unsellableItems"`
		UnsellablePizzaItems Pizza            `json:"unsellablePizzaItems"`
	}

	// UnsellableItem part of Category struct
	UnsellableItem struct {
		ID           string   `json:"id"`
		ProductID    string   `json:"productId"`
		Restrictions []string `json:"restrictions"`
	}

	// CategoryResponse from API when creating
	CategoryResponse struct {
		ID           string `json:"id"`
		Sequence     int    `json:"sequence"`
		Name         string `json:"name"`
		ExternalCode string `json:"externalCode"`
		Status       string `json:"status"`
		Items        []Item `json:"items"`
		Template     string `json:"template"`
		Pizza        Pizza  `json:"pizza"`
	}

	// Item product description
	Item struct {
		ID                  string   `json:"id"`
		Name                string   `json:"name"`
		Description         string   `json:"description"`
		ExternalCode        string   `json:"externalCode"`
		Status              string   `json:"status"`
		ProductID           string   `json:"productId"`
		Sequence            int      `json:"sequence"`
		MagePath            string   `json:"magePath"`
		Price               Price    `json:"price"`
		Shifts              []Shift  `json:"shifts"`
		Serving             string   `json:"serving"`
		DietaryRestrictions []string `json:"dietaryRestrictions"`
		Ean                 string   `json:"ean"`
		OptionGroups        []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			ExternalCode string `json:"externalCode"`
			Status       string `json:"status"`
			Sequence     int    `json:"sequence"`
			Min          int    `json:"min"`
			Max          int    `json:"max"`
			Options      struct {
				ID           string `json:"id"`
				Status       string `json:"status"`
				Sequence     int    `json:"sequence"`
				ProductID    string `json:"productId"`
				Name         string `json:"name"`
				Description  string `json:"description"`
				ExternalCode string `json:"externalCode"`
				ImagePath    string `json:"imagePath"`
				Price        Price  `json:"price"`
			} `json:"options"`
		} `json:"optionGroups"`
		// SellingOption struct {
		// } `json:"sellingOption"`
	}

	// CategoryItem linked product to a category
	CategoryItem struct {
		ID                  string    `json:"id"`
		Name                string    `json:"name"`
		Status              string    `json:"status"`
		ExternalCode        string    `json:"externalCode"`
		Template            string    `json:"template"`
		AcceptedFractions   []float64 `json:"acceptedFractions"`
		DietaryRestrictions []string  `json:"dietaryRestrictions"`
		Sequence            int       `json:"sequence"`
		Price               Price     `json:"price"`
		Shifts              []Shift   `json:"shifts"`
	}

	// CategoryCreateResponse create API response
	CategoryCreateResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ExternalCode string `json:"externalCode"`
		Status       string `json:"status"`
		Sequence     int    `json:"sequence"`
		Template     string `json:"template"`
	}

	// Price of a product object
	Price struct {
		Value         float64 `json:"value"`
		OriginalValue float64 `json:"originalValue"`
	}

	// Shift of a Item or CategoryItem
	Shift struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Monday    bool   `json:"monday"`
		Tuesday   bool   `json:"tuesday"`
		Wednesday bool   `json:"wednesday"`
		Thursday  bool   `json:"thursday"`
		Friday    bool   `json:"friday"`
		Saturday  bool   `json:"saturday"`
		Sunday    bool   `json:"sunday"`
	}

	apiError struct {
		Details struct {
			Code string `json:"code"`
		} `json:"details"`
	}
)
