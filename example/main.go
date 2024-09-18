package main

import (
	"github.com/YidiDev/echo-host-route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

func defineHost1Routes(group *echo.Group) {
	group.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host1")
	})
	group.GET("/hi", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi from host1")
	})
}

func defineHost2Routes(group *echo.Group) {
	group.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from host2")
	})
	group.GET("/hi", func(c echo.Context) error {
		log.Println("Important stuff")
		return c.String(http.StatusOK, "Hi from host2")
	})
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())

	hostConfigs := []hostroute.HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	// Setup host-based routes
	hostroute.SetupHostBasedRoutes(e, hostConfigs, genericHosts, true)

	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "No known route")
	})

	e.Start(":8080")
}
