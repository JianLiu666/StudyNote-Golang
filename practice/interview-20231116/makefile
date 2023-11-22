.PHONY: help init up down build local

help:
	@echo "Usage: make [commands]\n"
	@echo "Comands:"
	@echo "  init        初始化建置環境 (docker volume, build image, etc.)"
	@echo "  up          透過 docker-compose 啟動所有服務"
	@echo "  down        關閉 docker-cpmpose 所有服務"
	@echo "  build       建置 Linked List Server 映像檔"
	@echo "  local       本地執行 Linked List Server"

init:	
	go mod download
	go mod tidy
	make build

up:
	docker compose -f deployment/docker-compose.yaml down -v
	docker compose -f deployment/docker-compose.yaml up -d
	docker ps -a

down:
	docker compose -f deployment/docker-compose.yaml down -v
	docker ps -a

build:
	docker build -t listserver:latest .

local:
	go run main.go -f config/env.yaml server
