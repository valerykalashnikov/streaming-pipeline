PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test-unit:
	@echo "Running tests..."
	go test $(PACKAGES)

env-start:
	@echo "Running an environment..."
	docker compose up

build-fileemitter:
	@echo "Building file-emitter"
	go build -o bin/file-emitter github.com/valerykalashnikov/streaming-pipeline/cmd/fileemitter

build-publisher:
	@echo "Building publisher"
	go build -o ./bin/publisher github.com/valerykalashnikov/streaming-pipeline/cmd/publisher

build-consumer:
	@echo "Building consumer"
	go build -o ./bin/consumer github.com/valerykalashnikov/streaming-pipeline/cmd/consumer
consumer-setupdb:
setupdb:
	@echo "Setting up consumer database..."
	psql -U postgres stats < cmd/consumer/db/schema.sql