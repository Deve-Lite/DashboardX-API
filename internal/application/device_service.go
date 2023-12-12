package application

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
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
	es EventService
}

func NewDeviceService(dr repository.DeviceRepository, bs BrokerService, es EventService) DeviceService {
	return &deviceService{dr, bs, es}
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

	deviceID, err := d.dr.Create(ctx, device)
	if err != nil {
		return uuid.Nil, err
	}

	d.es.PublishDevices(ctx, enum.EntityCreatedAction, device.UserID, device.BrokerID.UUID, deviceID)

	return deviceID, nil
}

func (d *deviceService) Update(ctx context.Context, device *domain.UpdateDevice) error {
	if device.BrokerID.Set && !device.BrokerID.Null {
		if _, err := d.bs.Get(ctx, device.BrokerID.Value, device.UserID); err != nil {
			return err
		}
	}

	if err := d.dr.Update(ctx, device); err != nil {
		return err
	}

	{
		device, _ := d.Get(ctx, device.ID, device.UserID)
		d.es.PublishDevices(ctx, enum.EntityUpdatedAction, device.UserID, device.BrokerID.UUID, device.ID)
	}

	return nil
}

func (d *deviceService) Delete(ctx context.Context, deviceID uuid.UUID, userID uuid.UUID) error {
	device, _ := d.Get(ctx, deviceID, userID)

	if err := d.dr.Delete(ctx, deviceID, userID); err != nil {
		return err
	}

	d.es.PublishDevices(ctx, enum.EntityUpdatedAction, device.UserID, device.BrokerID.UUID, device.ID)

	return nil
}
