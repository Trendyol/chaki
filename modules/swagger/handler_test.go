package swagger

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_fiberWrapper(t *testing.T) {
	mockDocs := Docs{}
	app := fiber.New()
	wrappedApp := fiberWrapper(mockDocs)(app)

	assert.NotNil(t, wrappedApp)
	assert.Equal(t, uint32(2), wrappedApp.HandlersCount())
}

func Test_newRedirectMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(newRedirectMiddleware())

	tests := []struct {
		path       string
		expected   string
		statusCode int
	}{
		{"/", "/swagger/index.html", http.StatusFound},
		{"/swagger", "/swagger/index.html", http.StatusFound},
		{"/swagger.json", "/swagger/docs.json", http.StatusFound},
		{"/swagger/v1/swagger.json", "/swagger/docs.json", http.StatusFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		resp, _ := app.Test(req)

		assert.Equal(t, tt.statusCode, resp.StatusCode)
		assert.Equal(t, tt.expected, resp.Header.Get("Location"))
	}
}

func Test_newMiddleware(t *testing.T) {
	mockDocs := Docs{}.WithHost("localhost")
	fmt.Println(mockDocs)
	app := fiber.New()
	app.Use(newMiddleware(mockDocs))

	t.Run("Serve JSON for /swagger/docs.json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/docs.json", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Host")
	})

	t.Run("Serve static files for /swagger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
