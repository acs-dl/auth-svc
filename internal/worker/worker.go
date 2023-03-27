package worker

import (
	"context"

	"gitlab.com/distributed_lab/acs/auth/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewWorker(cfg).Run(ctx)
}
