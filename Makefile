.PHONY: build
build:
	go build -v ./cmd/app/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

MOCKGEN_BIN:=$(LOCAL_BIN)/mockgen
MOCKGEN_TAG:=1.5.0
LOCAL_BIN := $(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG:=1.39.0
INTERNAL_PATH := $(CURDIR)/internal
RELATIVE_INTERNAL_PATH := ./internal
GOOSE_BIN:=$(LOCAL_BIN)/goose
GOOSE_TAG:=2.6.0
DATABASE_URI=postgres://admin:password@localhost:6432/gopher_market?sslmode=disable
RUN_ADDRESS=localhost:10000
ACCRUAL_SYSTEM_ADDRESS=http://localhost:8080
export PATH := $(PATH):$(LOCAL_BIN)
export GO111MODULE=on

swagger:
	swag init -g ./cmd/app/main.go -o ./doc

install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	@echo "Downloading golangci-lint v$(GOLANGCI_TAG)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v$(GOLANGCI_TAG)
endif

lint:
	@echo "Running golangci-lint..."
	@PATH=$(PATH); golangci-lint run

docker-compose-up:
	@echo "Starting server in docker container..."
	@docker-compose -p infra -f build/dev/docker-compose.yml up -d

docker-compose-down:
	@echo "Down server in docker container..."
	@docker-compose -p infra -f build/dev/docker-compose.yml down -v

docker-compose-ps:
	@docker-compose -p infra -f build/dev/docker-compose.yml ps

docker-compose-logs:
	@docker-compose -p infra -f build/dev/docker-compose.yml logs -f

go-vendor:
	@go mod vendor

install-mockgen:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@echo "Downloading mockgen $(MOCKGEN_TAG)"
	tmp=$$(mktemp -d) && GOPATH=$$tmp go install github.com/golang/mock/mockgen@v$(MOCKGEN_TAG) && \
	mv $$tmp/bin/mockgen $(LOCAL_BIN)/mockgen && rm -rf $$tmp
endif

generate-mock:
	@PATH=$(PATH); mockgen -source=$(RELATIVE_INTERNAL_PATH)/handlers/balance.go -destination=$(RELATIVE_INTERNAL_PATH)/handlers/mock/mock.go -package=mock

POSTGRES_SETUP_TEST:=user=admin password=password dbname=gopher_market host=localhost port=6432 sslmode=disable

install-goose:
	go get -u github.com/pressly/goose/v3/cmd/goose

test-migrations-up:
	@goose -dir "$(INTERNAL_PATH)/storage/migrations" postgres "${POSTGRES_SETUP_TEST}" up

test-migrations-down:
	@goose -dir "$(INTERNAL_PATH)/storage/migrations" postgres "${POSTGRES_SETUP_TEST}" down

run:
	go run cmd/gophermart/main.go -d='postgres://admin:password@localhost:6432/gopher_market?sslmode=disable' -a=localhost:8081 -r=localhost:8080
