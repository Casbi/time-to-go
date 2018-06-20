package main

import (
  "log"
  "net/http"
)

func getTimeToGoHandler(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("testing - writen by getTimeToGoHandler"))
  log.Printf("finished handler")
}

func main() {
  http.HandleFunc("/getTimeToGo", getTimeToGoHandler)
  http.ListenAndServe(":1313", nil)

}
