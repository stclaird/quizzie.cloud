format:
	gofmt -w ./

build:
	go build -o api ./api/.

builddocker:
	docker build -t quizzie .

buildui:
	cd ui/ && yarn build

#===============================================
.DEFAULT_GOAL := start
.PHONY: start
start:
	go run ./api/.
