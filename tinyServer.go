package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("testing - writen by TimeToGoHandler\n"))

	log.Println(req.RequestURI)

	err := req.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %s", err)
	}
	log.Printf("from = %s & to = %s & departure time = %s", req.Form.Get("from"), req.Form.Get("to"), req.Form.Get("departureTime"))

	dirReq := &maps.DirectionsRequest{
		Origin:        req.Form.Get("from"),
		Destination:   req.Form.Get("to"),
		DepartureTime: req.Form.Get("departureTime"),
		Mode:          maps.TravelModeDriving,
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
			fmt.Println(routes[i].Summary)
			for j := range routes[i].Legs {
				fmt.Println(routes[i].Legs[j].DurationInTraffic)
			}
		}
	}

}

func main() {
	http.HandleFunc("/getTimeToGo", timeToGoHandler)
	log.Printf("Listening on localhost:1313/getTimeToGo")
	http.ListenAndServe(":1313", nil)
	// to test, use url:
	// http://localhost:1313/getTimeToGo?from=110 John Burke Drive, Porirua&to=91 Aitken Street, Wellington&departureTime=1530000000

}
