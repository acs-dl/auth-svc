package processor

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/acs/auth/internal/config"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/acs/auth/internal/data/postgres"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	serviceName = data.ModuleName + "-processor"

	//add needed actions for module
	GetModulesPermissionsAction = "get_modules_permissions"
	SetModulesPermissionsAction = "set_modules_permissions"
	RemoveModuleAction          = "remove_module"
)

type Processor interface {
	HandleNewMessage(msg data.ModulePayload) error
}

type processor struct {
	log          *logan.Entry
	modulesQ     data.Modules
	permissionsQ data.Permissions
	usersQ       data.Users
}

var handleActions = map[string]func(proc *processor, msg data.ModulePayload) error{
	SetModulesPermissionsAction: func(proc *processor, msg data.ModulePayload) error {
		return proc.handleSetModulesPermissionsAction(msg)
	},
	RemoveModuleAction: func(proc *processor, msg data.ModulePayload) error {
		return proc.handleRemoveModuleAction(msg)
	},
}

func NewProcessor(cfg config.Config) Processor {
	return &processor{
		log:          cfg.Log().WithField("service", serviceName),
		modulesQ:     postgres.NewModulesQ(cfg.DB()),
		permissionsQ: postgres.NewPermissionsQ(cfg.DB()),
		usersQ:       postgres.NewUsersQ(cfg.DB()),
	}
}

func (p *processor) HandleNewMessage(msg data.ModulePayload) error {
	p.log.Infof("handling message with id `%s`", msg.RequestId)

	err := validation.Errors{
		"action": validation.Validate(msg.Action, validation.Required, validation.In(SetModulesPermissionsAction, RemoveModuleAction)),
	}.Filter()
	if err != nil {
		p.log.WithError(err).Errorf("no such action to handle for message with id `%s`", msg.RequestId)
		return errors.Wrap(err, fmt.Sprintf("no such action `%s` to handle for message with id `%s`", msg.Action, msg.RequestId))
	}

	requestHandler := handleActions[msg.Action]
	if err = requestHandler(p, msg); err != nil {
		p.log.WithError(err).Errorf("failed to handle message with id `%s`", msg.RequestId)
		return err
	}

	p.log.Infof("finish handling message with id `%s`", msg.RequestId)
	return nil
}
