package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type trip struct {
	departureTime          time.Time
	arrivalTime            time.Time
	totalDurationInTraffic time.Duration
	totalDuration          time.Duration
}

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {
	var thisTrip trip

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "testing - writen by TimeToGoHandler")

	log.Println(req.RequestURI)

	err := req.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %s", err)
	}
	fmt.Fprintln(w, "from: " + req.Form.Get("from") + " to: " + req.Form.Get("to") + " departure time: " + req.Form.Get("departureTime"))

	dirReq := &maps.DirectionsRequest{
		Origin:        req.Form.Get("from"),
		Destination:   req.Form.Get("to"),
		DepartureTime: req.Form.Get("departureTime"),
		ArrivalTime: 	 req.Form.Get("arrivalTime"),
		Mode:          maps.TravelModeDriving,
	}

	departureTime, err := strconv.ParseInt(dirReq.DepartureTime, 10, 64)
	if err != nil {
		log.Printf("error parsing departure time from string to int64: %s", err)
	} else {
		thisTrip.departureTime = time.Unix(departureTime, 0)
		fmt.Fprintf(w, "Trip departureTime: %v\n", thisTrip.departureTime)
	}

	arrivalTime, err := strconv.ParseInt(dirReq.ArrivalTime, 10, 64)
	if err != nil {
		log.Printf("error parsing arrival time from string to int64: %s", err)
	} else {
		thisTrip.arrivalTime = time.Unix(arrivalTime, 0)
		fmt.Fprintf(w, "Trip arrivalTime: %v\n", thisTrip.arrivalTime)
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyAFgh1pAQpS59mEwuViE2ExOw7M_W-2rzQ"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	routes, _, err := mapsClient.Directions(context.Background(), dirReq)

	if err != nil {
		log.Printf("fatal error: %s", err)
	} else {
		for i := range routes {
			fmt.Println(w, routes[i].Summary)
			for j := range routes[i].Legs {
				if j == len(routes[i].Legs)-1 {
					thisTrip.arrivalTime = routes[i].Legs[j].ArrivalTime
					fmt.Fprintf(w, "ArrivalTime on last leg: %v\n", routes[i].Legs[j].ArrivalTime)
					fmt.Fprintf(w, "Trip arrivalTime: %v\n", thisTrip.arrivalTime)
				}
				thisTrip.totalDurationInTraffic += routes[i].Legs[j].DurationInTraffic
				thisTrip.totalDuration += routes[i].Legs[j].Duration

				fmt.Fprintf(w, "Leg DurationInTraffic: %v\n", routes[i].Legs[j].DurationInTraffic)
				fmt.Fprintf(w, "Leg Duration: %v\n", routes[i].Legs[j].Duration)

				fmt.Fprintf(w, "Trip totalDurationInTraffic: %v\n", thisTrip.totalDurationInTraffic)
				fmt.Fprintf(w, "Trip totalDuration: %v\n", thisTrip.totalDuration)

			}
		}
	}

	routesJSON, _ := json.Marshal(routes)
	w.Write([]byte(routesJSON))

}

func main() {
	http.HandleFunc("/getTimeToGo", timeToGoHandler)
	log.Printf("Listening on localhost:1313/getTimeToGo")
	http.ListenAndServe(":1313", nil)
	// to test, use url:
	// http://localhost:1313/getTimeToGo?from=110 John Burke Drive, Porirua&to=91 Aitken Street, Wellington&departureTime=1530000000

}
