#########################
# makefile for workshop #
# backend part          #
#########################

VERSION=$(shell git rev-parse HEAD)
DESC=$(shell git log -1 --pretty=%B)
BUCKET_NAME="deployments"
PROJECT_NAME="workshop"
REGION="us-east-2"
ENV_NAME="production"


build:
	go build ./cmd/...

workshop:
	env $(shell cat config/local) ./server
local-workshop:
	env $(shell cat config/local) MYSQL_DNS='root@/workshop?parseTime=true' ./server

deploy:
aws deploy create-deployment \
  --application-name Workshop \
  --deployment-config-name CodeDeployDefault.OneAtATime \
  --deployment-group-name Workshop-DepGrp \
  --description "Backend Website" \
  --github-location repository=repository,commitId=$(VERSION)

clean:
	rm server
	rm $(VERSION).zip
