# WebCommerce SaaS Platform - Product Description

## Latar Belakang

Saya ingin membangun platform WebCommerce SaaS yang memungkinkan UMKM, brand lokal, dan seller marketplace di Indonesia memiliki website e-commerce sendiri tanpa harus membangun sistem dari nol.

Platform ini bukan marketplace seperti Tokopedia atau Shopee, melainkan platform yang memungkinkan merchant membuat dan mengelola toko online mereka sendiri, mirip konsep Shopify.

Fokus utama platform adalah membantu merchant yang sudah berjualan di marketplace seperti Shopee, TikTok Shop, dan Lazada untuk memiliki aset digital sendiri berupa website e-commerce yang dapat dikontrol penuh oleh merchant.

Dalam jangka panjang platform akan berkembang menjadi Omnichannel Commerce Platform yang menghubungkan website, marketplace, inventory, customer management, dan marketing automation dalam satu sistem.

---

## Target Market

### Primary Market

* UMKM Indonesia
* Brand lokal
* Seller marketplace

### Secondary Market

* Distributor
* Reseller
* UKM yang ingin melakukan digitalisasi penjualan

---

## Business Model

Platform menggunakan model SaaS (Software as a Service).

Merchant membayar biaya berlangganan bulanan atau tahunan untuk menggunakan platform.

Rencana paket:

* Trial
* Starter
* Growth
* Enterprise

Setiap paket memiliki batasan fitur, jumlah pengguna, storage, dan kemampuan integrasi yang berbeda.

---

## User Roles

### Customer

Pengunjung toko yang melakukan pembelian produk.

### Merchant Owner

Pemilik toko yang memiliki akses penuh terhadap toko dan konfigurasi bisnis.

### Merchant Admin

Pengguna internal merchant yang membantu operasional toko.

### Merchant Staff

Pengguna dengan akses terbatas sesuai role dan permission yang diberikan.

### Super Admin

Tim internal platform yang mengelola seluruh tenant, subscription, billing, dan melakukan customer support.

---

## Aplikasi yang Akan Dibangun

### 1. Storefront Application

Digunakan customer untuk mengakses toko online merchant.

Fitur utama:

* Homepage
* Product Catalog
* Product Detail
* Search
* Cart
* Checkout
* Customer Account
* Order Tracking

---

### 2. Merchant Dashboard

Digunakan merchant untuk mengelola operasional toko.

Fitur utama:

* Product Management
* Category Management
* Inventory Management
* Order Management
* Customer Management
* Promotion Management
* Analytics
* Store Settings

---

### 3. Super Admin Dashboard

Digunakan internal perusahaan.

Fitur utama:

* Tenant Management
* Subscription Management
* Billing Management
* Feature Flag Management
* User Monitoring
* Audit Log
* Impersonation
* Platform Analytics

---

## Multi-Tenant Architecture

Platform harus dirancang sebagai Multi-Tenant SaaS.

Setiap merchant disebut Tenant.

Semua tenant menggunakan infrastruktur dan aplikasi yang sama, tetapi data harus terisolasi secara logis.

Target arsitektur harus mampu mendukung ribuan tenant tanpa perubahan besar pada desain sistem.

Setiap tenant memiliki:

* Toko
* Domain
* User
* Produk
* Pesanan
* Subscription

yang terpisah dari tenant lain.

---

## Domain Management

Platform harus mendukung dua jenis domain:

### Platform Subdomain

Contoh:

tokobunga.webcommerce.id

### Custom Domain

Contoh:

tokobunga.com

Merchant dapat menghubungkan domain miliknya sendiri ke platform.

Sistem harus dapat mengenali tenant berdasarkan domain yang digunakan customer saat mengakses toko.

---

## Identity and Access Management

Platform harus memiliki sistem otorisasi yang fleksibel.

Menggunakan Role Based Access Control (RBAC).

Contoh permission:

* product.view
* product.create
* product.update
* order.view
* order.manage
* customer.view
* subscription.manage

Permission harus dapat dikonfigurasi dan tidak hardcoded.

