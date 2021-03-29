package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arxdsilva/golang-ifood-sdk/adapters"
	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/kpango/glg"
)

const (
	v1Endpoint = "/v1.0/orders"
	v2Endpoint = "/v2.0/orders"
	v3Endpoint = "/v3.0/orders"
)

var (
	// CancelCodes are all valid iFood API cancellation codes
	CancelCodes = map[string]string{
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
	// Service determinates the order's interface
	Service interface {
		GetDetails(reference string) (OrderDetails, error)
		SetIntegrateStatus(reference string) error
		SetConfirmStatus(reference string) error
		SetDispatchStatus(reference string) error
		SetReadyToDeliverStatus(reference string) error
		SetCancelStatus(reference, code string) error
		ClientCancellationStatus(reference string, accepted bool) error
		Tracking(orderUUID string) (TrackingResponse, error)
		DeliveryInformation(orderUUID string) (DeliveryInformationResponse, error)
	}

	ordersService struct {
		adapter adapters.Http
		auth    auth.Service
	}
)

// New returns a new order service
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
	endpoint := fmt.Sprintf("%s/%s", v3Endpoint, orderReference)
	resp, status, err := o.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders GetDetails adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusOK {
		glg.Error("[SDK] Orders GetDetails status code: ", status)
		err = fmt.Errorf("Order reference '%s' could not retrieve details", orderReference)
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
	endpoint := fmt.Sprintf("%s/%s/statuses/integration", v1Endpoint, orderReference)
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
	endpoint := fmt.Sprintf("%s/%s/statuses/confirmation", v1Endpoint, orderReference)
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
	endpoint := fmt.Sprintf("%s/%s/statuses/dispatch", v1Endpoint, orderReference)
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
	endpoint := fmt.Sprintf("%s/%s/statuses/readyToDeliver", v2Endpoint, orderReference)
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
	endpoint := fmt.Sprintf("%s/%s/statuses/cancellationRequested", v3Endpoint, orderReference)
	detail := CancelCodes[code]
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
		err = fmt.Errorf(
			"Order reference '%s' could not be set as 'cancelled' code '%s', detail '%s'",
			orderReference, code, detail)
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
		glg.Error("[SDK] Orders ClientCancellationStatus: ", err.Error())
		return
	}
	if err = o.auth.Validate(); err != nil {
		glg.Error("[SDK] Orders ClientCancellationStatus auth.Validate: ", err.Error())
		return
	}
	cancelStatus := "consumerCancellationDenied"
	if accepted {
		cancelStatus = "consumerCancellationAccepted"
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/statuses/%s", v2Endpoint, orderReference, cancelStatus)
	_, status, err := o.adapter.DoRequest(http.MethodPost, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders ClientCancellationStatus adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders ClientCancellationStatus status code: ", status, " orderReference: ", orderReference)
		err = fmt.Errorf(
			"Order reference '%s' could not set 'client cancellation' status '%s'",
			orderReference, cancelStatus)
		glg.Error("[SDK] Orders ClientCancellationStatus err: ", err)
		return
	}
	return
}

// Tracking retorna a posicao do entregador
func (o *ordersService) Tracking(orderUUID string) (tr TrackingResponse, err error) {
	if orderUUID == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders Tracking: ", orderUUID, " err: ", err.Error())
		return
	}
	if err = o.auth.Validate(); err != nil {
		glg.Error("[SDK] Orders Tracking auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/tracking", v2Endpoint, orderUUID)
	resp, status, err := o.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders Tracking adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders Tracking status code: ", status, " order uuid: ", orderUUID)
		err = fmt.Errorf("Order reference '%s' could not get tracking information", orderUUID)
		glg.Error("[SDK] Orders Tracking err: ", err)
		return
	}
	return tr, json.Unmarshal(resp, &tr)
}

// DeliveryInformation retorna informacoes da entrega
func (o *ordersService) DeliveryInformation(orderUUID string) (di DeliveryInformationResponse, err error) {
	if orderUUID == "" {
		err = ErrOrderReferenceNotSpecified
		glg.Error("[SDK] Orders DeliveryInformation: ", orderUUID, " err: ", err.Error())
		return
	}
	if err = o.auth.Validate(); err != nil {
		glg.Error("[SDK] Orders DeliveryInformation auth.Validate: ", err.Error())
		return
	}
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", o.auth.GetToken())
	endpoint := fmt.Sprintf("%s/%s/delivery-information", v2Endpoint, orderUUID)
	resp, status, err := o.adapter.DoRequest(http.MethodGet, endpoint, nil, headers)
	if err != nil {
		glg.Error("[SDK] Orders DeliveryInformation adapter.DoRequest error: ", err.Error())
		return
	}
	if status != http.StatusAccepted {
		glg.Error("[SDK] Orders DeliveryInformation status code: ", status, " order uuid: ", orderUUID)
		err = fmt.Errorf("Order uuid '%s' could get delivery information", orderUUID)
		glg.Error("[SDK] Orders DeliveryInformation err: ", err)
		return
	}
	return di, json.Unmarshal(resp, &di)
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
	_, ok := CancelCodes[code]
	if !ok {
		err = fmt.Errorf(
			"cancel code '%s' is invalid, verify docs: https://developer.ifood.com.br/reference#pedido-de-cancelamento-30",
			code)
		return
	}
	return
}
