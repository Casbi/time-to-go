package main

import (
	//"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type trip struct {
	DepartureTime          time.Time
	ArrivalTime            time.Time
	TotalDurationInTraffic time.Duration
}

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {
	var thisTrip trip

	w.WriteHeader(http.StatusOK)
	//fmt.Fprintln(w, "testing - writen by TimeToGoHandler")

	log.Println(req.RequestURI)

	err := req.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %s", err)
	}
	//fmt.Fprintln(w, "from: " + req.Form.Get("from") + " to: " + req.Form.Get("to") + " departure time: " + req.Form.Get("departureTime"))

	dirReq := &maps.DirectionsRequest{
		Origin:        req.Form.Get("from"),
		Destination:   req.Form.Get("to"),
		DepartureTime: req.Form.Get("departureTime"),
		// ArrivalTime: 	 req.Form.Get("arrivalTime"),
		Mode:          maps.TravelModeDriving,
	}

	departureTime, err := strconv.ParseInt(dirReq.DepartureTime, 10, 64)
	if err != nil {
		log.Printf("error parsing departure time from string to int64: %s", err)
	} else {
		thisTrip.DepartureTime = time.Unix(departureTime, 0)
		//fmt.Fprintf(w, "Trip departureTime: %v\n", thisTrip.DepartureTime)
	}

	// initiate arrival time, will add leg(s) durations to it later to get the final arrival time
	thisTrip.ArrivalTime = thisTrip.DepartureTime

	mapsClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyAFgh1pAQpS59mEwuViE2ExOw7M_W-2rzQ"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	routes, _, err := mapsClient.Directions(context.Background(), dirReq)

	if err != nil {
		log.Printf("fatal error: %s", err)
	} else {
		for i := range routes {
			//fmt.Println(w, routes[i].Summary)
			for j := range routes[i].Legs {
				thisTrip.TotalDurationInTraffic += routes[i].Legs[j].DurationInTraffic

				//fmt.Fprintf(w, "Leg DurationInTraffic: %v\n", routes[i].Legs[j].DurationInTraffic)
				//fmt.Fprintf(w, "Leg Duration: %v\n", routes[i].Legs[j].Duration)
			}
		}

		thisTrip.ArrivalTime = thisTrip.ArrivalTime.Add(thisTrip.TotalDurationInTraffic)
		//fmt.Fprintf(w, "Trip arrivalTime: %v\n", thisTrip.ArrivalTime)

		//fmt.Fprintf(w, "Trip totalDurationInTraffic: %v\n", thisTrip.TotalDurationInTraffic)
	}

	//routesJSON, _ := json.Marshal(routes)
	//w.Write([]byte(routesJSON))

	thisTripJSON, err := json.Marshal(thisTrip)
	if err != nil {
		//fmt.Fprintf(w, "Error marshalling thisTrip: %v\n", err)
	}
	w.Write(thisTripJSON)

}

func main() {
	http.HandleFunc("/getTimeToGo", timeToGoHandler)
	log.Printf("Listening on localhost:1313/getTimeToGo")
	http.ListenAndServe(":1313", nil)
	// to test, use url:
	// http://localhost:1313/getTimeToGo?from=110 John Burke Drive, Porirua&to=91 Aitken Street, Wellington&departureTime=1530000000

}
