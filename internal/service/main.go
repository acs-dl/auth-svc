package service

import (
	"context"
	"sync"

	"github.com/acs-dl/auth-svc/internal/receiver"
	"github.com/acs-dl/auth-svc/internal/sender"
	"github.com/acs-dl/auth-svc/internal/service/api"
	"github.com/acs-dl/auth-svc/internal/worker"

	"github.com/acs-dl/auth-svc/internal/config"
	"github.com/acs-dl/auth-svc/internal/service/types"
)

var availableServices = map[string]types.Runner{
	"api":      api.Run,
	"sender":   sender.Run,
	"receiver": receiver.Run,
	"worker":   worker.Run,
}

func Run(cfg config.Config) {
	logger := cfg.Log().WithField("service", "main")
	ctx := context.Background()
	wg := new(sync.WaitGroup)

	logger.Info("Starting all available services...")

	for serviceName, service := range availableServices {
		wg.Add(1)

		go func(name string, runner types.Runner) {
			defer wg.Done()

			runner(ctx, cfg)

		}(serviceName, service)

		logger.WithField("service", serviceName).Info("Service started")
	}

	wg.Wait()

}
