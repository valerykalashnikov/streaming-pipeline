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