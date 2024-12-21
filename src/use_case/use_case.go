package use_case

import "boilerplate-demo/src/use_case/repository"

type UseCase struct {
	playerRepository repository.PlayerRepository
}

func New(playerRepository repository.PlayerRepository) *UseCase {
	return &UseCase{
		playerRepository: playerRepository,
	}
}
