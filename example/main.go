package main

import (
	"github.com/YidiDev/echo-host-route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

// defineHost1Routes sets up the routes specific to host1.com.
func defineHost1Routes(group *echo.Group) {
	// Route to handle request to root URL of host1, returns greeting message.
	group.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host1")
	})
	// Route to handle request to /hi URL of host1, returns a different greeting message.
	group.GET("/hi", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi from host1")
	})
}

// defineHost2Routes sets up the routes specific to host2.com.
func defineHost2Routes(group *echo.Group) {
	// Route to handle request to root URL of host2, returns greeting message.
	group.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host2")
	})
	// Route to handle request to /hi URL of host2, includes log statement and returns a greeting message.
	group.GET("/hi", func(c echo.Context) error {
		log.Println("Important stuff")
		return c.String(http.StatusOK, "Hi from host2")
	})
}

// routeNotFoundSpecifier defines behavior for unspecified routes.
func routeNotFoundSpecifier(_ string, group *echo.Group) error {
	// Route handler for unspecified routes, returns a not-found message.
	group.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "No known route")
	})
	return nil
}

// init function sets up logging output to standard output.
func init() {
	log.SetOutput(os.Stdout)
}

// main function initializes the Echo instance and sets up host-based routing.
func main() {
	e := echo.New()             // Create a new Echo instance.
	e.Use(middleware.Recover()) // Middleware to recover from panics.

	hostConfigs := []hostroute.HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	// Setup host-based routes and handle any errors during setup.
	err := hostroute.SetupHostBasedRoutes(e, hostConfigs, genericHosts, true, routeNotFoundSpecifier)
	if err != nil {
		log.Fatal(err) // Log fatal error if setup fails.
	}

	// Start the server on port 8080.
	e.Start(":8080")
}
