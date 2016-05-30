package main

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func shutdown() {

	fmt.Println("shutting down all zones")

	// force all zones low to minimize water bill
	initializePins()

	// close, per the hwio documentation
	hwio.CloseAll()
}

func main() {

	// for driver and zone cleanup
	defer shutdown()

	// handle signals and shut down the zones correctly
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-c
		shutdown()
		os.Exit(1)
	}()

	// Force the beaglebone driver since MatchesHardwareConfig looks in the wrong location for my distribution.
	// At some point a pull request should be submitted to address this.
	hwio.SetDriver(new(hwio.BeagleBoneBlackDriver))

	pins, err := initializePins()
	if err != nil {
		log.Fatal(err)
	}

	defaultDuration := 30 * time.Minute
	zoneFiveDuration := 15 * time.Minute

	for zone, pin := range pins {

		fmt.Printf("turning on zone %v\n", zone)

		hwio.DigitalWrite(pin, hwio.HIGH)

		switch zone {
		case 5:
			time.Sleep(zoneFiveDuration)
		default:
			time.Sleep(defaultDuration)
		}

		fmt.Printf("turning off zone %v\n", zone)

		hwio.DigitalWrite(pin, hwio.LOW)

		// pause between zones to allow the sprinkler heads to retract
		time.Sleep(30 * time.Second)
	}

}
