package goshoplazza

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	UserAgent = "goshoplazza/0.0.1"
)

var (
	// Shopify API version YYYY-MM - defaults to admin which uses the oldest stable version of the api
	globalApiPathPrefix string = "openapi"
)

// App represents basic app settings such as Api key, secret, scope, and redirect url.
// See oauth.go for OAuth related helper functions.
type App struct {
	ApiKey      string
	ApiSecret   string
	RedirectUrl string
	Scope       string
	Password    string
}

// Client manages communication with the Shopify API.
type Client struct {
	// HTTP client used to communicate with the DO API.
	Client *http.Client

	// App settings
	app App

	// Base URL for API requests.
	// This is set on a per-store basis which means that each store must have
	// its own client.
	baseURL *url.URL

	// A permanent access token
	token string

	// Services used for communicating with the API
	Product ProductService
	// CustomCollection           CustomCollectionService
	// SmartCollection            SmartCollectionService
	// Customer                   CustomerService
	// CustomerAddress            CustomerAddressService
	Order OrderService
	// DraftOrder                 DraftOrderService
	// Shop                       ShopService
	// Webhook                    WebhookService
	Variant VariantService
	Image   ImageService
	// Transaction                TransactionService
	// Theme                      ThemeService
	// Asset                      AssetService
	// ScriptTag                  ScriptTagService
	// RecurringApplicationCharge RecurringApplicationChargeService
	// UsageCharge                UsageChargeService
	// Metafield                  MetafieldService
	// Blog                       BlogService
	// ApplicationCharge          ApplicationChargeService
	// Redirect                   RedirectService
	// Page                       PageService
	// StorefrontAccessToken      StorefrontAccessTokenService
	// Collect                    CollectService
	// Location                   LocationService
	// DiscountCode               DiscountCodeService
	// InventoryItem              InventoryItemService
}

// A general response error that follows a similar layout to Shopify's response
// errors, i.e. either a single message or a list of messages.
type ResponseError struct {
	Status  int
	Message string
	Errors  []string
}

func (e ResponseError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	sort.Strings(e.Errors)
	s := strings.Join(e.Errors, ", ")

	if s != "" {
		return s
	}

	return "Unknown Error"
}

// ResponseDecodingError occurs when the response body from Shopify could
// not be parsed.
type ResponseDecodingError struct {
	Body    []byte
	Message string
	Status  int
}

func (e ResponseDecodingError) Error() string {
	return e.Message
}

// An error specific to a rate-limiting response. Embeds the ResponseError to
// allow consumers to handle it the same was a normal ResponseError.
type RateLimitError struct {
	ResponseError
	RetryAfter int
}

// Creates an API request. A relative URL can be provided in urlStr, which will
// be resolved to the BaseURL of the Client. Relative URLS should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body, options interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Make the full url based on the relative path
	u := c.baseURL.ResolveReference(rel)

	// Add custom options
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}

		for k, values := range u.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
		u.RawQuery = optionsQuery.Encode()
	}

	// A bit of JSON ceremony
	var js []byte = nil

	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)
	if c.token != "" {
		req.Header.Add("Token-Type", "Bearer")
		req.Header.Add("Access-Token", c.token)
	} else if c.app.Password != "" {
		req.SetBasicAuth(c.app.ApiKey, c.app.Password)
	}
	return req, nil
}

// Option is used to configure client with options
type Option func(c *Client)

