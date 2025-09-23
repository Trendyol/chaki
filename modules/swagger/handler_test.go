package swagger

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_fiberWrapper(t *testing.T) {
	mockDocs := Docs{}
	cfg := createTestConfig(map[string]interface{}{})
	app := fiber.New()
	wrappedApp := fiberWrapper(mockDocs, cfg)(app)

	assert.NotNil(t, wrappedApp)
	assert.Equal(t, uint32(3), wrappedApp.HandlersCount()) // host middleware + redirect + swagger
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
		defer resp.Body.Close()

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
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Host")
	})

	t.Run("Serve static files for /swagger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func createTestConfig(data map[string]interface{}) *config.Config {
	v := viper.New()
	for k, val := range data {
		v.Set(k, val)
	}
	return config.NewConfig(v, map[string]*viper.Viper{})
}

func Test_newHostAccessMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		config         map[string]interface{}
		host           string
		xForwardedHost string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "no swagger config - should allow all",
			config:         map[string]interface{}{},
			host:           "api.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			description:    "When no swagger config exists, all hosts should be allowed",
		},
		{
			name: "blocked host - should deny",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"api.example.com", "external.company.com"},
			},
			host:           "api.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "Host matching blocked list should be denied",
		},
		{
			name: "blocked host substring - should deny",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"example.com"},
			},
			host:           "api.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "Host containing blocked substring should be denied",
		},
		{
			name: "non-blocked host - should allow",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"external.company.com"},
			},
			host:           "internal.company.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			description:    "Host not in blocked list should be allowed",
		},
		{
			name: "allowed host only - should allow",
			config: map[string]interface{}{
				"server.swagger.allowedHostsContains": []string{"internal", "localhost"},
			},
			host:           "internal.company.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			description:    "Host matching allowed list should be allowed",
		},
		{
			name: "not in allowed host list - should deny",
			config: map[string]interface{}{
				"server.swagger.allowedHostsContains": []string{"internal", "localhost"},
			},
			host:           "api.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "Host not in allowed list should be denied",
		},
		{
			name: "blocked takes precedence over allowed for security",
			config: map[string]interface{}{
				"server.swagger.allowedHostsContains": []string{"internal"},
				"server.swagger.blockedHostsContains": []string{"example.com"},
			},
			host:           "internal.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "Blocked list should take precedence over allowed list for better security",
		},
		{
			name: "x-forwarded-host takes precedence",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"api.example.com"},
			},
			host:           "internal.company.com",
			xForwardedHost: "api.example.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusNotFound,
			description:    "X-Forwarded-Host header should be used instead of Host",
		},
		{
			name: "non-swagger path - should allow",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"api.example.com"},
			},
			host:           "api.example.com",
			path:           "/api/users",
			expectedStatus: http.StatusOK, // Should reach the catch-all handler
			description:    "Non-swagger paths should not be affected by host restrictions",
		},
		{
			name: "root path redirect blocked",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"external"},
			},
			host:           "external.api.com",
			path:           "/",
			expectedStatus: http.StatusNotFound,
			description:    "Root path should be blocked when host is in blocked list",
		},
		{
			name: "swagger.json blocked",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"external"},
			},
			host:           "external.api.com",
			path:           "/swagger.json",
			expectedStatus: http.StatusNotFound,
			description:    "swagger.json endpoint should be blocked",
		},
		{
			name: "localhost development allowed",
			config: map[string]interface{}{
				"server.swagger.allowedHostsContains": []string{"localhost", "127.0.0.1"},
			},
			host:           "localhost:8080",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			description:    "Localhost should be allowed for development",
		},
		{
			name: "empty blocked list item ignored",
			config: map[string]interface{}{
				"server.swagger.blockedHostsContains": []string{"", "api.example.com"},
			},
			host:           "any.host.com",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			description:    "Empty strings in blocked list should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createTestConfig(tt.config)
			middleware := newHostAccessMiddleware(cfg)

			app := fiber.New()
			app.Use(middleware)

			// Add a simple handler to test middleware behavior
			app.All("*", func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Host = tt.host
			if tt.xForwardedHost != "" {
				req.Header.Set("X-Forwarded-Host", tt.xForwardedHost)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, tt.description)
		})
	}
}

func Test_getEffectiveHost(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		xForwardedHost string
		expected       string
	}{
		{
			name:     "only host header",
			host:     "example.com",
			expected: "example.com",
		},
		{
			name:           "x-forwarded-host takes precedence",
			host:           "internal.com",
			xForwardedHost: "external.com",
			expected:       "external.com",
		},
		{
			name:           "empty x-forwarded-host falls back to host",
			host:           "example.com",
			xForwardedHost: "",
			expected:       "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			var actualHost string

			app.Get("/test", func(c *fiber.Ctx) error {
				actualHost = getEffectiveHost(c)
				return c.SendStatus(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Host = tt.host
			if tt.xForwardedHost != "" {
				req.Header.Set("X-Forwarded-Host", tt.xForwardedHost)
			}

			_, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actualHost)
		})
	}
}

func Test_isSwaggerRequest(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", true},
		{"/swagger", true},
		{"/swagger/", true},
		{"/swagger/index.html", true},
		{"/swagger/docs.json", true},
		{"/swagger.json", true},
		{"/swagger/v1/swagger.json", true},
		{"/swagger/assets/style.css", true},
		{"/api/users", false},
		{"/health", false},
		{"/metrics", false},
		{"/swaggerx", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			actual := isSwaggerRequest(tt.path)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_containsAny(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		subs     []string
		expected bool
	}{
		{
			name:     "exact match",
			s:        "example.com",
			subs:     []string{"example.com"},
			expected: true,
		},
		{
			name:     "substring match",
			s:        "api.example.com",
			subs:     []string{"example.com"},
			expected: true,
		},
		{
			name:     "multiple substrings - first matches",
			s:        "api.example.com",
			subs:     []string{"example.com", "test.com"},
			expected: true,
		},
		{
			name:     "multiple substrings - second matches",
			s:        "api.test.com",
			subs:     []string{"example.com", "test.com"},
			expected: true,
		},
		{
			name:     "no match",
			s:        "internal.company.com",
			subs:     []string{"example.com", "test.com"},
			expected: false,
		},
		{
			name:     "empty substring ignored",
			s:        "example.com",
			subs:     []string{"", "nomatch"},
			expected: false,
		},
		{
			name:     "empty string input",
			s:        "",
			subs:     []string{"example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := containsAny(tt.s, tt.subs)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
