package eventsourcing

// Event represents a single event object
type Event interface{}

// EventEmitter is the core interface to track changes in the event sourcing
// system.
type EventEmitter interface {
	Transition(event Event)
	AddChange(event Event)
	IncrementVersion()
}

// EventSource is a basic implementation of an event emitter.
type EventSource struct {
	Changes         []Event
	ExpectedVersion int
}

// TrackChange process an event into the given emitter to change its state
// add the event to the list of changes of the event source.
func TrackChange(src EventEmitter, event Event) {
	src.AddChange(event)
	src.Transition(event)
}

// LoadFromEvents reconstructs an object from a list of events.
func LoadFromEvents(src EventEmitter, events []Event) {
	for _, event := range events {
		src.Transition(event)
		src.IncrementVersion()
	}
}

func (src *EventSource) IncrementVersion() {
	src.ExpectedVersion++
}

func (src *EventSource) AddChange(event Event) {
	src.Changes = append(src.Changes, event)
}
