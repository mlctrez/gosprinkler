package main

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"log"
)

func initializePins() (pins []hwio.Pin, err error) {

	fmt.Println("initializePins")

	pinNames := []string{"P8.8", "P8.10", "P8.12", "P8.14", "P8.16", "P8.18"}

	pins = make([]hwio.Pin, len(pinNames))

	for i, pinName := range pinNames {
		myPin, err := hwio.GetPinWithMode(pinName, hwio.OUTPUT)
		if err != nil {
			return nil, err
		}

		err = hwio.DigitalWrite(myPin, hwio.LOW)

		if err != nil {
			return nil, err
		}

		pins[i] = myPin

	}

	return pins, nil

}

func main() {

	// force beaglebone driver since hwio.MatchesHardwareConfig for the beaglebone driver
	// fails to detect the driver.
	hwio.SetDriver(new(hwio.BeagleBoneBlackDriver))

	// for cleanup
	defer initializePins()

	defer hwio.CloseAll()

	pins, err := initializePins()
	if err != nil {
		log.Fatal(err)
	}

	_ = pins

	//time.Sleep(5 * time.Second)
	//
	//fmt.Println("turning on zone 0")
	//
	//hwio.DigitalWrite(pins[0], hwio.HIGH)
	//
	//time.Sleep(120 * time.Second)
	//
	//fmt.Println("turning off zone 0")
	//
	//hwio.DigitalWrite(pins[0], hwio.LOW)

	hwio.DebugPinMap()

}
