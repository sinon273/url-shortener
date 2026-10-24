-include user/.env
export

.PHONY: proto-gen up down \
        run-user run-url run-gateway run-analytics run-all \
        build \
        test-user test-url test-gateway test-analytics test-all \
        test-int \
        lint-user lint-url lint-gateway lint-analytics lint-all \
        migrate-up migrate-down

proto-gen:
	cd shared/proto && buf generate


up:
	docker compose up -d

down:
	docker compose down


run-user:
	cd user && go run ./cmd

run-url:
	cd url && go run ./cmd

run-gateway:
	cd gateway && go run ./cmd

run-analytics:
	cd analytics && go run ./cmd

run-all:
	$(MAKE) -j4 run-user run-url run-gateway run-analytics

build:
	docker compose build


test-user:
	cd user && go test -race -v ./...

test-url:
	cd url && go test -race -v ./...

test-gateway:
	cd gateway && go test -race -v ./...

test-analytics:
	cd analytics && go test -race -v ./...

test-all: test-user test-url test-gateway test-analytics

test-int:
	cd gateway && go test -tags=integration -v ./test/... -timeout 20m


lint-user:
	cd user && golangci-lint run ./...

lint-url:
	cd url && golangci-lint run ./...

lint-gateway:
	cd gateway && golangci-lint run ./...

lint-analytics:
	cd analytics && golangci-lint run ./...

lint-all: lint-user lint-url lint-gateway lint-analytics


migrate-up:
	migrate -path user/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path user/migrations -database "$(DATABASE_URL)" down 1