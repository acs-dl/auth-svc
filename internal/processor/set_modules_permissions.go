package processor

import (
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) handleSetModulesPermissionsAction(msg data.ModulePayload) error {
	p.log.Infof("started handling message with id `%s`", msg.RequestId)

	if msg.ModulePermissions == nil {
		p.log.Errorf("module permissions are empty")
		return errors.Errorf("module permissions are empty")
	}

	for module, statusPermissions := range msg.ModulePermissions {
		err := p.modulesQ.Upsert(data.Module{
			Name: module,
		})
		if err != nil {
			p.log.WithError(err).Errorf("failed to upsert module")
			return errors.Wrap(err, "failed to upsert module")
		}

		dbModule, err := p.modulesQ.GetByName(module)
		if err != nil {
			p.log.WithError(err).Errorf("failed to get module")
			return errors.Wrap(err, "failed to get module")
		}

		if dbModule == nil {
			p.log.Errorf("no such module `%s`", module)
			return errors.Errorf("no such module `%s`", module)
		}

		for status, permission := range statusPermissions {
			err = p.permissionsQ.Upsert(data.Permission{
				ModuleId: dbModule.Id,
				Name:     permission,
				Status:   data.UserStatus(status),
			})
			if err != nil {
				p.log.WithError(err).Errorf("failed to upsert permission")
				return errors.Wrap(err, "failed to upsert permission")
			}
		}
	}

	p.log.Infof("finished handling message with id `%s`", msg.RequestId)
	return nil
}
