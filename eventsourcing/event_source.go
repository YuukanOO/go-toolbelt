// Package eventsourcing exposes event sourcing basic interface and class.
// Based on the work of James Nugent available at http://jen20.com/2015/02/08/event-sourcing-in-go.html
// In order for this to work, I had to exposes interface members but I think its
// the best I can come up with.
package eventsourcing

// Event represents a single event object
type Event interface{}

// EventEmitter is the core interface to track changes in the event sourcing
// system.
type EventEmitter interface {
	// Transition process an event and update the state of the entity.
	// It should be implemented by domain objects.
	Transition(event Event)
	// AddChange adds a change to the entity without transitioning. Mostly a
	// persistence mechanism.
	AddChange(event Event)
	// IncrementVersion increments the version of the event source object.
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

// IncrementVersion increments the version. Used for concurrency.
func (src *EventSource) IncrementVersion() {
	src.ExpectedVersion++
}

// AddChange adds a change to the inner list of the event source. It will basically
// be called after an object has transited from a state to a new one.
func (src *EventSource) AddChange(event Event) {
	src.Changes = append(src.Changes, event)
}
