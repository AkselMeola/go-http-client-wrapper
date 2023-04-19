# Go http client wrapper

A convenience wrapper around go http client. 

## Install

`go get github.com/AkselMeola/go-http-client-wrapper`

## Example usage

```go
var responseData map[string]interface{}
handlerFunc := func(r *http.Response) error {
    return json.NewDecoder(r.Body).Decode(&responseData)
}

err := NewClient("https://api.domain", time.Second*15).
    Post(
        context.Background(),
        "/resource/path",
        strings.NewReader(`{"some": "json"}`),
        handlerFunc,
    )
```