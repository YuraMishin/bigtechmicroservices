// Code generated by ogen, DO NOT EDIT.

package order_v1

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (s *GenericErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Ref: #/components/schemas/bad_request_error
type BadRequestError struct {
	// HTTP-код ошибки.
	Code int `json:"code"`
	// Описание ошибки.
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *BadRequestError) GetCode() int {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *BadRequestError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *BadRequestError) SetCode(val int) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *BadRequestError) SetMessage(val string) {
	s.Message = val
}

func (*BadRequestError) cancelOrderByUUIDRes() {}
func (*BadRequestError) createNewOrderRes()    {}
func (*BadRequestError) getOrderByUUIDRes()    {}
func (*BadRequestError) payOrderRes()          {}

// CancelOrderByUUIDNoContent is response for CancelOrderByUUID operation.
type CancelOrderByUUIDNoContent struct{}

func (*CancelOrderByUUIDNoContent) cancelOrderByUUIDRes() {}

// Ref: #/components/schemas/conflict
type Conflict struct {
	// HTTP-код ошибки.
	Code int `json:"code"`
	// Описание ошибки.
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *Conflict) GetCode() int {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *Conflict) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *Conflict) SetCode(val int) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *Conflict) SetMessage(val string) {
	s.Message = val
}

func (*Conflict) cancelOrderByUUIDRes() {}

// Ref: #/components/schemas/create_order_request
type CreateOrderRequest struct {
	// UUID пользователя.
	UserUUID  uuid.UUID   `json:"user_uuid"`
	PartUuids []uuid.UUID `json:"part_uuids"`
}

// GetUserUUID returns the value of UserUUID.
func (s *CreateOrderRequest) GetUserUUID() uuid.UUID {
	return s.UserUUID
}

// GetPartUuids returns the value of PartUuids.
func (s *CreateOrderRequest) GetPartUuids() []uuid.UUID {
	return s.PartUuids
}

// SetUserUUID sets the value of UserUUID.
func (s *CreateOrderRequest) SetUserUUID(val uuid.UUID) {
	s.UserUUID = val
}

// SetPartUuids sets the value of PartUuids.
func (s *CreateOrderRequest) SetPartUuids(val []uuid.UUID) {
	s.PartUuids = val
}

// Ref: #/components/schemas/create_order_response
type CreateOrderResponse struct {
	// Уникальный идентификатор заказа.
	OrderUUID uuid.UUID `json:"order_uuid"`
	// Итоговая стоимость.
	TotalPrice float32 `json:"total_price"`
}

// GetOrderUUID returns the value of OrderUUID.
func (s *CreateOrderResponse) GetOrderUUID() uuid.UUID {
	return s.OrderUUID
}

// GetTotalPrice returns the value of TotalPrice.
func (s *CreateOrderResponse) GetTotalPrice() float32 {
	return s.TotalPrice
}

// SetOrderUUID sets the value of OrderUUID.
func (s *CreateOrderResponse) SetOrderUUID(val uuid.UUID) {
	s.OrderUUID = val
}

// SetTotalPrice sets the value of TotalPrice.
func (s *CreateOrderResponse) SetTotalPrice(val float32) {
	s.TotalPrice = val
}

func (*CreateOrderResponse) createNewOrderRes() {}

// Ref: #/components/schemas/generic_error
type GenericError struct {
	// HTTP-код ошибки.
	Code OptInt `json:"code"`
	// Описание ошибки.
	Message OptString `json:"message"`
}

// GetCode returns the value of Code.
func (s *GenericError) GetCode() OptInt {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *GenericError) GetMessage() OptString {
	return s.Message
}

// SetCode sets the value of Code.
func (s *GenericError) SetCode(val OptInt) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *GenericError) SetMessage(val OptString) {
	s.Message = val
}

// GenericErrorStatusCode wraps GenericError with StatusCode.
type GenericErrorStatusCode struct {
	StatusCode int
	Response   GenericError
}

// GetStatusCode returns the value of StatusCode.
func (s *GenericErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *GenericErrorStatusCode) GetResponse() GenericError {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *GenericErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *GenericErrorStatusCode) SetResponse(val GenericError) {
	s.Response = val
}

// Ref: #/components/schemas/internal_server_error
type InternalServerError struct {
	// HTTP-код ошибки.
	Code int `json:"code"`
	// Описание ошибки.
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *InternalServerError) GetCode() int {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *InternalServerError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *InternalServerError) SetCode(val int) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *InternalServerError) SetMessage(val string) {
	s.Message = val
}

func (*InternalServerError) cancelOrderByUUIDRes() {}
func (*InternalServerError) createNewOrderRes()    {}
func (*InternalServerError) getOrderByUUIDRes()    {}
func (*InternalServerError) payOrderRes()          {}

