package localbitcoins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.0.1"
	defaultBaseURL = "https://localbitcoins.com/"
	userAgent      = "go-localbitcoins/" + libraryVersion
)

// A Client manages communication with the LocalBitcoins API.
type Client struct {
	client *http.Client

	// Base URL used for API requests. Must end with a trailing slash.
	BaseURL *url.URL

	// User agent sent when communicating with the LocalBitcoins API.
	UserAgent string

	// Services for talking to different parts of the LocalBitcoins API.
	Accounts *AccountsService
}

// Adds the parameters in opt as URL query parameters to s. opt must be a
// struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Returns a new LocalBitcoins API client. If a nil httpClient is provided,
// http.DefaultClient will be used. To use API methods that require
// authentication (most, if not all, do), provide an http.Client that will
// perform the authentication for you (such as that provided by the goauth2
// library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Accounts = &AccountsService{client: c}
	return c
}

// Creates an API request. A relative URl can be provided in urlStr, in which
// case it is resolved relative to the BaseURL of the Client. Relative URLs
// should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string,
	body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Response is a LocalBitcoins API response. This wraps the standard
// http.Response returned from LocalBitcoins and provides convenient access to
// things like pagination links.
type Response struct {
	*http.Response
}

// Creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Sends an API request and returns the API response. The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even if there was an error, we still return the response body if the
		// user wants to inspect it further
		return response, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return response, err
}

// An ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	Response *http.Response
	Err      Error `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d - %v %d",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Err.Message, r.Err.Code)
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"error_code"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(`type %d error with message "%v"`, e.Code, e.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Represents data as returned by the LocalBitcoins API. Wraps the data and
// actions objects.
type ResponseData struct {
	Data    interface{} `json:"data"`
	Actions interface{} `json:"actions"`
}

// Bool is a helper function that allocates a new bool value to store v and
// returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper function that allocates a new int32 value to store v and
// returns a pointer to it, but unlike Int32 its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper function that allocates a new string value to store v and
// returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
