package mapper

import (
	"reflect"
	"unicode"

	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
)

type DeviceControlMapper interface {
	ModelToDTO(v *domain.DeviceControl) *dto.GetDeviceControlResponse
	CreateDTOToCreateModel(v *dto.CreateDeviceControlRequest) *domain.CreateDeviceControl
	UpdateDTOToUpdateModel(v *dto.UpdateDeviceControlRequest) *domain.UpdateDeviceControl
}

type deviceControlMapper struct{}

func NewDeviceControlMapper() DeviceControlMapper {
	return &deviceControlMapper{}
}

func (*deviceControlMapper) ModelToDTO(v *domain.DeviceControl) *dto.GetDeviceControlResponse {
	r := &dto.GetDeviceControlResponse{
		ID:       v.ID,
		DeviceID: v.DeviceID,
		Name:     v.Name,
		Icon: dto.Icon{
			Name:            v.IconName,
			BackgroundColor: v.IconBackgroundColor,
		},
		Type:                   v.Type,
		Topic:                  v.Topic,
		QoS:                    v.QoS,
		IsConfirmationRequired: v.IsConfirmationRequired,
		IsAvailable:            v.IsAvailable,
		CanNotifyOnPublish:     v.CanNotifyOnPublish,
		CanDisplayName:         v.CanDisplayName,
	}

	for k, e := range v.Attributes {
		if k == "colorFormat" {
			t := e.(string)
			r.Attributes.ColorFormat = &t
		}

		if k == "maxValue" {
			t := float32(e.(float64))
			r.Attributes.MaxValue = &t
		}

		if k == "minValue" {
			t := float32(e.(float64))
			r.Attributes.MinValue = &t
		}

		if k == "offPayload" {
			t := e.(string)
			r.Attributes.OffPayload = &t
		}

		if k == "onPayload" {
			t := e.(string)
			r.Attributes.OnPayload = &t
		}

		if k == "payload" {
			t := e.(string)
			r.Attributes.Payload = &t
		}

		if k == "payloadTemplate" {
			t := e.(string)
			r.Attributes.PayloadTemplate = &t
		}

		if k == "payloads" {
			t := e.(map[string]interface{})
			a := make(map[string]string)

			for tk, te := range t {
				a[tk] = te.(string)
			}

			r.Attributes.Payloads = &a
		}

		if k == "sendAsTicks" {
			t := e.(bool)
			r.Attributes.SendAsTicks = &t
		}

		if k == "value" {
			t := e.(string)
			r.Attributes.Value = &t
		}
	}

	return r
}

func (*deviceControlMapper) CreateDTOToCreateModel(v *dto.CreateDeviceControlRequest) *domain.CreateDeviceControl {
	d := &domain.CreateDeviceControl{
		Type:                   *v.Type,
		Name:                   v.Name,
		IconName:               v.Icon.Name,
		IconBackgroundColor:    v.Icon.BackgroundColor,
		QoS:                    *v.QoS,
		IsAvailable:            *v.IsAvailable,
		IsConfirmationRequired: *v.IsConfirmationRequired,
		CanDisplayName:         *v.CanDisplayName,
		CanNotifyOnPublish:     *v.CanNotifyOnPublish,
		Topic:                  v.Topic,
	}

	a := map[string]interface{}{}

	s := reflect.ValueOf(*v.Attributes)
	g := s.Type()

	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).IsNil() {
			continue
		}

		n := g.Field(i).Name
		k := string(unicode.ToLower([]rune(n)[0])) + n[1:]

		a[k] = s.Field(i).Interface()
	}

	d.Attributes = a

	return d
}

func (*deviceControlMapper) UpdateDTOToUpdateModel(v *dto.UpdateDeviceControlRequest) *domain.UpdateDeviceControl {
	d := &domain.UpdateDeviceControl{
		Type:                   v.Type,
		Name:                   v.Name,
		QoS:                    v.QoS,
		IsAvailable:            v.IsAvailable,
		IsConfirmationRequired: v.IsConfirmationRequired,
		CanDisplayName:         v.CanDisplayName,
		CanNotifyOnPublish:     v.CanNotifyOnPublish,
		Topic:                  v.Topic,
	}

	if v.Icon.Name.Set {
		d.IconName = &v.Icon.Name.String
	}

	if v.Icon.BackgroundColor.Set {
		d.IconBackgroundColor = &v.Icon.BackgroundColor.String
	}

	if v.Attributes != nil {
		a := map[string]interface{}{}

		s := reflect.ValueOf(*v.Attributes)
		g := s.Type()

		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).IsNil() {
				continue
			}

			n := g.Field(i).Name
			k := string(unicode.ToLower([]rune(n)[0])) + n[1:]

			a[k] = s.Field(i).Interface()
		}

		d.Attributes = a
	} else {
		if *v.Type == enum.ControlTextOut {
			d.Attributes = map[string]interface{}{}
		}
	}

	return d
}