// Ref: #/components/schemas/not_found_error
type NotFoundError struct {
	// HTTP-код ошибки.
	Code int `json:"code"`
	// Описание ошибки.
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *NotFoundError) GetCode() int {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *NotFoundError) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *NotFoundError) SetCode(val int) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *NotFoundError) SetMessage(val string) {
	s.Message = val
}

func (*NotFoundError) cancelOrderByUUIDRes() {}
func (*NotFoundError) getOrderByUUIDRes()    {}
func (*NotFoundError) payOrderRes()          {}

// NewOptInt returns new OptInt with value set to v.
func NewOptInt(v int) OptInt {
	return OptInt{
		Value: v,
		Set:   true,
	}
}

// OptInt is optional int.
type OptInt struct {
	Value int
	Set   bool
}

// IsSet returns true if OptInt was set.
func (o OptInt) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt) Reset() {
	var v int
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt) SetTo(v int) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt) Get() (v int, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/order_dto
type OrderDto struct {
	// Уникальный идентификатор заказа.
	OrderUUID uuid.UUID `json:"order_uuid"`
	// Уникальный идентификатор пользователя.
	UserUUID uuid.UUID `json:"user_uuid"`
	// Список идентификаторов деталей.
	PartUuids []uuid.UUID `json:"part_uuids"`
	// Общая стоимость заказа.
	TotalPrice float32 `json:"total_price"`
	// Идентификатор транзакции оплаты.
	TransactionUUID uuid.UUID `json:"transaction_uuid"`
	// Способ оплаты.
	PaymentMethod OrderDtoPaymentMethod `json:"payment_method"`
	// Статус закаща.
	Status OrderDtoStatus `json:"status"`
}

// GetOrderUUID returns the value of OrderUUID.
func (s *OrderDto) GetOrderUUID() uuid.UUID {
	return s.OrderUUID
}

// GetUserUUID returns the value of UserUUID.
func (s *OrderDto) GetUserUUID() uuid.UUID {
	return s.UserUUID
}

// GetPartUuids returns the value of PartUuids.
func (s *OrderDto) GetPartUuids() []uuid.UUID {
	return s.PartUuids
}

// GetTotalPrice returns the value of TotalPrice.
func (s *OrderDto) GetTotalPrice() float32 {
	return s.TotalPrice
}

// GetTransactionUUID returns the value of TransactionUUID.
func (s *OrderDto) GetTransactionUUID() uuid.UUID {
	return s.TransactionUUID
}

// GetPaymentMethod returns the value of PaymentMethod.
func (s *OrderDto) GetPaymentMethod() OrderDtoPaymentMethod {
	return s.PaymentMethod
}

// GetStatus returns the value of Status.
func (s *OrderDto) GetStatus() OrderDtoStatus {
	return s.Status
}

// SetOrderUUID sets the value of OrderUUID.
func (s *OrderDto) SetOrderUUID(val uuid.UUID) {
	s.OrderUUID = val
}

// SetUserUUID sets the value of UserUUID.
func (s *OrderDto) SetUserUUID(val uuid.UUID) {
	s.UserUUID = val
}

// SetPartUuids sets the value of PartUuids.
func (s *OrderDto) SetPartUuids(val []uuid.UUID) {
	s.PartUuids = val
}

// SetTotalPrice sets the value of TotalPrice.
func (s *OrderDto) SetTotalPrice(val float32) {
	s.TotalPrice = val
}

// SetTransactionUUID sets the value of TransactionUUID.
func (s *OrderDto) SetTransactionUUID(val uuid.UUID) {
	s.TransactionUUID = val
}

// SetPaymentMethod sets the value of PaymentMethod.
func (s *OrderDto) SetPaymentMethod(val OrderDtoPaymentMethod) {
	s.PaymentMethod = val
}

// SetStatus sets the value of Status.
func (s *OrderDto) SetStatus(val OrderDtoStatus) {
	s.Status = val
}

func (*OrderDto) getOrderByUUIDRes() {}

// Способ оплаты.
type OrderDtoPaymentMethod string

const (
	OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED   OrderDtoPaymentMethod = "PAYMENT_METHOD_UNSPECIFIED"
	OrderDtoPaymentMethodPAYMENTMETHODCARD          OrderDtoPaymentMethod = "PAYMENT_METHOD_CARD"
	OrderDtoPaymentMethodPAYMENTMETHODSBP           OrderDtoPaymentMethod = "PAYMENT_METHOD_SBP"
	OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD    OrderDtoPaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY OrderDtoPaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

// AllValues returns all OrderDtoPaymentMethod values.
func (OrderDtoPaymentMethod) AllValues() []OrderDtoPaymentMethod {
	return []OrderDtoPaymentMethod{
		OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED,
		OrderDtoPaymentMethodPAYMENTMETHODCARD,
		OrderDtoPaymentMethodPAYMENTMETHODSBP,
		OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD,
		OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s OrderDtoPaymentMethod) MarshalText() ([]byte, error) {
	switch s {
	case OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED:
		return []byte(s), nil
	case OrderDtoPaymentMethodPAYMENTMETHODCARD:
		return []byte(s), nil
	case OrderDtoPaymentMethodPAYMENTMETHODSBP:
		return []byte(s), nil
	case OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD:
		return []byte(s), nil
	case OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *OrderDtoPaymentMethod) UnmarshalText(data []byte) error {
	switch OrderDtoPaymentMethod(data) {
	case OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED:
		*s = OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
		return nil
	case OrderDtoPaymentMethodPAYMENTMETHODCARD:
		*s = OrderDtoPaymentMethodPAYMENTMETHODCARD
		return nil
	case OrderDtoPaymentMethodPAYMENTMETHODSBP:
		*s = OrderDtoPaymentMethodPAYMENTMETHODSBP
		return nil
	case OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD:
		*s = OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD
		return nil
	case OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY:
		*s = OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Статус закаща.
type OrderDtoStatus string

const (
	OrderDtoStatusPENDINGPAYMENT OrderDtoStatus = "PENDING_PAYMENT"
	OrderDtoStatusPAID           OrderDtoStatus = "PAID"
	OrderDtoStatusCANCELLED      OrderDtoStatus = "CANCELLED"
)

// AllValues returns all OrderDtoStatus values.
func (OrderDtoStatus) AllValues() []OrderDtoStatus {
	return []OrderDtoStatus{
		OrderDtoStatusPENDINGPAYMENT,
		OrderDtoStatusPAID,
		OrderDtoStatusCANCELLED,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s OrderDtoStatus) MarshalText() ([]byte, error) {
	switch s {
	case OrderDtoStatusPENDINGPAYMENT:
		return []byte(s), nil
	case OrderDtoStatusPAID:
		return []byte(s), nil
	case OrderDtoStatusCANCELLED:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *OrderDtoStatus) UnmarshalText(data []byte) error {
	switch OrderDtoStatus(data) {
	case OrderDtoStatusPENDINGPAYMENT:
		*s = OrderDtoStatusPENDINGPAYMENT
		return nil
	case OrderDtoStatusPAID:
		*s = OrderDtoStatusPAID
		return nil
	case OrderDtoStatusCANCELLED:
		*s = OrderDtoStatusCANCELLED
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/pay_order_request
type PayOrderRequest struct {
	// Способ оплаты.
	PaymentMethod PayOrderRequestPaymentMethod `json:"payment_method"`
}

// GetPaymentMethod returns the value of PaymentMethod.
func (s *PayOrderRequest) GetPaymentMethod() PayOrderRequestPaymentMethod {
	return s.PaymentMethod
}

// SetPaymentMethod sets the value of PaymentMethod.
func (s *PayOrderRequest) SetPaymentMethod(val PayOrderRequestPaymentMethod) {
	s.PaymentMethod = val
}

// Способ оплаты.
type PayOrderRequestPaymentMethod string

const (
	PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED   PayOrderRequestPaymentMethod = "PAYMENT_METHOD_UNSPECIFIED"
	PayOrderRequestPaymentMethodPAYMENTMETHODCARD          PayOrderRequestPaymentMethod = "PAYMENT_METHOD_CARD"
	PayOrderRequestPaymentMethodPAYMENTMETHODSBP           PayOrderRequestPaymentMethod = "PAYMENT_METHOD_SBP"
	PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD    PayOrderRequestPaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY PayOrderRequestPaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

// AllValues returns all PayOrderRequestPaymentMethod values.
func (PayOrderRequestPaymentMethod) AllValues() []PayOrderRequestPaymentMethod {
	return []PayOrderRequestPaymentMethod{
		PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED,
		PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
		PayOrderRequestPaymentMethodPAYMENTMETHODSBP,
		PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD,
		PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s PayOrderRequestPaymentMethod) MarshalText() ([]byte, error) {
	switch s {
	case PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		return []byte(s), nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return []byte(s), nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return []byte(s), nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return []byte(s), nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *PayOrderRequestPaymentMethod) UnmarshalText(data []byte) error {
	switch PayOrderRequestPaymentMethod(data) {
	case PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		*s = PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED
		return nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		*s = PayOrderRequestPaymentMethodPAYMENTMETHODCARD
		return nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		*s = PayOrderRequestPaymentMethodPAYMENTMETHODSBP
		return nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		*s = PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD
		return nil
	case PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		*s = PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/pay_order_response
type PayOrderResponse struct {
	// Уникальный идентификатор транзакции на оплату.
	TransactionUUID uuid.UUID `json:"transaction_uuid"`
}

// GetTransactionUUID returns the value of TransactionUUID.
func (s *PayOrderResponse) GetTransactionUUID() uuid.UUID {
	return s.TransactionUUID
}

// SetTransactionUUID sets the value of TransactionUUID.
func (s *PayOrderResponse) SetTransactionUUID(val uuid.UUID) {
	s.TransactionUUID = val
}

func (*PayOrderResponse) payOrderRes() {}
