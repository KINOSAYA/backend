BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker started!"

up_build:
# build_broker build_auth
	@echo Building docker containers...
	docker-compose up -d --build
	@echo Done!

down:
	@echo Stopping docker images...
	docker-compose down
	@echo Done!

swag:
	@echo Updating swagger documentation...
	cd ./broker-service/cmd/api&& swag init
#
#build_broker:
#	@echo Building broker service
#	cd ./broker-service&& set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BROKER_BINARY} ./cmd/api
#	@echo Done!
#
#build_auth:
#	@echo Building authentication service
#	cd ./authentication-service&& set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${AUTH_BINARY} ./cmd/api
#	@echo Done!
#
#restart: down up_build

