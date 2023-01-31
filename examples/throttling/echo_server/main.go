package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/go-awesome"
)

func main() {
	e := echo.New()

	t := awesome.NewThrottler()

	e.GET("/", func(ctx echo.Context) error {
		if err := t.Throttle(ctx.Request().URL.Path+"|"+ctx.RealIP(), 5*time.Second); err != nil {
			return ctx.String(http.StatusTooManyRequests, "Too Many Requests")
		}

		return ctx.String(http.StatusOK, "OK")
	})
	_ = e.Start(":8181")
}
