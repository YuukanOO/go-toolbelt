package eventsourcing

import (
	"reflect"
	"testing"
)

func TestDispatcher(t *testing.T) {
	obj := &ESObject{}
	TrackChange(obj, CreatedEvent{ID: 2})
	TrackChange(obj, AnotherEvent{})

	dispatcher := NewDispatcher()

	if len(dispatcher.handlers) != 0 {
		t.Error("The dispatcher should have any handlers for now")
	}

	handlerStack := []Event{}
	anotherHandlerStack := []Event{}

	handler := func(evt Event) {
		handlerStack = append(handlerStack, evt)
	}

	anotherHandler := func(evt Event) {
		anotherHandlerStack = append(anotherHandlerStack, evt)
	}

	dispatcher.AddHandlers(handler, anotherHandler)

	if len(dispatcher.handlers) != 2 {
		t.Error("We should have 2 handlers now")
	}

	dispatcher.Dispatch(obj)

	if len(anotherHandlerStack) != 2 || len(handlerStack) != 2 {
		t.Error("Every stack should contains 2 events")
	}

	evt, aevt := handlerStack[0], anotherHandlerStack[0]

	if evt != aevt {
		t.Error("Fire events should be the same in each handler")
	}

	if reflect.TypeOf(evt).Name() != "CreatedEvent" || reflect.TypeOf(aevt).Name() != "CreatedEvent" {
		t.Error("First event should be of type CreatedEvent")
	}

	evt, aevt = handlerStack[1], anotherHandlerStack[1]

	if reflect.TypeOf(evt).Name() != "AnotherEvent" || reflect.TypeOf(aevt).Name() != "AnotherEvent" {
		t.Error("First event should be of type AnotherEvent")
	}
}
