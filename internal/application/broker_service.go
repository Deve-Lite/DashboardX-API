package application

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain/repository"
	"github.com/google/uuid"
)

type BrokerService interface {
	Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error)
	Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error)
	Update(ctx context.Context, broker *domain.UpdateBroker) error
	Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error
}

type brokerService struct {
	br repository.BrokerRepository
}

func NewBrokerService(br repository.BrokerRepository) BrokerService {
	return &brokerService{br}
}

func (b *brokerService) Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error) {
	return b.br.Get(ctx, brokerID, userID)
}

func (b *brokerService) List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error) {
	return b.br.List(ctx, userID)
}

func (b *brokerService) Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error) {
	return b.br.Create(ctx, broker)
}

func (b *brokerService) Update(ctx context.Context, broker *domain.UpdateBroker) error {
	return b.br.Update(ctx, broker)
}

func (b *brokerService) Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error {
	return b.br.Delete(ctx, brokerID, userID)
}
