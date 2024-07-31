SERVICE_NAME=todo-service
ENGINE=main.go
BUILD_DIR=build
IMAGE_REPOSITORY=rizqi20/todo-service

.PHONY: build run build-run local docker-build docker-run docker-build-run compose-up compose-down check-if-present-env check-if-valid-env clean-proto proto-gen ssl-gen env-local env-docker migrate-create mock ssl-gen

CHECK_ENV := production|staging|local

proto-gen:	clean-proto
	@echo "Generating the stubs"
	./script/proto-gen.sh
	@echo "Success generate stubs. All stubs created are in the 'stubs/' directory"
	@echo "DO NOT EDIT ANY FILES STUBS!"

clean-proto:
	@echo "Delete all previous stubs ..."
	rm -rf stubs/*
	@echo "All stubs successfully deleted"

ssl-gen:
	@echo "Generating ssl configuration"
	./scripts/ssl-gen.sh
	@echo "Success generate ssl configuration. All SSL Configuration created in the 'cert/' directory"
	@echo "DO NOT EXPOSE SSL DIRECTORY!"

docker-build: check-if-present-env check-if-valid-env
	@echo "Building the app image....."
	docker build . --file Dockerfile --build-arg ENVIRONMENT=${ENV} --no-cache --tag ${SERVICE_NAME}
	@echo "Success build the app image"

docker-run:
	@echo "Running the docker container app....."
	docker run --name=${SERVICE_NAME} -p ${PORT}:${PORT} -d ${SERVICE_NAME}:latest
	@echo "Success run the container app"

docker-build-run: docker-build docker-run

docker-tag:
	@echo "Tagging the app image......"
	docker tag ${SERVICE_NAME} ${IMAGE_REPOSITORY}:latest 

docker-push:
	@echo "Pushing the image to the iamge repository....."
	docker push ${IMAGE_REPOSITORY}:latest
	@echo "Success push to image repository"

compose-up:
	@echo Starting docker compose
	docker compose -f docker-compose.yaml up -d --build

compose-down:
	@echo Stopping docker compose
	docker compose -f docker-compose.yaml down

build:
	@echo "Building app"
	go build -o ${BUILD_DIR}/${SERVICE_NAME} ${ENGINE}
	@echo "Success build app. Your app is ready to use in 'build/' directory."

run:
	@echo "Running app"
	./${BUILD_DIR}/${SERVICE_NAME}

build-run: build run

local:
	@go fmt ./...
	@go run .

mock: 
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go

env-local:
	cp config.local.yaml config.yaml

env-docker:
	cp .env.example .env

# Environment test/check
check-if-present-env:
	test $(ENV)
	
check-if-valid-env:
ifneq ($(ENV), $(filter $(ENV), production staging local))
	$(error "ENV=$(CHECK_ENV)" is required)
endif