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
)
