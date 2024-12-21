package use_case

import (
	"boilerplate-demo/src/entity/player"
	"context"
)

func (uc UseCase) CreatePlayer(ctx context.Context, player player.Player) (string, error) {
	playerId, err := uc.playerRepository.CreatePlayer(ctx, player)
	if err != nil {
		return "", err
	}
	return playerId, nil
}

func (uc UseCase) GetPlayer(ctx context.Context, pId int) (player.Player, error) {
	p, err := uc.playerRepository.GetPlayer(ctx, int64(pId))
	if err != nil {
		return player.Player{}, err
	}

	return p, nil
}
