// Package providerregistrysdk provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package providerregistrysdk

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Defines values for ConfigArgumentType.
const (
	String ConfigArgumentType = "string"
)

// Defines values for LogLevel.
const (
	ERROR   LogLevel = "ERROR"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
)

// Defines values for TargetKindType.
const (
	Object TargetKindType = "object"
)

// ConfigArgument defines model for ConfigArgument.
type ConfigArgument struct {
	Description string             `json:"description"`
	Secret      bool               `json:"secret"`
	Type        ConfigArgumentType `json:"type"`
}

// ConfigArgumentType defines model for ConfigArgument.Type.
type ConfigArgumentType string

// ConfigSchema defines model for ConfigSchema.
type ConfigSchema map[string]ConfigArgument

// DescribeResponse defines model for DescribeResponse.
type DescribeResponse struct {
	Config      map[string]interface{} `json:"config"`
	Diagnostics []DiagnosticLog        `json:"diagnostics"`
	Healthy     bool                   `json:"healthy"`

	// A registered provider version
	Provider Provider       `json:"provider"`
	Schema   ProviderSchema `json:"schema"`
}

// DiagnosticLog defines model for DiagnosticLog.
type DiagnosticLog struct {
	Level LogLevel `json:"level"`
	Msg   string   `json:"msg"`
}

// LogLevel defines model for LogLevel.
type LogLevel string

// Metadata about the schema
type MetaSchema struct {
	Framework string `json:"framework"`
}

// A registered provider version
type Provider struct {
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
	Version   string `json:"version"`
}

// A registered provider version
type ProviderDetail struct {
	CfnTemplateS3Arn string         `json:"cfnTemplateS3Arn"`
	LambdaAssetS3Arn string         `json:"lambdaAssetS3Arn"`
	Name             string         `json:"name"`
	Publisher        string         `json:"publisher"`
	Schema           ProviderSchema `json:"schema"`
	Version          string         `json:"version"`
}

// ProviderSchema defines model for ProviderSchema.
type ProviderSchema struct {
	Id     string        `json:"$id"`
	Schema string        `json:"$schema"`
	Config *ConfigSchema `json:"config,omitempty"`

	// Metadata about the schema
	Meta      *MetaSchema      `json:"meta,omitempty"`
	Resources *ResourcesSchema `json:"resources,omitempty"`
	Targets   *TargetSchema    `json:"targets,omitempty"`
}

// Resource defines model for Resource.
type Resource struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
	Type string `json:"type"`
}

// ResourceLoader defines model for ResourceLoader.
type ResourceLoader struct {
	Title string `json:"title"`
}

// ResourcesSchema defines model for ResourcesSchema.
type ResourcesSchema struct {
	Loaders map[string]ResourceLoader `json:"loaders"`
	Types   *map[string]Resource      `json:"types,omitempty"`
}

// Defines the metadata and data type for the argument
type TargetArgument struct {
	Description *string `json:"description,omitempty"`
	Id          string  `json:"id"`
	Resource    *string `json:"resource,omitempty"`
	Title       string  `json:"title"`
}

// TargetKind defines model for TargetKind.
type TargetKind struct {
	Properties map[string]TargetArgument `json:"properties"`
	Type       TargetKindType            `json:"type"`
}

// TargetKindType defines model for TargetKind.Type.
type TargetKindType string

// TargetSchema defines model for TargetSchema.
type TargetSchema map[string]TargetKind

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse defines model for HealthResponse.
type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

