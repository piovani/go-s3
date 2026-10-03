.PHONY: help build run docker-up docker-down clean

help:
	@echo "Comandos disponíveis:"
	@echo "  make run         - Sobe Docker Compose e executa a app"
	@echo "  make docker-up   - Sobe apenas o S3 Ninja"


run: docker-up build
	@echo "Iniciando servidor..."
	@echo "S3 Ninja rodando em http://localhost:9445"
	@echo "API rodando em http://localhost:8080"
	GIN_MODE=release ./bin/go-s3

docker-up:
	docker compose up -d
	@sleep 2
	@echo "✓ S3 Ninja iniciado"
