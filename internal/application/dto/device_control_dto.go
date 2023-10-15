package dto

import (
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/google/uuid"
)

type ControlAttributes struct {
	MaxValue        *float32           `json:"maxValue,omitempty"`
	MinValue        *float32           `json:"minValue,omitempty"`
	Value           *string            `json:"value,omitempty"`
	ColorFormat     *string            `json:"colorFormat,omitempty"`
	PayloadTemplate *string            `json:"payloadTemplate,omitempty"`
	Payload         *string            `json:"payload,omitempty"`
	Payloads        *map[string]string `json:"payloads,omitempty"`
	OnPayload       *string            `json:"onPayload,omitempty"`
	OffPayload      *string            `json:"offPayload,omitempty"`
	SecondSpan      *int               `json:"secondSpan,omitempty"`
	SendAsTicks     *bool              `json:"sendAsTicks,omitempty"`
}

type CreateDeviceControlRequest struct {
	Name                   string             `json:"name" binding:"required"`
	Type                   *enum.ControlType  `json:"type" binding:"required,control_type"`
	Attributes             *ControlAttributes `json:"attributes" binding:"control_attributes"`
	Topic                  string             `json:"topic" binding:"required"`
	Icon                   Icon               `json:"icon" binding:"required"`
	QoS                    *enum.QoSLevel     `json:"qualityOfService" binding:"qos_level"`
	IsConfirmationRequired *bool              `json:"isConfirmationRequired" binding:"required"`
	IsAvailable            *bool              `json:"isAvailable" binding:"required"`
	CanNotifyOnPublish     *bool              `json:"canNotifyOnPublish" binding:"required"`
	CanDisplayName         *bool              `json:"canDisplayName" binding:"required"`
}

type CreateDeviceControlResponse struct {
	ID uuid.UUID `json:"id" format:"uuid"`
}

type GetDeviceControlResponse struct {
	ID                     uuid.UUID         `json:"id" format:"uuid"`
	DeviceID               uuid.UUID         `json:"deviceId" format:"uuid"`
	Name                   string            `json:"name"`
	Type                   enum.ControlType  `json:"type"`
	Attributes             ControlAttributes `json:"attributes"`
	Topic                  string            `json:"topic"`
	Icon                   Icon              `json:"icon"`
	QoS                    enum.QoSLevel     `json:"qualityOfService"`
	IsConfirmationRequired bool              `json:"isConfirmationRequired"`
	IsAvailable            bool              `json:"isAvailable"`
	CanNotifyOnPublish     bool              `json:"canNotifyOnPublish"`
	CanDisplayName         bool              `json:"canDisplayName"`
}

type UpdateDeviceControlRequest struct {
	Name                   *string            `json:"name"`
	Type                   *enum.ControlType  `json:"type" binding:"omitempty,control_type"`
	Attributes             *ControlAttributes `json:"attributes" binding:"omitempty,control_attributes"`
	Topic                  *string            `json:"topic"`
	Icon                   IconOptional       `json:"icon"`
	QoS                    *enum.QoSLevel     `json:"qualityOfService" binding:"omitempty,qos_level"`
	IsConfirmationRequired *bool              `json:"isConfirmationRequired"`
	IsAvailable            *bool              `json:"isAvailable"`
	CanNotifyOnPublish     *bool              `json:"canNotifyOnPublish"`
	CanDisplayName         *bool              `json:"canDisplayName"`
}
