.PHONY: build up down clean status logs test jaeger

# Build all services
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Clean up containers and images
clean:
	docker-compose down -v --rmi all --remove-orphans

# Show status of all services
status:
	docker-compose ps

# Show logs
logs:
	docker-compose logs -f

# Test the services
test:
	@echo "Testing Book API..."
	curl -s "http://localhost:8080/v1/books/1" | jq .
	@echo "Testing Author API..."
	curl -s "http://localhost:8080/v1/authors/1" | jq .

# Open Jaeger UI
jaeger:
	@echo "Jaeger UI available at: http://localhost:16687" 
