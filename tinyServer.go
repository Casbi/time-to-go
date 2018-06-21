package main

import (
  "log"
  "net/http"

  "github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// var mapsClient maps.Client

// func directionsAPICall(req maps.DirectionsRequest) ([]maps.Route, []maps.GeocodedWaypoint, error)  {
//
// }

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("testing - writen by TimeToGoHandler\n"))

  log.Println(req.RequestURI)

  err := req.ParseForm()
  if err != nil {
    log.Panicf("error parsing form: %s", err)
  }
  log.Printf("from = %s to = %s", req.Form.Get("from"), req.Form.Get("to") )

  dirReq := &maps.DirectionsRequest{
		Origin:      req.Form.Get("from"),
		Destination: req.Form.Get("to"),
	}

  mapsClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyAFgh1pAQpS59mEwuViE2ExOw7M_W-2rzQ"))
  if err != nil {
	 	log.Fatalf("fatal error: %s", err)
	}

	routes, wayPoints, err := mapsClient.Directions(context.Background(), dirReq)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

  w.Write([]byte(pretty.Sprint(routes, wayPoints)))

}

func main() {

  http.HandleFunc("/getTimeToGo", timeToGoHandler)
  http.ListenAndServe(":1313", nil)

}
