format:
	gofmt -w ./

build:
	go build -o api ./api/.

builddocker:
	docker build -t quizzie.cloud .

buildui:
	cd ui/ && yarn build

# Docker Compose commands
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

compose-logs:
	docker-compose logs -f

compose-build:
	docker-compose up --build -d

#===============================================
.DEFAULT_GOAL := start
.PHONY: start
start:
	go run ./api/.
