package route

import (
	"boilerplate-demo/src/interface/fiber_server/config"
	"boilerplate-demo/src/interface/fiber_server/helper"
	"boilerplate-demo/src/use_case"

	system_spec "boilerplate-demo/src/interface/fiber_server/spec/system"

	"github.com/gofiber/fiber/v2"
)

type routeSystem struct {
	config  *config.ServerConfig
	useCase *use_case.UseCase
}

func (r routeSystem) GetLiveliness(c *fiber.Ctx) error {
	return c.Send([]byte(helper.OK))
}

func (r routeSystem) GetLiveness(c *fiber.Ctx) error {
	return c.Send([]byte(helper.OK))
}

func (r routeSystem) GetReadiness(c *fiber.Ctx) error {
	err := r.useCase.HealthCheck(c.Context())
	if err != nil {
		return helper.ErrorHandler(c, err)
	}

	return c.Send([]byte(helper.OK))
}

func (r routeSystem) GetVersion(c *fiber.Ctx) error {
	return c.Send([]byte(r.config.AppVersion))
}

func NewRouteSystem(config *config.ServerConfig, useCase *use_case.UseCase) system_spec.ServerInterface {
	return routeSystem{
		config:  config,
		useCase: useCase,
	}
}
