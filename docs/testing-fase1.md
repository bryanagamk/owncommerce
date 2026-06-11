# Manual Testing — Phase 1 (MVP Core Commerce)

Panduan testing API Phase 1: produk, katalog, cart, checkout, Midtrans, order management.

> **Prasyarat:** Phase 0 sudah jalan (`make infra-up`, `make api-dev`).  
> **Header wajib storefront:** `X-Tenant-Slug: <slug-toko>`

---

## Setup

```bash
# Pastikan Midtrans sandbox key di .env
MIDTRANS_SERVER_KEY=SB-Mid-server-...
MIDTRANS_CLIENT_KEY=SB-Mid-client-...
MIDTRANS_IS_PRODUCTION=false
```

Restart API setelah update `.env`.

---

## Flow Testing Lengkap

### 1. Register Merchant (jika belum ada)

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

```bash
export MERCHANT_TOKEN="<access_token>"
export TENANT_SLUG="tokobunga"
```

---

### 2. Merchant — Update Store Settings

```bash
curl -s -X PATCH http://localhost:8080/v1/merchant/store/settings \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Toko bunga segar",
    "contact_email": "hello@tokobunga.com",
    "contact_phone": "08123456789",
    "address": "Jl. Melati No. 1",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345"
  }' | jq
```

---

### 3. Merchant — CRUD Kategori

```bash
curl -s -X POST http://localhost:8080/v1/merchant/categories \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Bunga Segar", "slug": "bunga-segar", "is_active": true}' | jq
```

Simpan `category_id`, lalu:

```bash
curl -s http://localhost:8080/v1/merchant/categories \
  -H "Authorization: Bearer $MERCHANT_TOKEN" | jq
```

---

### 4. Merchant — CRUD Produk

```bash
curl -s -X POST http://localhost:8080/v1/merchant/products \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "<category_id>",
    "name": "Buket Mawar Merah",
    "slug": "buket-mawar-merah",
    "description": "12 batang mawar merah",
    "sku": "BMR-001",
    "price": 250000,
    "status": "active",
    "is_featured": true,
    "quantity": 50,
    "low_stock_threshold": 5
  }' | jq
```

Simpan `product_id`, lalu upload gambar:

```bash
curl -s -X POST "http://localhost:8080/v1/merchant/products/<product_id>/images" \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -F "image=@/path/to/image.jpg" | jq
```

---

### 5. Storefront — Katalog Publik

```bash
curl -s http://localhost:8080/v1/storefront/home \
  -H "X-Tenant-Slug: $TENANT_SLUG" | jq

curl -s http://localhost:8080/v1/storefront/products \
  -H "X-Tenant-Slug: $TENANT_SLUG" | jq

curl -s http://localhost:8080/v1/storefront/products/buket-mawar-merah \
  -H "X-Tenant-Slug: $TENANT_SLUG" | jq

curl -s "http://localhost:8080/v1/storefront/products?q=mawar" \
  -H "X-Tenant-Slug: $TENANT_SLUG" | jq
```

---

### 6. Customer — Register & Login

```bash
curl -s -X POST http://localhost:8080/v1/storefront/auth/register \
  -H "X-Tenant-Slug: $TENANT_SLUG" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@test.com",
    "password": "password123",
    "name": "Siti",
    "phone": "08129999999"
  }' | jq
```

```bash
export CUSTOMER_TOKEN="<customer_access_token>"
export CART_SESSION="guest-$(uuidgen)"   # untuk guest cart
```

---

### 7. Cart & Checkout

**Guest cart** (pakai `X-Cart-Session`):

```bash
curl -s -X POST http://localhost:8080/v1/storefront/cart/items \
  -H "X-Tenant-Slug: $TENANT_SLUG" \
  -H "X-Cart-Session: $CART_SESSION" \
  -H "Content-Type: application/json" \
  -d '{"product_id": "<product_id>", "quantity": 2}' | jq

curl -s http://localhost:8080/v1/storefront/cart \
  -H "X-Tenant-Slug: $TENANT_SLUG" \
  -H "X-Cart-Session: $CART_SESSION" | jq
```

**Checkout:**

