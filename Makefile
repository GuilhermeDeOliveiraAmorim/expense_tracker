#!make
.PHONY: up down rmv swag

build:
	docker-compose up -d --build

up:
	docker-compose -f ./docker-compose.yml up -d

run:
	go run ./main.go

down:
	docker-compose -f ./docker-compose.yml down

rmv:
	docker volume rm $$(docker volume ls -q)

swag:
	swag init -g ./main.go -o ./api

test:
	go test -race -covermode=atomic -coverprofile=coverage.out $$(go list ./... | grep domain)
