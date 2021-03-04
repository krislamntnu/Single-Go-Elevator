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
	//startFloor := elevio.GetFloor()
	startFloor := -1
	if startFloor == -1 {
		elev = elevator{floor: -1, dir: elevio.MdDown, behaviour: ebMoving}
	} else {
		elev = elevator{floor: startFloor, dir: elevio.MdStop, behaviour: ebIdle}
	}
	elevio.SetMotorDirection(elev.dir)

	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)

	for {
		select {
		case a := <-drvButtons:
			fmt.Printf("%+v\n", a)
			// elevio.SetButtonLamp(a.Button, a.Floor, true)
			elev.requests[a.Floor][a.Button] = true
			setAllLights(elev)

		case a := <-drvFloors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				elev.dir = elevio.MdDown
			} else if a == 0 {
				elev.dir = elevio.MdUp
			}
			elevio.SetMotorDirection(elev.dir)
			elev.floor = a
			elev = requestsClearCurrentFloor(elev)
			setAllLights(elev)

		case a := <-drvObstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MdStop)
			} else {
				elevio.SetMotorDirection(elev.dir)
			}
		}
	}
}
