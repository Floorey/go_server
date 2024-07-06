package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	r := mux.NewRouter()

	r.HandleFunc("/login", Login).Methods("POST")

	s := r.PathPrefix("/api").Subrouter()
	s.Use(Authenticate)
	s.HandleFunc("/vms", GetVMs).Methods("GET")
	s.HandleFunc("/vms", CreateVM).Methods("POST")

	http.Handle("/", r)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetVMs(w http.ResponseWriter, r *http.Request) {
	respChan := make(chan []VM)
	dbOps <- func() {
		var vms []VM
		db.Find(&vms)
		respChan <- vms
	}
	vms := <-respChan
	json.NewEncoder(w).Encode(vms)
}

func CreateVM(w http.ResponseWriter, r *http.Request) {
	respChan := make(chan VM)
	dbOps <- func() {
		var vm VM
		json.NewDecoder(r.Body).Decode(&vm)
		db.Create(&vm)
		respChan <- vm
	}
	createdVM := <-respChan
	json.NewEncoder(w).Encode(createdVM)
}
