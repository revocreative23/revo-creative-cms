# Revo Creative CMS — Backend

REST API backend untuk CMS Revo Creative (PT Rajawali Cakra Digdaya).
Mengelola: settings perusahaan, logo, portfolio, dan products.

**Stack**: Go 1.21+ · Gin · GORM · PostgreSQL · Session Cookie Auth · Local file upload

---

## Prerequisites

| Tool | Versi minimum | Cek |
|---|---|---|
| Go | 1.21 | `go version` |
| PostgreSQL | 14 | `psql --version` |
| Git | any | `git --version` |

PostgreSQL di Windows bisa diinstall dari [EnterpriseDB installer](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads).

---

## Setup

### 1. Buat database

Via `psql`:
```bash
psql -U postgres -h localhost -p 5433 -c "CREATE DATABASE revocreative_cms;"
```

Atau via pgAdmin: Servers → klik kanan Databases → Create → Database → nama `revocreative_cms`.

### 2. Copy & isi `.env`

```bash
cp .env.example .env
```

Edit `.env`:
```ini
DB_PASSWORD=<password_postgres_anda>
SESSION_SECRET=<string_acak_minimum_32_karakter>
SEED_ADMIN_EMAIL=admin@revocreative.local
SEED_ADMIN_PASSWORD=<password_admin_yang_kuat>
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Run server

```bash
go run ./cmd/server
```

Output sukses:
```
✓ connected ke postgres
→ menjalankan AutoMigrate...
✓ semua tabel ter-migrate
✓ admin user admin@revocreative.local berhasil dibuat
✓ seed setting: company_name
... (dst)
server jalan di http://localhost:8080
```

Tes: buka [http://localhost:8080/api/health](http://localhost:8080/api/health) — harus return `{"status":"ok"}`.

---

## Struktur Folder

```
backend/
├── cmd/server/main.go              # entry point
├── internal/
│   ├── config/                     # env loader & DB connection
│   ├── models/                     # GORM models + AutoMigrate
│   ├── handlers/                   # HTTP handlers (= Controllers)
│   ├── middleware/                 # auth, CORS
│   ├── routes/                     # URL → handler mapping
│   └── services/                   # business logic (seeder)
├── uploads/                        # file upload tersimpan di sini
├── docs/
│   ├── API.md                      # endpoint reference lengkap
│   ├── ARCHITECTURE.md             # diagram & flow
│   └── RevoCMS.postman_collection.json
├── .env.example
├── go.mod
└── README.md
```

Detail arsitektur: lihat [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## Endpoint Ringkasan

| Group | Endpoint | Auth |
|---|---|---|
| Health | `GET /api/health` | - |
| Auth | `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/auth/me` | mixed |
| Public | `GET /api/settings`, `GET /api/logos/active`, `GET /api/portfolio`, `GET /api/products` | - |
| Admin Settings | `GET/PUT /api/admin/settings[/:key]` | ✓ |
| Admin Logos | CRUD `/api/admin/logos` + `/activate` | ✓ |
| Admin Portfolio | CRUD `/api/admin/portfolio[/:id]` | ✓ |
| Admin Products | CRUD `/api/admin/products[/:id]` | ✓ |
| Admin Upload | `POST /api/admin/upload` (multipart) | ✓ |
| Static | `GET /uploads/<file>` | - |

Lihat detail request/response: [docs/API.md](docs/API.md).

---

## Testing dengan Postman

Import file [docs/RevoCMS.postman_collection.json](docs/RevoCMS.postman_collection.json) ke Postman:

1. Postman → **Import** → pilih file JSON
2. Aktifkan environment **Revo CMS Local** (otomatis ter-import)
3. Jalankan **Login** dulu (cookie otomatis tersimpan)
4. Request lain bisa dijalankan tanpa setup tambahan

---

## Environment Variables

| Variable | Default | Keterangan |
|---|---|---|
| `APP_ENV` | development | `development` atau `production` |
| `SERVER_PORT` | 8080 | Port HTTP server |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5433 | PostgreSQL port |
| `DB_USER` | postgres | DB user |
| `DB_PASSWORD` | - | **Wajib** diisi |
| `DB_NAME` | revocreative_cms | Nama database |
| `DB_SSLMODE` | disable | Set `require` di prod |
| `SESSION_SECRET` | - | **Wajib**, min 32 char acak |
| `SESSION_COOKIE_NAME` | revocms_session | Nama cookie |
| `SESSION_MAX_AGE_SECONDS` | 86400 | Durasi session (default 24 jam) |
| `CORS_ALLOWED_ORIGIN` | http://localhost:5173 | Origin SPA frontend (Vite default) |
| `UPLOAD_DIR` | ./uploads | Folder simpan file |
| `UPLOAD_MAX_BYTES` | 5242880 | Max ukuran file (5MB) |
| `SEED_ADMIN_EMAIL` | - | Email admin yang di-seed |
| `SEED_ADMIN_PASSWORD` | - | Password admin yang di-seed |

---

## Build untuk Production

```bash
go build -o revocms-server ./cmd/server
./revocms-server
```

Catatan production:
- Set `APP_ENV=production`
- Set `DB_SSLMODE=require` (atau verifikasi sesuai setup DB)
- Cookie akan otomatis pakai `Secure` flag (butuh HTTPS)
- Update `CORS_ALLOWED_ORIGIN` ke domain produksi
