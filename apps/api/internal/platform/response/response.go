package response

import "github.com/gofiber/fiber/v2"

type Body struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedBody struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Body{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Body{Success: true, Data: data})
}

func Message(c *fiber.Ctx, message string) error {
	return c.JSON(Body{Success: true, Message: message})
}

func Paginated(c *fiber.Ctx, data interface{}, total int64, limit, offset int) error {
	return c.JSON(PaginatedBody{
		Success: true,
		Data:    data,
		Meta: Meta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}

func Error(c *fiber.Ctx, status int, message string, err interface{}) error {
	return c.Status(status).JSON(Body{Success: false, Message: message, Error: err})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message, nil)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message, nil)
}
