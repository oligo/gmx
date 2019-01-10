package main

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
	"time"
)

var sessionStore = sessions.NewCookieStore([]byte("session-key-123456"))
const sessionName = "gmxui-session"

var stopper chan bool

func connect(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Parse form failed", http.StatusInternalServerError)
		return
	}

	addr := r.Form.Get("addr")

	if len(addr) == 0 {
		log.Println(err)
		http.Error(w, "Parse form failed", http.StatusInternalServerError)
		return
	}

	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}

	if connPool.HasAddr(addr){
		http.Error(w, "Host already connected", http.StatusOK)
		return 
	}

	if _, err := connPool.Push(addr); err != nil {
		log.Println(err)
		http.Error(w, "Connecting host failed", http.StatusInternalServerError)
		return
	}

	session.Values["current-addr"] =  addr
	session.Save(r, w)

	syncer := NewMetricsSyncer(connPool.Get(addr), 5 * time.Second, db)
	syncer.retrieveKeys()
	stopper = syncer.Run()
}

func allKeys(w http.ResponseWriter, r *http.Request) {
	_, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}
	 
	keys, err := db.getKeys()
	if err != nil {
		http.Error(w, "Fetch all keys failed", 500)
		return
	}

	keysByte, err := json.Marshal(keys)
	if err != nil {
		http.Error(w, "Marshalling data failed", 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(keysByte)
	
}


func keyValue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parse form failed", http.StatusInternalServerError)
		return
	}

	key := r.Form.Get("key")
	if len(key) == 0 {
		http.Error(w, "No GMX key", http.StatusNotFound)
		return 
	}

	_, err = sessionStore.Get(r, "gmxui-session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}
	 

	metrics, err := db.getMetrics(key)

	if err != nil {
		http.Error(w, "Query metrics failed", 500)
		return
	}

	metricBytes, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, "Marshalling data failed", 500)
		return
	}

	
	w.Header().Set("content-type", "application/json")
	w.Write(metricBytes)
	
}
