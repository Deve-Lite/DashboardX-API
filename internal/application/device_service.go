package application

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	"github.com/google/uuid"
)

type DeviceService interface {
	Get(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) (*domain.Device, error)
	List(ctx context.Context, filters *domain.ListDeviceFilters) ([]*domain.Device, error)
	Create(ctx context.Context, device *domain.CreateDevice) (uuid.UUID, error)
	Update(ctx context.Context, device *domain.UpdateDevice) error
	Delete(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error
}

type deviceService struct {
	dr repository.DeviceRepository
	bs BrokerService
}

func NewDeviceService(dr repository.DeviceRepository, bs BrokerService) DeviceService {
	return &deviceService{dr, bs}
}

func (d *deviceService) Get(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) (*domain.Device, error) {
	return d.dr.Get(ctx, deviceID, userID)
}

func (d *deviceService) List(ctx context.Context, filters *domain.ListDeviceFilters) ([]*domain.Device, error) {
	return d.dr.List(ctx, filters)
}

func (d *deviceService) Create(ctx context.Context, device *domain.CreateDevice) (uuid.UUID, error) {
	if device.BrokerID.Valid {
		if _, err := d.bs.Get(ctx, device.BrokerID.UUID, device.UserID); err != nil {
			return uuid.Nil, err
		}
	}

	return d.dr.Create(ctx, device)
}

func (d *deviceService) Update(ctx context.Context, device *domain.UpdateDevice) error {
	if device.BrokerID.Set && !device.BrokerID.Null {
		if _, err := d.bs.Get(ctx, device.BrokerID.Value, device.UserID); err != nil {
			return err
		}
	}

	return d.dr.Update(ctx, device)
}

func (d *deviceService) Delete(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error {
	return d.dr.Delete(ctx, deviceID, userID)
}
