package sender

import (
	"context"

	"github.com/acs-dl/auth-svc/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewSender(cfg).Run(ctx)
}
