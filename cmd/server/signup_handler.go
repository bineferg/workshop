package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/workshop/lib/repository"
	"github.com/workshop/lib/workshop"
)

type SignupHandler struct {
	workshopRepo repository.WorkshopDB
}

type SignUpListResponse struct {
	WorkshopID string `json:"workshop_id"`
	SignUps    []SignUp `json:"sign_ups"`
}

type SignUp struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	Message string `json:"Message"`
}

func createSignup(su SignUp, id string) workshop.SignUp {
	return workshop.SignUp{
		WorkshopID: id,
		FirstName:  su.FirstName,
		LastName:   su.LastName,
		Email:      su.Email,
		Message: su.Message,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (h SignupHandler) GetSignupsByWorkshopID(w http.ResponseWriter, r *http.Request) error {
	urlVars := mux.Vars(r)
	workshopID := urlVars["workshop_id"]
	signups, err := h.workshopRepo.GetSignUpsByWorkshopID(workshopID)
	if err != nil {
		return err
	}
	var sResp []SignUp
	for _, s := range signups {
		sResp = append(sResp, SignUp{
			FirstName:        s.FirstName,
			LastName: 	s.LastName,
			Email: s.Email,
		})
	}
	resp := SignUpListResponse{SignUps: sResp, WorkshopID: workshopID}
	json.NewEncoder(w).Encode(resp)

	return nil
}

func (h SignupHandler) GetSignups(w http.ResponseWriter, r *http.Request) error {
	urlVars := mux.Vars(r)
	if urlVars["workshop_id"] != "all" {
		return h.GetSignupsByWorkshopID(w, r)
	}
	table, err := h.workshopRepo.GetAllSignUps()
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(table)
	return nil
}

func (h SignupHandler) CreateSignup(w http.ResponseWriter, r *http.Request) error {
	var signup SignUp
	urlVars := mux.Vars(r)
	workshopID := urlVars["workshop_id"]
	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	su := createSignup(signup, workshopID)
	if err := h.workshopRepo.SignUp(su); err != nil {
		return err
	}
	io.WriteString(w, "OK")
	return nil

}

func (h SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := h.GetSignups(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "POST":
		err := h.CreateSignup(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		http.Error(w, "not a valid request", http.StatusBadRequest)
	}

}
