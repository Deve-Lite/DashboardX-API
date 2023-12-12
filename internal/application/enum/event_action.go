package enum

type EventAction string

const (
	ChannelOpenedAction EventAction = "CHANNEL_OPENED"
	ChannelClosedAction EventAction = "CHANNEL_CLOSED"
	EntityCreatedAction EventAction = "ENTITY_CREATED"
	EntityUpdatedAction EventAction = "ENTITY_UPDATED"
	EntityDeletedAction EventAction = "ENTITY_DELETED"
)
