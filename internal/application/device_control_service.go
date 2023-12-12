package application

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/google/uuid"
)

type DeviceControlService interface {
	List(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) ([]*domain.DeviceControl, error)
	Create(ctx context.Context, userID uuid.UUID, control *domain.CreateDeviceControl) (uuid.UUID, error)
	Update(ctx context.Context, userID uuid.UUID, control *domain.UpdateDeviceControl) error
	Delete(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID, controlID uuid.UUID) error
}

type deviceControlService struct {
	dcr repository.DeviceControlRepository
	ds  DeviceService
	es  EventService
}

func NewDeviceControlService(dcr repository.DeviceControlRepository, ds DeviceService, es EventService) DeviceControlService {
	return &deviceControlService{dcr, ds, es}
}

func (dc *deviceControlService) List(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID) ([]*domain.DeviceControl, error) {
	if _, err := dc.ds.Get(ctx, deviceID, userID); err != nil {
		return nil, err
	}

	return dc.dcr.ListByDevice(ctx, deviceID)
}

func (dc *deviceControlService) Create(ctx context.Context, userID uuid.UUID, control *domain.CreateDeviceControl) (uuid.UUID, error) {
	device, err := dc.ds.Get(ctx, control.DeviceID, userID)
	if err != nil {
		return uuid.Nil, err
	}

	if control.Type == enum.ControlState {
		r, err := dc.dcr.Exist(ctx, &domain.DeviceControlFilters{DeviceID: control.DeviceID, Type: control.Type})
		if err != nil {
			return uuid.Nil, err
		}
		if r {
			return uuid.Nil, ae.ErrControlStateExists
		}
	}

	controlID, err := dc.dcr.Create(ctx, control)
	if err != nil {
		return uuid.Nil, err
	}

	dc.es.PublishDeviceControls(ctx, enum.EntityCreatedAction, userID, device.BrokerID.UUID, device.ID, controlID)

	return controlID, nil
}

func (dc *deviceControlService) Update(ctx context.Context, userID uuid.UUID, control *domain.UpdateDeviceControl) error {
	device, err := dc.ds.Get(ctx, control.DeviceID, userID)
	if err != nil {
		return err
	}

	if control.Type != nil && *control.Type == enum.ControlState {
		ci, err := dc.dcr.ListByType(ctx, &domain.DeviceControlFilters{
			DeviceID: control.DeviceID,
			Type:     *control.Type,
		})
		if err != nil {
			return err
		}

		if len(ci) == 1 && ci[0].ID != control.ID {
			return ae.ErrControlStateExists
		}
	}

	err = dc.dcr.Update(ctx, control)
	if err != nil {
		return err
	}

	dc.es.PublishDeviceControls(ctx, enum.EntityCreatedAction, userID, device.BrokerID.UUID, device.ID, control.ID)

	return nil
}

func (dc *deviceControlService) Delete(ctx context.Context, userID uuid.UUID, deviceID uuid.UUID, controlID uuid.UUID) error {
	device, err := dc.ds.Get(ctx, deviceID, userID)
	if err != nil {
		return err
	}

	err = dc.dcr.Delete(ctx, deviceID, controlID)
	if err != nil {
		return err
	}

	dc.es.PublishDeviceControls(ctx, enum.EntityCreatedAction, userID, device.BrokerID.UUID, deviceID, controlID)

	return nil
}
