package player_repository

import (
	"boilerplate-demo/src/entity/player"
	"strconv"
)

func playerToGormPlayer(p player.Player) (postgresGormPlayerModel, error) {
	var playerID uint
	if p.ID != "" {
		playerID64, err := strconv.ParseUint(p.ID, 10, 0)
		if err != nil {
			return postgresGormPlayerModel{}, err
		}

		playerID = uint(playerID64)
	}

	return postgresGormPlayerModel{
		PlayerId: playerID,
		Name:     p.Name,
		Address:  p.Address,
	}, nil
}
