# Revo Creative CMS

**Content Management System** untuk website [Revo Creative](https://revocreative.co.id/) (PT Rajawali Cakra Digdaya) — agensi digital & software house Indonesia.

Mengelola konten website company profile secara mandiri tanpa edit HTML: ganti logo (light/dark untuk transisi navbar), tambah/ubah/hapus item portfolio & products, update info perusahaan (alamat, telp, email) yang muncul di semua halaman.

---

## Stack Teknologi

| Bagian | Teknologi |
|---|---|
| **Backend** | Go 1.21+ · Gin · GORM · PostgreSQL |
| **Auth** | Session cookie httpOnly |
| **File storage** | Local disk (`backend/uploads/`) |
| **Frontend (lama, statis)** | HTML · CSS · JavaScript vanilla |
| **Frontend (rencana, SPA)** | React · Vite · TypeScript |
| **Database** | PostgreSQL 16/18 |

---

## Struktur Project

```
revo-creative-cms/
├── README.md                 ← Anda di sini
├── backend/                  ← REST API Go (✅ Tahap 1-5 selesai)
│   ├── cmd/server/main.go
│   ├── internal/             (config, models, handlers, middleware, routes, services)
│   ├── uploads/              (file gambar yang di-upload)
│   ├── docs/
│   │   ├── API.md            (referensi 24 endpoint)
│   │   ├── ARCHITECTURE.md   (diagram & penjelasan layer)
│   │   └── RevoCMS.postman_collection.json
│   └── README.md             ← detail setup backend
│
├── frontend/                 ← HTML statis lama (referensi visual, akan di-port)
│   ├── index.html
│   ├── about.html
│   ├── portfolio.html
│   ├── products.html
│   ├── contact.html
│   ├── css/ · js/ · assets/
│
└── web/                      ← SPA React (belum dibuat, rencana Tahap 6+)
```

---

## Quick Start (untuk yang baru clone)

### 1. Prerequisites

| Tool | Versi minimum | Cek |
|---|---|---|
| Go | 1.21 | `go version` |
| PostgreSQL | 14 | `psql --version` |
| Node.js (untuk SPA nanti) | 20 | `node --version` |
| Git | any | `git --version` |

### 2. Jalankan Backend

```bash
cd backend
cp .env.example .env
# edit .env — isi DB_PASSWORD, SESSION_SECRET, SEED_ADMIN_PASSWORD
go mod download
go run ./cmd/server
```

Detail lengkap: **[backend/README.md](backend/README.md)**.

Backend akan jalan di `http://localhost:8080`.

### 3. Jalankan Frontend Statis (lama)

Buka `frontend/index.html` langsung di browser, atau pakai live server VSCode.

> Setelah SPA dibuat (Tahap 6+), frontend statis ini akan diarsipkan.

---

## Status Development

| Tahap | Komponen | Status |
|---|---|---|
| 1 | Setup backend Go (modules, struktur, koneksi DB) | ✅ Selesai |
| 2 | Models GORM + AutoMigrate + Seeder | ✅ Selesai |
| 3 | Auth (session cookie, login/logout/me) | ✅ Selesai |
| 4 | REST API CRUD (settings, logos, portfolio, products) | ✅ Selesai |
| 5 | File upload (multipart, validasi mime+size) | ✅ Selesai |
| 6 | Setup SPA React + Vite + TypeScript | 🔜 Berikutnya |
| 7 | Port public pages ke SPA (Home, About, Portfolio, dll) | ⏳ Pending |
| 8 | Admin panel SPA (login, dashboard, CRUD UI) | ⏳ Pending |
| 9 | Polish (rate-limit, validasi extra, README final, deploy script) | ⏳ Pending |

---

## Testing Backend

Cara tercepat — import Postman collection:

1. **Postman** → **Import** → pilih `backend/docs/RevoCMS.postman_collection.json`
2. Edit variabel `admin_email` & `admin_password` di Collection → tab **Variables** supaya match `.env` Anda
3. Jalankan request `Auth → POST /api/auth/login` dulu (cookie auto-tersimpan)
4. 24 request siap dipakai

Detail lengkap: **[backend/docs/API.md](backend/docs/API.md)**.

---

## Dokumentasi

| Topik | File |
|---|---|
| Setup & menjalankan backend | [backend/README.md](backend/README.md) |
| Endpoint reference lengkap | [backend/docs/API.md](backend/docs/API.md) |
| Arsitektur, flow request, peran tiap folder | [backend/docs/ARCHITECTURE.md](backend/docs/ARCHITECTURE.md) |
| Postman collection (import-able) | [backend/docs/RevoCMS.postman_collection.json](backend/docs/RevoCMS.postman_collection.json) |

---

## Konvensi Git

- Branch utama: `main`
- Commit message: singkat, deskriptif, Bahasa Indonesia atau Inggris konsisten
- `.env` **tidak** di-commit (sudah di `.gitignore`); selalu update `.env.example` saat tambah variable baru
- `backend/uploads/*` **tidak** di-commit (kecuali `.gitkeep`)

---

## Kontak

- **Brand**: Revo Creative
- **Legal**: PT Rajawali Cakra Digdaya
- **Email**: rajawalicakradigdaya@gmail.com
- **Website**: https://revocreative.co.id/
