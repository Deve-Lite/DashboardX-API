package enum

type ControlType string

const (
	ControlButton   ControlType = "button"
	ControlColor    ControlType = "color"
	ControlDateTime ControlType = "date-time"
	ControlSwitch   ControlType = "switch"
	ControlSlider   ControlType = "slider"
	ControlState    ControlType = "state"
	ControlRadio    ControlType = "radio"
	ControlTextOut  ControlType = "text-out"
)
