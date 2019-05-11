all: build run

backend: build-backend run

build:
	npm run --prefix frontend build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o calendar-server main.go
	docker build --tag=calendar .

build-backend:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o calendar-server main.go
	docker build --tag=calendar .

run:
	docker-compose up calendar db cache

clean:
	docker container stop calendar_calendar_1
	docker container stop calendar_db_1
	docker container stop calendar_cache_1
	docker container rm calendar_calendar_1
	docker container rm calendar_db_1
	docker container rm calendar_cache_1

.PHONY: build run backend
