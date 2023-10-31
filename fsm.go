package fsm

import "fmt"

// State is the condition that an item can be in at a given time.
type State string

// Event is the idea of something happening to trigger a state transition.
type Event string

// Transition represents a path from one state to another through a specific event.
// A transition can have extra validations passed to it via a BeforeFn (Example: Make sure a balance is above $0).
// A transition can have side effects passed to it via an AfterFn (Example: save to a database).
type Transition[T any] struct {
	From  State
	Event Event
	To    State
	// BeforeFn can be used to run extra validations for state transitions.
	BeforeFn func(*T) error
	// AfterFn can be used for any actions that need to happen after a state transition.
	// For example saving to a database, logging, etc.
	AfterFn func(*T) error
}

// FSM models a finite state machine.
// It holds an item or entity that has multiple states and a set of allowable transitions between the states.
type FSM[T any] struct {
	item *T

	currentState State
	transitions  []Transition[T]
}

// NewFSM constructs a new state machine given the current state of an item and a list of possible transitions.
func NewFSM[T any](currentState State, item *T, transitions ...Transition[T]) *FSM[T] {
	return &FSM[T]{
		item:         item,
		currentState: currentState,
		transitions:  transitions,
	}
}

// State returns the current state of the item associated with the FSM.
func (f *FSM[T]) State() State {
	return f.currentState
}

// HandleEvent handles an event that should trigger a transition between states.
func (f *FSM[T]) HandleEvent(event Event) error {
	currentState := f.currentState
	for _, t := range f.transitions {
		if t.From == currentState && t.Event == event {
			if t.BeforeFn != nil {
				if err := t.BeforeFn(f.item); err != nil {
					return err
				}
			}

			f.currentState = t.To
			if t.AfterFn != nil {
				if err := t.AfterFn(f.item); err != nil {
					return err
				}
			}

			return nil
		}
	}

	return &TransitionError{from: currentState, event: event}
}

// TransitionError indicates an invalid transition between states.
type TransitionError struct {
	from  State
	event Event
}

func (te *TransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s with event %s", te.from, te.event)
}
