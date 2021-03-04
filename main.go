package main

import (
	"fmt"

	"github.com/krislamntnu/Single-Go-Elevator.git/elevio"
)

func setAllLights(elev elevator) {
	for floor := 0; floor < numFloors; floor++ {
		for btn := elevio.ButtonType(0); btn < numButtons; btn++ {
			elevio.SetButtonLamp(btn, floor, elev.requests[floor][btn])
		}
	}
}

func main() {
	elevio.Init("localhost:15657", numFloors)

	// Initialize and start the elevator
	var elev elevator
	startFloor := elevio.GetFloor()
	if startFloor == -1 {
		elev = elevator{floor: -1, dir: elevio.MdDown, behaviour: ebMoving}
	} else {
		elev = elevator{floor: startFloor, dir: elevio.MdStop, behaviour: ebIdle}
	}
	elevio.SetMotorDirection(elev.dir)

	buttonEvents := make(chan elevio.ButtonEvent)
	hitFloorEvents := make(chan int)
	openDoor := make(chan bool)
	doorClosed := make(chan bool)

	go elevio.PollButtons(buttonEvents)
	go elevio.PollFloorSensor(hitFloorEvents)
	go doorFsm(doorClosed, openDoor)

	for {
		select {
		case button := <-buttonEvents:
			fmt.Printf("ButtonEvent: %+v\n", button)

			switch elev.behaviour {
			case ebDoorOpen:
				if elev.floor == button.Floor {
					openDoor <- true
				} else {
					elev.requests[button.Floor][button.Button] = true
				}

			case ebMoving:
				elev.requests[button.Floor][button.Button] = true

			case ebIdle:
				if elev.floor == button.Floor {
					openDoor <- true
					elev.behaviour = ebDoorOpen
				} else {
					elev.requests[button.Floor][button.Button] = true
					elev.dir = requestsChooseDirection(elev)
					elevio.SetMotorDirection(elev.dir)
					elev.behaviour = ebMoving
				}

			default:

			}
			setAllLights(elev)

		case newFloor := <-hitFloorEvents:
			fmt.Printf("New floor: %+v\n", newFloor)

			elev.floor = newFloor
			elevio.SetFloorIndicator(newFloor)

			switch elev.behaviour {
			case ebMoving:
				if requestsShouldStop(elev) {
					elevio.SetMotorDirection(elevio.MdStop)
					openDoor <- true
					elev = requestsClearCurrentFloor(elev)
					setAllLights(elev)
					elev.behaviour = ebDoorOpen
				}
			default:
			}

		case <-doorClosed:
			fmt.Printf("Door closed\n")

			switch elev.behaviour {
			case ebDoorOpen:
				elev.dir = requestsChooseDirection(elev)
				elevio.SetMotorDirection(elev.dir)

				if elev.dir == elevio.MdStop {
					elev.behaviour = ebIdle
				} else {
					elev.behaviour = ebMoving
				}

			default:

			}
		}

	}
}
