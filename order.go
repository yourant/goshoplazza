package goshoplazza

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const ordersBasePath = "orders"
const ordersResourceName = "orders"

// OrderService is an interface for interfacing with the orders endpoints of
// the Shopify API.
// See: https://help.shopify.com/api/reference/order
type OrderService interface {
	List(interface{}) ([]Order, error)
	Count(interface{}) (int, error)
	Get(string, interface{}) (*Order, error)
	Create(Order) (*Order, error)
	Update(Order) (*Order, error)

	// MetafieldsService used for Order resource to communicate with Metafields resource
	// MetafieldsService

	// FulfillmentsService used for Order resource to communicate with Fulfillments resource
	FulfillmentsService
}

// OrderServiceOp handles communication with the order related methods of the
// Shopify API.
type OrderServiceOp struct {
	client *Client
}

// A struct for all available order count options
type OrderCountOptions struct {
	Page              int       `url:"page,omitempty"`
	Limit             int       `url:"limit,omitempty"`
	SinceID           int64     `url:"since_id,omitempty"`
	CreatedAtMin      time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax      time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin      time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax      time.Time `url:"updated_at_max,omitempty"`
	Order             string    `url:"order,omitempty"`
	Fields            string    `url:"fields,omitempty"`
	Status            string    `url:"status,omitempty"`
	FinancialStatus   string    `url:"financial_status,omitempty"`
	FulfillmentStatus string    `url:"fulfillment_status,omitempty"`
}

// A struct for all available order list options.
// See: https://help.shopify.com/api/reference/order#index
type OrderListOptions struct {
	Page              int       `url:"page,omitempty"`
	Limit             int       `url:"limit,omitempty"`
	SinceID           int64     `url:"since_id,omitempty"`
	Status            string    `url:"status,omitempty"`
	FinancialStatus   string    `url:"financial_status,omitempty"`
	FulfillmentStatus string    `url:"fulfillment_status,omitempty"`
	CreatedAtMin      time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax      time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin      time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax      time.Time `url:"updated_at_max,omitempty"`
	ProcessedAtMin    time.Time `url:"processed_at_min,omitempty"`
	ProcessedAtMax    time.Time `url:"processed_at_max,omitempty"`
	Fields            string    `url:"fields,omitempty"`
	Order             string    `url:"order,omitempty"`
}

