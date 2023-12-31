package lichess

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://lichess.org/"
)

// NewClient returns a new Lichess API client. If a nil httpClient is
// provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}
	c.common.client = c
	c.Games = (*GamesService)(&c.common)
	c.Puzzles = (*PuzzlesService)(&c.common)

	return c
}

// WithAuthToken returns a copy of the client configured to use the
// provided token for the Authorization header.
func (c *Client) WithAuthToken(token string) *Client {
	transport := c.client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	c.client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return transport.RoundTrip(req)
		},
	)

	return c
}

// A Client manages communication with the Lichess API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public GitHub API, but can be
	// set to a domain endpoint. BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Lichess API.
	Games   *GamesService
	Puzzles *PuzzlesService
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	return c.newRequest(ctx, method, urlStr, nil)
}

// RequestBody is a data structure that holds the request body as well as
// the type of the content, which will be used to set the Content-Type header.
type RequestBody struct {
	Bytes io.Reader
	Type  string
}

// NewRequestWithBody creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequestWithBody(
	ctx context.Context,
	method, urlStr string,
	body RequestBody,
) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, urlStr, body.Bytes)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", body.Type)

	return req, nil
}

func (c *Client) newRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	switch typeOfResponse(method, req.URL.Path) {
	case jsonResponseType:
		req.Header.Set("Accept", "application/json")
	case ndJsonResponseType:
		req.Header.Set("Accept", "application/x-ndjson")
	}

	return req, nil
}

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise, you
// are supposed to read and close the response's Body. If rate limit is exceeded
// and reset time is in the future, BareDo returns *RateLimitError immediately
// without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(req *http.Request) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}

		return nil, err
	}

	response := newResponse(resp)

	return response, err
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error happens, the response is returned as is.
// If rate limit is exceeded and reset time is in the future, Do returns
// *RateLimitError immediately without making a network API call.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	res, err := c.BareDo(req)
	if err != nil {
		return res, err
	}

	// We only close the response body, when the
	// caller expects it to parse the response body.
	// In other words, when v is not nil.
	// That's to also have support for streaming.
	defer func() {
		if v != nil {
			// Explicit ignore error.
			// We might want to revisit this later.
			_ = res.Body.Close()
		}
	}()

	err = c.decodeResponse(req, res, v)

	return res, err
}

func (c *Client) decodeResponse(req *http.Request, res *Response, v interface{}) error {
	var err error

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		switch typeOfResponse(req.Method, req.URL.Path) {
		case jsonResponseType:
			decErr := json.NewDecoder(res.Body).Decode(v)
			// Ignore EOF errors caused by empty response body
			if decErr != nil && !errors.Is(decErr, io.EOF) {
				err = decErr
			}
		case ndJsonResponseType:
			err = c.decodeNdJson(res, v)
		}
	}

	return err
}

func (c *Client) decodeNdJson(res *Response, v interface{}) error {
	if reflect.ValueOf(v).Elem().Kind() != reflect.Slice {
		return errors.New("v is not a pointer to a slice")
	}

	itemType := reflect.ValueOf(v).Elem().Type().Elem()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		item := reflect.New(itemType).Interface()
		if err := json.Unmarshal(scanner.Bytes(), item); err != nil {
			return err
		}

		reflect.ValueOf(v).Elem().Set(reflect.Append(reflect.ValueOf(v).Elem(), reflect.ValueOf(item).Elem()))
	}

	return nil
}

type service struct {
	client *Client
}

// Response is a Lichess API response. This wraps the standard http.Response
// returned from Lichess and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

type responseType uint8

const (
	jsonResponseType responseType = iota

	ndJsonResponseType // An array of this length will be able to contain all rate limit categories.
)

// typeOfResponse returns the response type of the endpoint, determined by HTTP method and Request.URL.Path.
func typeOfResponse(_, path string) responseType {
	switch {
	default:
		return jsonResponseType
	case strings.Contains(path, "api/games/user/"),
		strings.Contains(path, "api/stream/games-by-users"),
		strings.Contains(path, "api/puzzle/activity"):
		return ndJsonResponseType
	}
}

// roundTripperFunc creates a RoundTripper (transport).
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
