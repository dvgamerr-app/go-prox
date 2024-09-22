package notice

// type NotifyPayload struct {
// 	Message string `json:"msg"`
// }

// // PUT /notify/:serviceName/:roomName with fiber router
// func HandlePutNotify(c *fiber.Ctx) error {
// 	var req *NotifyPayload = &NotifyPayload{
// 		Message: c.Query("msg"),
// 	}

// 	if req.Message != "" {
// 		var err error
// 		if req.Message, err = url.QueryUnescape(req.Message); err != nil {
// 			return fmt.Errorf("Error:", err)
// 		}
// 	} else {
// 		if err := c.BodyParser(&req); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 	}

// 	serviceName := c.Params("serviceName")
// 	roomName := c.Params("roomName")

// 	if serviceName == "" || roomName == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "serviceName and roomName are required",
// 		})
// 	}

// 	if err := h.service.PutNotify(c.Context(), serviceName, roomName, req.Message); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "success",
// 	})
// }
