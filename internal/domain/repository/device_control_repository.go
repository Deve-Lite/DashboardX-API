package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/google/uuid"
)

type DeviceControlRepository interface {
	List(ctx context.Context, filters *domain.ListDeviceControlFilters) ([]*domain.DeviceControl, error)
	Create(ctx context.Context, control *domain.CreateDeviceControl) (uuid.UUID, error)
	Exist(ctx context.Context, filters *domain.ExistDeviceControlFilters) (bool, error)
	Update(ctx context.Context, control *domain.UpdateDeviceControl) error
	Delete(ctx context.Context, deviceID uuid.UUID, controlID uuid.UUID) error
}
