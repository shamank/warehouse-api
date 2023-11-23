include .env
export $(shell sed 's/=.*//' .env)
.PHONY:
.SILENT:


up:
	docker-compose build && docker-compose up -d

run:
	go run cmd/warehouse/main.go -config=./configs/local.yaml


tests:
	go test ./...