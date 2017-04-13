# go-toolbelt

This is my personal toolbelt for the golang awesome programming language. It is built around ddd practices.

## Installation

`go get github.com/YuukanOO/go-toolbelt`

## Usage

### Errors

A `DomainError` struct is defined and encapsulates a Domain error. A domain error is an expected error in case some input were wrong and will eventually be displayed to the user.

```go
err := errors.NewDomainError("AConstantCode", "A friendly message for the developper", errors.New("Any number of inner errors"))
```

### Validation

This is a very simple fluent like API that uses [the go-playground validator](https://github.com/go-playground/validator) under the hood.

```go
err := validation.Validate("User").
  Field("username", "mytoolongusername", "required,max=10,min=1").
  Field("password", "aS3cretP@ssw0rd", "required,min=10").
  Errors() // Will trigger the evaluation

// If it has an error, a domain error will be returned
domErr := err.(*errors.DomainError)

// And inner Errors will be of type FieldError
fieldErr := domErr.Errors[0].(*validation.FieldError)
```

Don't hesitate to check the tests for more examples.

### Event sourcing

I know I shouldn't have to expose `Transition` and other methods but I had to for this to work.

```go
type User struct {
  eventsourcing.EventSource
  ID int
}

type UserCreated struct { ID int }

func NewUser() *User {
  usr := &User{}
  eventsourcing.TrackChange(usr, UserCreated{ ID: 1 })
  return usr
}

func NewUserFromStore(events []Event) *User {
  usr := &User{}
  eventsourcing.LoadFromEvents(usr, events)
  return usr
}

func (u *User) Transition(evt eventsourcing.Event) {
  switch e := evt.(type) {
    case UserCreated:
      u.ID = e.ID
      break
  }
}

u := NewUser()

// len(u.Changes) == 1
// u.Changes[0] == UserCreated{ 1 }
// u.ID == 1

evts := []Event{
  UserCreated{ ID: 6 }
}

us := NewUserFromStore(evts)

// len(us.Changes) == 0 since it has been reconstructed from the store
// u.ID == 6
```

This toolbelt also implements an event dispatcher:

```go
// Using same user has above

dispatcher := eventsourcing.NewDispatcher()

func handler(evt Event) {
  fmt.Println(reflect.TypeOf(evt).Name())
}

// Add one or more handlers for this dispatcher
dispatcher.AddHandlers(handler)

// len(dispatcher.handlers) == 1

// Dispatch one or more event emitters
dispatcher.Dispatch(u) // u being a NewUser()

// handler will be called so => "UserCreated" will be printed out
```