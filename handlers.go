package main

import (
	"context"
	"net/url"
	"prox/pgsql"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	TELEGRAM   = "telegram"
	SLACK      = "slack"
	MSTEAM     = "msteam"
	LINE       = "line"
	LINENOTIFY = "line-notify"
	EMAIL      = "email"
	WEBHOOK    = "webhook"
	NATIVE     = "native"
)

type NotifyPayload struct {
	Message string `json:"msg"`
}

func fiberThrowError(c *fiber.Ctx, code int, err error) error {
	return c.Status(code).JSON(fiber.Map{"code": code, "error": err.Error()})
}

func handlerPutNotify(ctx context.Context, pgx *pgsql.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req *NotifyPayload = &NotifyPayload{
			Message: c.Query("msg"),
		}

		var err error
		if req.Message != "" {
			if req.Message, err = url.QueryUnescape(req.Message); err != nil {
				return fiberThrowError(c, fiber.StatusBadRequest, err)
			}
		} else {
			if err = c.BodyParser(&req); err != nil {
				return fiberThrowError(c, fiber.StatusBadRequest, err)
			}
		}
		epochMilliseconds := time.Now().UnixMilli()

		// result, err := pgx.DB.ExecContext(ctx, `
		// 	INSERT INTO "notice"."history" () VALUES
		// `)
		// err = stx.Execute(fmt.Sprintf(`
		// 	INSERT INTO "app"."notice_history" ("notice_room_id", "o_sender", "b_sended")
		// 	VALUES %s;
		// `, strings.Join(historyInserted, ",")))

		// if db.IsRollback(err, stx) {
		// 	return api.ErrorHandlerThrow(c, fiber.StatusInternalServerError, err)
		// }

		estimated := time.Now().UnixMilli() - epochMilliseconds

		return c.JSON(fiber.Map{"code": fiber.StatusOK, "used": estimated})
	}
}
