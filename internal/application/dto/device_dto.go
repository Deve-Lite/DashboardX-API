package dto

import (
	"time"

	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
)

type DeviceParams struct {
	DeviceID string `uri:"deviceId" binding:"required,uuid"`
}

type DeviceQuery struct {
	BrokerID *string `form:"brokerId" format:"uuid"`
}

type DeviceControlParams struct {
	DeviceID  string `uri:"deviceId" binding:"required,uuid"`
	ControlID string `uri:"controlId" binding:"required,uuid"`
}

type CreateDeviceRequest struct {
	BrokerID uuid.NullUUID `json:"brokerId" binding:"emptyuuid" swaggertype:"string" format:"uuid"`
	Name     string        `json:"name" binding:"required"`
	Icon     Icon          `json:"icon" binding:"required"`
	Placing  t.String      `json:"placing" swaggertype:"string"`
	BasePath t.String      `json:"basePath" swaggertype:"string"`
}

type CreateDeviceResponse struct {
	ID uuid.UUID `json:"id" format:"uuid"`
}

type UpdateDeviceRequest struct {
	BrokerID t.Nullable[uuid.UUID] `json:"brokerId" swaggertype:"string" format:"uuid" extensions:"x-nullable"`
	Name     t.String              `json:"name" swaggertype:"string"`
	Icon     IconOptional          `json:"icon"`
	Placing  t.String              `json:"placing" swaggertype:"string" extensions:"x-nullable"`
	BasePath t.String              `json:"basePath" swaggertype:"string" extensions:"x-nullable"`
}

type GetDeviceResponse struct {
	ID        uuid.UUID     `json:"id" format:"uuid"`
	BrokerID  uuid.NullUUID `json:"brokerId" swaggertype:"string" format:"uuid"`
	Name      string        `json:"name"`
	Icon      Icon          `json:"icon"`
	Placing   *string       `json:"placing"`
	BasePath  *string       `json:"basePath"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
