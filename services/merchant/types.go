package merchant

type (
	Merchant struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Phones  []string `json:"phones"`
		Address Address  `json:"address"`
	}

	Address struct {
		Formattedaddress string `json:"formattedAddress"`
		Country          string `json:"country"`
		State            string `json:"state"`
		City             string `json:"city"`
		Neighborhood     string `json:"neighborhood"`
		Streetname       string `json:"streetName"`
		Streetnumber     string `json:"streetNumber"`
		Postalcode       string `json:"postalCode"`
	}

	Unavailabilities []Unavailability

	Unavailability struct {
		ID          string `json:"id"`
		Storeid     string `json:"storeId"`
		Description string `json:"description"`
		Authorid    string `json:"authorId"`
		Start       string `json:"start"`
		End         string `json:"end"`
	}

	unavailability struct {
		Description string `json:"description"`
		Minutes     int32  `json:"minutes"`
	}

	UnavailabilityResponse struct {
		ID          string `json:"id"`
		Storeid     string `json:"storeId"`
		Description string `json:"description"`
		Authorid    string `json:"authorId"`
		Start       string `json:"start"`
		End         string `json:"end"`
	}

	AvailabilityResponse []Availability
	Availability         struct {
		Context    string `json:"context"`
		Available  bool   `json:"available"`
		State      string `json:"state"`
		Reopenable struct {
			// Identifier interface{} `json:"identifier"`
			// Type       interface{} `json:"type"`
			Reopenable bool `json:"reopenable"`
		} `json:"reopenable"`
		Validations []struct {
			ID      string `json:"id"`
			Code    string `json:"code"`
			State   string `json:"state"`
			Message struct {
				Title       string `json:"title"`
				Subtitle    string `json:"subtitle"`
				Description string `json:"description"`
				Priority    int    `json:"priority"`
			} `json:"message"`
		} `json:"validations"`
		Message struct {
			Title       string `json:"title"`
			Subtitle    string `json:"subtitle"`
			Description string `json:"description"`
			Priority    int    `json:"priority"`
		} `json:"message"`
	}
)
