package repository

import (
	"boilerplate-demo/src/entity/player"
	"context"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
	Name() string
}

type PlayerRepository interface {
	HealthChecker
	CreatePlayer(ctx context.Context, player player.Player) (string, error)
	GetPlayer(ctx context.Context, pId int64) (player.Player, error)
}
