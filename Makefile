.PHONY: infra-up infra-down infra-logs db-create api-dev api-build api-test

db-create:
	createdb owncommerce 2>/dev/null || echo "Database 'owncommerce' sudah ada atau gagal dibuat — cek koneksi PostgreSQL lokal"

infra-up:
	docker compose -f infra/docker/docker-compose.dev.yml up -d

infra-down:
	docker compose -f infra/docker/docker-compose.dev.yml down

infra-logs:
	docker compose -f infra/docker/docker-compose.dev.yml logs -f

api-dev:
	$(MAKE) -C apps/api dev

api-build:
	$(MAKE) -C apps/api build

api-test:
	$(MAKE) -C apps/api test

dev: infra-up api-dev
