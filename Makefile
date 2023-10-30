
up_build:
	@echo Starting docker containers...
	docker-compose up -d --build
	@echo Done!

down:
	@echo Stopping docker images...
	docker-compose down
	@echo Done!