test:
	@echo "Starting unit tests..."
	go test --count=1 --short -v ./...
	@echo "Starting integration tests..."
	docker compose -f tests/docker-compose.yml up -d
	@echo "Waiting for services to start..."
	@sleep 30
	@echo "Running tests..."
	@trap 'docker compose -f tests/docker-compose.yml down -v' EXIT; go test -v -count=1 ./tests

lint:
	@echo "Linting code..."
	golangci-lint run --fix
	@echo "Linting completed"

money-tracker-deploy:
	helm upgrade --install moneytracker ./deploy/helm/moneytracker \
	   	-n default \
	   	-f ./deploy/helm/moneytracker/values.yaml \
	   	-f ./deploy/helm/moneytracker/prod-values-secrets.yaml


upload-images:
	docker buildx build --platform linux/amd64 -t moon1it/money_tracker:latest --push ./
