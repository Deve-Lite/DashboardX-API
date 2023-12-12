package domain

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type EventChannels struct {
	Channels map[string]chan Event
	Mutex    sync.Mutex
}

type Event struct {
	ID   uuid.UUID `json:"id"`
	Data EventData `json:"data"`
}

func (e *Event) String() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Panicf("%s", errors.Wrap(err, "domain.Event.String.Marshal"))
	}
	return string(bytes)
}

type EventData struct {
	Action  enum.EventAction `json:"action"`
	Entity  *EventEntity     `json:"entity"`
	Related *[]EventEntity   `json:"related"`
}

type EventEntity struct {
	ID   uuid.UUID        `json:"id"`
	Name enum.EventEntity `json:"name"`
}
