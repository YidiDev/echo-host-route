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

			if _, specificHostExists := hostConfigMap[host]; specificHostExists {
				return next(c)
			}

			if _, isGenericHost := genericHosts[host]; isGenericHost {
				return next(c)
			}

			if secureAgainstUnknownHosts {
				return c.String(http.StatusNotFound, "Unknown host")
			}

			return next(c)
		}
	}
}

func SetupHostBasedRoutes(e *echo.Echo, hostConfigs []HostConfig, genericHosts []string, noRouteFactory func(*echo.Group), secureAgainstUnknownHosts bool) {
	hostConfigMap := make(map[string]*HostConfig)
	genericHostsMap := stringSliceToMap(genericHosts)

	for _, hostConfig := range hostConfigs {
		hostGroup := e.Host(hostConfig.Host).Group("")
		hostConfig.RouterFactory(hostGroup)
		noRouteFactory(hostGroup)

		hostConfigMap[hostConfig.Host] = &hostConfig
	}

	for _, genericHost := range genericHosts {
		genericGroup := e.Host(genericHost).Group("")

		for _, hostConfig := range hostConfigs {
			if hostConfig.Prefix != "" {
				prefixedGroup := genericGroup.Group(fmt.Sprintf("/%s", hostConfig.Prefix))
				hostConfig.RouterFactory(prefixedGroup)
			}
		}

		noRouteFactory(genericGroup)
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
