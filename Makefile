-include user/.env
export

.PHONY: proto-gen run test lint up down migrate-up migrate-down run-url test-url lint-url run-gateway test-int run-analytics

proto-gen:
	cd shared/proto && buf generate

run:
	cd user && go run ./cmd

test:
	cd user && go test -race -v ./...

lint:
	cd user && golangci-lint run ./...

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	migrate -path user/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path user/migrations -database "$(DATABASE_URL)" down 1

run-url:
	go run ./url/cmd

test-url:
	go test -race -v ./url/...

lint-url:
	golangci-lint run ./url/...

run-gateway:
	go run ./gateway/cmd

test-int:
	cd gateway && go test -tags=integration -v ./test/... -timeout 20m

run-analytics:
	go run ./analytics/cmd