VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

start:
	docker-compose up -d

stop:
	docker-compose stop

build:
	docker-compose build

show-version:
	echo ${VERSION}

rebuild:
	go build && swag init && make fix-swagger-models

run-with-swagger:
	go build && swag init && make fix-swagger-models && go run main.go

# fix-swagger-models
# swaggo/swag has a bug that will prevent renaming of Models from "model.Account" ino "Account"
# we are going to fix this generation with this command

fix-swagger-models:
	./fix-swagger-files.sh

build-dev:
	swag init
	./fix-swagger-files.sh
	go build -v main.go