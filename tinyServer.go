package main

import (
	//"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	// "google.golang.org/genproto/googleapis/type/latlng"
)

var googleMapAPIKey = os.Getenv("googleMapAPIKey")

type trip struct {
	DepartureTime          	time.Time
	ArrivalTime            	time.Time
	TotalDurationInTraffic 	time.Duration
	OriginLatLng						maps.LatLng
	DestinationLatLng 			maps.LatLng
}

func getTripFromAPI(req *http.Request) trip {
	var thisTrip trip

	err := req.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %s", err)
	}

	dirReq := &maps.DirectionsRequest{
		Origin:        req.Form.Get("from"),
		Destination:   req.Form.Get("to"),
		DepartureTime: req.Form.Get("departureTime"),
		Mode:          maps.TravelModeDriving,
	}

	departureTime, err := strconv.ParseInt(dirReq.DepartureTime, 10, 64)
	if err != nil {
		log.Printf("error parsing departure time from string to int64: %s", err)
	} else {
		thisTrip.DepartureTime = time.Unix(departureTime, 0)
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey(googleMapAPIKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	routes, _, err := mapsClient.Directions(context.Background(), dirReq)

	if err != nil {
		log.Printf("fatal error: %s", err)
	} else {
		for i := range routes {
			for j := range routes[i].Legs {
				thisTrip.TotalDurationInTraffic += routes[i].Legs[j].DurationInTraffic
				if (i == 0 && j == 0) {
					thisTrip.OriginLatLng = routes[i].Legs[j].StartLocation
				}
				if (i == len(routes)-1 && j == len(routes[i].Legs)-1) {
					thisTrip.DestinationLatLng = routes[i].Legs[j].EndLocation
				}
			}
		}
		thisTrip.ArrivalTime = thisTrip.DepartureTime
		thisTrip.ArrivalTime = thisTrip.ArrivalTime.Add(thisTrip.TotalDurationInTraffic)
	}
	return thisTrip
}

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	thisTripJSON, err := json.Marshal(getTripFromAPI(req))
	if err != nil {
		log.Printf("Error marshalling thisTrip: %v\n", err)
	}
	w.Write(thisTripJSON)

}

func main() {
	http.HandleFunc("/getTimeToGo", timeToGoHandler)
	log.Printf("Listening on localhost:1313/getTimeToGo")
	http.ListenAndServe(":1313", nil)
	// to test, use url:
	// http://localhost:1313/getTimeToGo?from=110 John Burke Drive, Porirua&to=91 Aitken Street, Wellington&departureTime=1600000000

}
