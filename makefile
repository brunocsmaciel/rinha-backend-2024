.PHONY: build run clean

APP_NAME := rinha-backend-2024
DOCKER_IMAGE := brunocsmaciel/rinha-backend-2024

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-push:
	@echo "Pushing Docker image..."
	@docker push ${DOCKER_IMAGE}

docker-run:
	@echo "Running Docker container..."
	@docker compose -f Docker-compose.yml up -d

clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)

