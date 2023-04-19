package http_client_wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Example_buildURL() {
	fmt.Println(buildUrlWoErr("http://www.example.com", "foo/bar"))
	fmt.Println(buildUrlWoErr("http://www.example.com", "foo/bar?baz=1"))
	fmt.Println(buildUrlWoErr("http://www.example.com", "/foo/bar"))
	fmt.Println(buildUrlWoErr("http://www.example.com/", "/double-slash"))
	// Output:
	// http://www.example.com/foo/bar
	// http://www.example.com/foo/bar?baz=1
	// http://www.example.com/foo/bar
	// http://www.example.com/double-slash
}

func buildUrlWoErr(basePath string, path string) string {
	url, _ := NewClient(basePath, 60).buildRequestUrl(path)
	return url
}

type HttpClientTestSuite struct {
	suite.Suite
}

func TestHttpClientTestSuite(t *testing.T) {
	suite.Run(t, new(HttpClientTestSuite))
}

func (s *HttpClientTestSuite) SetupTest() {

}

func (s *HttpClientTestSuite) TestHttpClientPost() {
	s.Run("Post form", func() {
		var responseData map[string]interface{}
		handler := func(r *http.Response) error {
			s.Equal(http.StatusOK, r.StatusCode)
			return json.NewDecoder(r.Body).Decode(&responseData)
		}

		err := NewClient("https://postman-echo.com", time.Second*60).
			PostForm(
				context.Background(),
				"post",
				map[string]string{
					"foo": "bar",
				},
				handler,
				WithQueryParams(map[string]string{
					"zoo": "zar",
				}),
			)

		if s.NoError(err) {
			s.Equal(
				map[string]interface{}{
					"zoo": "zar",
				},
				responseData["args"],
			)
			s.Equal(
				map[string]interface{}{
					"foo": "bar",
				},
				responseData["form"],
			)
		}
	})
	s.Run("Post data", func() {
		var responseData map[string]interface{}
		handler := func(r *http.Response) error {
			s.Equal(http.StatusOK, r.StatusCode)
			return json.NewDecoder(r.Body).Decode(&responseData)
		}

		err := NewClient("https://postman-echo.com", time.Second*60).
			Post(
				context.Background(),
				"post",
				strings.NewReader("my data of string kind"),
				handler,
			)

		if s.NoError(err) {
			s.Equal(
				map[string]interface{}{},
				responseData["args"],
			)
			s.Equal(
				map[string]interface{}{},
				responseData["form"],
			)
			s.Equal(
				"my data of string kind",
				responseData["data"],
			)
		}
	})
}
