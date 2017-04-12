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

// AddHandlers adds one or more handlers to this dispatcher's instance.
func (d *Dispatcher) AddHandlers(handlers ...EventHandler) {
	d.handlers = append(d.handlers, handlers...)
}

func (d *Dispatcher) Dispatch(emitters ...EventEmitter) {

}
