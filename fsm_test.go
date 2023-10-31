package fsm

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewFSM(t *testing.T) {
	t.Run("no transitions", func(t *testing.T) {
		item := "testing string"
		fsm := NewFSM(State("test"), &item)

		if len(fsm.transitions) != 0 {
			t.Errorf("expected 0 transitions got %d", len(fsm.transitions))
		}
		if fsm.currentState != State("test") {
			t.Errorf("expected current state to be %s, got %s", State("test"), fsm.currentState)
		}
	})

	t.Run("one transition", func(t *testing.T) {
		item := "testing string"
		fsm := NewFSM(State("test"), &item, Transition[string]{
			From:  "test",
			Event: "run",
			To:    "running",
		})

		if len(fsm.transitions) != 1 {
			t.Errorf("expected 1 transitions got %d", len(fsm.transitions))
		}
		if fsm.currentState != State("test") {
			t.Errorf("expected current state to be %s, got %s", State("test"), fsm.currentState)
		}
	})
}

func TestFSM_State(t *testing.T) {
	x := "hi"
	fsm := NewFSM(State("test"), &x)

	if fsm.State() != State("test") {
		t.Errorf("expected current state to be %s, got %s", State("test"), fsm.State())
	}
}

func TestFSM_HandleEvent(t *testing.T) {
	t.Run("no before/after func", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
		})

		if err := fsm.HandleEvent(Event("send")); err != nil {
			t.Error(err)
		}

		if fsm.State() != State("sent") {
			t.Errorf("expected current state to be %s, got %s", State("sent"), fsm.State())
		}
	})

	t.Run("before func", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			BeforeFn: func(s *string) error {
				*s = "ive been changed"
				return nil
			},
		})

		if err := fsm.HandleEvent(Event("send")); err != nil {
			t.Error(err)
		}

		if item != "ive been changed" {
			t.Errorf("expected beforeFn to run but it didn't")
		}
		if fsm.State() != State("sent") {
			t.Errorf("expected current state to be %s, got %s", State("sent"), fsm.State())
		}
	})

	t.Run("after func", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			AfterFn: func(s *string) error {
				*s = "ive been changed"
				return nil
			},
		})

		if err := fsm.HandleEvent(Event("send")); err != nil {
			t.Error(err)
		}

		if item != "ive been changed" {
			t.Errorf("expected beforeFn to run but it didn't")
		}
		if fsm.State() != State("sent") {
			t.Errorf("expected current state to be %s, got %s", State("sent"), fsm.State())
		}
	})

	t.Run("before and after func", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			BeforeFn: func(s *string) error {
				*s = "ive been changed"
				return nil
			},
			AfterFn: func(s *string) error {
				*s = "ive been changed again"
				return nil
			},
		})

		if err := fsm.HandleEvent(Event("send")); err != nil {
			t.Error(err)
		}

		if item != "ive been changed again" {
			t.Errorf("expected beforeFn to run but it didn't")
		}
		if fsm.State() != State("sent") {
			t.Errorf("expected current state to be %s, got %s", State("sent"), fsm.State())
		}
	})

	t.Run("no transition for event", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			BeforeFn: func(s *string) error {
				*s = "ive been changed"
				return nil
			},
			AfterFn: func(s *string) error {
				*s = "ive been changed again"
				return nil
			},
		})

		err := fsm.HandleEvent(Event("collect"))
		if err == nil {
			t.Errorf("expected error got nil")
		}

		var te *TransitionError
		if !errors.As(err, &te) {
			t.Errorf("expected error to be type TransitionError, got %T", err)
		}
	})

	t.Run("before func error", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			BeforeFn: func(s *string) error {
				return fmt.Errorf("before func error")
			},
			AfterFn: func(s *string) error {
				*s = "ive been changed again"
				return nil
			},
		})

		err := fsm.HandleEvent(Event("send"))
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("after func error", func(t *testing.T) {
		item := "pending"
		fsm := NewFSM(State("pending"), &item, Transition[string]{
			From:  "pending",
			Event: "send",
			To:    "sent",
			BeforeFn: func(s *string) error {
				*s = "ive been changed"
				return nil
			},
			AfterFn: func(s *string) error {
				return fmt.Errorf("after func error")
			},
		})

		err := fsm.HandleEvent(Event("send"))
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestTransitionError_Error(t *testing.T) {
	type fields struct {
		from  State
		event Event
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "transition error",
			fields: fields{
				from:  State("test"),
				event: Event("collect"),
			},
			want: "invalid transition from test with event collect",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := &TransitionError{
				from:  tt.fields.from,
				event: tt.fields.event,
			}
			if got := te.Error(); got != tt.want {
				t.Errorf("TransitionError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
