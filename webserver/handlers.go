package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// // map[jobName]jobReadyNumber
// var jobMember map[string]int

func writeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,HEAD,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type,Accept,Origin")
}

func setTaskReady(w http.ResponseWriter, r *http.Request) {
	writeHeader(w)
	lock.Lock()
	defer lock.Unlock()
	vars := mux.Vars(r)
	jobName := vars["jobName"]
	jobMember[jobName]++
	if err := json.NewEncoder(w).Encode(jobMember); err != nil {
		log.Printf("Add task %s fail.", jobName)
	}
}

func checkJobReady(w http.ResponseWriter, r *http.Request) {
	writeHeader(w)
	// lock.RLock()
	// defer lock.RUnlock()
	vars := mux.Vars(r)
	jobName := vars["jobName"]
	if err := json.NewEncoder(w).Encode(jobMember[jobName]); err != nil {
		log.Printf("Check job %s Member fail", jobName)
	}
}
