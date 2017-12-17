package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/workshop/lib/repository"
)

func main() {
	port := flag.String("PORT", "8000", "listening-port")
	dbDNS := flag.String("MYSQL_DNS", os.Getenv("MYSQL_DNS"), "dns string for workshop db")

	flag.Parse()
	if *port == "" {
		log.Fatal("couldnt parse port")
	}

	if *dbDNS == "" {
		log.Fatal("dns string not found")
	}

	var workshopDB repository.WorkshopDB
	workshopDB, err := repository.NewWorkshopDB(*dbDNS)
	if err != nil {
		log.Fatalf("%v", err)
	}

	eventHandler := EventHandler{workshopRepo: workshopDB}
	workshopHandler := WorkshopHandler{workshopRepo: workshopDB}
	signupHandler := SignupHandler{workshopRepo: workshopDB}	

	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "OK") })
	router.Handle("/events", eventHandler)
	router.Handle("/workshops", workshopHandler)
	router.Handle("/signup/{workshop_id}", signupHandler)
	log.Printf("listening on port %s", *port)
	go func() {
		if err := http.ListenAndServe(":"+*port, router); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals

}
