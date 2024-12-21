package middleware

import (
	"boilerplate-demo/src/interface/fiber_server/helper"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Recovery(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				helper.ErrorHandler(c, fmt.Errorf("%v", r))
			} else {
				helper.ErrorHandler(c, err)
			}
		}
	}()

	return c.Next()
}
