package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/devcharmander/habit-tracker/database"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	router       *httprouter.Router
	resourcePath string
	habits       []*database.Habit
	dbURL        string // TODO: URLs should come from some sort of configuration
}

//NewServer creates a new instance of the server
func newServer() *server {
	return &server{
		router:       httprouter.New(),
		resourcePath: "/etc/timetable/resources/",
		dbURL:        "mongodb://localhost:27017",
	}
}

func (srv *server) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filePath := srv.resourcePath + "index.html"
	t, err := template.ParseFiles(filePath)
	if err != nil {
		log.Fatalf("unable to parse the template at location: %s Error %v", filePath, err)
	}
	t.Execute(w, srv)
}

func (srv *server) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	startDate := queryValues.Get("sd")
	endDate := queryValues.Get("ed")
	if startDate != "" && endDate != "" {
		client := database.NewClient("mongo", srv.dbURL)
		habits := client.Get(startDate, endDate)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(habits)
	} else {
		client := database.NewClient("mongo", srv.dbURL)
		habits := client.Get()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(habits)
	}
}

func (srv *server) Put(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data []*database.Habit
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Print("Error decoding the request. Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := database.NewClient("mongo", srv.dbURL)
	if status := client.Add(data); status {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *server) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data *database.Habit
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Print("Error decoding the request. Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := database.NewClient("mongo", srv.dbURL)
	if status := client.Update(data); status {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *server) Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data *database.Habit
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("error decoding body. Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := database.NewClient("mongo", srv.dbURL)
	if client.Remove(data) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
