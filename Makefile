format:
	gofmt -w ./

build:
	go build -o api ./api/.

builddocker:
	docker build -t quizzie .

buildui:
	cd api/ui/ && yarn build

# install:
# 	go get -u "github.com/tidwall/sjson"

#===============================================
.DEFAULT_GOAL := start
.PHONY: start
start:
	go run ./api/.