```bash
curl -s -X POST http://localhost:8080/v1/storefront/checkout \
  -H "X-Tenant-Slug: $TENANT_SLUG" \
  -H "X-Cart-Session: $CART_SESSION" \
  -H "Content-Type: application/json" \
  -d '{
    "recipient_name": "Siti",
    "recipient_phone": "08129999999",
    "shipping_address": "Jl. Anggrek No. 5",
    "shipping_city": "Jakarta",
    "shipping_province": "DKI Jakarta",
    "shipping_postal_code": "12345",
    "customer_email": "customer@test.com",
    "shipping_cost": 15000
  }' | jq
```

Simpan `order_id` dari response.

---

### 8. Pembayaran Midtrans

```bash
curl -s -X POST "http://localhost:8080/v1/storefront/orders/<order_id>/pay" \
  -H "X-Tenant-Slug: $TENANT_SLUG" | jq
```

**Expected:** `snap_token` — gunakan di frontend Midtrans Snap, atau test via Midtrans sandbox dashboard.

> Webhook Midtrans: `POST /v1/webhooks/midtrans`  
> Untuk local dev, gunakan ngrok/cloudflared agar Midtrans bisa kirim notifikasi.

---

### 9. Merchant — Kelola Order

```bash
curl -s http://localhost:8080/v1/merchant/orders \
  -H "Authorization: Bearer $MERCHANT_TOKEN" | jq

curl -s http://localhost:8080/v1/merchant/orders/<order_id> \
  -H "Authorization: Bearer $MERCHANT_TOKEN" | jq

curl -s -X PATCH "http://localhost:8080/v1/merchant/orders/<order_id>/status" \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "processing", "note": "Sedang diproses"}' | jq

curl -s -X PATCH "http://localhost:8080/v1/merchant/orders/<order_id>/status" \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "shipped", "note": "Dikirim via JNE"}' | jq

curl -s -X PATCH "http://localhost:8080/v1/merchant/orders/<order_id>/status" \
  -H "Authorization: Bearer $MERCHANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed", "note": "Selesai"}' | jq
```

---

### 10. Customer — Riwayat Order

```bash
curl -s http://localhost:8080/v1/storefront/orders \
  -H "X-Tenant-Slug: $TENANT_SLUG" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" | jq
```

---

## Checklist Phase 1

| # | Fitur | Pass? |
|---|---|---|
| 1 | Store settings update | ☐ |
| 2 | CRUD kategori | ☐ |
| 3 | CRUD produk + upload gambar | ☐ |
| 4 | Storefront home & katalog | ☐ |
| 5 | Customer register/login | ☐ |
| 6 | Cart (guest & logged-in) | ☐ |
| 7 | Checkout → order created | ☐ |
| 8 | Midtrans snap token | ☐ |
| 9 | Merchant kelola order | ☐ |
| 10 | Customer riwayat order | ☐ |

---

## API Endpoints Phase 1

### Merchant (`Authorization: Bearer <merchant_token>`)

| Method | Endpoint |
|---|---|
| PATCH | `/v1/merchant/store/settings` |
| GET/POST | `/v1/merchant/categories` |
| PATCH/DELETE | `/v1/merchant/categories/:id` |
| GET/POST | `/v1/merchant/products` |
| GET/PATCH/DELETE | `/v1/merchant/products/:id` |
| PATCH | `/v1/merchant/products/:id/inventory` |
| POST | `/v1/merchant/products/:id/images` |
| GET | `/v1/merchant/orders` |
| GET | `/v1/merchant/orders/:id` |
| PATCH | `/v1/merchant/orders/:id/status` |
| POST | `/v1/merchant/orders/:id/cancel` |

### Storefront (`X-Tenant-Slug` wajib)

| Method | Endpoint |
|---|---|
| GET | `/v1/storefront/home` |
| GET | `/v1/storefront/categories` |
| GET | `/v1/storefront/products` |
| GET | `/v1/storefront/products/:slug` |
| POST | `/v1/storefront/auth/register` |
| POST | `/v1/storefront/auth/login` |
| GET | `/v1/storefront/cart` |
| POST | `/v1/storefront/cart/items` |
| POST | `/v1/storefront/checkout` |
| POST | `/v1/storefront/orders/:id/pay` |
| GET | `/v1/storefront/orders/:id` |

### Webhook

| Method | Endpoint |
|---|---|
| POST | `/v1/webhooks/midtrans` |
