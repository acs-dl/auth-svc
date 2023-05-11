package processor

import (
	"github.com/acs-dl/auth-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) handleRemoveModuleAction(msg data.ModulePayload) error {
	p.log.Infof("started handling message with id `%s`", msg.RequestId)

	if msg.ModuleName == "" {
		p.log.Errorf("module name is empty")
		return errors.Errorf("module name is empty")
	}

	err := p.modulesQ.FilterByNames(msg.ModuleName).Delete()
	if err != nil {
		p.log.WithError(err).Errorf("failed to delete module")
		return errors.Wrap(err, "failed to delete module")
	}

	p.log.Infof("finished handling message with id `%s`", msg.RequestId)
	return nil
}
