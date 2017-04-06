package eventsourcing

import "testing"
import "reflect"

type ESObject struct {
	EventSource
	ID              int
	processedEvents int
}

type CreatedEvent struct{ ID int }
type AnotherEvent struct{}

func (obj *ESObject) Transition(evt Event) {
	switch e := evt.(type) {
	case CreatedEvent:
		obj.processedEvents++
		obj.ID = e.ID
		break
	case AnotherEvent:
		obj.processedEvents++
		break
	}
}

func TestEventSource(t *testing.T) {
	obj := &ESObject{}

	if len(obj.Changes) != 0 {
		t.Error("Changes should be empty for now")
	}

	TrackChange(obj, CreatedEvent{ID: 1337})

	if len(obj.Changes) == 0 {
		t.Error("Changes should contains the CreatedEvent now")
	}

	if obj.processedEvents != 1 || obj.ID != 1337 {
		t.Error("Processed events should have been incremented and ID should be set")
	}

	eventsFromStore := []Event{
		CreatedEvent{ID: 6},
		AnotherEvent{},
	}

	newObj := &ESObject{}

	LoadFromEvents(newObj, eventsFromStore)

	if newObj.ID != 6 || newObj.processedEvents != 2 {
		t.Error("Processed events should have been set and ID should be set")
	}
}

func TestPopChange(t *testing.T) {
	obj := &ESObject{}
	TrackChange(obj, CreatedEvent{ID: 2})
	TrackChange(obj, AnotherEvent{})

	if len(obj.Changes) != 2 {
		t.Error("Object should have 2 pending changes")
	}

	evt := obj.PopChange()
	evtType := reflect.TypeOf(evt)

	if evtType.Name() != "CreatedEvent" {
		t.Error("First event should be of type CreatedEvent")
	}

	evt = obj.PopChange()
	evtType = reflect.TypeOf(evt)

	if evtType.Name() != "AnotherEvent" {
		t.Error("Second event should be of type AnotherEvent")
	}

	evt = obj.PopChange()

	if evt != nil {
		t.Error("There should be no event anymore")
	}
}
