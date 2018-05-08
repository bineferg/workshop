package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/workshop/lib/repository"
)

func main() {
	port := flag.String("PORT", "8000", "listening-port")
	dbDNS := flag.String("MYSQL_DNS", os.Getenv("MYSQL_DNS"), "dns string for workshop db")
	awsKeyID := flag.String("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"), "aws access key for ses session")
	awsSecretKey := flag.String("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"), "aws access key for ses session")
	awsRegion := flag.String("AWS_REGION", os.Getenv("AWS_REGION"), "region for aws")
	uploadBucket := flag.String("S3_UPLOAD_BUCKET", os.Getenv("S3_UPLOAD_BUCKET"), "bucket for photos")

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
	if *uploadBucket == "" {
		log.Fatal("uploadBucket string not found")
	}

	var workshopDB repository.WorkshopDB
	workshopDB, err := repository.NewWorkshopDB(*dbDNS)
	if err != nil {
		log.Fatalf("%v", err)
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(*awsRegion),
		Credentials: credentials.NewStaticCredentials(*awsKeyID, *awsSecretKey, ""),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	sesSession := ses.New(awsSession)

	eventHandler := EventHandler{workshopRepo: workshopDB}
	workshopHandler := WorkshopHandler{workshopRepo: workshopDB}
	signupHandler := SignupHandler{workshopRepo: workshopDB}
	mailHandler := MailHandler{ses: sesSession}
	uploadHandler := UploadHandler{s3Cli: s3.New(awsSession), bucket: *uploadBucket}

	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "OK") })
	router.Handle("/events", eventHandler)
	router.Handle("/workshops", workshopHandler)
	router.Handle("/signup/{workshop_id}", signupHandler)
	router.Handle("/mail", mailHandler)
	router.HandleFunc("/upload/{folder}/{key}", uploadHandler.SignURL)
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
