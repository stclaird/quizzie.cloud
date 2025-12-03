format:
	gofmt -w ./

build:
	go build -o api ./api/.

builddocker:
	docker build -t quizzie.cloud .

buildui:
	cd ui/ && yarn build

# Testing commands
test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

integration-test:
	go test -tags=integration ./api/... -v

integration-test-docker:
	./scripts/integration-test.sh

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
