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
)

func main() {
	port := flag.String("PORT", "8000", "listening-port")
	fixtures := flag.String("TEST_EVENTS", "fixtures/events.txt", "test-events")

	flag.Parse()
	if *port == "" {
		log.Fatal("couldnt parse port")
	}

	if *fixtures == "" {
		log.Fatal("could not get events path")
	}
	eventHandler := EventHandler{eventsList: *fixtures}

	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "OK") })
	router.Handle("/events", eventHandler)
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
