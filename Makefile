default: vet test build

.PHONY: build
build:
	go build -o bin/terraform-provider-postgresreplication

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	cd dockercompose && docker-compose up wait_for
	go test -count 1 -v ./...
