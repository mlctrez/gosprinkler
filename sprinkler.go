package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/mlctrez/gosprinkler/beagle"
	"github.com/mlctrez/gosprinkler/dashbutton"
	"github.com/mlctrez/gosprinkler/sighandler"
)

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

func runHttpServer(api *beagle.Api, sig *sighandler.Api) {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("go away"))
	})

	http.HandleFunc("/stop", func(rw http.ResponseWriter, r *http.Request) {
		sig.Interrupt()
	})

	http.HandleFunc("/sprinkler", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(sprinklerControlHTML))
	})

	http.HandleFunc("/sprinkler/api", func(rw http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		api.ChangePin(r.FormValue("pin"), r.FormValue("cmd"))

		rw.Write([]byte(sprinklerControlHTML))
	})

	http.ListenAndServe(":9090", nil)
}

func main() {

	api := beagle.New()
	sig := sighandler.New()
	sig.RegisterHandler(api.Shutdown)

	// for driver and zone cleanup on main exit
	defer api.Shutdown()

	httpOnly := flag.Bool("http", false, "if set, run only the http server")
	flag.Parse()

	go dashbutton.New(sig.Interrupt)

	if *httpOnly {
		runHttpServer(api, sig)
		return
	}

	defaultDuration := 30 * time.Minute
	zoneFiveDuration := 15 * time.Minute

	for _, zone := range strings.Split("0,1,2,3,4,5", ",") {

		api.ChangePin(zone, "on")

		switch zone {
		case "5":
			time.Sleep(zoneFiveDuration)
		default:
			time.Sleep(defaultDuration)
		}

		api.ChangePin(zone, "off")

		// pause between zones to allow the sprinkler heads to retract
		time.Sleep(30 * time.Second)
	}

}
