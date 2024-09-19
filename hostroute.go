package hostroute

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HostConfig holds the configuration for each host.
type HostConfig struct {
	Host          string              // The specific hostname (e.g., "host1.com").
	Prefix        string              // Prefix for route paths (e.g., "1" or "2").
	RouterFactory func(e *echo.Group) // Function to define routes for the host.
}

// SecureAgainstUnknownHosts returns a middleware function to secure the server against requests from unknown hosts.
func SecureAgainstUnknownHosts(knownHosts map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := c.Request().Host

			if _, known := knownHosts[host]; !known {
				// If host is not recognized, return a 404 Not Found response.
				return c.String(http.StatusNotFound, "Unknown host")
			}
			return next(c)
		}
	}
}

// SetupHostBasedRoutes configures routing based on hostnames.
func SetupHostBasedRoutes(e *echo.Echo, hostConfigs []HostConfig, genericHosts []string, secureAgainstUnknownHosts bool, additionalHostConfig ...func(string, *echo.Group) error) error {
	var allHosts []string

	for _, hostConfig := range hostConfigs {
		hostGroup := e.Host(hostConfig.Host) // Create a host-specific route group
		hostConfig.RouterFactory(hostGroup)  // Set up routes using the provided factory function
		if secureAgainstUnknownHosts {
			allHosts = append(allHosts, hostConfig.Host) // Keep track of known hosts
		}
		if len(additionalHostConfig) > 0 {
			for _, config := range additionalHostConfig {
				err := config(hostConfig.Host, hostGroup) // Apply additional configurations if provided
				if err != nil {
					return err // Return on error
				}
			}
		}
	}

	for _, genericHost := range genericHosts {
		genericGroup := e.Host(genericHost) // Create a group for generic hosts

		for _, hostConfig := range hostConfigs {
			if hostConfig.Prefix != "" {
				// Create prefixed routes for generic hosts
				prefixedGroup := genericGroup.Group(fmt.Sprintf("/%s", hostConfig.Prefix))
				hostConfig.RouterFactory(prefixedGroup) // Set up routes for each prefix
			}
		}

		if secureAgainstUnknownHosts {
			allHosts = append(allHosts, genericHost) // Track generic known hosts
		}

		if len(additionalHostConfig) > 0 {
			for _, config := range additionalHostConfig {
				err := config(genericHost, genericGroup) // Apply additional configurations
				if err != nil {
					return err // Return on error
				}
			}
		}
	}

	if secureAgainstUnknownHosts {
		// Apply the security middleware against unknown hosts
		e.Use(SecureAgainstUnknownHosts(stringSliceToMap(allHosts)))
	}

	return nil
}

// stringSliceToMap converts a slice of strings to a map with the string as the key and true as the value.
func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range slice {
		result[s] = true
	}
	return result
}
