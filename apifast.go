package apifast

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Header struct {
	Tag   string
	Value interface{}
}

type Auth struct {
	Username string
	Password string
	Token    string
}

// RequestOptions represents optional parameters for making API requests
type RequestOptions struct {
	Timeout time.Duration // Request timeout duration
	payload []byte
	Headers []Header
	Auth    Auth
}

type FastBuilder struct {
	method  string
	url     string
	options RequestOptions
	result  interface{}
}

type Response struct {
	Code int    // HTTP code
	Msg  string // Status message
	Body interface{}
}

// Build initializes a new FastBuilder instance
func Build() *FastBuilder {
	return &FastBuilder{}
}

// Uri sets the request URL
func (b *FastBuilder) Uri(url string) *FastBuilder {
	b.url = url
	return b
}

// Timeout sets the request timeout
func (b *FastBuilder) Timeout(timeout time.Duration) *FastBuilder {
	b.options.Timeout = timeout
	return b
}

// Auth sets the authentication options
func (b *FastBuilder) Auth(auth Auth) *FastBuilder {
	b.options.Auth = auth
	return b
}

// Headers sets custom headers for the request
func (b *FastBuilder) Headers(headers []Header) *FastBuilder {
	b.options.Headers = headers
	return b
}

// Payload sets the request payload (body)
func (b *FastBuilder) Payload(payload []byte) *FastBuilder {
	b.options.payload = payload
	return b
}

// Result specifies where to store the response result
func (b *FastBuilder) Result(result interface{}) *FastBuilder {
	b.result = result
	return b
}

// Get initiates a GET request
func (b *FastBuilder) Get() (*Response, error) {
	b.method = "GET"
	return b.makeRequest()
}

// Post initiates a POST request
func (b *FastBuilder) Post() (*Response, error) {
	b.method = "POST"
	return b.makeRequest()
}

// Patch initiates a PATCH request
func (b *FastBuilder) Patch() (*Response, error) {
	b.method = "PATCH"
	return b.makeRequest()
}

// Delete initiates a DELETE request
func (b *FastBuilder) Delete() (*Response, error) {
	b.method = "DELETE"
	return b.makeRequest()
}

// makeRequest handles sending the request and receiving the response
func (b *FastBuilder) makeRequest() (*Response, error) {
	// Create a context with timeout if specified
	var ctx context.Context
	var cancel context.CancelFunc
	if b.options.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), b.options.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	// Create a fasthttp client with the given timeout
	client := &fasthttp.Client{
		ReadTimeout:  b.options.Timeout,
		WriteTimeout: b.options.Timeout,
	}

	// Prepare the request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Set custom headers if provided
	for _, h := range b.options.Headers {
		req.Header.Set(h.Tag, fmt.Sprintf("%v", h.Value))
	}

	// Add Basic or Bearer authentication if provided
	if b.options.Auth.Username != "" && b.options.Auth.Password != "" {
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(b.options.Auth.Username+":"+b.options.Auth.Password))
		req.Header.Set("Authorization", authHeader)
	} else if b.options.Auth.Token != "" {
		authHeader := "Bearer " + b.options.Auth.Token
		req.Header.Set("Authorization", authHeader)
	}

	// Set the request URI and method
	req.SetRequestURI(b.url)
	req.Header.SetMethod(b.method)

	// Set the request body if payload is provided
	if b.options.payload != nil {
		req.SetBody(b.options.payload)
	}

	// Create a fasthttp response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Send the request
	err := client.DoTimeout(req, resp, b.options.Timeout)
	if err != nil {
		// Check if the error is due to a timeout
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out")
		}
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Get the response body
	body := resp.Body()

	// Map response body to the result if provided
	if b.result != nil {
		if err := mapper(body, b.result); err != nil {
			return nil, err
		}
	}

	// Return the response
	return &Response{
		Code: resp.StatusCode(),
		Msg:  resp.String(),
		Body: body,
	}, nil
}

// mapper function unmarshals the JSON response into the provided destination
func mapper(source []byte, dest interface{}) error {
	return json.Unmarshal(source, dest)
}
