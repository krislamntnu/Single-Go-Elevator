package main

import (
	"fmt"
	"time"

	"github.com/krislamntnu/Single-Go-Elevator.git/elevio"
)

const doorOpenDurationSec = 3

type doorState int

const (
	dsClosed doorState = iota
	dsOpen
	dsObstructed
)

func doorFsm(doorClosed chan<- bool, openDoor <-chan bool) {
	state := dsClosed
	elevio.SetDoorOpenLamp(false)

	obstructionEvents := make(chan bool)
	go elevio.PollObstructionSwitch(obstructionEvents)

	for {
		select {
		case shouldOpen := <-openDoor:
			if shouldOpen {
				fmt.Println("Door open")
				elevio.SetDoorOpenLamp(true)
				state = dsOpen

				// For now a sleep will do. Fixing a real timer later
				time.Sleep(doorOpenDurationSec * time.Second)
				elevio.SetDoorOpenLamp(false)
				doorClosed <- true
				state = dsClosed
			}
		case obstruction := <-obstructionEvents:
			// Some nonesense. Fix later
			if state == dsClosed {
				fmt.Printf("%+v\n", obstruction)
			}
		}
	}
}
