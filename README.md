# FSM

[![GoDoc Widget]][GoDoc]

`fsm` is a small library to model a finite state machine.

## Install

`go get -u github.com/actatum/fsm`

## Examples

See [\_examples/](https://github.com/actatum/fsm/blob/master/_examples/) for various examples.

** Simple Example **

```go
package main

import (
	"fmt"
	"log"

	"github.com/actatum/fsm"
)

type LightSwitch struct {
	State string
	Name  string
}

func (ls *LightSwitch) SetState(s fsm.State) {
	ls.State = string(s)
}

func main() {
	ls := LightSwitch{
		State: "off",
		Name:  "living room",
	}
	m := fsm.NewFSM[*LightSwitch](fsm.State("off"), &ls, []fsm.Transition[*LightSwitch]{
		{
			From:  fsm.State("off"),
			Event: fsm.Event("flip_switch"),
			To:    fsm.State("on"),
		},
		{
			From:  fsm.State("on"),
			Event: fsm.Event("flip_switch"),
			To:    fsm.State("off"),
		},
	}...)

	if err := m.HandleEvent(fsm.Event("flip_switch")); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m.State(), ls.State)
}
```

[GoDoc]: https://pkg.go.dev/github.com/actatum/fsm
[GoDoc Widget]: https://godoc.org/github.com/actatum/fsm?status.svg
