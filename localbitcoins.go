package localbitcoins

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/url"
  "reflect"

  "github.com/google/go-querystring/query"
)

const (
  libraryVersion = "0.0.1"
  defaultBaseURL = "https://localbitcoins.com/"
  userAgent = "go-localbitcoins/" + libraryVersion
)

// A Client manages communication with the LocalBitcoins API.
type Client struct {
  client *http.Client

  // Base URL used for API requests. Must end with a trailing slash.
  BaseURL *url.URL

  // User agent sent when communicating with the LocalBitcoins API.
  UserAgent string

  // Services for talking to different parts of the LocalBitcoins API.
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

// String is a helper routine that allocates a new string value to store v and
// returns a pointer to it.
func String(v string) *string {
  p := new(string)
  *p = v
  return p
}