---

## Subscription System

Platform menggunakan sistem berlangganan.

Fitur yang harus didukung:

* Trial
* Monthly Subscription
* Yearly Subscription
* Upgrade Plan
* Downgrade Plan
* Subscription Renewal
* Grace Period
* Subscription Expiration

Akses fitur ditentukan berdasarkan paket berlangganan yang aktif.

---

## Feature Flag System

Fitur tertentu hanya tersedia pada paket tertentu.

Contoh:

* Custom Domain
* Multi Staff
* Marketplace Integration
* Advanced Analytics
* AI Features

Sistem harus mampu membatasi akses fitur berdasarkan paket yang dimiliki tenant.

---

## Billing System

Platform harus memiliki billing engine internal.

Fitur:

* Invoice Generation
* Payment Tracking
* Subscription Billing
* Renewal Reminder

Integrasi pembayaran direncanakan menggunakan Midtrans dan penyedia pembayaran lokal lainnya.

---

## Impersonation

Super Admin harus dapat melakukan Login As Merchant untuk kebutuhan customer support dan debugging.

Seluruh aktivitas impersonation harus dicatat pada audit log.

---

## Audit Log

Semua aktivitas penting harus direkam.

Contoh:

* Login
* Update Product
* Delete Product
* Change Domain
* Upgrade Subscription
* Impersonation

Audit log harus dapat digunakan untuk kebutuhan investigasi dan debugging.

---

# OwnCommerce - Rencana Arsitektur Aplikasi

## Gambaran Umum

OwnCommerce adalah platform WebCommerce SaaS Multi-Tenant yang memungkinkan merchant memiliki website e-commerce sendiri, mengelola pelanggan sendiri, dan membangun brand sendiri tanpa bergantung sepenuhnya pada marketplace.

Platform ini dirancang untuk berkembang dari WebCommerce menjadi Omnichannel Commerce Platform yang mengintegrasikan website, marketplace, inventory, customer management, analytics, dan AI dalam satu ekosistem.

Prinsip utama platform:

* Own Your Store
* Own Your Brand
* Own Your Customers

---

# Arsitektur Tingkat Tinggi

Platform terdiri dari tiga aplikasi utama yang menggunakan backend dan database yang sama.

## Storefront Application

Digunakan oleh customer untuk berbelanja.

Fitur utama:

* Homepage
* Product Catalog
* Product Detail
* Search
* Cart
* Checkout
* Customer Account
* Order Tracking

Target pengguna:

* Customer
* Visitor

---

## Merchant Dashboard

Digunakan oleh merchant untuk mengelola toko.

Fitur utama:

* Product Management
* Category Management
* Inventory Management
* Order Management
* Customer Management
* Promotion Management
* Analytics
* Store Configuration

Target pengguna:

* Merchant Owner
* Merchant Admin
* Merchant Staff

---

## Super Admin Dashboard

Digunakan oleh tim internal OwnCommerce.

Fitur utama:

* Tenant Management
* Subscription Management
* Billing Management
* Feature Management
* Impersonation
* Monitoring
* Audit Log
* Platform Analytics

Target pengguna:

* Super Admin
* Customer Support
* Internal Operations

---

# Arsitektur Backend

Backend dibangun menggunakan pendekatan Modular Monolith.

Alasan:

* Lebih cepat dikembangkan dibanding microservices
* Lebih mudah dikelola oleh tim kecil
* Memiliki batas modul yang jelas
* Dapat dipecah menjadi microservices di masa depan jika diperlukan

Teknologi:

* Golang
* Fiber
* PostgreSQL
* Redis
* Asynq

Struktur aplikasi:

Core Platform
├── Auth
├── Tenant
├── IAM
├── Subscription
├── Billing
├── Audit
│
Commerce
├── Product
├── Category
├── Inventory
├── Customer
├── Cart
├── Order
├── Payment
├── Shipment
│
Omnichannel
├── Marketplace
├── Sync Engine
├── Channel Manager
│
Future
├── CRM
├── Marketing
├── Loyalty
├── AI

