package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/workshop/lib/repository"
	"github.com/workshop/lib/workshop"
)

type EventHandler struct {
	workshopRepo repository.WorkshopDB
}

type EventListResponse struct {
	Events []Event `json:"events"`
}
type Event struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StartTime   string  `json:"time"`
	Cost        float64 `json:"cost"`
	Location    string  `json:"location"`
}

func createEvent(e Event) (workshop.Event, error) {
	//TODO validate workshop entry later
	t, err := time.Parse(time.RFC3339, e.StartTime)
        if err !=nil {
                log.Printf("invalid date!")
        }
	return workshop.Event{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   t, //Deal e.th this later
		Cost:        float64(e.Cost),
		Location:    e.Location,
	}, nil
}

func (h EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) error {

	today := time.Now()
	events, err := h.workshopRepo.GetEventsAfterDate(today)
	if err != nil {
		return nil
	}
	var eResp []Event
	for _, e := range events {
		eResp = append(eResp, Event{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			StartTime:   fmt.Sprintf("%v", e.StartTime),
			Cost:        e.Cost,
			Location:    e.Location,
		})
	}
	resp := EventListResponse{Events: eResp}
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
	e, _ := createEvent(event)
	if err := h.workshopRepo.InsertEvent(e); err != nil {
		return err
	}
	log.Printf("event created %v", e)
	io.WriteString(w, "OK")
	return nil

}

func (h EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) error {
	var event workshop.Event
	err  := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return err
	}	
	defer r.Body.Close()
	if err = h.workshopRepo.UpdateEvent(event); err != nil {
		return err
	}
	io.WriteString(w, "OK")
	return nil

}

func (h EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) error {
	v := r.URL.Query()
	eventID := v.Get("event_id")
	if err := h.workshopRepo.DeleteEvent(eventID); err != nil {
		return err
	}
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
	case "PUT":
		err := h.UpdateEvent(w,r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "DELETE":
		err := h.DeleteEvent(w,r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "not a valid request", http.StatusBadRequest)
	}

}
