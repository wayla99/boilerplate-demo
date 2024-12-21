package route

import (
	player_spec "boilerplate-demo/src/interface/fiber_server/spec/player"
	"boilerplate-demo/src/use_case"
)

type routeV1 struct {
	useCase *use_case.UseCase
}

func NewRouteV1(useCase *use_case.UseCase) player_spec.ServerInterface {
	return routeV1{useCase: useCase}
}
