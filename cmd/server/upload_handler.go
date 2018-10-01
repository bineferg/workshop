package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

type UploadHandler struct {
	s3Cli  *s3.S3
	bucket string
}

type URLResponse struct {
	URL     string   `json:"url"`
}

const (
	workshopsFolder = "workshops"
	eventsFodler    = "events"
)

func (u UploadHandler) SignURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	folder := vars["folder"]
	req, _ := u.s3Cli.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(fmt.Sprintf("/%s/%s", folder, key)),
	})
	url, _, _ := req.PresignRequest(30 * time.Minute)
	resp := URLResponse{
		URL: url,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