---

# Multi Tenant Architecture

Menggunakan Shared Database Shared Schema.

Seluruh merchant menggunakan database yang sama.

Setiap data bisnis wajib memiliki tenant_id.

Contoh:

products
orders
customers
inventories

Semua data selalu diisolasi berdasarkan tenant.

Keuntungan:

* Biaya infrastruktur rendah
* Mudah scaling
* Mudah maintenance
* Cocok untuk ribuan merchant

Target awal:

10.000+ tenant aktif

---

# Domain Management

Platform mendukung dua tipe domain.

## Platform Domain

Contoh:

tokobunga.owncommerce.id

Merchant langsung dapat menggunakan subdomain yang disediakan platform.

---

## Custom Domain

Contoh:

tokobunga.com

Merchant dapat menghubungkan domain pribadi ke platform.

Sistem akan melakukan:

* Domain Verification
* DNS Validation
* SSL Provisioning
* Domain Routing

Tenant akan dikenali berdasarkan host yang digunakan customer saat mengakses website.

---

# Identity & Access Management

Menggunakan Role Based Access Control (RBAC).

Role standar:

* Super Admin
* Merchant Owner
* Merchant Admin
* Merchant Staff

Permission bersifat dinamis.

Contoh:

product.view
product.create
product.update

order.view
order.manage

customer.view

subscription.manage

Role dan permission dapat dikembangkan tanpa perubahan kode aplikasi.

---

# Authentication Strategy

Menggunakan JWT Authentication.

Tipe token:

* Access Token
* Refresh Token

Mendukung:

* Multi Device Login
* Session Revocation
* Refresh Token Rotation

Seluruh API menggunakan Bearer Authentication.

---

# Subscription Architecture

Platform menggunakan model SaaS Subscription.

Paket awal:

* Trial
* Starter
* Growth
* Enterprise

Setiap paket memiliki:

* Batas fitur
* Batas user
* Batas storage
* Batas integrasi

Lifecycle:

Trial
→ Active
→ Renewal
→ Expired
→ Suspended

---

# Feature Flag System

Fitur ditentukan berdasarkan paket subscription.

Contoh:

Custom Domain
Marketplace Sync
Advanced Analytics
AI Features

Sistem dapat mengaktifkan atau menonaktifkan fitur tanpa deployment aplikasi.

---

# Billing Architecture

Modul billing terpisah dari subscription.

Komponen:

* Invoice
* Payment
* Tax
* Billing History

Integrasi pembayaran:

* Midtrans
* Xendit

Melalui abstraction layer agar provider dapat diganti tanpa mengubah business logic.

---

# Impersonation System

Super Admin dapat masuk ke akun merchant untuk kebutuhan support dan debugging.

Alur:

Super Admin
→ Login As Merchant
→ Temporary Session
→ Stop Impersonation

Seluruh aktivitas tercatat pada audit log.

---

# Audit Architecture

Seluruh aktivitas penting dicatat.

Contoh:

* Login
* Logout
* Product Update
* Product Delete
* Order Cancel
* Subscription Upgrade
* Domain Change
* Impersonation

Audit log menjadi sumber utama untuk investigasi dan troubleshooting.

---

# Background Processing

Menggunakan Asynq sebagai job queue.

Digunakan untuk:

* Email
* Notification
* Invoice Generation
* Marketplace Sync
* AI Processing
* Analytics Calculation

Worker dapat di-scale secara terpisah dari API.

---

# Storage Architecture

Database:

* PostgreSQL

Cache:

* Redis

File Storage:

* S3 Compatible Storage

Digunakan untuk:

* Product Image
* Theme Asset
* Document
* Invoice

---

# Struktur Repo

owncommerce/
│
├── apps/
│   ├── web/
│   ├── seller/
│   ├── superadmin/
│   └── api/
│
├── packages/
│   ├── ui/
│   ├── auth/
│   ├── types/
│   ├── sdk/
│   └── shared/
│
├── docs/
│
└── infra/
│
└── prompt/
