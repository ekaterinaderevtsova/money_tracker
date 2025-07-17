test:
	@echo "Starting unit tests..."
	go test --count=1 --short -v ./...
	@echo "Starting integration tests..."
	docker-compose -f tests/docker-compose.yml up -d
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Running tests..."
	@trap 'docker-compose -f tests/docker-compose.yml down -v' EXIT; go test -v -count=1 ./tests
