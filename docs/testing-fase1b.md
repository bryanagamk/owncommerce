# Manual Testing — Phase 1B (Frontend MVP)

Panduan testing UI end-to-end: merchant dashboard (`seller`) + customer storefront (`web`).

> **Prasyarat:** Phase 1A API jalan, infra Docker aktif. Lihat [testing-fase1.md](./testing-fase1.md) untuk setup API.

---

## Setup

### 1. Environment

```bash
# Root — API & database
cp .env.example .env

# Frontend apps
cp apps/seller/.env.example apps/seller/.env
cp apps/web/.env.example apps/web/.env
```

Isi `VITE_MIDTRANS_CLIENT_KEY` di `apps/web/.env` dengan sandbox client key Midtrans (sama dengan `MIDTRANS_CLIENT_KEY` di root `.env`).

### 2. Jalankan Services

Terminal 1 — infra + API:

```bash
make infra-up
make api-dev
```

Terminal 2 — merchant dashboard:

```bash
npm install
npm run dev:seller
```

Buka http://localhost:5173

Terminal 3 — storefront:

```bash
npm run dev:web
```

Buka http://localhost:5174

---

## Flow Testing Lengkap (UI)

### A. Merchant — Daftar & Setup Toko

1. Buka http://localhost:5173/register
2. Daftar merchant baru (email, password, nama, nama toko, slug)
3. Setelah login, buka **Settings** → isi deskripsi, kontak, alamat → simpan
4. Buka **Categories** → tambah kategori (mis. "Bunga Segar")
5. Buka **Products** → **Tambah Produk**:
   - Nama, harga, stok, status **Active**
   - Upload gambar produk
   - Centang **Featured** untuk tampil di homepage
6. Verifikasi produk muncul di tabel dengan gambar dan stok

### B. Customer — Belanja (Guest)

1. Buka http://localhost:5174 (tenant default: `tokobunga` via `VITE_TENANT_SLUG`)
2. Homepage menampilkan nama toko + produk unggulan
3. Klik produk → **Tambah ke Keranjang**
4. Buka **Keranjang** (ikon cart) → ubah qty / hapus item
5. **Checkout** → isi alamat pengiriman → **Buat Pesanan & Bayar**
6. Jendela Midtrans Snap terbuka (sandbox)
7. Selesaikan pembayaran sandbox → redirect ke `/payment/success`

> Jika Snap tidak terbuka: cek `VITE_MIDTRANS_CLIENT_KEY` dan console browser.

### C. Merchant — Proses Pesanan

1. Kembali ke http://localhost:5173/orders
2. Pesanan baru muncul dengan status `pending_payment` atau `paid` (setelah webhook)
3. Klik pesanan → update status: **processing** → **shipped** → **completed**

### D. Customer — Akun Terdaftar

1. Buka http://localhost:5174/account/register
2. Daftar akun customer
3. Login ulang di `/account/login`
4. Buka **Profil** → update nama/telepon, tambah alamat pengiriman
5. Belanja lagi (cart terikat ke akun)
6. Checkout → pilih alamat tersimpan dari dropdown
7. Buka **Pesanan** di header → lihat riwayat order

### E. Merchant — Onboarding Wizard

1. Daftar merchant baru di http://localhost:5173/register
2. Setelah daftar, diarahkan ke `/onboarding`
3. Lengkapi profil toko → (opsional) kategori → (opsional) produk pertama
4. Selesai → masuk dashboard
5. Login merchant lama → langsung ke dashboard (skip onboarding)

---

## Checklist Cepat

| # | Skenario | Seller | Web |
|---|---|:---:|:---:|
| 1 | Register / Login | ✅ | ✅ |
| 2 | CRUD kategori | ✅ | - |
| 3 | CRUD produk + upload gambar | ✅ | - |
| 4 | Homepage + katalog | - | ✅ |
| 5 | Search produk | - | ✅ |
| 6 | Add to cart (guest) | - | ✅ |
| 7 | Checkout + buat order | - | ✅ |
| 8 | Midtrans Snap payment | - | ✅ |
| 9 | Daftar & kelola pesanan merchant | ✅ | - |
| 10 | Customer orders history | - | ✅ |
| 11 | Customer profile & alamat | - | ✅ |
| 12 | Onboarding wizard merchant | ✅ | - |
| 13 | Filter kategori di katalog | - | ✅ |

---

## Troubleshooting

| Masalah | Penyebab | Solusi |
|---|---|---|
| Halaman kosong / API error | API tidak jalan | `make api-dev`, cek http://localhost:8080/v1/health |
| Toko tidak ditemukan | Slug salah | Pastikan `VITE_TENANT_SLUG` sama dengan slug saat register |
| Gambar tidak tampil | URL relatif | `VITE_API_URL` harus `http://localhost:8080` |
| 401 di seller | Token expired | Logout & login ulang |
| Snap tidak load | Client key kosong | Isi `VITE_MIDTRANS_CLIENT_KEY` di `apps/web/.env` |
| CORS error | Origin tidak diizinkan | Pastikan API dev mode mengizinkan localhost:5173/5174 |

---

## Port Reference

| Service | URL |
|---|---|
| API | http://localhost:8080 |
| Seller (merchant) | http://localhost:5173 |
| Web (storefront) | http://localhost:5174 |
| PostgreSQL (Docker) | localhost:5433 |
| Redis (Docker) | localhost:6380 |
