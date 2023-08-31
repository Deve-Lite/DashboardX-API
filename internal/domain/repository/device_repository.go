package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	Get(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) (*domain.Device, error)
	List(ctx context.Context, filters *domain.ListDeviceFilters) ([]*domain.Device, error)
	Create(ctx context.Context, device *domain.CreateDevice) (uuid.UUID, error)
	Update(ctx context.Context, device *domain.UpdateDevice) error
	Delete(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error
}
