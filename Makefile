integration_tests:
	@echo "Starting integration tests..."
	docker-compose -f tests/docker-compose.yml up -d db_test redis_test
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Running tests..."
	@trap 'docker-compose -f tests/docker-compose.yml down -v' EXIT; \
	docker-compose -f tests/docker-compose.yml run --rm app_test go test -v -count=1 ./tests
