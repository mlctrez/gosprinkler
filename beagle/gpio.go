package beagle

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Api struct {
	GpioPaths []string
}

func New() *Api {

	api := &Api{
		GpioPaths: []string{
			// 6 0 P8.8 /sys/class/gpio/gpio67
			"/sys/class/gpio/gpio67",
			// 8 1 P8.10 /sys/class/gpio/gpio68
			"/sys/class/gpio/gpio68",
			// 10 2 P8.12 /sys/class/gpio/gpio44
			"/sys/class/gpio/gpio44",
			// 12 3 P8.14 /sys/class/gpio/gpio26
			"/sys/class/gpio/gpio26",
			// 14 4 P8.16 /sys/class/gpio/gpio46
			"/sys/class/gpio/gpio46",
			// 16 5 P8.18 /sys/class/gpio/gpio65
			"/sys/class/gpio/gpio65",
		},
	}

	api.InitializePins()

	return api
}

func writeString(path, value string) (err error) {

	var f *os.File

	mode := os.O_WRONLY | os.O_TRUNC

	if f, err = os.OpenFile(path, mode, 0666); err != nil {
		return
	}

	defer f.Close()

	_, err = f.Write([]byte(value))

	return
}

func (a *Api) InitializePins() error {

	log.Println("InitializePins")

	for _, path := range a.GpioPaths {
		err := writeString(path+"/direction", "out")
		if err != nil {
			log.Println(err)
		}
	}

	a.PinsOff()

	return nil
}

func (a *Api) ChangePin(pin, state string) {
	thePin, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	if thePin >= 0 && thePin < len(a.GpioPaths) {
		switch strings.ToLower(state) {
		case "on", "true":
			log.Printf("turning on pin %d\n", thePin)
			writeString(a.GpioPaths[thePin]+"/value", "1")
		default:
			log.Printf("turning off pin %d\n", thePin)
			writeString(a.GpioPaths[thePin]+"/value", "0")
		}
	}
}

func (a *Api) PinsOff() {
	for _, p := range a.GpioPaths {
		writeString(p+"/value", "0")
	}
}

func (a *Api) Shutdown() {
	log.Println("shutting down all zones")
	// force all zones low to minimize water bill
	a.PinsOff()
}
