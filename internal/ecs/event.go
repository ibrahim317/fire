package ecs

// EventType identifies the type of event.
type EventType int

const (
	// EventPlayerJump is fired when the player initiates a jump.
	EventPlayerJump EventType = iota
	// EventPlayerLand is fired when the player lands on the ground.
	EventPlayerLand
	// EventCollision is fired when two entities collide.
	EventCollision
	// EventDamage is fired when an entity takes damage.
	EventDamage
	// EventDeath is fired when an entity's health reaches zero.
	EventDeath
	// EventCoinCollected is fired when a coin is collected.
	EventCoinCollected
)

// Event represents a game event with source, target, and data.
type Event struct {
	Type   EventType
	Source EntityID
	Target EntityID
	Data   interface{}
}

// EventHandler is a function that handles an event.
type EventHandler func(Event)

// EventBus manages event publication and subscription.
type EventBus struct {
	handlers map[EventType][]EventHandler
	queue    []Event
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
		queue:    make([]Event, 0),
	}
}

// Subscribe registers a handler for a specific event type.
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish queues an event for processing.
func (eb *EventBus) Publish(event Event) {
	eb.queue = append(eb.queue, event)
}

// Process processes all queued events and clears the queue.
func (eb *EventBus) Process() {
	for _, event := range eb.queue {
		if handlers, ok := eb.handlers[event.Type]; ok {
			for _, handler := range handlers {
				handler(event)
			}
		}
	}
	eb.queue = eb.queue[:0]
}

// PublishImmediate publishes and immediately processes an event.
func (eb *EventBus) PublishImmediate(event Event) {
	if handlers, ok := eb.handlers[event.Type]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// Clear removes all handlers and queued events.
func (eb *EventBus) Clear() {
	eb.handlers = make(map[EventType][]EventHandler)
	eb.queue = eb.queue[:0]
}
