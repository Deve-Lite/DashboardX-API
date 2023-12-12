package application

import (
	"context"
	"log"
	"sync"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type EventService interface {
	Subscribe(ctx context.Context, userID uuid.UUID) (*domain.EventChannels, uuid.UUID)
	Unsubscribe(ctx context.Context, userID, channelID uuid.UUID, cause string)
	Publish(ctx context.Context, event domain.Event, userID, channelID uuid.UUID)
	PublishBrokers(ctx context.Context, action enum.EventAction, userID, brokerID uuid.UUID)
	PublishDevices(ctx context.Context, action enum.EventAction, userID, brokerID, deviceID uuid.UUID)
	PublishDeviceControls(ctx context.Context, action enum.EventAction, userID, brokerID, deviceID, deviceControlID uuid.UUID)
}

type eventService struct {
	users map[string]*domain.EventChannels
}

func NewEventService() EventService {
	return &eventService{
		users: make(map[string]*domain.EventChannels),
	}
}

func (s *eventService) Subscribe(ctx context.Context, userID uuid.UUID) (*domain.EventChannels, uuid.UUID) {
	u, ok := s.users[userID.String()]
	if !ok {
		u = &domain.EventChannels{
			Channels: make(map[string]chan domain.Event),
			Mutex:    sync.Mutex{},
		}
		s.users[userID.String()] = u
		log.Printf("eventService.Subscribe: user %s, created a new event channels", userID)
	}

	channelID := uuid.New()

	u.Mutex.Lock()
	u.Channels[channelID.String()] = make(chan domain.Event)
	u.Mutex.Unlock()

	log.Printf("eventService.Subscribe: user %s, added a new channel %s", userID, channelID)
	return u, channelID
}

func (s *eventService) Publish(ctx context.Context, event domain.Event, userID, channelID uuid.UUID) {
	u, ok := s.users[userID.String()]
	if !ok {
		log.Printf("eventService.Publish: user %s, action %s, no event channels has been found", userID, event.Data.Action)
		return
	}

	u.Mutex.Lock()
	if channelID == uuid.Nil {
		for _, ch := range u.Channels {
			ch <- event
		}
		log.Printf("eventService.Publish: user %s, action %s, published to all channels (len: %d)", userID, event.Data.Action, len(u.Channels))
	} else {
		u.Channels[channelID.String()] <- event
		log.Printf("eventService.Publish: user %s, action %s, published to channel %s", userID, event.Data.Action, channelID)
	}
	u.Mutex.Unlock()
}

func (s *eventService) Unsubscribe(ctx context.Context, userID, channelID uuid.UUID, cause string) {
	u, ok := s.users[userID.String()]
	if !ok {
		log.Printf("eventService.Unsubscribe: user %s, no event channels has been found", userID)
		return
	}

	u.Mutex.Lock()
	close(u.Channels[channelID.String()])
	delete(u.Channels, channelID.String())
	u.Mutex.Unlock()

	log.Printf("eventService.Unsubscribe: user %s, channel %s -> %s", userID, channelID, cause)

	if len(u.Channels) == 0 {
		delete(s.users, userID.String())
		log.Printf("eventService.Unsubscribe: user %s, removed - no channels left", userID)
	}
}

func (s *eventService) PublishBrokers(ctx context.Context, action enum.EventAction, userID, brokerID uuid.UUID) {
	s.Publish(ctx, domain.Event{
		ID: uuid.New(),
		Data: domain.EventData{
			Action: action,
			Entity: &domain.EventEntity{
				ID:   brokerID,
				Name: enum.BrokersEntity,
			},
		},
	}, userID, uuid.Nil)
}

func (s *eventService) PublishDevices(ctx context.Context, action enum.EventAction, userID, brokerID, deviceID uuid.UUID) {
	related := []domain.EventEntity{}

	if brokerID != uuid.Nil {
		related = append(related, domain.EventEntity{
			ID:   brokerID,
			Name: enum.BrokersEntity,
		})
	}

	s.Publish(ctx, domain.Event{
		ID: uuid.New(),
		Data: domain.EventData{
			Action: action,
			Entity: &domain.EventEntity{
				ID:   deviceID,
				Name: enum.DevicesEntity,
			},
			Related: &related,
		},
	}, userID, uuid.Nil)
}

func (s *eventService) PublishDeviceControls(ctx context.Context, action enum.EventAction, userID, brokerID, deviceID, deviceControlID uuid.UUID) {
	related := []domain.EventEntity{
		{
			ID:   deviceID,
			Name: enum.DevicesEntity,
		},
	}

	if brokerID != uuid.Nil {
		related = append(related, domain.EventEntity{
			ID:   brokerID,
			Name: enum.BrokersEntity,
		})
	}

	s.Publish(ctx, domain.Event{
		ID: uuid.New(),
		Data: domain.EventData{
			Action: action,
			Entity: &domain.EventEntity{
				ID:   deviceID,
				Name: enum.DeviceControlsEntity,
			},
			Related: &related,
		},
	}, userID, uuid.Nil)
}
