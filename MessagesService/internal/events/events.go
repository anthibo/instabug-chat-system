package events

type EventQueue struct {
	Name string
	Type string
}

const (
	MessageCreatedQueue           = "message_created_queue"
	MessageCreationRequestedQueue = "message_creation_requested_queue"
)

var EventQueues = map[string]EventQueue{
	MessageCreatedQueue:           {Name: "message_created_queue", Type: "direct"},
	MessageCreationRequestedQueue: {Name: "message_creation_requested_queue", Type: "direct"},
}

func GetEventQueue(name string) (EventQueue, bool) {
	queue, exists := EventQueues[name]
	return queue, exists
}
