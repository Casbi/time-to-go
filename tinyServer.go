package main

import (
  "log"
  "net/http"

  "github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func timeToGoHandler(w http.ResponseWriter, req *http.Request) {

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("testing - writen by TimeToGoHandler\n"))

  log.Println(req.RequestURI)

  err := req.ParseForm()
  if err != nil {
    log.Panicf("error parsing form: %s", err)
  }
  log.Printf("from = %s to = %s departure time = %s", req.Form.Get("from"), req.Form.Get("to"), req.Form.Get("departureTime") )

  dirReq := &maps.DirectionsRequest {
		Origin:      req.Form.Get("from"),
		Destination: req.Form.Get("to"),
    DepartureTime: req.Form.Get("departureTime"),
    Mode: maps.TravelModeDriving,
  }

  mapsClient, err := maps.NewClient(maps.WithAPIKey("AIzaSyAFgh1pAQpS59mEwuViE2ExOw7M_W-2rzQ"))
  if err != nil {
    log.Fatalf("fatal error: %s", err)
  }

	routes, wayPoints, err := mapsClient.Directions(context.Background(), dirReq)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

  for key := range routes {
    w.Write([]byte(pretty.Sprint(routes[key])))
  }

  w.Write([]byte(pretty.Sprint(wayPoints)))

}

func main() {
  http.HandleFunc("/getTimeToGo", timeToGoHandler)
  http.ListenAndServe(":1313", nil)
}
