package mapper

import (
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
)

type BrokerMapper interface {
	ModelToDTO(v *domain.Broker) *dto.GetBrokerResponse
	ModelToCredentialsDTO(v *domain.Broker) *dto.GetBrokerCredentialsResponse
	CreateDTOToCreateModel(v *dto.CreateBrokerRequest) *domain.CreateBroker
	UpdateDTOToUpdateModel(v *dto.UpdateBrokerRequest) *domain.UpdateBroker
	SetCredentialsDTOToUpdateModel(v *dto.SetBrokerCredentialsRequest) *domain.UpdateBroker
}

type brokerMapper struct{}

func NewBrokerMapper() BrokerMapper {
	return &brokerMapper{}
}

func (*brokerMapper) ModelToDTO(v *domain.Broker) *dto.GetBrokerResponse {
	r := &dto.GetBrokerResponse{
		ID:        v.ID,
		Name:      v.Name,
		Server:    v.Server,
		Port:      v.Port,
		KeepAlive: v.KeepAlive,
		Icon: dto.Icon{
			Name:            v.IconName,
			BackgroundColor: v.IconBackgroundColor,
		},
		IsSSL:     v.IsSSL,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}

	if v.ClientID.Null {
		r.ClientID = nil
	} else {
		r.ClientID = &v.ClientID.String
	}

	return r
}

func (*brokerMapper) ModelToCredentialsDTO(v *domain.Broker) *dto.GetBrokerCredentialsResponse {
	r := &dto.GetBrokerCredentialsResponse{
		ID: v.ID,
	}

	if v.Password.Null {
		r.Password = nil
	} else {
		r.Password = &v.Password.String
	}

	if v.Username.Null {
		r.Username = nil
	} else {
		r.Username = &v.Username.String
	}

	return r
}

func (*brokerMapper) CreateDTOToCreateModel(v *dto.CreateBrokerRequest) *domain.CreateBroker {
	return &domain.CreateBroker{
		Name:                v.Name,
		Server:              v.Server,
		Port:                *v.Port,
		KeepAlive:           *v.KeepAlive,
		IconName:            v.Icon.Name,
		IconBackgroundColor: v.Icon.BackgroundColor,
		IsSSL:               *v.IsSSL,
		ClientID:            v.ClientID,
	}
}

func (*brokerMapper) UpdateDTOToUpdateModel(v *dto.UpdateBrokerRequest) *domain.UpdateBroker {
	r := &domain.UpdateBroker{
		Name:      v.Name,
		Server:    v.Server,
		Port:      v.Port,
		KeepAlive: v.KeepAlive,
		IsSSL:     v.IsSSL,
		ClientID:  v.ClientID,
	}

	if v.Icon.Name.Set {
		r.IconName = v.Icon.Name
	}

	if v.Icon.BackgroundColor.Set {
		r.IconBackgroundColor = v.Icon.BackgroundColor
	}

	return r
}

func (*brokerMapper) SetCredentialsDTOToUpdateModel(v *dto.SetBrokerCredentialsRequest) *domain.UpdateBroker {
	return &domain.UpdateBroker{
		Username: v.Username,
		Password: v.Password,
	}
}
