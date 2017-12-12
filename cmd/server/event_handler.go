package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type EventHandler struct {
	eventsList string
}

type EventListResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	Description string `json:"description"`
	Time        string `json:"time"`
	Cost        int    `json:"cost"`
	Cap         int    `json:"person_cap"`
}

func createEvent(description string) Event {
	defaultTime := time.Now()
	defaultCost := 5
	defaultCap := 15
	return Event{Description: description, Time: fmt.Sprintf("%v", defaultTime), Cost: defaultCost, Cap: defaultCap}
}

func (h EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) error {
	file, err := os.Open(h.eventsList)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer file.Close()
	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		events = append(events, createEvent(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Printf(err.Error())
		return err
	}

	resp := EventListResponse{Events: events}
	json.NewEncoder(w).Encode(resp)

	return nil
}

func (h EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) error {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	log.Printf("event created %v", event)
	io.WriteString(w, "OK")
	return nil

}

func (h EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := h.GetEvents(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "POST":
		err := h.CreateEvent(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		http.Error(w, "not a valid request", http.StatusBadRequest)
	}

}
