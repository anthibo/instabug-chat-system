package messaging

import "message_service/internal/events"

type EventPublisherManager interface {
	PublishEvent(event events.EventQueue, eventData interface{}) error
}
