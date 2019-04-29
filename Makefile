all: build run

build:
	npm run --prefix frontend build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o calendar-server main.go
	docker build --tag=calendar .

run:
	docker-compose up calendar db cache

.PHONY: build run
