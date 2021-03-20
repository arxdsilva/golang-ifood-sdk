package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/arxdsilva/golang-ifood-sdk/services/merchant"
	"github.com/kpango/glg"
)

const (
	V1Endpoint = "/v1.0/orders"
	V2Endpoint = "/v2.0/orders"
	V3Endpoint = "/v3.0/orders"
)

var (
	ErrOrderReferenceNotSpecified = errors.New("Order reference not specified")
	ErrCancelCodeNotSpecified     = errors.New("Order cancel code not specified")
	cancelCodes                   = map[string]string{
		"501": "PROBLEMAS DE SISTEMA",
		"502": "PEDIDO EM DUPLICIDADE",
		"503": "ITEM INDISPONÍVEL",
		"504": "RESTAURANTE SEM MOTOBOY",
		"505": "CARDÁPIO DESATUALIZADO",
		"506": "PEDIDO FORA DA ÁREA DE ENTREGA",
		"507": "CLIENTE GOLPISTA / TROTE",
		"508": "FORA DO HORÁRIO DO DELIVERY",
		"509": "DIFICULDADES INTERNAS DO RESTAURANTE",
		"511": "ÁREA DE RISCO",
		"512": "RESTAURANTE ABRIRÁ MAIS TARDE",
		"513": "RESTAURANTE FECHOU MAIS CEDO",
		"803": "ITEM INDISPONÍVEL",
		"805": "RESTAURANTE SEM MOTOBOY",
		"801": "PROBLEMAS DE SISTEMA",
		"804": "CADASTRO DO CLIENTE INCOMPLETO - CLIENTE NÃO ATENDE",
		"807": "PEDIDO FORA DA ÁREA DE ENTREGA",
		"808": "CLIENTE GOLPISTA / TROTE",
		"809": "FORA DO HORÁRIO DO DELIVERY",
		"815": "DIFICULDADES INTERNAS DO RESTAURANTE",
		"818": "TAXA DE ENTREGA INCONSISTENTE",
		"820": "ÁREA DE RISCO",
	}
)

type (
	Service interface {
		GetDetails(reference string) (OrderDetails, error)
		SetIntegrateStatus(reference string) error
		SetConfirmStatus(reference string) error
		SetDispatchStatus(reference string) error
		SetReadyToDeliverStatus(reference string) error
		SetCancelStatus(reference, code string) error
		ClientCancellationStatus(reference string, accepted bool) error
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

	cancelOrder struct {
		Code    string `json:"cancellationCode"`
		Details string `json:"details"`
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
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders GetDetails: ", err.Error())
		return
	}
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
		glg.Error("[SDK] Orders GetDetails status code: ", status)
		err = fmt.Errorf("Order reference %s could not retrieve details", orderReference)
		glg.Error("[SDK] Orders GetDetails err: ", err)
		return
	}
	return od, json.Unmarshal(resp, &od)
}

func (o *ordersService) SetIntegrateStatus(orderReference string) (err error) {
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders SetIntegrateStatus: ", err.Error())
		return
	}
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
		glg.Error("[SDK] Orders SetIntegrateStatus err: ", err)
		return
	}
	return
}

func (o *ordersService) SetConfirmStatus(orderReference string) (err error) {
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders SetConfirmStatus: ", err.Error())
		return
	}
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
		err = fmt.Errorf("Order reference '%s' could not be confirmed", orderReference)
		glg.Error("[SDK] Orders SetConfirmStatus err: ", err)
		return
	}
	return
}

func (o *ordersService) SetDispatchStatus(orderReference string) (err error) {
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders SetDispatchStatus: ", err.Error())
		return
	}
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
		err = fmt.Errorf("Order reference '%s' could not be dispatched", orderReference)
		glg.Error("[SDK] Orders SetDispatchStatus err: ", err)
		return
	}
	return
}

func (o *ordersService) SetReadyToDeliverStatus(orderReference string) (err error) {
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders SetReadyToDeliverStatus: ", err.Error())
		return
	}
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders SetReadyToDeliverStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/readyToDeliver", V2Endpoint, orderReference)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders SetReadyToDeliverStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders SetReadyToDeliverStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference '%s' could not be set as 'ready to deliver'", orderReference)
		glg.Error("[SDK] Orders SetReadyToDeliverStatus err: ", err)
		return
	}
	return
}

func (o *ordersService) SetCancelStatus(orderReference, code string) (err error) {
	if err = verifyCancel(orderReference, code); err != nil {
		glg.Error("[SDK] Orders SetCancelStatus verifyCancel: ", err.Error())
		return
	}
	err = o.auth.Validate()
	if err != nil {
		glg.Error("[SDK] Orders SetCancelStatus auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/cancellationRequested", V3Endpoint, orderReference)
	detail := cancelCodes[code]
	co := cancelOrder{Code: code, Details: detail}
	reader, err := httpadapter.NewJsonReader(co)
	if err != nil {
		glg.Error("[SDK] Orders SetCancelStatus NewJsonReader error: ", err.Error())
		return
	}
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, reader, headers)
	if err != nil {
		glg.Error("[SDK] Orders SetCancelStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders SetCancelStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference '%s' could not be set as 'ready to deliver'", orderReference)
		glg.Error("[SDK] Orders SetCancelStatus err: ", err)
		return
	}
	return
}

// ClientCancellationStatus lida com o cancelamento do pedido por parte do cliente
//
// link: https://developer.ifood.com.br/reference#handshake-cancelamento
//
// reference: order reference id
// accepted: aceitacao pelo e-PDV do cancelamento do pedido
func (o *ordersService) ClientCancellationStatus(orderReference string, accepted bool) (err error) {
	if orderReference == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders SetReadyToDeliverStatus: ", err.Error())
		return
	}
	if err = o.auth.Validate(); err != nil {
		glg.Error("[SDK] Orders AcceptCancelStatus auth.Validate: ", err.Error())
		return
	}
	cancelStatus := "consumerCancellationDenied"
	if accepted {
		cancelStatus = "consumerCancellationAccepted"
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/%s", V2Endpoint, orderReference, cancelStatus)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders AcceptCancelStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders AcceptCancelStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf("Order reference '%s' could not be set as 'ready to deliver'", orderReference)
		glg.Error("[SDK] Orders AcceptCancelStatus err: ", err)
		return
	}
	return
}

func verifyCancel(reference, code string) (err error) {
	if reference == "" {
		err = ErrOrderReferenceNotSpecified
		return
	}
	if code == "" {
		err = ErrCancelCodeNotSpecified
		return
	}
	_, ok := cancelCodes[code]
	if !ok {
		err = fmt.Errorf(
			"cancel code '%s' is invalid, verify docs: https://developer.ifood.com.br/reference#pedido-de-cancelamento-30",
			code)
		return
	}
	return
}
