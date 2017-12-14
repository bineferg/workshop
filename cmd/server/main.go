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
	fixtures := flag.String("TEST_EVENTS", "fixtures/events.txt", "test-events")
	dbDNS := flag.String("MYSQL_DNS", os.Getenv("MYSQL_DNS"), "dns string for workshop db")

	flag.Parse()
	if *port == "" {
		log.Fatal("couldnt parse port")
	}

	if *fixtures == "" {
		log.Fatal("could not get events path")
	}

	if *dbDNS == "" {
		log.Fatal("dns string not found")
	}

	var workshopDB repository.WorkshopDB
	workshopDB, err := repository.NewWorkshopDB(*dbDNS)
	if err != nil {
		log.Fatalf("%v", err)
	}

	eventHandler := EventHandler{eventsList: *fixtures}
	workshopHandler := WorkshopHandler{workshopRepo: workshopDB}

	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "OK") })
	router.Handle("/events", eventHandler)
	router.Handle("/workshops", workshopHandler)
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
