package api

import (
	daas "github.com/touno-io/goasa"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func ThrowTraceHandler(c *fiber.Ctx, code int, err error, ignoreCapture ...bool) error {
	daas.Error(err)
	if len(ignoreCapture) > 0 {
		sentry.CaptureException(err)
	}
	return c.Status(code).JSON((HTTP{Code: code, Error: err.Error()}))
}
