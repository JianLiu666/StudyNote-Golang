.PHONY: redis push-message receive-message

redis:
	docker compose -f ./docker-compose.yaml --compatibility up -d

push-message:
	go run ./cmd/push-message/*.go

receive-message:
	go run ./cmd/receive-message/*.go -addr :6379