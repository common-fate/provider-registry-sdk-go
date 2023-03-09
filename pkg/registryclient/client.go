package registryclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

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
	HTTPClient Doer
	// Logger to use. If not provided,
	// we default to the global logger with zap.S().
	//
	// The registry client writes a debug log message
	// when being configured, with the registry URL in use.
	Logger *zap.SugaredLogger
}

// New creates a new client. By default, it uses
// https://api.registry.commonfate.io as the registry URL.
//
// By default, if the COMMON_FATE_PROVIDER_REGISTRY_URL environment variables is set,
// it will be used as the registry URL.
//
// If you want to change the registry local, set this environment variable.
//
// Alternatively, and not recommended, you can use the registryclient.NewWithURL()
// to provide a custom URL defined in Go code.
func New(ctx context.Context, opts ...func(co *ClientOpts)) (*Client, error) {
	url := os.Getenv("COMMON_FATE_PROVIDER_REGISTRY_URL")

	if url == "" {
		// default to the Common Fate production registry URL.
		url = "https://api.registry.commonfate.io"
	}

	return NewWithURL(ctx, url, opts...)
}

// NewWithURL allows a custom registry URL to be provided. This method is not recommended.
// Most of the time you won't need this - just call registryclient.New() and set the
// COMMON_FATE_PROVIDER_REGISTRY_URL environment variable.
func NewWithURL(ctx context.Context, url string, opts ...func(co *ClientOpts)) (*Client, error) {
	co := &ClientOpts{
		HTTPClient: &ErrorHandlingClient{Client: http.DefaultClient},
		Logger:     zap.S(),
	}

	for _, o := range opts {
		o(co)
	}

	co.Logger.Debugw("configuring provider registry client", "url", url)

	return providerregistrysdk.NewClientWithResponses(url, providerregistrysdk.WithHTTPClient(co.HTTPClient))
}

// Client is an alias for the exported Go SDK client type
type Client = providerregistrysdk.ClientWithResponses
