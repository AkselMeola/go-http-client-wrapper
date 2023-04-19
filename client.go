package http_client_wrapper

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestOption = func(r *http.Request)

type ResponseHandler func(r *http.Response) error

type HttpClient struct {
	client   *http.Client
	basePath string
}

func NewClient(basePath string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		basePath: strings.TrimRight(basePath, "/"),
	}
}

func (c *HttpClient) Post(
	ctx context.Context,
	path string,
	body io.Reader,
	handler ResponseHandler,
	opts ...RequestOption,
) error {
	requestUrl, err := c.buildRequestUrl(path)
	if err != nil {
		return err
	}

	req, err := c.createRequest(ctx, http.MethodPost, requestUrl, body, opts...)
	if err != nil {
		return NewError("create request error: "+err.Error(), err)
	}

	return c.executeRequest(req, handler)
}

func (c *HttpClient) PostForm(
	ctx context.Context,
	path string,
	formValues map[string]string,
	handler ResponseHandler,
	opts ...RequestOption,
) error {
	formData := url.Values{}
	for key, value := range formValues {
		formData.Set(key, value)
	}

	opts = append(opts, WithContentType("application/x-www-form-urlencoded"))

	return c.Post(ctx, path, strings.NewReader(formData.Encode()), handler, opts...)
}

func (c *HttpClient) Get(ctx context.Context, path string, handler ResponseHandler, opts ...RequestOption) error {
	requestUrl, err := c.buildRequestUrl(path)
	if err != nil {
		return err
	}

	req, err := c.createRequest(ctx, http.MethodGet, requestUrl, nil, opts...)
	if err != nil {
		return NewError("create request error: "+err.Error(), err)
	}

	return c.executeRequest(req, handler)
}

func (c *HttpClient) buildRequestUrl(path string) (string, error) {
	requestUrl, err := url.Parse(c.basePath + "/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return "", NewError("url parse error: "+err.Error(), err)
	}

	return requestUrl.String(), nil
}

func (c *HttpClient) createRequest(ctx context.Context, method string, url string, body io.Reader, opts ...RequestOption) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *HttpClient) executeRequest(req *http.Request, handler ResponseHandler) (err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return NewError("request error: "+err.Error(), err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return handler(resp)
}
