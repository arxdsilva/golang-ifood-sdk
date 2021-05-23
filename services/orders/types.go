package orders

import (
	"time"

	"github.com/arxdsilva/golang-ifood-sdk/services/merchant"
)

type (
	// OrderDetails endpoint return
	OrderDetails struct {
		ID                       string            `json:"id"`
		Reference                string            `json:"reference"`
		Shortreference           string            `json:"shortReference"`
		Createdat                string            `json:"createdAt"`
		Type                     string            `json:"type"`
		Merchant                 merchant.Merchant `json:"merchant"`
		Payments                 []Payment         `json:"payments"`
		Customer                 Customer          `json:"customer"`
		Items                    []Item            `json:"items"`
		Subtotal                 string            `json:"subTotal"`
		Totalprice               string            `json:"totalPrice"`
		Deliveryfee              string            `json:"deliveryFee"`
		Deliveryaddress          DeliveryAddress   `json:"deliveryAddress"`
		Deliverydatetime         string            `json:"deliveryDateTime"`
		Preparationtimeinseconds string            `json:"preparationTimeInSeconds"`
	}

	// Payment details
	Payment struct {
		Name      string `json:"name"`
		Code      string `json:"code"`
		Value     string `json:"value"`
		Prepaid   string `json:"prepaid"`
		Issuer    string `json:"issuer"`
		Collector string `json:"collector,omitempty"`
	}

	// Customer details
	Customer struct {
		ID                           string `json:"id"`
		UUID                         string `json:"uuid"`
		Name                         string `json:"name"`
		Taxpayeridentificationnumber string `json:"taxPayerIdentificationNumber"`
		Phone                        string `json:"phone"`
		Orderscountonrestaurant      string `json:"ordersCountOnRestaurant"`
	}

	V2Customer struct {
		Phone struct {
			Number              string    `json:"number"`
			Localizer           string    `json:"localizer"`
			Localizerexpiration time.Time `json:"localizerExpiration"`
		} `json:"phone"`
		Documentnumber        string `json:"documentNumber"`
		Name                  string `json:"name"`
		Orderscountonmerchant int    `json:"ordersCountOnMerchant"`
		ID                    string `json:"id"`
	}

	// DeliveryAddress from customer
	DeliveryAddress struct {
		Formattedaddress string      `json:"formattedAddress"`
		Country          string      `json:"country"`
		State            string      `json:"state"`
		City             string      `json:"city"`
		Coordinates      Coordinates `json:"coordinates"`
		Neighborhood     string      `json:"neighborhood"`
		Streetname       string      `json:"streetName"`
		Streetnumber     string      `json:"streetNumber"`
		Postalcode       string      `json:"postalCode"`
		Reference        string      `json:"reference"`
		Complement       string      `json:"complement"`
	}

	// Coordinates of a delivery
	Coordinates struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	// Item of the order
	Item struct {
		Name          string    `json:"name"`
		Quantity      string    `json:"quantity"`
		Price         string    `json:"price"`
		Subitemsprice string    `json:"subItemsPrice"`
		Totalprice    string    `json:"totalPrice"`
		Discount      string    `json:"discount"`
		Addition      string    `json:"addition"`
		Externalcode  string    `json:"externalCode,omitempty"`
		Subitems      []Subitem `json:"subItems,omitempty"`
		Observations  string    `json:"observations,omitempty"`
	}

	// Subitem of the order
	Subitem struct {
		Name         string `json:"name"`
		Quantity     string `json:"quantity"`
		Price        string `json:"price"`
		Totalprice   string `json:"totalPrice"`
		Discount     string `json:"discount"`
		Addition     string `json:"addition"`
		Externalcode string `json:"externalCode"`
	}

	// TrackingResponse API response of tracking
	TrackingResponse struct {
		Date             interface{} `json:"date"`
		DeliveryTime     time.Time   `json:"deliveryTime"`
		Eta              int         `json:"eta"`
		EtaToDestination int         `json:"etaToDestination"`
		EtaToOrigin      int         `json:"etaToOrigin"`
		Latitude         float64     `json:"latitude"`
		Longitude        float64     `json:"longitude"`
		OrderID          string      `json:"orderId"`
		TrackDate        time.Time   `json:"trackDate"`
	}

	// DeliveryInformationResponse API response of the delivery
	DeliveryInformationResponse struct {
		ExternalID         string      `json:"externalId"`
		OrderStatus        string      `json:"orderStatus"`
		WorkerName         string      `json:"workerName"`
		WorkerPhone        string      `json:"workerPhone"`
		WorkerPhoto        string      `json:"workerPhoto"`
		VehicleType        string      `json:"vehicleType"`
		VehiclePlateNumber interface{} `json:"vehiclePlateNumber"`
		LogisticCompany    string      `json:"logisticCompany"`
		Latitude           float64     `json:"latitude"`
		Longitude          float64     `json:"longitude"`
		Eta                int         `json:"eta"`
	}
	cancelOrder struct {
		Code    string `json:"cancellationCode"`
		Details string `json:"details"`
	}

	v2CancelOrder struct {
		Reason           string `json:"reason"`
		CancellationCode string `json:"cancellationCode"`
	}

	V2OrderDetails struct {
		Benefits                 []V2Benefit           `json:"benefits"`
		Ordertype                string                `json:"orderType"`
		Payments                 V2Payments            `json:"payments"`
		Merchant                 V2MerchantInfos       `json:"merchant"`
		SalesChannel             string                `json:"salesChannel"`
		OrderTiming              string                `json:"orderTiming"`
		CreatedAt                time.Time             `json:"createdAt"`
		Total                    V2OrderValues         `json:"total"`
		PreparationStartDatetime time.Time             `json:"preparationStartDateTime"`
		ID                       string                `json:"id"`
		DisplayID                string                `json:"displayId"`
		Items                    []V2Item              `json:"items"`
		Customer                 V2Customer            `json:"customer"`
		ExtraInfo                string                `json:"extraInfo"`
		Delivery                 V2DeliveryInformation `json:"delivery"`
		Schedule                 V2Schedule            `json:"schedule"`
		Indoor                   V2Indoor              `json:"indoor"`
		Takeout                  V2Takeout             `json:"takeout"`
	}

	V2Item struct {
		Unitprice    int    `json:"unitPrice"`
		Quantity     int    `json:"quantity"`
		Externalcode string `json:"externalCode"`
		Totalprice   int    `json:"totalPrice"`
		Index        int    `json:"index"`
		Unit         string `json:"unit"`
		Ean          string `json:"ean"`
		Price        int    `json:"price"`
		Observations string `json:"observations"`
		Name         string `json:"name"`
		Options      []struct {
			Unitprice    int    `json:"unitPrice"`
			Unit         string `json:"unit"`
			Ean          string `json:"ean"`
			Quantity     int    `json:"quantity"`
			Externalcode string `json:"externalCode"`
			Price        int    `json:"price"`
			Name         string `json:"name"`
			Index        int    `json:"index"`
			ID           string `json:"id"`
		} `json:"options"`
		ID           string `json:"id"`
		Optionsprice int    `json:"optionsPrice"`
	}

	V2DeliveryInformation struct {
		Mode            string `json:"mode"`
		Deliveredby     string `json:"deliveredBy"`
		Deliveryaddress struct {
			Reference        string `json:"reference"`
			Country          string `json:"country"`
			Streetname       string `json:"streetName"`
			Formattedaddress string `json:"formattedAddress"`
			Streetnumber     string `json:"streetNumber"`
			City             string `json:"city"`
			Postalcode       string `json:"postalCode"`
			Coordinates      struct {
				Latitude  int `json:"latitude"`
				Longitude int `json:"longitude"`
			} `json:"coordinates"`
			Neighborhood string `json:"neighborhood"`
			State        string `json:"state"`
			Complement   string `json:"complement"`
		} `json:"deliveryAddress"`
		Deliverydatetime time.Time `json:"deliveryDateTime"`
	}

	V2Payments struct {
		Methods []struct {
			Wallet struct {
				Name string `json:"name"`
			} `json:"wallet"`
			Method   string `json:"method"`
			Prepaid  bool   `json:"prepaid"`
			Currency string `json:"currency"`
			Type     string `json:"type"`
			Value    int    `json:"value"`
			Cash     struct {
				Changefor int `json:"changeFor"`
			} `json:"cash"`
			Card struct {
				Brand string `json:"brand"`
			} `json:"card"`
		} `json:"methods"`
		Pending int `json:"pending"`
		Prepaid int `json:"prepaid"`
	}

	V2OrderValues struct {
		Benefits    int `json:"benefits"`
		Deliveryfee int `json:"deliveryFee"`
		Orderamount int `json:"orderAmount"`
		Subtotal    int `json:"subTotal"`
	}

	V2Benefit struct {
		Targetid          string `json:"targetId"`
		Sponsorshipvalues []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"sponsorshipValues"`
		Value  int    `json:"value"`
		Target string `json:"target"`
	}

	V2MerchantInfos struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	V2Schedule struct {
		Deliverydatetimestart time.Time `json:"deliveryDateTimeStart"`
		Deliverydatetimeend   time.Time `json:"deliveryDateTimeEnd"`
	}

	V2Indoor struct {
		Mode             string    `json:"mode"`
		Deliverydatetime time.Time `json:"deliveryDateTime"`
		Table            string    `json:"table"`
	}

	V2Takeout struct {
		Mode            string    `json:"mode"`
		Takeoutdatetime time.Time `json:"takeoutDateTime"`
	}
)
