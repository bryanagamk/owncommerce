# Manual Testing — Phase 0

Panduan manual testing untuk fondasi OwnCommerce (Auth, Tenant, IAM, Audit) dengan setup **Docker** (PostgreSQL port **5433**, Redis port **6380**).

---

## Persiapan

### Terminal 1 — Infra Docker

```bash
cd /Volumes/SSD500/AgamTech/OwnCommerce
make infra-up
```

Cek container jalan:

```bash
docker ps --filter name=owncommerce
```

Harus ada:

| Container | Port host | Status |
|---|---|---|
| `owncommerce-postgres` | 5433 | healthy |
| `owncommerce-redis` | 6380 | healthy |

### Terminal 2 — API

```bash
cd /Volumes/SSD500/AgamTech/OwnCommerce
make api-dev
```

Tunggu log:

```
OwnCommerce API listening on :8080 (env=development)
```

### Terminal 3 — Testing (curl)

Opsional — pretty print JSON:

```bash
brew install jq
```

Pastikan `.env` sudah dikonfigurasi untuk Docker:

```env
DATABASE_URL=postgres://owncommerce:owncommerce@localhost:5433/owncommerce?sslmode=disable
REDIS_URL=redis://localhost:6380/0
```

---

## Test 1 — Health Check

```bash
curl -s http://localhost:8080/health | jq
curl -s http://localhost:8080/v1/health | jq
```

**Expected:**

```json
{
  "success": true,
  "data": {
    "status": "ok",
    "database": "ok",
    "service": "owncommerce-api"
  }
}
```

Jika `database: "down"` → jalankan `make infra-up` dan pastikan `DATABASE_URL` memakai port **5433**.

---

## Test 2 — Register Merchant

```bash
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@tokobunga.com",
    "password": "password123",
    "name": "Budi",
    "store_name": "Toko Bunga",
    "slug": "tokobunga"
  }' | jq
```

**Expected:**

- `success: true`
- `user` (email, name)
- `tenant` (slug: `tokobunga`)
- `tokens.access_token` & `tokens.refresh_token`
- `roles`: `["merchant_owner"]`
- `permissions`: berisi `product.view`, `order.manage`, dll.

**Simpan token untuk test berikutnya:**

```bash
export ACCESS_TOKEN="paste_access_token_disini"
export REFRESH_TOKEN="paste_refresh_token_disini"
```

---

## Test 3 — Register Duplikat (negative)

```bash
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@tokobunga.com",
    "password": "password123",
    "name": "Budi",
    "store_name": "Toko Lain",
    "slug": "tokolain"
  }' | jq
```

**Expected:** `success: false` — email sudah terdaftar.

---

## Test 4 — Login

```bash
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@tokobunga.com",
    "password": "password123"
  }' | jq
```

**Expected:** token baru + roles/permissions.

Update token:

```bash
export ACCESS_TOKEN="access_token_baru"
export REFRESH_TOKEN="refresh_token_baru"
```

---

## Test 5 — Login Password Salah (negative)

```bash
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "owner@tokobunga.com", "password": "salah"}' | jq
```

**Expected:** HTTP 401, pesan invalid credentials.

---

## Test 6 — Profil User (`/me`)

```bash
curl -s http://localhost:8080/v1/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

**Expected:** data user + `tenant_id` + `roles` + `permissions`.

Tanpa token (harus gagal):

```bash
curl -s http://localhost:8080/v1/me | jq
```

**Expected:** HTTP 401.

---

## Test 7 — Tenant Saat Ini

```bash
curl -s http://localhost:8080/v1/tenants/current \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

**Expected:**

- `tenant.name`: Toko Bunga
- `tenant.slug`: tokobunga
- `domains` berisi `tokobunga.localhost`

---

## Test 8 — Tenant Resolver (simulasi storefront)

```bash
# Dengan header tenant (local dev)
curl -s http://localhost:8080/v1/store \
  -H "X-Tenant-Slug: tokobunga" | jq

# Tanpa header (harus gagal)
curl -s http://localhost:8080/v1/store | jq
```

**Expected:**

- Dengan `X-Tenant-Slug` → data tenant + domains
- Tanpa header → HTTP 404, tenant not found

---

## Test 9 — Refresh Token

```bash
curl -s -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq
```

**Expected:** access token & refresh token **baru**.

Update token:

```bash
export ACCESS_TOKEN="access_token_baru"
export REFRESH_TOKEN="refresh_token_baru"
```

Coba refresh token **lama**:

```bash
curl -s -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "TOKEN_LAMA"}' | jq
```

**Expected:** HTTP 401 (token rotation — token lama sudah tidak valid).

---

## Test 10 — Logout

```bash
curl -s -X POST http://localhost:8080/v1/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq
```

**Expected:** `success: true`, message `logged out`.

Refresh token setelah logout → harus HTTP 401.

**Login ulang** sebelum Test 11–12:

```bash
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "owner@tokobunga.com", "password": "password123"}' | jq
```

Update `ACCESS_TOKEN` lagi.

---

## Test 11 — RBAC (Roles & Permissions)

```bash
curl -s http://localhost:8080/v1/iam/roles \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

curl -s http://localhost:8080/v1/iam/permissions \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

**Expected:** daftar role (`super_admin`, `merchant_owner`, dll.) dan permission (`product.view`, `order.manage`, dll.).

---

## Test 12 — Audit Logs

```bash
curl -s "http://localhost:8080/v1/audit-logs?limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

**Expected:** log aktivitas seperti:

- `auth.register`
- `auth.login`
- `auth.refresh`
- `auth.logout`

---

## Verifikasi Database (opsional)

```bash
docker exec -it owncommerce-postgres psql -U owncommerce -d owncommerce
```

```sql
SELECT id, name, slug, status FROM tenants;
SELECT domain, type, is_primary FROM tenant_domains;
SELECT id, email, name FROM users;

SELECT u.email, r.name AS role
FROM user_roles ur
JOIN users u ON u.id = ur.user_id
JOIN roles r ON r.id = ur.role_id;

SELECT action, entity_type, created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT 10;

\q
```

---

## Checklist

| # | Test | Expected | Pass? |
|---|---|---|---|
| 1 | Health check | `database: ok` | ☐ |
| 2 | Register | Tenant + user + token dibuat | ☐ |
| 3 | Register duplikat | Error email taken | ☐ |
| 4 | Login | Token baru | ☐ |
| 5 | Login salah password | HTTP 401 | ☐ |
| 6 | `GET /v1/me` | User + permissions | ☐ |
| 7 | `GET /v1/tenants/current` | Tenant + domain | ☐ |
| 8 | `GET /v1/store` + `X-Tenant-Slug` | Tenant resolved | ☐ |
| 9 | Refresh token | Token baru, token lama invalid | ☐ |
| 10 | Logout | Refresh token revoked | ☐ |
| 11 | IAM roles & permissions | Daftar role & permission | ☐ |
| 12 | Audit logs | Riwayat auth events | ☐ |

---

## Reset Data (mulai dari awal)

```bash
make infra-down
docker volume rm owncommerce_postgres_data owncommerce_redis_data
make infra-up
make api-dev
```

AutoMigrate + seed akan jalan ulang dari database kosong.

---

## Troubleshooting

| Masalah | Solusi |
|---|---|
| `connection refused` port 8080 | API belum jalan → `make api-dev` |
| `database: down` | `make infra-up`, tunggu ~10 detik |
| Docker error | Pastikan Docker Desktop running |
| `401` di semua endpoint | Token expired (15 menit) → login ulang |
| Port bentrok | Docker pakai port **5433/6380**, bukan 5432/6379 |
