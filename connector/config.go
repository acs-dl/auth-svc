package connector

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AuthConfiger interface {
	AuthConnector() *Connector
}

type authConnectorConfigurator struct {
	once   comfig.Once
	getter kv.Getter
}

func NewAuthConnectorConfigurator(getter kv.Getter) AuthConfiger {
	return &authConnectorConfigurator{getter: getter}
}

type ModuleConnectorConfig struct {
	ServiceUrl string `fig:"url"`
}

func (c *authConnectorConfigurator) AuthConnector() *Connector {
	return c.once.Do(func() interface{} {
		cfg := ModuleConnectorConfig{}

		raw := kv.MustGetStringMap(c.getter, "auth_connector")

		if err := figure.
			Out(&cfg).
			From(raw).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out auth connector"))
		}

		return NewConnector(cfg.ServiceUrl)
	}).(*Connector)
}
