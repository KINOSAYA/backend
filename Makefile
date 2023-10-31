
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker started!"

up_build:
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