PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test-unit:
	@echo "Running tests..."
	go test $(PACKAGES)

