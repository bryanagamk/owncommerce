# OwnCommerce

Platform WebCommerce SaaS multi-tenant untuk UMKM Indonesia.

## Prasyarat

- Go 1.23+
- Docker & Docker Compose (PostgreSQL + Redis)
- Node.js 20+ (untuk monorepo frontend nanti)

## Quick Start (Local Dev)

### 1. Environment

```bash
cp .env.example .env
```

### 2. Jalankan PostgreSQL & Redis (Docker)

Port **5433** (Postgres) dan **6380** (Redis) di host — tidak bentrok dengan PostgreSQL/Redis lokal Mac.

```bash
make infra-up
```

`.env` sudah dikonfigurasi untuk Docker secara default.

### 3. Jalankan API

```bash
make api-dev
```

API tersedia di `http://localhost:8080`

### 4. Health Check

```bash
curl http://localhost:8080/health
curl http://localhost:8080/v1/health
```

## API Endpoints (Phase 0)

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| GET | `/health` | - | Health check |
| GET | `/v1/health` | - | Health check v1 |
| POST | `/v1/auth/register` | - | Daftar merchant + buat tenant |
| POST | `/v1/auth/login` | - | Login |
| POST | `/v1/auth/refresh` | - | Refresh token |
| POST | `/v1/auth/logout` | Bearer | Logout |
| GET | `/v1/me` | Bearer | Profil user |
| GET | `/v1/tenants/current` | Bearer | Tenant merchant saat ini |
| GET | `/v1/store` | - | Info toko (butuh tenant context) |
| GET | `/v1/audit-logs` | Bearer + `audit.view` | Audit log tenant |
| GET | `/v1/iam/roles` | Bearer + `staff.manage` | Daftar role |
| GET | `/v1/iam/permissions` | Bearer + `staff.manage` | Daftar permission |

## Local Dev: Tenant Context

Karena belum pakai domain production, gunakan header untuk resolve tenant:

```bash
# Register merchant
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@tokobunga.com",
    "password": "password123",
    "name": "Budi",
    "store_name": "Toko Bunga",
    "slug": "tokobunga"
  }'

# Resolve tenant store (local dev)
curl http://localhost:8080/v1/store \
  -H "X-Tenant-Slug: tokobunga"
```

## Struktur Proyek

```
apps/api/          → Backend Go (Fiber + GORM)
packages/          → Shared packages (placeholder)
infra/docker/      → Docker Compose dev
storage/           → Local file storage
docs/              → Dokumentasi teknis & testing
prompt/            → Dokumentasi produk & rencana
```

## Manual Testing

Lihat [docs/testing-fase0.md](./docs/testing-fase0.md) untuk panduan testing Phase 0.

## Perintah Berguna

```bash
make infra-up        # Start PostgreSQL & Redis via Docker (port 5433/6380)
make infra-down      # Stop infra Docker
make api-dev         # Run API
make api-build       # Build API binary
make api-test        # Run Go tests
```

## Troubleshooting

| Error | Penyebab | Solusi |
|---|---|---|
| `port 5432/6379 already in use` | Port default bentrok dengan service lokal | Docker sudah pakai port **5433/6380** — pastikan `make infra-up` |
| `database: down` di health check | Container belum jalan / URL salah | `make infra-up`, cek `DATABASE_URL` pakai port **5433** |
| `connection refused` | Docker belum jalan | Start Docker Desktop, lalu `make infra-up` |
