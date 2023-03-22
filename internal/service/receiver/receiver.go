package receiver

import (
	"context"

	"gitlab.com/distributed_lab/acs/auth/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewReceiver(cfg).Run(ctx)
}
