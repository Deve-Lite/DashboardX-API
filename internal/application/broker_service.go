package application

import (
	"context"
	"log"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
)

type BrokerService interface {
	Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error)
	Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error)
	Update(ctx context.Context, broker *domain.UpdateBroker) error
	Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error
	GetCredentials(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error)
	SetCredentials(ctx context.Context, broker *domain.UpdateBroker) error
}

type brokerService struct {
	c  *config.Config
	br repository.BrokerRepository
	cs CryptoService
	es EventService
}

func NewBrokerService(c *config.Config, br repository.BrokerRepository, cs CryptoService, es EventService) BrokerService {
	return &brokerService{c, br, cs, es}
}

func (b *brokerService) Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error) {
	return b.br.Get(ctx, brokerID, userID)
}

func (b *brokerService) List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error) {
	return b.br.List(ctx, userID)
}

func (b *brokerService) Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error) {
	if broker.Password.Set || broker.Username.Set {
		log.Print("can not set credential, use brokerService.SetCredentials instead")
		broker.Password = t.NewString("", false, false)
		broker.Username = t.NewString("", false, false)
	}

	brokerID, err := b.br.Create(ctx, broker)
	if err != nil {
		return uuid.Nil, err
	}

	// enum.EntityCreatedAction, broker.UserID, brokerID

	b.es.Publish(ctx, domain.Event{
		ID: uuid.New(),
		Data: domain.EventData{
			Action: enum.EntityCreatedAction,
			Entity: &domain.EventEntity{
				ID:   brokerID,
				Name: enum.BrokersEntity,
			},
		},
	}, broker.UserID, uuid.Nil)

	return brokerID, nil
}

func (b *brokerService) Update(ctx context.Context, broker *domain.UpdateBroker) error {
	if broker.Password.Set || broker.Username.Set {
		log.Print("can not set credential, use brokerService.SetCredentials instead")
		broker.Password = t.NewString("", false, false)
		broker.Username = t.NewString("", false, false)
	}

	if err := b.br.Update(ctx, broker); err != nil {
		return err
	}

	b.es.PublishBrokers(ctx, enum.EntityUpdatedAction, broker.UserID, broker.ID)

	return nil
}

func (b *brokerService) Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error {
	if err := b.br.Delete(ctx, brokerID, userID); err != nil {
		return err
	}

	b.es.PublishBrokers(ctx, enum.EntityDeletedAction, userID, brokerID)

	return nil
}

func (b *brokerService) GetCredentials(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error) {
	broker, err := b.br.Get(ctx, brokerID, userID)
	if err != nil {
		return nil, err
	}

	if broker.Username.Set && !broker.Username.Null {
		username, err := b.cs.Decrypt(broker.Username.String, enum.CryptoBrokerKey)
		if err != nil {
			return nil, err
		}

		broker.Username = t.NewString(username, false, true)
	}

	if broker.Password.Set && !broker.Password.Null {
		password, err := b.cs.Decrypt(broker.Password.String, enum.CryptoBrokerKey)
		if err != nil {
			return nil, err
		}

		broker.Password = t.NewString(password, false, true)
	}

	return broker, nil
}

func (b *brokerService) SetCredentials(ctx context.Context, broker *domain.UpdateBroker) error {
	if !broker.Username.Set || !broker.Password.Set {
		return ae.ErrNoBrokerCredentials
	}

	if !broker.Username.Null {
		username, err := b.cs.Encrypt(broker.Username.String, enum.CryptoBrokerKey)
		if err != nil {
			return err
		}

		broker.Username = t.NewString(username, false, true)
	}

	if !broker.Password.Null {
		password, err := b.cs.Encrypt(broker.Password.String, enum.CryptoBrokerKey)
		if err != nil {
			return err
		}

		broker.Password = t.NewString(password, false, true)
	}

	if err := b.br.Update(ctx, broker); err != nil {
		return err
	}

	b.es.PublishBrokers(ctx, enum.EntityUpdatedAction, broker.UserID, broker.ID)

	return nil
}
