package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/google/uuid"
)

type BrokerRepository interface {
	Get(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) (*domain.Broker, error)
	List(ctx context.Context, userID uuid.UUID) ([]*domain.Broker, error)
	Create(ctx context.Context, broker *domain.CreateBroker) (uuid.UUID, error)
	Update(ctx context.Context, broker *domain.UpdateBroker) error
	Delete(ctx context.Context, brokerID uuid.UUID, userID uuid.UUID) error
}
