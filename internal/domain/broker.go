package domain

import (
	"time"

	t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
	"github.com/google/uuid"
)

type Broker struct {
	ID                  uuid.UUID `db:"id"`
	UserID              uuid.UUID `db:"user_id"`
	Name                string    `db:"name"`
	Server              string    `db:"server"`
	Port                uint16    `db:"port"`
	KeepAlive           uint16    `db:"keep_alive"`
	IconName            string    `db:"icon_name"`
	IconBackgroundColor string    `db:"icon_background_color"`
	IsSSL               bool      `db:"is_ssl"`
	Username            t.String  `db:"username"`
	Password            t.String  `db:"password"`
	ClientID            t.String  `db:"client_id"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

type CreateBroker struct {
	UserID              uuid.UUID `db:"user_id"`
	Name                string    `db:"name"`
	Server              string    `db:"server"`
	Port                uint16    `db:"port"`
	KeepAlive           uint16    `db:"keep_alive"`
	IconName            string    `db:"icon_name"`
	IconBackgroundColor string    `db:"icon_background_color"`
	IsSSL               bool      `db:"is_ssl"`
	Username            t.String  `db:"username"`
	Password            t.String  `db:"password"`
	ClientID            t.String  `db:"client_id"`
}

type UpdateBroker struct {
	ID                  uuid.UUID `db:"user_id"`
	UserID              uuid.UUID `db:"user_id"`
	Name                t.String  `db:"name"`
	Server              t.String  `db:"server"`
	Port                t.Uint16  `db:"port"`
	KeepAlive           t.Uint16  `db:"keep_alive"`
	IconName            t.String  `db:"icon_name"`
	IconBackgroundColor t.String  `db:"icon_background_color"`
	IsSSL               t.Bool    `db:"is_ssl"`
	Username            t.String  `db:"username"`
	Password            t.String  `db:"password"`
	ClientID            t.String  `db:"client_id"`
}
