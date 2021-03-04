package main

import "github.com/krislamntnu/Single-Go-Elevator.git/elevio"

const (
	numFloors  = 4
	numButtons = 3
)

type elevatorBehaviour int

const (
	ebIdle elevatorBehaviour = iota
	ebDoorOpen
	ebMoving
)

type elevator struct {
	floor     int
	dir       elevio.MotorDirection // Just up or down
	requests  [numFloors][numButtons]bool
	behaviour elevatorBehaviour // State on which det fsm switches
}
