.PHONY: help build run docker-up docker-down clean

help:
	@echo "Comandos disponíveis:"
	@echo "  make run         - Sobe Docker Compose e executa a app"
	@echo "  make build       - Compila o binário em bin/go-s3"
	@echo "  make docker-up   - Sobe apenas o S3 Ninja"
	@echo "  make docker-down - Derruba o Docker Compose"
	@echo "  make clean       - Remove o binário compilado"

build:
	go build -o bin/go-s3 ./cmd


run: docker-up build
	@echo "Iniciando servidor..."
	@echo "S3 Ninja rodando em http://localhost:9445"
	@echo "API rodando em http://localhost:8080"
	GIN_MODE=release ./bin/go-s3

docker-up:
	docker compose up -d
	@sleep 2
	@echo "✓ S3 Ninja iniciado"

docker-down:
	docker compose down

clean:
	rm -f bin/go-s3
