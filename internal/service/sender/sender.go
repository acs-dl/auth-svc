package sender

import (
	"context"

	"gitlab.com/distributed_lab/acs/auth/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewSender(cfg).Run(ctx)
}
