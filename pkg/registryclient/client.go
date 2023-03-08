package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/common-fate/useragent"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ErrorHandlingClient checks the response status code
// and creates an error if the API returns greater than 300.
type ErrorHandlingClient struct {
	Client Doer
	Logger *zap.SugaredLogger
}

func (rd *ErrorHandlingClient) Do(req *http.Request) (*http.Response, error) {
	// add a user agent to the request
	ua := useragent.FromContext(req.Context())
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	res, err := rd.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 300 {
		// response is ok
		return res, nil
	}

	// if we get here, the API has returned an error
	// surface this as a Go error so we don't need to handle it everywhere in our CLI codebase.
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res, errors.Wrap(err, "reading error response body")
	}

	// if we get here, the API has returned a response code > 300.
	// we treat this as an error.
	return res, fmt.Errorf("the Provider Registry returned an error (code %v): %s", res.StatusCode, string(body))
}

type ClientOpts struct {
	APIURL     string
	HTTPClient Doer
}

// WithAPIURL overrides the API URL.
// If the url is empty, it is not overriden and the regular
// API URL from aws-exports.json is used instead.
//
// This can be used for local development to provider a localhost URL.
func WithAPIURL(url string) func(co *ClientOpts) {
	return func(co *ClientOpts) {
		co.APIURL = url
	}
}

// New creates a new client. By default, it uses
// https://api.registry.commonfate.io as the registry URL.
func New(ctx context.Context, opts ...func(co *ClientOpts)) (*Client, error) {
	co := &ClientOpts{
		APIURL:     "https://api.registry.commonfate.io",
		HTTPClient: http.DefaultClient,
	}

	for _, o := range opts {
		o(co)
	}

	httpClient := &ErrorHandlingClient{Client: co.HTTPClient}

	return providerregistrysdk.NewClientWithResponses(co.APIURL, providerregistrysdk.WithHTTPClient(httpClient))
}

// Client is an alias for the exported Go SDK client type
type Client = providerregistrysdk.ClientWithResponses
