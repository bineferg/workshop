#########################
# makefile for workshop #
# backend part          #
#########################

build:
	go build ./cmd/...

workshop:
	env $(shell cat config/local) ./server

clean:
	rm server
