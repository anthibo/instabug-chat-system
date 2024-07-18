package messaging

type EventPublisherManager interface {
	PublishEvent(event string, body []byte) error
}
