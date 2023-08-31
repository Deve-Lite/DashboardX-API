package mapper

import (
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/auth"
	n "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
)

type UserMapper interface {
	ModelToDTO(v *domain.User) *dto.GetUserResponse
	LoginDTOToModel(v *dto.LoginUserRequest) *domain.User
	CreateDTOToCreateModel(v *dto.CreateUserRequest) *domain.CreateUser
	UpdateDTOToUpdateModel(v *dto.UpdateUserRequest) *domain.UpdateUser
	TokenModelToTokenDTO(v *auth.Tokens) *dto.LoginUserResponse
}

type userMapper struct{}

func NewUserMapper() UserMapper {
	return &userMapper{}
}

func (*userMapper) ModelToDTO(v *domain.User) *dto.GetUserResponse {
	return &dto.GetUserResponse{
		ID:       v.ID,
		Name:     v.Name,
		Email:    v.Email,
		Theme:    v.Theme,
		Language: v.Language,
	}
}

func (*userMapper) LoginDTOToModel(v *dto.LoginUserRequest) *domain.User {
	return &domain.User{
		Password: v.Password,
		Email:    v.Email,
	}
}

func (*userMapper) CreateDTOToCreateModel(v *dto.CreateUserRequest) *domain.CreateUser {
	return &domain.CreateUser{
		Name:     v.Name,
		Password: v.Password,
		Email:    v.Email,
		IsAdmin:  false,
		Language: nil,
		Theme:    nil,
	}
}

func (*userMapper) UpdateDTOToUpdateModel(v *dto.UpdateUserRequest) *domain.UpdateUser {
	return &domain.UpdateUser{
		Name:     v.Name,
		Email:    v.Email,
		Password: n.NewString("", false, false),
		IsAdmin:  n.NewBool(false, false, false),
		Language: v.Language,
		Theme:    v.Theme,
	}
}

func (*userMapper) TokenModelToTokenDTO(v *auth.Tokens) *dto.LoginUserResponse {
	return &dto.LoginUserResponse{
		AccessToken:  v.AccessToken,
		RefreshToken: v.RefreshToken,
	}
}
