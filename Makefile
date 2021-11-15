PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test-unit:
	@echo "Running tests..."
	go test $(PACKAGES)

test-e2e:
	@echo "Running e2e test..."
	./e2e.sh

env-start:
	@echo "Running an environment..."
	docker compose up -d
env-stop:
	@echo "Stop all containers"
	docker compose down

dashboard-start:
	@echo "Running dashboard"
	go build -o ./bin/dashboard github.com/valerykalashnikov/streaming-pipeline/cmd/dashboard
	./bin/dashboard

build-fileemitter:
	@echo "Building file-emitter"
	go build -o bin/file-emitter github.com/valerykalashnikov/streaming-pipeline/cmd/fileemitter

build-publisher:
	@echo "Building publisher"
	go build -o ./bin/publisher github.com/valerykalashnikov/streaming-pipeline/cmd/publisher

build-consumer:
	@echo "Building consumer"
	go build -o ./bin/consumer github.com/valerykalashnikov/streaming-pipeline/cmd/consumer

build-all: build-fileemitter build-publisher build-consumer