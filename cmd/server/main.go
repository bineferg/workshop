package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/workshop/lib/repository"
)

func main() {
	port := flag.String("PORT", "8000", "listening-port")
	dbDNS := flag.String("MYSQL_DNS", os.Getenv("MYSQL_DNS"), "dns string for workshop db")
	awsKeyID := flag.String("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"), "aws access key for ses session")
	awsSecretKey := flag.String("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"), "aws access key for ses session")
	awsRegion := flag.String("AWS_REGION", os.Getenv("AWS_REGION"), "region for aws")

	flag.Parse()
	if *port == "" {
		log.Fatal("couldnt parse port")
	}

	if *dbDNS == "" {
		log.Fatal("dns string not found")
	}
	if *awsKeyID == "" {
		log.Fatal("awsKeyID string not found")
	}
	if *awsSecretKey == "" {
		log.Fatal("awsSecretKey string not found")
	}
	if *awsRegion == "" {
		log.Fatal("awsRegion string not found")
	}

	var workshopDB repository.WorkshopDB
	workshopDB, err := repository.NewWorkshopDB(*dbDNS)
	if err != nil {
		log.Fatalf("%v", err)
	}
	awsSession := session.New(&aws.Config{
		Region:      aws.String(*awsRegion),
		Credentials: credentials.NewStaticCredentials(*awsKeyID, *awsSecretKey, ""),
	})

	sesSession := ses.New(awsSession)

	eventHandler := EventHandler{workshopRepo: workshopDB}
	workshopHandler := WorkshopHandler{workshopRepo: workshopDB}
	signupHandler := SignupHandler{workshopRepo: workshopDB}
	mailHandler := MailHandler{ses: sesSession}

	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "OK") })
	router.Handle("/events", eventHandler)
	router.Handle("/workshops", workshopHandler)
	router.Handle("/signup/{workshop_id}", signupHandler)
	router.Handle("/mail", mailHandler)
	log.Printf("listening on port %s", *port)
	go func() {
		if err := http.ListenAndServe(":"+*port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(router)); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals

}
