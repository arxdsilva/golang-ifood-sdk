package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/arxdsilva/golang-ifood-sdk/services/merchant"
	"github.com/kpango/glg"
)

const (
	V1Endpoint = "/v1.0/orders"
	V3Endpoint = "/v3.0/orders"
)

type (
	Service interface {
		GetDetails(reference string) (OrderDetails, error)
		SetIntegrateStatus(reference string) error
		SetConfirmStatus(reference string) error
		SetDispatchStatus(reference string) error
	}

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

	ordersService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

func New(adapter adapters.Http, authService auth.Service) *ordersService {
	return &ordersService{adapter, authService}
}

func (o *ordersService) GetDetails(orderReference string) (od OrderDetails, err error) {
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders GetDetails auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s", V3Endpoint, orderReference)
	resp, status, err := o.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders GetDetails adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Warn("[SDK] Orders GetDetails status code: ", status)
		err = fmt.Errorf("Order reference %s could not retrieve details", orderReference)
		return
	}
	return od, json.Unmarshal(resp, &od)
}

func (o *ordersService) SetIntegrateStatus(orderReference string) (err error) {
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders SetIntegrateStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/integration", V1Endpoint, orderReference)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders SetIntegrateStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders SetIntegrateStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference %s could not be integrated", orderReference)
		return
	}
	return
}

func (o *ordersService) SetConfirmStatus(orderReference string) (err error) {
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders SetConfirmStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/confirmation", V1Endpoint, orderReference)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders SetConfirmStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders SetConfirmStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference %s could not be confirmed", orderReference)
		return
	}
	return
}

func (o *ordersService) SetDispatchStatus(orderReference string) (err error) {
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders SetDispatchStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/dispatch", V1Endpoint, orderReference)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders SetDispatchStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders SetDispatchStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference %s could not be confirmed", orderReference)
		return
	}
	return
}
