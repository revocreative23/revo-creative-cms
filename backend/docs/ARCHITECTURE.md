# Architecture — Revo Creative CMS Backend

## Pola Arsitektur

Project ini pakai **Layered Architecture** (= variasi MVC tanpa View, karena backend REST API).

| Layer MVC klasik | Di project kita | Catatan |
|---|---|---|
| **Model** | `internal/models/` | GORM struct → tabel PostgreSQL |
| **View** | tidak ada di backend | Backend cuma return JSON. "View" sebenarnya ada di SPA React (rencana Tahap 7) |
| **Controller** | `internal/handlers/` | Konvensi Go → disebut **Handler** |

Tambahan layer di luar MVC klasik:
- `services/` — business logic kompleks (saat ini: seeder)
- `middleware/` — cross-cutting concerns (auth, CORS, session)
- `routes/` — URL → handler mapping terpusat
- `config/` — env loader & DB connection

---

## Flow Request

```
HTTP request
    │
    ▼
┌─────────────────────────────────────────────────┐
│ cmd/server/main.go                              │
│  - Setup Gin server                             │
│  - CORS, Session middleware                     │
│  - Static files /uploads                        │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│ internal/routes/routes.go                       │
│  - URL → handler mapping                        │
│  - Apply middleware (RequireAuth, dll)          │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│ internal/middleware/auth.go                     │
│  - Cek session cookie                           │
│  - Block 401 kalau tidak login                  │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│ internal/handlers/*.go        (= "Controller")  │
│  - Parse & validasi request                     │
│  - Panggil model / service                      │
│  - Format response JSON                         │
└────────┬──────────────────────────┬─────────────┘
         │                          │
         ▼                          ▼
┌──────────────────┐    ┌──────────────────────┐
│ services/*.go    │    │ models/*.go          │
│ (business logic) │    │ (struct + GORM tag)  │
│ - seeder         │    │ - User, Logo, dll    │
│ - (kompleks bs.) │    │ - migrasi otomatis   │
└──────────────────┘    └──────────┬───────────┘
                                   │
                                   ▼
                       ┌──────────────────────┐
                       │ config/database.go   │
                       │ → GORM → PostgreSQL  │
                       └──────────────────────┘
```

---

## Peran Tiap Folder

| Folder | Konsep | Tugas | Boleh import dari |
|---|---|---|---|
| `cmd/server/` | Entry point | `main.go`, bootstrap | semua |
| `internal/config/` | Infrastructure | Load env, koneksi DB | godotenv, GORM |
| `internal/models/` | **Model** | Struct + GORM tag + AutoMigrate | GORM |
| `internal/services/` | **Service Layer** | Business logic kompleks | models, GORM |
| `internal/handlers/` | **Controller** | HTTP → call service/model → JSON | services, models, gin |
| `internal/middleware/` | Cross-cutting | Auth, logging, CORS | gin, sessions |
| `internal/routes/` | Routing | URL → handler mapping | handlers, middleware |

**Aturan dependency** (one-way, atas → bawah):
```
cmd → routes → middleware → handlers → services → models → config
```
Jangan ada cyclic import (Go akan complain compile-time).

---

## Models

5 model utama (lihat [internal/models/](../internal/models/)):

| Model | Tabel | Tujuan |
|---|---|---|
| `User` | `users` | Admin login (1 user untuk MVP, scalable ke multi-admin) |
| `SiteSetting` | `site_settings` | Key-value untuk info perusahaan (alamat, telp, dll) |
| `Logo` | `logos` | Per-type logo (light/dark/favicon), 1 aktif per type |
| `PortfolioItem` | `portfolio_items` | Item portfolio dengan tags JSON |
| `Product` | `products` | Product dengan features JSON |

### Kenapa SiteSetting pakai Key-Value?

Alternatifnya: bikin 1 tabel `companies` dengan kolom `name`, `address`, `phone`, dll. Tapi:
- ✅ Key-value: nambah field baru (mis. `tiktok_url`) **tanpa migrasi DB**, cukup insert row baru via seeder atau admin panel
- ❌ Per-kolom: setiap field baru = migrasi schema = restart server

Trade-off:
- Key-value: type-safety lemah (semua string), butuh seeder untuk default
- Per-kolom: type-safety kuat, struktur jelas di DB

Untuk CMS yang field-nya bisa berubah seiring waktu, key-value lebih fleksibel.

### Soft Delete

