package route

import (
	"boilerplate-demo/src/interface/fiber_server/helper"

	player_spec "boilerplate-demo/src/interface/fiber_server/spec/player"

	"github.com/gofiber/fiber/v2"
)

func (r routeV1) PostPlayer(c *fiber.Ctx) error {
	p := player_spec.PlayerRequest{}
	if err := c.BodyParser(&p); err != nil {
		return helper.ErrorHandler(c, err)
	}

	pEntity := toEntityPlayer(p)

	playerId, err := r.useCase.CreatePlayer(c.Context(), pEntity)
	if err != nil {
		return helper.ErrorHandler(c, err)
	}

	return c.SendString(playerId)
}

func (r routeV1) GetPlayerPlayerId(c *fiber.Ctx, pid int) error {
	p, err := r.useCase.GetPlayer(c.Context(), pid)
	if err != nil {
		return helper.ErrorHandler(c, err)
	}
	return c.JSON(toPlayer(p))
}