// Order represents a Shopify order
type Order struct {
	ID            string     `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Email         string     `json:"email,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	// Customer              *Customer        `json:"customer,omitempty"`
	BillingAddress        *Address         `json:"billing_address,omitempty"`
	ShippingAddress       *Address         `json:"shipping_address,omitempty"`
	Currency              string           `json:"currency,omitempty"`
	TotalPrice            *decimal.Decimal `json:"total_price,omitempty"`
	SubtotalPrice         *decimal.Decimal `json:"subtotal_price,omitempty"`
	TotalDiscounts        *decimal.Decimal `json:"total_discount,omitempty"`
	TotalShipping         *decimal.Decimal `json:"total_shipping,omitempty"`
	TotalLineItemsPrice   *decimal.Decimal `json:"total_line_items_price,omitempty"`
	TaxesIncluded         bool             `json:"taxes_included,omitempty"`
	TotalTax              *decimal.Decimal `json:"total_tax,omitempty"`
	TaxLines              []TaxLine        `json:"tax_lines,omitempty"`
	TotalWeight           int              `json:"total_weight,omitempty"`
	FinancialStatus       string           `json:"financial_status,omitempty"`
	Fulfillments          []Fulfillment    `json:"fulfillments,omitempty"`
	FulfillmentStatus     string           `json:"fulfillment_status,omitempty"`
	Token                 string           `json:"token,omitempty"`
	CartToken             string           `json:"cart_token,omitempty"`
	Number                string           `json:"number,omitempty"`
	OrderNumber           int              `json:"order_number,omitempty"`
	Note                  string           `json:"note,omitempty"`
	Test                  bool             `json:"test,omitempty"`
	BrowserIp             string           `json:"browser_ip,omitempty"`
	BuyerAcceptsMarketing bool             `json:"buyer_accepts_marketing,omitempty"`
	CancelReason          string           `json:"cancel_reason,omitempty"`
	NoteAttributes        []NoteAttribute  `json:"note_attributes,omitempty"`
	DiscountCodes         []DiscountCode   `json:"discount_codes,omitempty"`
	LineItems             []LineItem       `json:"line_items,omitempty"`
	ShippingLine          *ShippingLine    `json:"shipping_line,omitempty"`
	Transactions          []Transaction    `json:"transactions,omitempty"`
	AppID                 int              `json:"app_id,omitempty"`
	CustomerLocale        string           `json:"customer_locale,omitempty"`
	LandingSite           string           `json:"landing_site,omitempty"`
	ReferringSite         string           `json:"referring_site,omitempty"`
	SourceName            string           `json:"source_name,omitempty"`
	ClientDetails         *ClientDetails   `json:"client_details,omitempty"`
	Tags                  string           `json:"tags,omitempty"`
	LocationId            int64            `json:"location_id,omitempty"`
	PaymentGatewayNames   []string         `json:"payment_gateway_names,omitempty"`
	ProcessingMethod      string           `json:"processing_method,omitempty"`
	Refunds               []Refund         `json:"refunds,omitempty"`
	UserId                int64            `json:"user_id,omitempty"`
	OrderStatusUrl        string           `json:"order_status_url,omitempty"`
	Gateway               string           `json:"gateway,omitempty"`
	Confirmed             bool             `json:"confirmed,omitempty"`
	TotalPriceUSD         *decimal.Decimal `json:"total_price_usd,omitempty"`
	CheckoutToken         string           `json:"checkout_token,omitempty"`
	Reference             string           `json:"reference,omitempty"`
	SourceIdentifier      string           `json:"source_identifier,omitempty"`
	SourceURL             string           `json:"source_url,omitempty"`
	DeviceID              int64            `json:"device_id,omitempty"`
	Phone                 string           `json:"phone,omitempty"`
	LandingSiteRef        string           `json:"landing_site_ref,omitempty"`
	CheckoutID            int64            `json:"checkout_id,omitempty"`
	ContactEmail          string           `json:"contact_email,omitempty"`
	// Metafields            []Metafield      `json:"metafields,omitempty"`
}

type Address struct {
	ID           int64   `json:"id,omitempty"`
	Address1     string  `json:"address1,omitempty"`
	Address2     string  `json:"address2,omitempty"`
	City         string  `json:"city,omitempty"`
	Company      string  `json:"company,omitempty"`
	Country      string  `json:"country,omitempty"`
	CountryCode  string  `json:"country_code,omitempty"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	Latitude     float64 `json:"latitude,omitempty"`
	Longitude    float64 `json:"longitude,omitempty"`
	Name         string  `json:"name,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Province     string  `json:"province,omitempty"`
	ProvinceCode string  `json:"province_code,omitempty"`
	Zip          string  `json:"zip,omitempty"`
}

type DiscountCode struct {
	Amount *decimal.Decimal `json:"amount,omitempty"`
	Code   string           `json:"code,omitempty"`
	Type   string           `json:"type,omitempty"`
}

type LineItem struct {
	ID                         string           `json:"id,omitempty"`
	ProductID                  string           `json:"product_id,omitempty"`
	VariantID                  string           `json:"variant_id,omitempty"`
	Quantity                   int              `json:"quantity,omitempty"`
	Price                      *decimal.Decimal `json:"price,omitempty"`
	TotalDiscount              *decimal.Decimal `json:"total_discount,omitempty"`
	ProductTitle               string           `json:"product_title,omitempty"`
	VariantTitle               string           `json:"variant_title,omitempty"`
	Name                       string           `json:"name,omitempty"`
	SKU                        string           `json:"sku,omitempty"`
	Vendor                     string           `json:"vendor,omitempty"`
	GiftCard                   bool             `json:"gift_card,omitempty"`
	Taxable                    bool             `json:"taxable,omitempty"`
	FulfillmentService         string           `json:"fulfillment_service,omitempty"`
	RequiresShipping           bool             `json:"requires_shipping,omitempty"`
	VariantInventoryManagement string           `json:"variant_inventory_management,omitempty"`
	PreTaxPrice                *decimal.Decimal `json:"pre_tax_price,omitempty"`
	Properties                 []NoteAttribute  `json:"properties,omitempty"`
	ProductExists              bool             `json:"product_exists,omitempty"`
	FulfillableQuantity        int              `json:"fulfillable_quantity,omitempty"`
	Grams                      int              `json:"grams,omitempty"`
	FulfillmentStatus          string           `json:"fulfillment_status,omitempty"`
	TaxLines                   []TaxLine        `json:"tax_lines,omitempty"`
	OriginLocation             *Address         `json:"origin_location,omitempty"`
	DestinationLocation        *Address         `json:"destination_location,omitempty"`
	// AppliedDiscount            *AppliedDiscount `json:"applied_discount,omitempty"`
}

