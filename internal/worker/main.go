package worker

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/acs/auth/internal/config"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/data/postgres"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const serviceName = data.ModuleName + "-worker"

type Worker interface {
	Run(ctx context.Context)
}

type worker struct {
	logger         *logan.Entry
	refreshTokensQ data.RefreshTokens
}

func NewWorker(cfg config.Config) Worker {
	return &worker{
		logger:         cfg.Log().WithField("runner", serviceName),
		refreshTokensQ: postgres.NewRefreshTokensQ(cfg.DB()),
	}
}

func (w *worker) Run(ctx context.Context) {
	running.WithBackOff(
		ctx,
		w.logger,
		serviceName,
		w.processWork,
		30*time.Minute,
		30*time.Minute,
		30*time.Minute,
	)
}

func (w *worker) processWork(_ context.Context) error {
	w.logger.Info("started worker")

	err := w.removeExpiredTokens()
	if err != nil {
		return errors.Wrap(err, " failed to remove expired refresh tokens")
	}

	w.logger.Info("finished worker")
	return nil
}

func (w *worker) removeExpiredTokens() error {
	w.logger.Info("started removing expired tokens")

	tokens, err := w.refreshTokensQ.FilterByValidTill(time.Now().Unix()).Select()
	if err != nil {
		return errors.Wrap(err, " failed to select refresh tokens")
	}

	w.logger.Infof("found `%d` tokens to remove", len(tokens))

	for _, token := range tokens {
		err = w.refreshTokensQ.Delete(token.Token)
		if err != nil {
			return errors.Wrap(err, " failed to delete refresh token")
		}
	}

	w.logger.Info("finished removing expired tokens")
	return nil
}
