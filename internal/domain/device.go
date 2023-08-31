package domain

import (
	"time"

	t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
	"github.com/google/uuid"
)

type Device struct {
	ID                  uuid.UUID     `db:"id"`
	UserID              uuid.UUID     `db:"user_id"`
	BrokerID            uuid.NullUUID `db:"broker_id"`
	Name                string        `db:"name"`
	IconName            string        `db:"icon_name"`
	IconBackgroundColor string        `db:"icon_background_color"`
	Placing             *string       `db:"placing"`
	BasePath            *string       `db:"base_path"`
	CreatedAt           time.Time     `db:"created_at"`
	UpdatedAt           time.Time     `db:"updated_at"`
}

type CreateDevice struct {
	UserID              uuid.UUID     `db:"user_id"`
	BrokerID            uuid.NullUUID `db:"broker_id"`
	Name                string        `db:"name"`
	IconName            string        `db:"icon_name"`
	IconBackgroundColor string        `db:"icon_background_color"`
	Placing             t.String      `db:"placing"`
	BasePath            t.String      `db:"base_path"`
}

type UpdateDevice struct {
	ID                  uuid.UUID             `db:"id"`
	UserID              uuid.UUID             `db:"user_id"`
	BrokerID            t.Nullable[uuid.UUID] `db:"broker_id"`
	Name                t.String              `db:"name"`
	IconName            t.String              `db:"icon_name"`
	IconBackgroundColor t.String              `db:"icon_background_color"`
	Placing             t.String              `db:"placing"`
	BasePath            t.String              `db:"base_path"`
}

type ListDeviceFilters struct {
	UserID   uuid.UUID
	BrokerID uuid.NullUUID
}
