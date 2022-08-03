BINARY_NAME=backend
DSN=postgres://devuser:password@localhost/go_books?sslmode=disable

## build: will build the application
build:
	@echo "building the application"
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/api
	@echo "built!"

## run: builds and runs the app
run: build
	@echo "Starting..."
	@env DSN=${DSN} ./${BINARY_NAME} &
	@echo "Started"

## start: alias to run
start: run

## stop: stops the server
stop:
	@echo "Stopping server..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart will restart the server
restart: stop start