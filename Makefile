BINARY_NAME=backend
DSN=postgres://devuser:password@localhost/go_books?sslmode=disable

## build: will build the application
build:
	@echo "building the application"
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/${BINARY_NAME} ./cmd/api
	@echo "built!"

## run: builds and runs the app
run: build
	@echo "Starting..."
	@env DSN=${DSN} ./bin/${BINARY_NAME} &
	@echo "Started"

## start: alias to run
start: run

## stop: stops the server
stop:
	@echo "Stopping server..."
	@-pkill -SIGTERM -f "./bin/${BINARY_NAME}"
	@echo "Stopped!"

## restart will restart the server
restart: stop start


## docker-build: builds the docker-compose
docker-build: build
	@echo "building with docker-compose"
	docker-compose build
	@echo "built!"

## docker-run: builds and starts docker-compose
docker-run: docker-build
	@echo "Starting docker-compose..."
	docker-compose up -d
	@echo "Started"

## docker-start: alias to run
docker-start: docker-run

## docker-stop: stops docker-compose
docker-stop:
	@echo "Stopping docker-compose"
	docker-compose down
	@echo "Stopped!"

## docker-restart: restart docker-compose
docker-restart: docker-stop docker-start