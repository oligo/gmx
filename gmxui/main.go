package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

)

var connPool *GMXConnPool
var sessionStore = sessions.NewCookieStore([]byte("session-key-123456"))


func main() {
	connPool = NewGMXConnPool()
	r := mux.NewRouter()	

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/subscribe", subscribe)
	api.HandleFunc("/allkeys", allKeys)
	api.HandleFunc("/key", keyValue)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3305", nil))
}


func subscribe(w http.ResponseWriter, r *http.Request) {
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

	session, err := sessionStore.Get(r, "gmxui-session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}

	if connPool.HasAddr(addr){
		http.Error(w, "address connected", http.StatusOK)
		return 
	}

	if err := connPool.Push(addr); err != nil {
		log.Println(err)
		http.Error(w, "Parse form failed", http.StatusInternalServerError)
		return
	}

	session.Values["addr"] =  addr
	session.Save(r, w)
}

func allKeys(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "gmxui-session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}
	 
	addr := session.Values["addr"].(string)
	if len(addr) == 0 {
		http.Error(w, "No GMX process address", http.StatusInternalServerError)
		return 
	}

	conn := connPool.Get(addr)
	if conn == nil {
		http.Error(w, "No live connection", http.StatusInternalServerError)
		return
	}
	
	keys, err := json.Marshal(conn.FetchKeys())
	if err != nil {
		http.Error(w, "Marshalling data failed", 500)
		return
	}

	
	w.Header().Set("content-type", "application/json")
	w.Write(keys)
	
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

	session, err := sessionStore.Get(r, "gmxui-session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return 
	}
	 
	addr := session.Values["addr"].(string)
	if len(addr) == 0 {
		http.Error(w, "No GMX process address", http.StatusInternalServerError)
		return 
	}

	conn := connPool.Get(addr)
	if conn == nil {
		http.Error(w, "No live connection", http.StatusInternalServerError)
		return
	}
	
	keyValue, err := json.Marshal(conn.GetValues([]string{key}))
	if err != nil {
		http.Error(w, "Marshalling data failed", 500)
		return
	}

	
	w.Header().Set("content-type", "application/json")
	w.Write(keyValue)
	
}

