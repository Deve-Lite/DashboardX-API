package validate

import (
	"fmt"
	"reflect"

	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var EmptyMin validator.Func = func(fl validator.FieldLevel) bool {
	v1, ok1 := fl.Field().Interface().(t.String)
	if ok1 {
		if err := validate.Var(v1.String, fmt.Sprintf("omitempty,min=%s", fl.Param())); err != nil {
			return false
		}

		return true
	}

	v2, ok2 := fl.Field().Interface().(t.Int)
	if ok2 {
		if err := validate.Var(v2.Int, fmt.Sprintf("omitempty,min=%s", fl.Param())); err != nil {
			return false
		}

		return true
	}

	return false
}

var EmptyEmail validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(t.String)
	if !ok {
		return false
	}

	if err := validate.Var(v.String, "omitempty,email"); err != nil {
		return false
	}

	return true
}

var EmptyUUID validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(uuid.NullUUID)
	if !ok {
		return false
	}

	if err := validate.Var(v, "omitempty,uuid"); err != nil {
		return false
	}

	return true
}

var EmptyHexColor validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(t.String)
	if !ok {
		return false
	}

	if err := validate.Var(v.String, "omitempty,hexcolor"); err != nil {
		return false
	}

	return true
}

var RequiredNullString validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(t.String)
	if !ok {
		return false
	}

	if !v.Set {
		return false
	}

	if v.Null {
		return true
	}

	if err := validate.Var(v.String, "required"); err != nil {
		return false
	}

	return true
}

var ControlType validator.Func = func(fl validator.FieldLevel) bool {
	_, ok := fl.Parent().FieldByName("Type").Interface().(*enum.ControlType)
	if !ok {
		return false
	}

	v, ok := fl.Field().Interface().(enum.ControlType)
	if !ok {
		return false
	}

	if fl.Parent().FieldByName("Attributes").IsNil() && v != enum.ControlTextOut {
		return false
	}

	switch v {
	case enum.ControlButton:
		fallthrough
	case enum.ControlColor:
		fallthrough
	case enum.ControlDateTime:
		fallthrough
	case enum.ControlRadio:
		fallthrough
	case enum.ControlSlider:
		fallthrough
	case enum.ControlState:
		fallthrough
	case enum.ControlSwitch:
		fallthrough
	case enum.ControlTextOut:
		return true
	}

	return false
}

func onlyRequiredAttributes(d dto.ControlAttributes, v ...string) bool {
	e := []string{}

	s := reflect.ValueOf(d)
	g := s.Type()

	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).IsNil() {
			continue
		}

		e = append(e, g.Field(i).Name)
	}

	if len(v) == 0 {
		return len(v) == len(e)
	}

	less := func(a, b string) bool { return a < b }
	return cmp.Diff(v, e, cmpopts.SortSlices(less)) == ""
}

var ControlAttributes validator.Func = func(fl validator.FieldLevel) bool {
	if fl.Parent().FieldByName("Type").IsNil() {
		return false
	}

	t, ok := fl.Parent().FieldByName("Type").Interface().(*enum.ControlType)
	if !ok {
		return false
	}

	v, ok := fl.Field().Interface().(dto.ControlAttributes)
	if !ok {
		return false
	}

	switch *t {
	case enum.ControlButton:
		return onlyRequiredAttributes(v, "Payload")
	case enum.ControlColor:
		return onlyRequiredAttributes(v, "PayloadTemplate", "ColorFormat")
	case enum.ControlDateTime:
		return onlyRequiredAttributes(v, "PayloadTemplate", "SendAsTicks")
	case enum.ControlRadio:
		return onlyRequiredAttributes(v, "Payloads")
	case enum.ControlSlider:
		return onlyRequiredAttributes(v, "PayloadTemplate", "MinValue", "MaxValue")
	case enum.ControlState:
		return onlyRequiredAttributes(v, "SecondSpan")
	case enum.ControlSwitch:
		return onlyRequiredAttributes(v, "OnPayload", "OffPayload")
	case enum.ControlTextOut:
		return onlyRequiredAttributes(v)
	}

	return false
}

var QoSLevel validator.Func = func(fl validator.FieldLevel) bool {
	if v, ok := fl.Field().Interface().(enum.QoSLevel); ok {
		return v >= enum.QoSZero && v <= enum.QoSTwo
	}

	return false
}
