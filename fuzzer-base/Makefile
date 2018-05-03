all: dep clean build

dep:
	dep ensure

clean:
	@rm -rf ./bin

test:
	@echo "=============="
	@echo "== UNIT TESTS:"
	MAXFUZZ_ENV="test" go test ./internal/helpers -v -tags=unit
	@echo "=============="

build:
	go build -o ./bin/monitor ./cmd/monitor
	go build -o ./bin/go-monitor ./cmd/go-monitor
	go build -o ./bin/pre-sync ./cmd/pre-sync
	go build -o ./bin/reproduce ./cmd/reproduce

install:
	mv -t ${GOBIN} ./bin/monitor ./bin/pre-sync ./bin/reproduce ./bin/go-monitor

teardown:
	docker-compose down

deploy: teardown
	SYNC_DIR=$(shell pwd)/sync docker-compose build
	SYNC_DIR=$(shell pwd)/sync docker-compose up
