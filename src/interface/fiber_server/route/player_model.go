package route

import (
	"boilerplate-demo/src/entity/player"
	player_spec "boilerplate-demo/src/interface/fiber_server/spec/player"
	"strconv"
)

func toEntityPlayer(p player_spec.PlayerRequest) player.Player {
	return player.Player{
		Name:    p.Name,
		Address: *p.Address,
	}
}

func toPlayer(p player.Player) player_spec.PlayerResponse {
	pId, _ := strconv.Atoi(p.ID)

	return player_spec.PlayerResponse{
		Address:  &p.Address,
		Name:     &p.Name,
		PlayerId: &pId,
	}
}
