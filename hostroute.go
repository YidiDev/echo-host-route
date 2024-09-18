package hostroute

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HostConfig struct {
	Host          string
	Prefix        string
	RouterFactory func(e *echo.Group)
}

func createHostBasedRoutingMiddleware(hostConfigMap map[string]*HostConfig, genericHosts map[string]bool, secureAgainstUnknownHosts bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := c.Request().Host

			if _, exists := hostConfigMap[host]; exists {
				return next(c)
			}

			if _, exists := genericHosts[host]; exists {
				return next(c)
			}

			if secureAgainstUnknownHosts {
				return c.String(http.StatusNotFound, "Unknown host")
			}

			return next(c)
		}
	}
}

func SetupHostBasedRoutes(e *echo.Echo, hostConfigs []HostConfig, genericHosts []string, secureAgainstUnknownHosts bool) {
	hostConfigMap := make(map[string]*HostConfig)
	genericHostsMap := stringSliceToMap(genericHosts)

	for i := range hostConfigs {
		group := e.Host(hostConfigs[i].Host)
		hostConfigs[i].RouterFactory(group)

		if hostConfigs[i].Prefix != "" {
			group = e.Group(fmt.Sprintf("/%s", hostConfigs[i].Prefix))
			hostConfigs[i].RouterFactory(group)
		}

		hostConfigMap[hostConfigs[i].Host] = &hostConfigs[i]
	}

	e.Use(createHostBasedRoutingMiddleware(hostConfigMap, genericHostsMap, secureAgainstUnknownHosts))

}

func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range slice {
		result[s] = true
	}
	return result
}
