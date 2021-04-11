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

	timeout := time.NewTimer(doorOpenDurationSec * time.Second)

	// Solves the special case with obstruction while the door is closed.
	obstructionSwitch := false

	// Finite state machine
	for {
		select {
		case <-timeout.C:
			switch state {
			case dsOpen:
				elevio.SetDoorOpenLamp(false)
				doorClosed <- true
				state = dsClosed
			default:
			}

		case <-openDoor:
			fmt.Println("Door open")
			switch state {
			case dsClosed:
				elevio.SetDoorOpenLamp(true)
				if !obstructionSwitch {
					timeout.Reset(doorOpenDurationSec * time.Second)
					state = dsOpen
				} else {
					state = dsObstructed
				}
			case dsOpen:
				timeout.Reset(doorOpenDurationSec * time.Second)
			default:
			}

		case obstruction := <-obstructionEvents:
			fmt.Println("Obstruction:", obstruction)
			switch state {
			case dsOpen:
				if obstruction {
					state = dsObstructed
					obstructionSwitch = true
				}
			case dsObstructed:
				if !obstruction {
					state = dsOpen
					timeout.Reset(doorOpenDurationSec * time.Second)
					obstructionSwitch = false
				}
			case dsClosed:
				obstructionSwitch = obstruction
			}
		}
	}
}
