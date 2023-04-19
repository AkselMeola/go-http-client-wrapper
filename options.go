package http_client_wrapper

import "net/http"

type RequestOptionFunc func(*http.Request)

func WithQueryParams(params map[string]string) RequestOptionFunc {
	return func(req *http.Request) {
		values := req.URL.Query()
		for key, value := range params {
			values.Add(key, value)
		}

		req.URL.RawQuery = values.Encode()
	}
}

func WithHeaders(headers map[string]string) RequestOptionFunc {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
}

func WithContentType(contentType string) RequestOptionFunc {
	return WithHeaders(map[string]string{
		"Content-Type": contentType,
	})
}

func WithAcceptHeader(contentType string) RequestOptionFunc {
	return WithHeaders(map[string]string{
		"Accept": contentType,
	})
}
