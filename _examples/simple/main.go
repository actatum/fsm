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

func main() {
	ls := LightSwitch{
		State: "off",
		Name:  "living room",
	}
	m := fsm.NewFSM[LightSwitch](fsm.State("off"), &ls, []fsm.Transition[LightSwitch]{
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

	fmt.Println(m.State())
}
