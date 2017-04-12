package eventsourcing

import "testing"

func TestDispatcher(t *testing.T) {
	obj := &ESObject{}
	TrackChange(obj, CreatedEvent{ID: 2})
	TrackChange(obj, AnotherEvent{})

	dispatcher := NewDispatcher()

	if len(dispatcher.handlers) != 0 {
		t.Error("The dispatcher should have any handlers for now")
	}

	handler := func(evt Event) {

	}

	dispatcher.AddHandler(handler)
}
