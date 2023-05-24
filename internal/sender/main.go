package sender

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/auth-svc/internal/config"
	"github.com/acs-dl/auth-svc/internal/data"
	"github.com/acs-dl/auth-svc/internal/processor"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const serviceName = data.ModuleName + "-sender"

type Sender struct {
	publisher *amqp.Publisher
	log       *logan.Entry
	topics    map[string]string
}

func NewSender(cfg config.Config) *Sender {
	return &Sender{
		publisher: cfg.Amqp().Publisher,
		log:       logan.New().WithField("service", serviceName),
		topics: map[string]string{
			data.ModuleName:       cfg.Amqp().Topic,
			data.OrchestratorName: cfg.Amqp().Orchestrator,
		},
	}
}

func (s *Sender) Run(ctx context.Context) {
	go running.WithBackOff(ctx, s.log,
		serviceName,
		s.processMessages,
		30*time.Second,
		30*time.Second,
		30*time.Second,
	)
}

func (s *Sender) processMessages(ctx context.Context) error {
	s.log.Infof("start sending message to get module permissions")
	msg, err := s.createMessage()
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}

	err = s.publisher.Publish(s.topics[data.OrchestratorName], msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message: "+msg.UUID)
	}

	s.log.Infof("finish sending message to get module permissions")

	return nil
}

func (s *Sender) createMessage() (*message.Message, error) {
	payload := data.ModulePayload{
		RequestId: uuid.New().String(),
		Action:    processor.GetModulesPermissionsAction,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		s.log.WithError(err).Errorf("failed to marshal payload")
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	return &message.Message{
		UUID:     payload.RequestId,
		Metadata: nil,
		Payload:  payloadJson,
	}, nil
}
