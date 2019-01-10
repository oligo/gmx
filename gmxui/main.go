package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"

)

var connPool *GMXConnPool


func main() {
	connPool = NewGMXConnPool()
	r := mux.NewRouter()	

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/connect", connect)
	api.HandleFunc("/allkeys", allKeys)
	api.HandleFunc("/key", keyValue)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3305", nil))
}