type LineItemProperty struct {
	Message string `json:"message"`
}

type NoteAttribute struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// Represents the result from the orders/X endpoint
type OrderResource struct {
	Order *Order `json:"order"`
}

// Represents the result from the orders endpoint
type OrdersResource struct {
	Orders []Order `json:"orders"`
}

type PaymentDetails struct {
	AVSResultCode     string `json:"avs_result_code,omitempty"`
	CreditCardBin     string `json:"credit_card_bin,omitempty"`
	CVVResultCode     string `json:"cvv_result_code,omitempty"`
	CreditCardNumber  string `json:"credit_card_number,omitempty"`
	CreditCardCompany string `json:"credit_card_company,omitempty"`
}

type ShippingLine struct {
	Name string `json:"name,omitempty"`
}

type TaxLine struct {
	Title string           `json:"title,omitempty"`
	Price *decimal.Decimal `json:"price,omitempty"`
	Rate  *decimal.Decimal `json:"rate,omitempty"`
}

type Transaction struct {
	ID             int64            `json:"id,omitempty"`
	OrderID        int64            `json:"order_id,omitempty"`
	Amount         *decimal.Decimal `json:"amount,omitempty"`
	Kind           string           `json:"kind,omitempty"`
	Gateway        string           `json:"gateway,omitempty"`
	Status         string           `json:"status,omitempty"`
	Message        string           `json:"message,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
	Test           bool             `json:"test,omitempty"`
	Authorization  string           `json:"authorization,omitempty"`
	Currency       string           `json:"currency,omitempty"`
	LocationID     *int64           `json:"location_id,omitempty"`
	UserID         *int64           `json:"user_id,omitempty"`
	ParentID       *int64           `json:"parent_id,omitempty"`
	DeviceID       *int64           `json:"device_id,omitempty"`
	ErrorCode      string           `json:"error_code,omitempty"`
	SourceName     string           `json:"source_name,omitempty"`
	PaymentDetails *PaymentDetails  `json:"payment_details,omitempty"`
}

type ClientDetails struct {
	AcceptLanguage string `json:"accept_language,omitempty"`
	BrowserHeight  int    `json:"browser_height,omitempty"`
	BrowserIp      string `json:"browser_ip,omitempty"`
	BrowserWidth   int    `json:"browser_width,omitempty"`
	SessionHash    string `json:"session_hash,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
}

type Refund struct {
	Id              int64            `json:"id,omitempty"`
	OrderId         int64            `json:"order_id,omitempty"`
	CreatedAt       *time.Time       `json:"created_at,omitempty"`
	Note            string           `json:"note,omitempty"`
	Restock         bool             `json:"restock,omitempty"`
	UserId          int64            `json:"user_id,omitempty"`
	RefundLineItems []RefundLineItem `json:"refund_line_items,omitempty"`
	Transactions    []Transaction    `json:"transactions,omitempty"`
}

type RefundLineItem struct {
	Id         int64            `json:"id,omitempty"`
	Quantity   int              `json:"quantity,omitempty"`
	LineItemId int64            `json:"line_item_id,omitempty"`
	LineItem   *LineItem        `json:"line_item,omitempty"`
	Subtotal   *decimal.Decimal `json:"subtotal,omitempty"`
	TotalTax   *decimal.Decimal `json:"total_tax,omitempty"`
}

