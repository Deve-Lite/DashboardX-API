package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type DeviceControlRepository interface {
	ListByType(ctx context.Context, filters *domain.DeviceControlFilters) ([]*domain.DeviceControl, error)
	ListByDevice(ctx context.Context, deviceID uuid.UUID) ([]*domain.DeviceControl, error)
	Create(ctx context.Context, control *domain.CreateDeviceControl) (uuid.UUID, error)
	Exist(ctx context.Context, filters *domain.DeviceControlFilters) (bool, error)
	Update(ctx context.Context, control *domain.UpdateDeviceControl) error
	Delete(ctx context.Context, deviceID uuid.UUID, controlID uuid.UUID) error
}