// WithVersion optionally sets the api-version if the passed string is valid
func WithVersion(apiVersion string) Option {
	return func(c *Client) {
		var rxPat = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}$`)
		if len(apiVersion) > 0 && rxPat.MatchString(apiVersion) {
			globalApiPathPrefix = fmt.Sprintf("admin/api/%s", apiVersion)
		} else {
			globalApiPathPrefix = "admin"
		}
	}
}

// NewClient returns a new Shopify API client with an already authenticated shopname and
// token. The shopName parameter is the shop's myshoplazza domain,
// e.g. "theshop.myshoplazza.com", or simply "theshop"
// a.NewClient(shopName, token, opts) is equivalent to NewClient(a, shopName, token, opts)
func (a App) NewClient(shopName, token string, opts ...Option) *Client {
	return NewClient(a, shopName, token, opts...)
}

// Returns a new Shopify API client with an already authenticated shopname and
// token. The shopName parameter is the shop's myshoplazza domain,
// e.g. "theshop.myshoplazza.com", or simply "theshop"
func NewClient(app App, shopName, token string, opts ...Option) *Client {
	httpClient := http.DefaultClient

	baseURL, _ := url.Parse(ShopBaseUrl(shopName))

	c := &Client{Client: httpClient, app: app, baseURL: baseURL, token: token}
	c.Product = &ProductServiceOp{client: c}
	// c.CustomCollection = &CustomCollectionServiceOp{client: c}
	// c.SmartCollection = &SmartCollectionServiceOp{client: c}
	// c.Customer = &CustomerServiceOp{client: c}
	// c.CustomerAddress = &CustomerAddressServiceOp{client: c}
	c.Order = &OrderServiceOp{client: c}
	// c.DraftOrder = &DraftOrderServiceOp{client: c}
	// c.Shop = &ShopServiceOp{client: c}
	// c.Webhook = &WebhookServiceOp{client: c}
	c.Variant = &VariantServiceOp{client: c}
	c.Image = &ImageServiceOp{client: c}
	// c.Transaction = &TransactionServiceOp{client: c}
	// c.Theme = &ThemeServiceOp{client: c}
	// c.Asset = &AssetServiceOp{client: c}
	// c.ScriptTag = &ScriptTagServiceOp{client: c}
	// c.RecurringApplicationCharge = &RecurringApplicationChargeServiceOp{client: c}
	// c.Metafield = &MetafieldServiceOp{client: c}
	// c.Blog = &BlogServiceOp{client: c}
	// c.ApplicationCharge = &ApplicationChargeServiceOp{client: c}
	// c.Redirect = &RedirectServiceOp{client: c}
	// c.Page = &PageServiceOp{client: c}
	// c.StorefrontAccessToken = &StorefrontAccessTokenServiceOp{client: c}
	// c.UsageCharge = &UsageChargeServiceOp{client: c}
	// c.Collect = &CollectServiceOp{client: c}
	// c.Location = &LocationServiceOp{client: c}
	// c.DiscountCode = &DiscountCodeServiceOp{client: c}
	// c.InventoryItem = &InventoryItemServiceOp{client: c}

	// apply any options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Do sends an API request and populates the given interface with the parsed
// response. It does not make much sense to call Do without a prepared
// interface instance.
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = CheckResponseError(resp)
	if err != nil {
		return err
	}

	if v != nil {
		// content, e := ioutil.ReadAll(resp.Body)
		// if e != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(string(content))

		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&v)
		if err != nil {
			return err
		}
	}

	return nil
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
	if err.Status == 429 {
		f, _ := strconv.ParseFloat(r.Header.Get("retry-after"), 64)
		return RateLimitError{
			ResponseError: err,
			RetryAfter:    int(f),
		}
	}
	if err.Status == 406 {
		err.Message = "Not acceptable"
	}
	return err
}

func CheckResponseError(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}

	// Create an anonoymous struct to parse the JSON data into.
	shoplazzaError := struct {
		Error  string      `json:"error"`
		Errors interface{} `json:"errors"`
	}{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// empty body, this probably means shoplazza returned an error with no body
	// we'll handle that error in wrapSpecificError()
	if len(bodyBytes) > 0 {
		err := json.Unmarshal(bodyBytes, &shoplazzaError)
		if err != nil {
			return ResponseDecodingError{
				Body:    bodyBytes,
				Message: err.Error(),
				Status:  r.StatusCode,
			}
		}
	}

	// Create the response error from the Shopify error.
	responseError := ResponseError{
		Status:  r.StatusCode,
		Message: shoplazzaError.Error,
	}

	// If the errors field is not filled out, we can return here.
	if shoplazzaError.Errors == nil {
		return wrapSpecificError(r, responseError)
	}

	// Shopify errors usually have the form:
	// {
	//   "errors": {
	//     "title": [
	//       "something is wrong"
	//     ]
	//   }
	// }
	// This structure is flattened to a single array:
	// [ "title: something is wrong" ]
	//
	// Unfortunately, "errors" can also be a single string so we have to deal
	// with that. Lots of reflection :-(
	switch reflect.TypeOf(shoplazzaError.Errors).Kind() {
	case reflect.String:
		// Single string, use as message
		responseError.Message = shoplazzaError.Errors.(string)
	case reflect.Slice:
		// An array, parse each entry as a string and join them on the message
		// json always serializes JSON arrays into []interface{}
		for _, elem := range shoplazzaError.Errors.([]interface{}) {
			responseError.Errors = append(responseError.Errors, fmt.Sprint(elem))
		}
		responseError.Message = strings.Join(responseError.Errors, ", ")
	case reflect.Map:
		// A map, parse each error for each key in the map.
		// json always serializes into map[string]interface{} for objects
		for k, v := range shoplazzaError.Errors.(map[string]interface{}) {
			// Check to make sure the interface is a slice
			// json always serializes JSON arrays into []interface{}
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				for _, elem := range v.([]interface{}) {
					// If the primary message of the response error is not set, use
					// any message.
					if responseError.Message == "" {
						responseError.Message = fmt.Sprintf("%v: %v", k, elem)
					}
					topicAndElem := fmt.Sprintf("%v: %v", k, elem)
					responseError.Errors = append(responseError.Errors, topicAndElem)
				}
			}
		}
	}

	return wrapSpecificError(r, responseError)
}

// General list options that can be used for most collections of entities.
type ListOptions struct {
	Page         int       `url:"page,omitempty"`
	Limit        int       `url:"limit,omitempty"`
	SinceID      int64     `url:"since_id,omitempty"`
	CreatedAtMin time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax time.Time `url:"updated_at_max,omitempty"`
	Order        string    `url:"order,omitempty"`
	Fields       string    `url:"fields,omitempty"`
	Vendor       string    `url:"vendor,omitempty"`
	IDs          []int64   `url:"ids,omitempty,comma"`
}

// General count options that can be used for most collection counts.
type CountOptions struct {
	CreatedAtMin time.Time `url:"created_at_min,omitempty"`
	CreatedAtMax time.Time `url:"created_at_max,omitempty"`
	UpdatedAtMin time.Time `url:"updated_at_min,omitempty"`
	UpdatedAtMax time.Time `url:"updated_at_max,omitempty"`
}

func (c *Client) Count(path string, options interface{}) (int, error) {
	resource := struct {
		Count int `json:"count"`
	}{}
	err := c.Get(path, &resource, options)
	return resource.Count, err
}

// CreateAndDo performs a web request to Shopify with the given method (GET,
// POST, PUT, DELETE) and relative path (e.g. "/admin/orders.json").
// The data, options and resource arguments are optional and only relevant in
// certain situations.
// If the data argument is non-nil, it will be used as the body of the request
// for POST and PUT requests.
// The options argument is used for specifying request options such as search
// parameters like created_at_min
// Any data returned from Shopify will be marshalled into resource argument.
func (c *Client) CreateAndDo(method, path string, data, options, resource interface{}) error {
	req, err := c.NewRequest(method, path, data, options)
	if err != nil {
		return err
	}

	err = c.Do(req, resource)
	if err != nil {
		return err
	}

	return nil
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *Client) Get(path string, resource, options interface{}) error {
	return c.CreateAndDo("GET", path, nil, options, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *Client) Post(path string, data, resource interface{}) error {
	return c.CreateAndDo("POST", path, data, nil, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *Client) Put(path string, data, resource interface{}) error {
	return c.CreateAndDo("PUT", path, data, nil, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(path string) error {
	return c.CreateAndDo("DELETE", path, nil, nil, nil)
}
