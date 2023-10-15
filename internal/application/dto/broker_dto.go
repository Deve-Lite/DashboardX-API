package dto

import (
	"time"

	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
)

type BrokerParams struct {
	BrokerID string `uri:"brokerId" binding:"required,uuid"`
}

type CreateBrokerRequest struct {
	Name      string   `json:"name" binding:"required"`
	Server    string   `json:"server" binding:"required"`
	Port      *uint16  `json:"port" binding:"required"`
	KeepAlive *uint16  `json:"keepAlive" binding:"required"`
	Icon      Icon     `json:"icon" binding:"required"`
	IsSSL     *bool    `json:"isSsl" binding:"required"`
	ClientID  t.String `json:"clientId" swaggertype:"string" extensions:"x-nullable"`
}

type CreateBrokerResponse struct {
	ID uuid.UUID `json:"id" binding:"required,uuid"`
}

type UpdateBrokerRequest struct {
	Name      t.String     `json:"name" swaggertype:"string"`
	Server    t.String     `json:"server" swaggertype:"string"`
	Port      t.Uint16     `json:"port" swaggertype:"integer"`
	KeepAlive t.Uint16     `json:"keepAlive" swaggertype:"integer"`
	Icon      IconOptional `json:"icon"`
	IsSSL     t.Bool       `json:"isSsl" swaggertype:"boolean"`
	ClientID  t.String     `json:"clientId" swaggertype:"string" extensions:"x-nullable"`
}

type GetBrokerResponse struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	Name      string    `json:"name"`
	Server    string    `json:"server"`
	Port      uint16    `json:"port"`
	KeepAlive uint16    `json:"keepAlive"`
	Icon      Icon      `json:"icon"`
	IsSSL     bool      `json:"isSsl"`
	ClientID  *string   `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetBrokerCredentialsResponse struct {
	ID       uuid.UUID `json:"id" format:"uuid"`
	Username *string   `json:"username"`
	Password *string   `json:"password"`
}

type SetBrokerCredentialsRequest struct {
	Username t.String `json:"username" binding:"requirednullstring" swaggertype:"string" extensions:"x-nullable"`
	Password t.String `json:"password" binding:"requirednullstring" swaggertype:"string" extensions:"x-nullable"`
}
