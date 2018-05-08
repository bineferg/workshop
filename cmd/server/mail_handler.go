	package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MailHandler struct {
	ses *ses.SES
}

type MailRequest struct {
	Email     string `json:"Email"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Message   string `json"Message"`
}

func (h MailHandler) SendMail(w http.ResponseWriter, r *http.Request) error {
	recipient := "workshop-on-forster@gmx.de"
	from := "sacre.kool@gmail.com"
	var mr MailRequest
	err := json.NewDecoder(r.Body).Decode(&mr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	emailBody := fmt.Sprintf("Message: %s\r\n\tSent By: %s %s %s", mr.Message, mr.FirstName, mr.LastName, mr.Email)
	sesEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(emailBody)},
			},
			Subject: &ses.Content{
				Data: aws.String("Workshop-Web ContactUs"),
			},
		},
		Source: aws.String(from),
		ReplyToAddresses: []*string{
			aws.String(from),
		},
	}
	_, err = h.ses.SendEmail(sesEmailInput)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (h MailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := h.SendMail(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		http.Error(w, "not a valid request", http.StatusBadRequest)
	}

}