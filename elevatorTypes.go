package main

import "github.com/krislamntnu/Single-Go-Elevator.git/elevio"

const (
	numFloors           = 4
	numButtons          = 3 // Cab, HallUp and HallDown
	doorOpenDurationSec = 3
)

type elevatorBehaviour int

const (
	ebIdle elevatorBehaviour = iota
	ebDoorOpen
	ebMoving
)

type elevator struct {
	floor     int
	dir       elevio.MotorDirection
	requests  [numFloors][numButtons]bool
	behaviour elevatorBehaviour
}
