# Host Route Library

A high-performance Echo middleware library for routing based on the host. This library facilitates the configuration of different routes and behaviors for distinct hostnames, enhancing the ability to host multi-tenant applications on a single server.

[![go report card](https://goreportcard.com/badge/github.com/YidiDev/echo-host-route "go report card")](https://goreportcard.com/report/github.com/YidiDev/echo-host-route)
[![Go](https://github.com/YidiDev/echo-host-route/actions/workflows/tests.yml/badge.svg)](https://github.com/YidiDev/echo-host-route/actions/workflows/tests.yml)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/YidiDev/echo-host-route?tab=doc)

## Installation

Add the module to your project by running:

```sh
go get github.com/YidiDev/echo-host-route
```

## Usage

The following example demonstrates how to use the library to define various routes based on the host name. This helps in setting up multiple applications or APIs served from the same server, each with its specific routing configuration.

### Example

```go
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
```

## Configuration Options

### `HostConfig`
The `HostConfig` struct is used to define the configuration for a specific host:
- `Host`: The hostname for which the configuration is defined.
- `Prefix`: A prefix to use for routes specific to this host when accessed on a generic host.
- `RouterFactory`: A function that sets up routing for this host, taking an `*echo.Group` instance to define its routes.

### Generic Hosts
Generic hosts are hosts that will have access to all routes defined in all the host configs and any others defined on the default router. This is useful for:
- **Local Testing**: to be able to access all routes without changing the host. 
- **Consolidated Access**: Handle routes from multiple applications on a single host. For example:
  - You have two applications hosted on one Go server: one at `application1.example.com` and the other at `application2.example.com`. However, you also want people to be able to access both applications by going to `example.com/application1` or `example.com/application2`.

### Secure Against Unknown Hosts
The `secureAgainstUnknownHosts` boolean flag controls how the middleware handles requests from unknown hosts:
- `true`: Requests from unknown hosts will receive a 404 Not Found Response. This is useful for securing your application against unexpected or unauthorized hosts.
- `false`: Requests from unknown hosts will be passed through the primary router. This is useful if you want to catch and handle such requests manually.

### Additional Host Config
This param is optional and allows for unlimited inputs. Each input should be a `func(*echo.Group) error`. This is meant for specifying functions that `SetupHostBasedRoutes` should run on every host group after creating it. Common use cases of this are:
- Configuring a `RouteNotFound` Handler. 
- Configuring Host Specific Middleware. This can be done in the `HostConfig` in the `RouterFactory`. Alternatively, it could be done here. This may be useful if you want to centralize a lot of the host-specific middleware.

### Handling Different Hosts

1. **Host-specific Routes**:
   Routes are defined uniquely for each host using a specific `RouterFactory`. The `HostConfig` struct includes the hostname, path prefix, and a function to define routes for that host.

    ```go
    hostConfigs := []hostroute.HostConfig{
        {Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
        {Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
    }
    ```

2. **Generic Hosts**:
   Generic hosts allow for fallback to common routes defined in the primary router.

    ```go
    genericHosts := []string{"host3.com", "host4.com"}
    ```

3. **Secure Against Unknown Hosts**:
   Secure your application by handling unknown hosts, preventing them from accessing unintended routes.

    ```go
    hostroute.SetupHostBasedRoutes(r, hostConfigs, genericHosts, true)
    ```
   
## Sister Project
This project has a sister project for Gin framework users. If you are using Gin, check out the [Echo Host Route Library](https://github.com/YidiDev/gin-host-route) for similar functionality.

## Contributing
Contributions are always welcome! If you're interested in contributing to the project, please take a look at our [Contributing Guidelines](CONTRIBUTING.md) file for guidelines on how to get started. We appreciate your help in improving the library!
