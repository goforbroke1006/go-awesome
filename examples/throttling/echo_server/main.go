package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/go-awesome"
)

func main() {
	e := echo.New()

	throttler := awesome.NewThrottlerFirstEntry(5 * time.Second)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			const throttlingForPrefix = "/api"

			if strings.HasPrefix(ctx.Request().URL.Path, throttlingForPrefix) {
				if err := throttler.Throttle(ctx.RealIP()); err != nil {
					fmt.Println("throttle", ctx.RealIP(), time.Now().UTC().Format(time.RFC3339))
					return ctx.String(http.StatusTooManyRequests, "Too Many Requests")
				}
			}

			return next(ctx)
		}
	})

	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "pong")
	})
	e.GET("/api/hello", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	e.GET("/api/world", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})
	_ = e.Start(":8181")

	// http://localhost:8181/ping
	// http://localhost:8181/api/hello
	// http://localhost:8181/api/world
}
