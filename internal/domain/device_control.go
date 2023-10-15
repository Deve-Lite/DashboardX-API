package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/google/uuid"
)

type DeviceControl struct {
	ID                     uuid.UUID         `db:"id"`
	DeviceID               uuid.UUID         `db:"device_id"`
	Name                   string            `db:"name"`
	Type                   enum.ControlType  `db:"type"`
	QoS                    enum.QoSLevel     `db:"quality_of_service"`
	IconName               string            `db:"icon_name"`
	IconBackgroundColor    string            `db:"icon_background_color"`
	IsAvailable            bool              `db:"is_available"`
	IsConfirmationRequired bool              `db:"is_confirmation_required"`
	CanNotifyOnPublish     bool              `db:"can_notify_on_publish"`
	CanDisplayName         bool              `db:"can_display_name"`
	Topic                  string            `db:"topic"`
	Attributes             ControlAttributes `db:"attributes"`
}

type CreateDeviceControl struct {
	DeviceID               uuid.UUID         `db:"device_id"`
	Name                   string            `db:"name"`
	Type                   enum.ControlType  `db:"type"`
	QoS                    enum.QoSLevel     `db:"quality_of_service"`
	IconName               string            `db:"icon_name"`
	IconBackgroundColor    string            `db:"icon_background_color"`
	IsAvailable            bool              `db:"is_available"`
	IsConfirmationRequired bool              `db:"is_confirmation_required"`
	CanNotifyOnPublish     bool              `db:"can_notify_on_publish"`
	CanDisplayName         bool              `db:"can_display_name"`
	Topic                  string            `db:"topic"`
	Attributes             ControlAttributes `db:"attributes"`
}

type UpdateDeviceControl struct {
	ID                     uuid.UUID         `db:"id"`
	DeviceID               uuid.UUID         `db:"device_id"`
	Name                   *string           `db:"name"`
	Type                   *enum.ControlType `db:"type"`
	QoS                    *enum.QoSLevel    `db:"quality_of_service"`
	IconName               *string           `db:"icon_name"`
	IconBackgroundColor    *string           `db:"icon_background_color"`
	IsAvailable            *bool             `db:"is_available"`
	IsConfirmationRequired *bool             `db:"is_confirmation_required"`
	CanNotifyOnPublish     *bool             `db:"can_notify_on_publish"`
	CanDisplayName         *bool             `db:"can_display_name"`
	Topic                  *string           `db:"topic"`
	Attributes             ControlAttributes `db:"attributes"`
}

type ExistDeviceControlFilters struct {
	DeviceID uuid.UUID        `db:"device_id"`
	Type     enum.ControlType `db:"type"`
}

type ListDeviceControlFilters struct {
	DeviceID uuid.UUID `db:"device_id"`
}

type ControlAttributes map[string]interface{}

func (a *ControlAttributes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ControlAttributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
