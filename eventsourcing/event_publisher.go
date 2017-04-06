package eventsourcing

// EventHandler represents a delegate that should be called when a new event has
// been dispatched.
type EventHandler func(evt Event)

// Dispatcher represents an event dispatcher which holds handlers.
type Dispatcher struct {
	handlers []EventHandler
}

// NewDispatcher instantiates a new dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}
