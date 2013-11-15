package localbitcoins

import (
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
)

var (
  // HTTP request multiplayer used with test server.
  mux *http.ServeMux

  // localbitcoins client being used
  client *Client

  // test HTTP that provides mock API responses
  server *httptest.Server
)

// Sets up a test HTTP server along with a localbitcoins.Client that is
// configured to talk to the test server. Tests should register handlers on mux
// which provide mock responses for the API method being tested.
func setup() {
  // test server
  mux = http.NewServeMux()
  server = httptest.NewServer(mux)

  // localbitcoins client configured to use test server
  client = NewClient(nil)
  client.BaseURL, _ = url.Parse(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
  server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
  if want != r.Method {
    t.Errorf("Request method = %v, want %v", r.Method, want)
  }
}

func TestNewClient(t *testing.T) {
  c := NewClient(nil)

  if c.BaseURL.String() != defaultBaseURL {
    t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(),
             defaultBaseURL)
  }
  if c.UserAgent != userAgent {
    t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, userAgent)
  }
}
