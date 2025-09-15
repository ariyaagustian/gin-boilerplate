run:
	APP_PORT=8081 go run ./cmd/server

dev:
	reflex -r '\.go$$' -s -- sh -c 'go run ./cmd/server'

test:
	go test ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down
