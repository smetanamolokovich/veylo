.PHONY: dev api web db seed down reset

# Start everything: db -> migrate -> seed -> api + web (parallel)
dev: db seed
	@$(MAKE) -j2 api web

# Start Postgres and run migrations
db:
	docker compose up -d db migrate
	@echo "Waiting for migrations..."
	@docker compose wait migrate

# Insert test data
seed: db
	go run ./cmd/seed

# Backend only (Go)
api:
	go run ./cmd/api

# Frontend only (Next.js)
web:
	cd web && npm run dev

# Stop Docker services
down:
	docker compose down

# Full reset: drop volumes and rebuild
reset:
	docker compose down -v
	@$(MAKE) dev
