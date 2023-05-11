package types

import (
	"context"

	"github.com/acs-dl/auth-svc/internal/config"
)

type Runner = func(context context.Context, config config.Config)
