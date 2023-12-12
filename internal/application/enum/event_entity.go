package enum

type EventEntity string

const (
	UserEntity           EventEntity = "USER"
	BrokersEntity        EventEntity = "BROKERS"
	DevicesEntity        EventEntity = "DEVICES"
	DeviceControlsEntity EventEntity = "DEVICE_CONTROLS"
)
