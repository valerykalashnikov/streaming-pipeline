PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test-unit:
	@echo "Running tests..."
	go test $(PACKAGES)

build-fileemitter:
	@echo "Building file-emitter"
	go build -o bin/file-emitter github.com/valerykalashnikov/streaming-pipeline/cmd/fileemitter