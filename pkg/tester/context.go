package tester

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manta-coder/golang-serverless-example/pkg/server"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
)

type ContextOptions func(req *http.Request, c echo.Context)

func NopHandlerFunc(c echo.Context) error {
	return nil
}

func NewContext(opts ...ContextOptions) (echo.Context, *httptest.ResponseRecorder) {
	e := server.NewEcho(GetLogger())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	for _, opt := range opts {
		opt(req, c)
	}

	return c, res
}

// WithJwtClaims adds JWT token with the provided claims to the context
func WithJwtClaims(claims jwt.MapClaims) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		c.Set(middleware.DefaultJWTConfig.ContextKey, token)
	}
}

// WithUser adds a user to the context
func WithUser(user interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		// key comes from odin/middleware.go
		c.Set("AuthenticatedUser", user)
	}
}

// WithContact adds a contact to the context
func WithContact(contact interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		// key comes from odin/middleware.go
		c.Set("AuthenticatedContact", contact)
	}
}

// WithProvider adds a provider to the context
func WithProvider(provider interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		// key comes from odin/middleware.go
		c.Set("AuthenticatedProvider", provider)
	}
}

// WithJsonBody marshals v into the body of the request
func WithJsonBody(v interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		body, _ := json.Marshal(v)

		r := bytes.NewReader(body)

		req.ContentLength = int64(r.Len())
		req.Body = ioutil.NopCloser(r)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
}

func WithParams(keyvalues ...string) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		var names []string
		var values []string

		for i := 0; i < len(keyvalues); i = i + 2 {
			names = append(names, keyvalues[i])
			values = append(values, keyvalues[i+1])
		}

		c.SetParamNames(names...)
		c.SetParamValues(values...)
	}
}

func WithQueryParams(keyvalues ...string) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		q := req.URL.Query()

		for i := 0; i < len(keyvalues); i = i + 2 {
			key := keyvalues[i]
			value := keyvalues[i+1]

			q[key] = append(q[key], value)
		}

		req.URL.RawQuery = q.Encode()
	}
}

// WithData adds data to the context
func WithData(key string, data interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		c.Set(key, data)
	}
}

// WithMultipartForm adds a file and form fields to the context
func WithMultipartForm(key string, file []byte, filename string, fields map[string]interface{}) ContextOptions {
	return func(req *http.Request, c echo.Context) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		defer writer.Close()
		part, _ := writer.CreateFormFile(key, filename)
		part.Write(file)

		for i, f := range fields {
			writer.WriteField(i, f.(string))
		}

		req.Body = ioutil.NopCloser(body)
		c.FormParams()
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	}
}