Semua model (kecuali `SiteSetting`) punya `DeletedAt gorm.DeletedAt`. Saat `db.Delete()`, GORM tidak hapus row — hanya isi `deleted_at`. Query default sudah exclude row yang ter-delete. Bisa recovery dengan `db.Unscoped()`.

---

## Authentication Flow

```
1. Client POST /api/auth/login { email, password }
   ↓
2. Handler:
   - Lookup user by email (GORM)
   - bcrypt.CompareHashAndPassword
   - Kalau cocok: simpan user_id, email, role di session
   ↓
3. Session middleware encrypt & set cookie `revocms_session` (httpOnly)
   ↓
4. Client browser auto-attach cookie di request berikutnya
   ↓
5. Request ke /api/auth/me atau /api/admin/*:
   - middleware/auth.go cek session.Get("user_id")
   - Kalau nil → 401 Unauthorized
   - Kalau ada → c.Set("user_id", uid) untuk dipakai handler
   ↓
6. Handler ambil c.GetUint("user_id") untuk query owner-specific data
```

**Kenapa session cookie, bukan JWT?**

Untuk admin panel internal (kasus kita):
- ✅ Session cookie httpOnly = lebih aman dari XSS
- ✅ Auto-managed browser, tidak perlu logic di SPA
- ✅ Logout server-side (clear session)
- ❌ Tidak portable ke mobile app / cross-domain — tapi kita tidak butuh

JWT lebih cocok kalau:
- API dipakai mobile native
- Multiple frontend di domain berbeda
- Microservices stateless

---

## File Upload

```
1. Client POST /api/admin/upload (multipart, field: file, subdir)
   ↓
2. UploadHandler:
   - Cek size (UPLOAD_MAX_BYTES)
   - Cek ekstensi (whitelist)
   - Cek magic bytes (validasi isi file)
   - Sanitize subdir (anti path traversal)
   - Generate nama uuid + ekstensi asli
   - Simpan ke ./uploads/<subdir>/<uuid>.<ext>
   ↓
3. Return file_path: "/uploads/<subdir>/<uuid>.<ext>"
   ↓
4. Client POST /api/admin/logos (atau portfolio/products)
   dengan thumbnail_path/file_path = file_path tadi
   ↓
5. Saat publik akses GET /uploads/... :
   Gin static file serving → return file dari disk
```

**Kenapa upload terpisah dari create entity?**

Alternative: 1 endpoint multipart yang sekaligus upload + create portfolio.

Pilihan kita (upload dulu, lalu create):
- ✅ UI bisa show preview gambar dulu sebelum user submit form
- ✅ Replace gambar tanpa create entity baru (cukup update field path)
- ✅ Re-use upload endpoint untuk semua entity (1 implementasi)
- ❌ 2 request, butuh "garbage collection" untuk file yang ter-upload tapi tidak pernah dipakai (nice-to-have, bukan blocker)

---

## Database

PostgreSQL, schema di-manage via GORM `AutoMigrate` (lihat [internal/models/migrate.go](../internal/models/migrate.go)).

**Strategi migrasi:**
- Dev: AutoMigrate jalan otomatis tiap server start. Aman untuk ADD column. **Tidak aman** untuk DROP/RENAME column — perlu migrasi manual.
- Prod: idealnya pakai tool migrasi terpisah (golang-migrate, atlas). Untuk MVP, AutoMigrate cukup.

**Connection pool**: default GORM. Bisa di-tuning via `sqlDB.SetMaxOpenConns()` di `config/database.go` saat traffic naik.

---

## Configuration

Pakai env file `.env` di-load via `godotenv` saat startup. Di production, env vars di-set di shell/OS, file `.env` tidak perlu ada.

Semua config terpusat di [internal/config/config.go](../internal/config/config.go) sebagai struct `Config`. Handler/service ambil dari sini, tidak `os.Getenv()` langsung — supaya mudah test.

---

## Kapan Refactor?

Saat ini handler langsung `h.db.Find(...)` — OK untuk CRUD sederhana. Refactor ke pattern lebih kompleks ketika:

1. **Repository pattern** — kalau query sering di-reuse di multiple handler / butuh mock untuk testing
2. **Service layer expansion** — kalau handler punya >50 baris logic atau orchestrate banyak operasi
3. **DTO terpisah** — kalau request/response shape jauh dari model (sekarang sudah dipisah via `xxxRequest` structs)
4. **Domain-Driven Design** — overkill untuk CMS, jangan dulu

**Aturan main**: tunggu sampai ada 2-3 instance yang butuh abstraksi sama baru extract. Premature abstraction = boilerplate tanpa benefit.
