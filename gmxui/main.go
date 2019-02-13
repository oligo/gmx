package main

import (
	"net/http"
	"log"
	"flag"
	"time"
	"github.com/gorilla/mux"

)

const GMX_VERSION = 0

// var connPool *GMXConnPool
var (
	remoteProcess string
	localProcess int
)

var (
	metricsSyncer MetricsSyncer
	db *DataStore
)

func init() {
}

func init() {
	flag.StringVar(&remoteProcess, "remote", "127.0.0.1:9997", "remote IP and port of the GMX server")
	flag.IntVar(&localProcess, "pid", 0, "Local process PID of the GMX server")

	db = NewDataStore("./metrics.db")
}

func main() {
	flag.Parse()

	var address string
	if localProcess > 0 {
		address = findUnixSocket(localProcess)
	} else if remoteProcess != "" {
		address = remoteProcess
	} else {
		log.Fatalln("No GMX server address provided!")
	}

	gmxServerConn, err := dial(address)
	if err != nil {
		log.Fatalf("Cannot connect to GMX server: %s", address)
	}	

	metricsSyncer = NewMetricsSyncer(gmxServerConn, 10 * time.Second, db)
	metricsSyncer.Run()
	defer metricsSyncer.Stop()

	r := mux.NewRouter()	

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/connect", connect)
	api.HandleFunc("/allkeys", allKeys)
	api.HandleFunc("/key", keyValue)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3305", nil))
}