// ListProvidersResponse defines model for ListProvidersResponse.
type ListProvidersResponse struct {
	Next      *string          `json:"next"`
	Providers []ProviderDetail `json:"providers"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ListAllProviders request
	ListAllProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProvider request
	GetProvider(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProviderSetupDocs request
	GetProviderSetupDocs(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProviderUsageDoc request
	GetProviderUsageDoc(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListAllProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListAllProvidersRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProvider(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProviderRequest(c.Server, publisher, name, version)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProviderSetupDocs(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProviderSetupDocsRequest(c.Server, publisher, name, version)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProviderUsageDoc(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProviderUsageDocRequest(c.Server, publisher, name, version)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListAllProvidersRequest generates requests for ListAllProviders
func NewListAllProvidersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1alpha1/providers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProviderRequest generates requests for GetProvider
func NewGetProviderRequest(server string, publisher string, name string, version string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "publisher", runtime.ParamLocationPath, publisher)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "version", runtime.ParamLocationPath, version)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1alpha1/providers/%s/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProviderSetupDocsRequest generates requests for GetProviderSetupDocs
func NewGetProviderSetupDocsRequest(server string, publisher string, name string, version string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "publisher", runtime.ParamLocationPath, publisher)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "version", runtime.ParamLocationPath, version)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1alpha1/providers/%s/%s/%s/setup", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProviderUsageDocRequest generates requests for GetProviderUsageDoc
func NewGetProviderUsageDocRequest(server string, publisher string, name string, version string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "publisher", runtime.ParamLocationPath, publisher)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "name", runtime.ParamLocationPath, name)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "version", runtime.ParamLocationPath, version)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1alpha1/providers/%s/%s/%s/usage", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ListAllProviders request
	ListAllProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListAllProvidersResponse, error)

	// GetProvider request
	GetProviderWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderResponse, error)

	// GetProviderSetupDocs request
	GetProviderSetupDocsWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderSetupDocsResponse, error)

	// GetProviderUsageDoc request
	GetProviderUsageDocWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderUsageDocResponse, error)
}

type ListAllProvidersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Next      *string          `json:"next"`
		Providers []ProviderDetail `json:"providers"`
	}
	JSON500 *struct {
		Error string `json:"error"`
	}
}

// Status returns HTTPResponse.Status
func (r ListAllProvidersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListAllProvidersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProviderResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ProviderDetail
	JSON404      *struct {
		Error string `json:"error"`
	}
	JSON500 *struct {
		Error string `json:"error"`
	}
}

// Status returns HTTPResponse.Status
func (r GetProviderResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProviderResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProviderSetupDocsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]string
}

// Status returns HTTPResponse.Status
func (r GetProviderSetupDocsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProviderSetupDocsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProviderUsageDocResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]string
}

// Status returns HTTPResponse.Status
func (r GetProviderUsageDocResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProviderUsageDocResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListAllProvidersWithResponse request returning *ListAllProvidersResponse
func (c *ClientWithResponses) ListAllProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListAllProvidersResponse, error) {
	rsp, err := c.ListAllProviders(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListAllProvidersResponse(rsp)
}

// GetProviderWithResponse request returning *GetProviderResponse
func (c *ClientWithResponses) GetProviderWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderResponse, error) {
	rsp, err := c.GetProvider(ctx, publisher, name, version, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProviderResponse(rsp)
}

// GetProviderSetupDocsWithResponse request returning *GetProviderSetupDocsResponse
func (c *ClientWithResponses) GetProviderSetupDocsWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderSetupDocsResponse, error) {
	rsp, err := c.GetProviderSetupDocs(ctx, publisher, name, version, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProviderSetupDocsResponse(rsp)
}

// GetProviderUsageDocWithResponse request returning *GetProviderUsageDocResponse
func (c *ClientWithResponses) GetProviderUsageDocWithResponse(ctx context.Context, publisher string, name string, version string, reqEditors ...RequestEditorFn) (*GetProviderUsageDocResponse, error) {
	rsp, err := c.GetProviderUsageDoc(ctx, publisher, name, version, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProviderUsageDocResponse(rsp)
}

// ParseListAllProvidersResponse parses an HTTP response from a ListAllProvidersWithResponse call
func ParseListAllProvidersResponse(rsp *http.Response) (*ListAllProvidersResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListAllProvidersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Next      *string          `json:"next"`
			Providers []ProviderDetail `json:"providers"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetProviderResponse parses an HTTP response from a GetProviderWithResponse call
func ParseGetProviderResponse(rsp *http.Response) (*GetProviderResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProviderResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ProviderDetail
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetProviderSetupDocsResponse parses an HTTP response from a GetProviderSetupDocsWithResponse call
func ParseGetProviderSetupDocsResponse(rsp *http.Response) (*GetProviderSetupDocsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProviderSetupDocsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetProviderUsageDocResponse parses an HTTP response from a GetProviderUsageDocWithResponse call
func ParseGetProviderUsageDocResponse(rsp *http.Response) (*GetProviderUsageDocResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProviderUsageDocResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List Providers
	// (GET /v1alpha1/providers)
	ListAllProviders(w http.ResponseWriter, r *http.Request)
	// Get Provider
	// (GET /v1alpha1/providers/{publisher}/{name}/{version})
	GetProvider(w http.ResponseWriter, r *http.Request, publisher string, name string, version string)
	// Get Provider Setup Docs
	// (GET /v1alpha1/providers/{publisher}/{name}/{version}/setup)
	GetProviderSetupDocs(w http.ResponseWriter, r *http.Request, publisher string, name string, version string)
	// Get Provider Usage Doc
	// (GET /v1alpha1/providers/{publisher}/{name}/{version}/usage)
	GetProviderUsageDoc(w http.ResponseWriter, r *http.Request, publisher string, name string, version string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// ListAllProviders operation middleware
func (siw *ServerInterfaceWrapper) ListAllProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListAllProviders(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetProvider operation middleware
func (siw *ServerInterfaceWrapper) GetProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "publisher" -------------
	var publisher string

	err = runtime.BindStyledParameter("simple", false, "publisher", chi.URLParam(r, "publisher"), &publisher)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "publisher", Err: err})
		return
	}

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// ------------- Path parameter "version" -------------
	var version string

	err = runtime.BindStyledParameter("simple", false, "version", chi.URLParam(r, "version"), &version)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "version", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProvider(w, r, publisher, name, version)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetProviderSetupDocs operation middleware
func (siw *ServerInterfaceWrapper) GetProviderSetupDocs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "publisher" -------------
	var publisher string

	err = runtime.BindStyledParameter("simple", false, "publisher", chi.URLParam(r, "publisher"), &publisher)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "publisher", Err: err})
		return
	}

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// ------------- Path parameter "version" -------------
	var version string

	err = runtime.BindStyledParameter("simple", false, "version", chi.URLParam(r, "version"), &version)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "version", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProviderSetupDocs(w, r, publisher, name, version)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetProviderUsageDoc operation middleware
func (siw *ServerInterfaceWrapper) GetProviderUsageDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "publisher" -------------
	var publisher string

	err = runtime.BindStyledParameter("simple", false, "publisher", chi.URLParam(r, "publisher"), &publisher)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "publisher", Err: err})
		return
	}

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", chi.URLParam(r, "name"), &name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// ------------- Path parameter "version" -------------
	var version string

	err = runtime.BindStyledParameter("simple", false, "version", chi.URLParam(r, "version"), &version)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "version", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProviderUsageDoc(w, r, publisher, name, version)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1alpha1/providers", wrapper.ListAllProviders)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1alpha1/providers/{publisher}/{name}/{version}", wrapper.GetProvider)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1alpha1/providers/{publisher}/{name}/{version}/setup", wrapper.GetProviderSetupDocs)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1alpha1/providers/{publisher}/{name}/{version}/usage", wrapper.GetProviderUsageDoc)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"FG8eAEco22V/BwAA//8Zj4DRtRsAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
