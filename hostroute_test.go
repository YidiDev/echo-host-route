package hostroute

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func defineHost1Routes(g *echo.Group) {
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host1")
	})
	g.GET("/hi", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi from host1")
	})
}

func defineHost2Routes(rg *echo.Group) {
	rg.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host2")
	})
	rg.GET("/hi", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi from host2")
	})
}

func noRouteHandler(c echo.Context) error {
	return c.String(http.StatusNotFound, "No known route")
}

func TestHostBasedRouting(t *testing.T) {
	r := echo.New()
	r.RouteNotFound("/*", noRouteHandler)

	hostConfigs := []HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	SetupHostBasedRoutes(r, hostConfigs, genericHosts, true)

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		host, path string
		expected   string
		statusCode int
	}{
		{"host1.com", "/", "Hello from host1", http.StatusOK},
		{"host1.com", "/hi", "Hi from host1", http.StatusOK},
		{"host1.com", "/unknown", "No known route", http.StatusNotFound},

		{"host2.com", "/", "Hello from host2", http.StatusOK},
		{"host2.com", "/hi", "Hi from host2", http.StatusOK},
		{"host2.com", "/unknown", "No known route", http.StatusNotFound},

		{"host3.com", "/1", "Hello from host1", http.StatusOK},
		{"host3.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host3.com", "/2", "Hello from host2", http.StatusOK},
		{"host3.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host3.com", "/unknown", "No known route", http.StatusNotFound},

		{"host4.com", "/1", "Hello from host1", http.StatusOK},
		{"host4.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host4.com", "/2", "Hello from host2", http.StatusOK},
		{"host4.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host4.com", "/unknown", "No known route", http.StatusNotFound},

		{"unknown.com", "/", "Unknown host", http.StatusNotFound},
	}

	client := &http.Client{}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", server.URL+tt.path, nil)
		req.Host = tt.host
		resp, err := client.Do(req)

		assert.NoError(t, err)

		body := make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(body)
		assert.NoError(t, err)

		assert.Equal(t, tt.statusCode, resp.StatusCode)
		assert.Equal(t, tt.expected, string(body))
		err = resp.Body.Close()
		assert.NoError(t, err)
	}
}

func TestHostBasedRoutingWithoutSecureAgainstUnknownHosts(t *testing.T) {
	r := echo.New()
	r.RouteNotFound("/*", noRouteHandler)

	hostConfigs := []HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	SetupHostBasedRoutes(r, hostConfigs, genericHosts, false)

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		host, path string
		expected   string
		statusCode int
	}{
		{"host1.com", "/", "Hello from host1", http.StatusOK},
		{"host1.com", "/hi", "Hi from host1", http.StatusOK},
		{"host1.com", "/unknown", "No known route", http.StatusNotFound},

		// Host 2 specific routes
		{"host2.com", "/", "Hello from host2", http.StatusOK},
		{"host2.com", "/hi", "Hi from host2", http.StatusOK},
		{"host2.com", "/unknown", "No known route", http.StatusNotFound},

		// Generic Host Routes
		{"host3.com", "/1", "Hello from host1", http.StatusOK},
		{"host3.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host3.com", "/2", "Hello from host2", http.StatusOK},
		{"host3.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host3.com", "/unknown", "No known route", http.StatusNotFound},

		{"host4.com", "/1", "Hello from host1", http.StatusOK},
		{"host4.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host4.com", "/2", "Hello from host2", http.StatusOK},
		{"host4.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host4.com", "/unknown", "No known route", http.StatusNotFound},

		{"unknown.com", "/", "No known route", http.StatusNotFound},
	}

	client := &http.Client{}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", server.URL+tt.path, nil)
		req.Host = tt.host
		resp, err := client.Do(req)

		assert.NoError(t, err)

		body := make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(body)
		assert.NoError(t, err)

		assert.Equal(t, tt.statusCode, resp.StatusCode)
		assert.Equal(t, tt.expected, string(body))
		err = resp.Body.Close()
		assert.NoError(t, err)

	}
}
