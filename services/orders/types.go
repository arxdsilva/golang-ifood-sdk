package orders

import "github.com/arxdsilva/golang-ifood-sdk/services/merchant"

type (
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

	Payment struct {
		Name      string `json:"name"`
		Code      string `json:"code"`
		Value     string `json:"value"`
		Prepaid   string `json:"prepaid"`
		Issuer    string `json:"issuer"`
		Collector string `json:"collector,omitempty"`
	}

	Customer struct {
		ID                           string `json:"id"`
		UUID                         string `json:"uuid"`
		Name                         string `json:"name"`
		Taxpayeridentificationnumber string `json:"taxPayerIdentificationNumber"`
		Phone                        string `json:"phone"`
		Orderscountonrestaurant      string `json:"ordersCountOnRestaurant"`
	}

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

	Coordinates struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	Item struct {
		Name          string     `json:"name"`
		Quantity      string     `json:"quantity"`
		Price         string     `json:"price"`
		Subitemsprice string     `json:"subItemsPrice"`
		Totalprice    string     `json:"totalPrice"`
		Discount      string     `json:"discount"`
		Addition      string     `json:"addition"`
		Externalcode  string     `json:"externalCode,omitempty"`
		Subitems      []Subitems `json:"subItems,omitempty"`
		Observations  string     `json:"observations,omitempty"`
	}

	Subitems struct {
		Name         string `json:"name"`
		Quantity     string `json:"quantity"`
		Price        string `json:"price"`
		Totalprice   string `json:"totalPrice"`
		Discount     string `json:"discount"`
		Addition     string `json:"addition"`
		Externalcode string `json:"externalCode"`
	}

	cancelOrder struct {
		Code    string `json:"cancellationCode"`
		Details string `json:"details"`
	}
)
