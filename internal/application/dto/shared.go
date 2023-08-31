package dto

import t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"

type Icon struct {
	Name            string `json:"name" binding:"required"`
	BackgroundColor string `json:"backgroundColor" binding:"required,hexcolor"`
}

type IconOptional struct {
	Name            t.String `json:"name" swaggertype:"string"`
	BackgroundColor t.String `json:"backgroundColor" binding:"nullhexcolor" swaggertype:"string"`
}
