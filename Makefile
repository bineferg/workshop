#########################
# makefile for workshop #
# backend part          #
#########################

VERSION=$(shell git rev-parse --short HEAD)
DESC=$(shell git log -1 --pretty=%B)
BUCKET_NAME="deployments"
PROJECT_NAME="workshop"
REGION="us-east-2"
ENV_NAME="production"


build:
	go build ./cmd/...

workshop:
	env $(shell cat config/local) ./server

clean:
	rm server
	rm $(VERSION).zip
