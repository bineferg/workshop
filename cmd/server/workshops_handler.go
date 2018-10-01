package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/workshop/lib/repository"
	"github.com/workshop/lib/workshop"
)

type WorkshopHandler struct {
	workshopRepo repository.WorkshopDB
}

type WorkshopListResponse struct {
	Workshops []workshop.Workshop `json:"workshops"`
}

type Workshop struct {
	WorkshopID  string `json:"workshop_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Duration    string `json:"duration"`
	Caption     string `json:"caption"`
	Cost        string `json:"cost"`
	Cap         int    `json:"cap"`
	IsFull      bool   `json:"isFull"`
	Location    string `json:"location"`
	Level       string `json:"level"`
}

func createWorkshop(w Workshop) (workshop.Workshop, error) {
	//TODO validate workshop entry later
	return workshop.Workshop{
		WorkshopID:  w.WorkshopID,
		Name:        w.Name,
		Caption:     w.Caption,
		Description: w.Description,
		Time:        w.Time,
		Cost:        w.Cost,
		Cap:         w.Cap,
		Location:    w.Location,
		Level:       w.Level,
	}, nil

}

// Get all workshops that start after TODAY
func (h WorkshopHandler) GetWorkshops(w http.ResponseWriter, r *http.Request) error {
	workshops, err := h.workshopRepo.GetWorkshops()
	if err != nil {
		return err
	}
	var wsResp []Workshop
	for _, w := range workshops {
		wsResp = append(wsResp, Workshop{
			Name:        w.Name,
			Description: w.Description,
			Time:        w.Time,
			Duration:    "", //TODO deal with this later
			Cost:        w.Cost,
			Cap:         w.Cap,
			Caption:     w.Caption,
			Location:    w.Location,
			Level:       w.Level,
		})
	}
	resp := WorkshopListResponse{Workshops: workshops}
	json.NewEncoder(w).Encode(resp)
	return nil
}

func (h WorkshopHandler) CreateWorkshop(w http.ResponseWriter, r *http.Request) error {
	var workshop Workshop
	err := json.NewDecoder(r.Body).Decode(&workshop)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	ws, _ := createWorkshop(workshop)
	if err := h.workshopRepo.InsertWorkshop(ws); err != nil {
		return err
	}
	log.Printf("event created %v", ws)
	io.WriteString(w, "OK")
	return nil

}

func (h WorkshopHandler) UpdateWorkshop(w http.ResponseWriter, r *http.Request) error {
	var ws workshop.Workshop
	err := json.NewDecoder(r.Body).Decode(&ws)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err = h.workshopRepo.UpdateWorkshop(ws); err != nil {
		return err
	}
	io.WriteString(w, "OK")
	return nil

}

func (h WorkshopHandler) DeleteWorkshop(w http.ResponseWriter, r *http.Request) error {
	v := r.URL.Query()
	workshopID := v.Get("workshop_id")
	if err := h.workshopRepo.DeleteWorkshop(workshopID); err != nil {
		return err
	}
	return nil

}

func (h WorkshopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := h.GetWorkshops(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "POST":
		err := h.CreateWorkshop(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "PUT":
		err := h.UpdateWorkshop(w, r) //h.UpdateWorkshop(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "DELETE":
		err := h.DeleteWorkshop(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "not a valid request", http.StatusBadRequest)
		return
	}

}
