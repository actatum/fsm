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
	db := make(map[string]*LightSwitch)
	ls := LightSwitch{
		State: "off",
		Name:  "living room",
	}
	m := fsm.NewFSM[*LightSwitch](fsm.State("off"), &ls, []fsm.Transition[*LightSwitch]{
		{
			From:  fsm.State("off"),
			Event: fsm.Event("flip_switch"),
			To:    fsm.State("on"),
			BeforeFn: func(ls *LightSwitch) error {
				if _, ok := db[ls.Name]; !ok {
					db[ls.Name] = ls
				}
				return nil
			},
			AfterFn: func(ls *LightSwitch) error {
				db[ls.Name] = ls
				return nil
			},
		},
		{
			From:  fsm.State("on"),
			Event: fsm.Event("flip_switch"),
			To:    fsm.State("off"),
			BeforeFn: func(ls *LightSwitch) error {
				if _, ok := db[ls.Name]; !ok {
					db[ls.Name] = ls
				}
				return nil
			},
			AfterFn: func(ls *LightSwitch) error {
				db[ls.Name] = ls
				return nil
			},
		},
	}...)

	if err := m.HandleEvent(fsm.Event("flip_switch")); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m.State(), ls.State)
	fmt.Printf("%+v\n", db)
}
