package worker

import (
	"context"

	"github.com/acs-dl/auth-svc/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewWorker(cfg).Run(ctx)
}