// List orders
func (s *OrderServiceOp) List(options interface{}) ([]Order, error) {
	path := fmt.Sprintf("%s/%s", globalApiPathPrefix, ordersBasePath)
	resource := new(OrdersResource)
	err := s.client.Get(path, resource, options)
	return resource.Orders, err
}

// Count orders
func (s *OrderServiceOp) Count(options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%s/count", globalApiPathPrefix, ordersBasePath)
	return s.client.Count(path, options)
}

// Get individual order
func (s *OrderServiceOp) Get(orderID string, options interface{}) (*Order, error) {
	path := fmt.Sprintf("%s/%s/%s", globalApiPathPrefix, ordersBasePath, orderID)
	resource := new(OrderResource)
	err := s.client.Get(path, resource, options)
	return resource.Order, err
}

// Create order
func (s *OrderServiceOp) Create(order Order) (*Order, error) {
	path := fmt.Sprintf("%s/%s", globalApiPathPrefix, ordersBasePath)
	wrappedData := OrderResource{Order: &order}
	resource := new(OrderResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Order, err
}

// Update order
func (s *OrderServiceOp) Update(order Order) (*Order, error) {
	path := fmt.Sprintf("%s/%s/%s", globalApiPathPrefix, ordersBasePath, order.ID)
	wrappedData := OrderResource{Order: &order}
	resource := new(OrderResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Order, err
}

// List metafields for an order
// func (s *OrderServiceOp) ListMetafields(orderID int64, options interface{}) ([]Metafield, error) {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.List(options)
// }

// Count metafields for an order
// func (s *OrderServiceOp) CountMetafields(orderID int64, options interface{}) (int, error) {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.Count(options)
// }

// Get individual metafield for an order
// func (s *OrderServiceOp) GetMetafield(orderID int64, metafieldID int64, options interface{}) (*Metafield, error) {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.Get(metafieldID, options)
// }

// Create a new metafield for an order
// func (s *OrderServiceOp) CreateMetafield(orderID int64, metafield Metafield) (*Metafield, error) {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.Create(metafield)
// }

// // Update an existing metafield for an order
// func (s *OrderServiceOp) UpdateMetafield(orderID int64, metafield Metafield) (*Metafield, error) {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.Update(metafield)
// }

// // Delete an existing metafield for an order
// func (s *OrderServiceOp) DeleteMetafield(orderID int64, metafieldID int64) error {
// 	metafieldService := &MetafieldServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
// 	return metafieldService.Delete(metafieldID)
// }

// List fulfillments for an order
func (s *OrderServiceOp) ListFulfillments(orderID string, options interface{}) ([]Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.List(options)
}

// Count fulfillments for an order
func (s *OrderServiceOp) CountFulfillments(orderID string, options interface{}) (int, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Count(options)
}

// Get individual fulfillment for an order
func (s *OrderServiceOp) GetFulfillment(orderID string, fulfillmentID string, options interface{}) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Get(fulfillmentID, options)
}

// Create a new fulfillment for an order
func (s *OrderServiceOp) CreateFulfillment(orderID string, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Create(fulfillment)
}

// Update an existing fulfillment for an order
func (s *OrderServiceOp) UpdateFulfillment(orderID string, fulfillment Fulfillment) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Update(fulfillment)
}

// Complete an existing fulfillment for an order
func (s *OrderServiceOp) CompleteFulfillment(orderID string, fulfillmentID string) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Complete(fulfillmentID)
}

// Transition an existing fulfillment for an order
func (s *OrderServiceOp) TransitionFulfillment(orderID string, fulfillmentID string) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Transition(fulfillmentID)
}

// Cancel an existing fulfillment for an order
func (s *OrderServiceOp) CancelFulfillment(orderID string, fulfillmentID string) (*Fulfillment, error) {
	fulfillmentService := &FulfillmentServiceOp{client: s.client, resource: ordersResourceName, resourceID: orderID}
	return fulfillmentService.Cancel(fulfillmentID)
}
