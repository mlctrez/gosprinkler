package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mlctrez/hwio"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func initializePins() (pins []hwio.Pin, err error) {

	log.Println("initializePins")

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

	return pins, err

}

func shutdown() {

	log.Println("shutting down all zones")

	// force all zones low to minimize water bill
	initializePins()

	// close, per the hwio documentation
	hwio.CloseAll()
}

var sprinklerControlHTML = `
<html>
<head>
</head>
<body>
<span style="font-size: 72px">
Pin 0 <a href="/sprinkler/api?pin=0&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=0">OFF</a><br>
Pin 1 <a href="/sprinkler/api?pin=1&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=1">OFF</a><br>
Pin 2 <a href="/sprinkler/api?pin=2&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=2">OFF</a><br>
Pin 3 <a href="/sprinkler/api?pin=3&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=3">OFF</a><br>
Pin 4 <a href="/sprinkler/api?pin=4&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=4">OFF</a><br>
Pin 5 <a href="/sprinkler/api?pin=5&cmd=on">ON</a>&nbsp;&nbsp;<a href="/sprinkler/api?pin=5">OFF</a><br>
</span>
</body>
`

func runHttpServer() {

	hwio.SetDriver(new(hwio.BeagleBoneBlackDriver))

	pins, err := initializePins()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("go away"))
	})

	http.HandleFunc("/sprinkler", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(sprinklerControlHTML))
	})

	http.HandleFunc("/sprinkler/api", func(rw http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pin, err := strconv.Atoi(r.FormValue("pin"))
		if err != nil {
			return
		}
		if pin >= 0 && pin < len(pins) {
			if r.FormValue("cmd") == "on" {
				hwio.DigitalWrite(pins[pin], hwio.HIGH)
			} else {
				hwio.DigitalWrite(pins[pin], hwio.LOW)
			}
		}
		rw.Write([]byte(sprinklerControlHTML))

	})

	http.ListenAndServe(":9090", nil)
}

func main() {

	// for driver and zone cleanup
	defer shutdown()

	if len(os.Args) == 2 {
		if os.Args[1] == "http" {
			runHttpServer()
			return
		}
	}

	if len(os.Args) == 2 {
		if os.Args[1] == "stop" {
			return
		}
	}

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

		log.Printf("turning on zone %v\n", zone)

		hwio.DigitalWrite(pin, hwio.HIGH)

		switch zone {
		case 5:
			time.Sleep(zoneFiveDuration)
		default:
			time.Sleep(defaultDuration)
		}

		log.Printf("turning off zone %v\n", zone)

		hwio.DigitalWrite(pin, hwio.LOW)

		// pause between zones to allow the sprinkler heads to retract
		time.Sleep(30 * time.Second)
	}

}
