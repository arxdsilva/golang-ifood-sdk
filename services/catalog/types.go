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

	CategoryResponse struct {
		ID           string `json:"id"`
		Sequence     int    `json:"sequence"`
		Name         string `json:"name"`
		ExternalCode string `json:"externalCode"`
		Status       string `json:"status"`
		Items        []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Description  string `json:"description"`
			ExternalCode string `json:"externalCode"`
			Status       string `json:"status"`
			ProductID    string `json:"productId"`
			Sequence     int    `json:"sequence"`
			MagePath     string `json:"magePath"`
			Price        struct {
				Value         int `json:"value"`
				OriginalValue int `json:"originalValue"`
			} `json:"price"`
			Shifts []struct {
				StartTime string `json:"startTime"`
				EndTime   string `json:"endTime"`
				Monday    bool   `json:"monday"`
				Tuesday   bool   `json:"tuesday"`
				Wednesday bool   `json:"wednesday"`
				Thursday  bool   `json:"thursday"`
				Friday    bool   `json:"friday"`
				Saturday  bool   `json:"saturday"`
				Sunday    bool   `json:"sunday"`
			} `json:"shifts"`
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
					Price        struct {
						Value         int `json:"value"`
						OriginalValue int `json:"originalValue"`
					} `json:"price"`
				} `json:"options"`
			} `json:"optionGroups"`
			SellingOption struct {
			} `json:"sellingOption"`
		} `json:"items"`
		Template string `json:"template"`
		Pizza    struct {
			ID    string `json:"id"`
			Sizes []struct {
				ID                string `json:"id"`
				Name              string `json:"name"`
				Sequence          int    `json:"sequence"`
				Status            string `json:"status"`
				ExternalCode      string `json:"externalCode"`
				Slices            int    `json:"slices"`
				AcceptedFractions []int  `json:"acceptedFractions"`
				Price             struct {
					Value         int `json:"value"`
					OriginalValue int `json:"originalValue"`
				} `json:"price"`
			} `json:"sizes"`
			Crusts []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				Status       string `json:"status"`
				ExternalCode string `json:"externalCode"`
				Sequence     int    `json:"sequence"`
				Price        struct {
					Value         int `json:"value"`
					OriginalValue int `json:"originalValue"`
				} `json:"price"`
			} `json:"crusts"`
			Edges []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				Status       string `json:"status"`
				ExternalCode string `json:"externalCode"`
				Sequence     int    `json:"sequence"`
				Price        struct {
					Value         int `json:"value"`
					OriginalValue int `json:"originalValue"`
				} `json:"price"`
			} `json:"edges"`
			Toppings []struct {
				ID                  string   `json:"id"`
				ExternalCode        string   `json:"externalCode"`
				Name                string   `json:"name"`
				Description         string   `json:"description"`
				Image               string   `json:"image"`
				Status              string   `json:"status"`
				DietaryRestrictions []string `json:"dietaryRestrictions"`
				Sequence            int      `json:"sequence"`
				Prices              struct {
					AdditionalProp struct {
						Value         int `json:"value"`
						OriginalValue int `json:"originalValue"`
					} `json:"additionalProp"`
				} `json:"prices"`
			} `json:"toppings"`
			Shifts []struct {
				StartTime string `json:"startTime"`
				EndTime   string `json:"endTime"`
				Monday    bool   `json:"monday"`
				Tuesday   bool   `json:"tuesday"`
				Wednesday bool   `json:"wednesday"`
				Thursday  bool   `json:"thursday"`
				Friday    bool   `json:"friday"`
				Saturday  bool   `json:"saturday"`
				Sunday    bool   `json:"sunday"`
			} `json:"shifts"`
		} `json:"pizza"`
	}
)
