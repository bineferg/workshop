package main

import (
	"encoding/json"
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
	URL string `json:"url"`
}

func (u UploadHandler) SignURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	req, _ := u.s3Cli.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})
	str, err := req.Presign(15 * time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := URLResponse{
		URL: str,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
