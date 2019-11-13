package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/google/go-querystring/query"

	"github.com/rs/xid"
)

const (
	defaultHTTPTimeout            = 15
	defaultMaxConcurrentRequest   = 100
	defaultRequestVolumeThreshold = 20
	defaultSleepWindow            = 5
	defaultErrorPercentThreshold  = 50
)

var (
	sharedClient = &http.Client{}
)

// Client defines http client with configured timeout
type Client struct {
	// the underlying http.Client object
	client *http.Client

	// http  timeout for all request
	httpTimeout time.Duration

	// if true, client will return err on non 2xx status code,
	// when that happens, resp object will be nill, thus could not be used
	errNon2xx bool

	// below fields are hystrix related fields
	name                  string
	maxConcurrentReq      int
	reqVolThreshold       int
	sleepWindow           time.Duration
	errorPercentThreshold int
}

// NewClient creates new Client object with given options
func NewClient() *Client {
	c := Client{
		name:                  xid.New().String(),
		httpTimeout:           defaultHTTPTimeout * time.Second,
		maxConcurrentReq:      defaultMaxConcurrentRequest,
		reqVolThreshold:       defaultRequestVolumeThreshold,
		sleepWindow:           defaultSleepWindow * time.Second,
		errorPercentThreshold: defaultErrorPercentThreshold,
	}

	c.client = sharedClient

	hystrix.ConfigureCommand(c.name, hystrix.CommandConfig{
		Timeout:                int(c.httpTimeout.Nanoseconds()) / 1e6,
		MaxConcurrentRequests:  c.maxConcurrentReq,
		RequestVolumeThreshold: c.reqVolThreshold,
		SleepWindow:            int(c.sleepWindow.Nanoseconds()) / 1e6,
		ErrorPercentThreshold:  c.errorPercentThreshold,
	})

	return &c
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	err = hystrix.Do(c.name, func() error {
		resp, err = c.client.Do(req)
		return err
	}, nil)

	if err != nil {
		return nil, err
	}

	// receive non 2xx status code when errNon2xx=true
	if c.errNon2xx && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return nil, fmt.Errorf("got non 2xx status code: %v", resp.StatusCode)
	}

	return resp, nil
}

// Do executes the given http request and returns the http response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

// DoJSON executes the given http request and unmarshall the response body
// into the given `data`
// The returned response Body is already closed
func (c *Client) DoJSON(req *http.Request, data interface{}) (*http.Response, error) {
	// do request
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal and return
	err = json.Unmarshal(b, data)
	return resp, err
}

// Get executes http GET request with the given headers
func (c *Client) Get(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	return c.doMethod(ctx, http.MethodGet, url, nil, headers)
}

// GetJSON executes http GET request with the given headers and parse the json response body
// into the given data.
// The returned response Body is already closed
func (c *Client) GetJSON(ctx context.Context, url string, headers http.Header, data interface{}) (*http.Response, error) {
	return c.doMethodJSON(ctx, http.MethodGet, url, nil, headers, data)
}

// GetJSONWithReqBody executes http GET request with the given headers, request body and parse the json response body
// into the given data.
// The returned response Body is already closed
func (c *Client) GetJSONWithReqBody(ctx context.Context, url string, headers http.Header, requestBody interface{}, data interface{}) (*http.Response, error) {
	return c.doMethodJSON(ctx, http.MethodGet, url, requestBody, headers, data)
}

// Post executes http POST request with the given body & headers.
// The supported data types for body are:
// - string
// - []byte
// - struct/map/slice : will be marshalled to JSON
// - io.Reader
func (c *Client) Post(ctx context.Context, url string, headers http.Header, body interface{}) (*http.Response, error) {
	return c.doMethod(ctx, http.MethodPost, url, body, headers)
}

// PostJSON executes http POST request with the given body & headers and parse the json response body
// into the given data.
// The supported data types for body are:
// - string
// - []byte
// - struct/map/slice : will be marshalled to JSON
// - io.Reader
// The returned response Body is already closed
func (c *Client) PostJSON(ctx context.Context, url string, headers http.Header, body, data interface{}) (*http.Response, error) {
	return c.doMethodJSON(ctx, http.MethodPost, url, body, headers, data)
}

// do specific request
func (c *Client) doMethod(ctx context.Context, method, url string, body interface{},
	headers http.Header) (*http.Response, error) {

	bodyReader, err := toIoReader(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	req = req.WithContext(ctx)
	return c.Do(req)
}

func (c *Client) doMethodJSON(ctx context.Context, method, url string, body interface{}, headers http.Header, data interface{}) (*http.Response, error) {
	bodyReader, err := toIoReader(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	req = req.WithContext(ctx)

	return c.DoJSON(req, data)
}

func toIoReader(body interface{}) (io.Reader, error) {
	var bodyReader io.Reader

	switch body := body.(type) {
	case nil:
	case io.Reader:
		bodyReader = body
	case string:
		bodyReader = bytes.NewBufferString(body)
	case []byte:
		bodyReader = bytes.NewBuffer(body)
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(b)
	}
	return bodyReader, nil
}

// DoURLQueryJSON function do http request with given param :
// method : http method to perform
// path : path to url
// param : struct with giving url tag, for reference : "github.com/google/go-querystring" * this is nullable
// response : pointer for response struct
// this function by default content type to form-url-encoded
func (c *Client) DoURLQueryJSON(ctx context.Context, method string, path string, param interface{}, response interface{}) (resp *http.Response, err error) {

	requestURL, err := url.Parse(path)
	if err != nil {
		return
	}

	urlQuery, err := query.Values(param)
	if err != nil {
		return
	}

	requestURL.RawQuery = urlQuery.Encode()
	req, err := http.NewRequest(method, requestURL.String(), nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = c.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, response)

	return
}
