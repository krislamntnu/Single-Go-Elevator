package main

import "github.com/krislamntnu/Single-Go-Elevator.git/elevio"

func requestsChooseDirection(elev elevator) elevio.MotorDirection {
	switch elev.dir {
	case elevio.MdUp:
		if requestsAbove(elev) {
			return elevio.MdUp
		} else if requestsBelow(elev) {
			return elevio.MdDown
		} else {
			return elevio.MdStop
		}
	case elevio.MdDown:
		if requestsBelow(elev) {
			return elevio.MdDown
		} else if requestsAbove(elev) {
			return elevio.MdUp
		} else {
			return elevio.MdStop
		}
	case elevio.MdStop:
		if requestsBelow(elev) {
			return elevio.MdDown
		} else if requestsAbove(elev) {
			return elevio.MdUp
		} else {
			return elevio.MdStop
		}
	default:
		return elevio.MdStop
	}
}

func requestsShouldStop(elev elevator) bool {
	switch elev.dir {
	case elevio.MdUp:
		return elev.requests[elev.floor][elevio.BtHallUp] ||
			elev.requests[elev.floor][elevio.BtCab] ||
			!requestsAbove(elev)
	case elevio.MdDown:
		return elev.requests[elev.floor][elevio.BtHallDown] ||
			elev.requests[elev.floor][elevio.BtCab] ||
			!requestsBelow(elev)
	default:
		return true
	}
}

func requestsClearCurrentFloor(elev elevator) elevator {
	for btn := elevio.ButtonType(0); btn < numButtons; btn++ {
		elev.requests[elev.floor][btn] = false
	}
	return elev
}

func requestsAbove(elev elevator) bool {
	for floor := elev.floor + 1; floor < numFloors; floor++ {
		for btn := elevio.ButtonType(0); btn < numButtons; btn++ {
			if elev.requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(elev elevator) bool {
	for floor := 0; floor < elev.floor; floor++ {
		for btn := elevio.ButtonType(0); btn < numButtons; btn++ {
			if elev.requests[floor][btn] {
				return true
			}
		}
	}
	return false
}
