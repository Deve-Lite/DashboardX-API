package mapper

import (
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
)

type DeviceMapper interface {
	ModelToDTO(v *domain.Device) *dto.GetDeviceResponse
	CreateDTOToCreateModel(v *dto.CreateDeviceRequest) *domain.CreateDevice
	UpdateDTOToUpdateModel(v *dto.UpdateDeviceRequest) *domain.UpdateDevice
}

type deviceMapper struct{}

func NewDeviceMapper() DeviceMapper {
	return &deviceMapper{}
}

func (*deviceMapper) ModelToDTO(v *domain.Device) *dto.GetDeviceResponse {
	r := &dto.GetDeviceResponse{
		ID:       v.ID,
		BrokerID: v.BrokerID,
		Name:     v.Name,
		Icon: dto.Icon{
			Name:            v.IconName,
			BackgroundColor: v.IconBackgroundColor,
		},
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		Placing:   v.Placing,
		BasePath:  v.BasePath,
	}

	return r
}

func (*deviceMapper) CreateDTOToCreateModel(v *dto.CreateDeviceRequest) *domain.CreateDevice {
	return &domain.CreateDevice{
		BrokerID:            v.BrokerID,
		Name:                v.Name,
		IconName:            v.Icon.Name,
		IconBackgroundColor: v.Icon.BackgroundColor,
		Placing:             v.Placing,
		BasePath:            v.BasePath,
	}
}

func (*deviceMapper) UpdateDTOToUpdateModel(v *dto.UpdateDeviceRequest) *domain.UpdateDevice {
	return &domain.UpdateDevice{
		BrokerID:            v.BrokerID,
		Name:                v.Name,
		IconName:            v.Icon.Name,
		IconBackgroundColor: v.Icon.BackgroundColor,
		Placing:             v.Placing,
		BasePath:            v.BasePath,
	}
}
