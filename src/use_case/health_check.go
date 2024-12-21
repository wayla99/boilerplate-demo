package use_case

import (
	"boilerplate-demo/src/use_case/repository"
	"context"
	"fmt"
)

func (uc UseCase) HealthCheck(ctx context.Context) error {
	return healthChecks(ctx, uc.playerRepository)
}

func healthChecks(ctx context.Context, healthCheckers ...repository.HealthChecker) error {
	errCh := make(chan error, len(healthCheckers))

	for _, checker := range healthCheckers {
		go func(checker repository.HealthChecker) {
			if err := checker.HealthCheck(ctx); err != nil {
				errCh <- fmt.Errorf("health check failed for [%s]: %w", checker.Name(), err)
			} else {
				errCh <- nil
			}
		}(checker)
	}

	for i := 0; i < len(healthCheckers); i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}
