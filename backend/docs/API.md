# API Reference — Revo Creative CMS

Base URL: `http://localhost:8080`

Semua response berupa JSON. Status code mengikuti standard REST:
- `200 OK` — sukses GET/PUT/DELETE
- `201 Created` — sukses POST
- `400 Bad Request` — validasi gagal / format request salah
- `401 Unauthorized` — perlu login / session expired
- `404 Not Found` — resource tidak ditemukan
- `500 Internal Server Error` — error server

Auth pakai **session cookie** (`revocms_session`, httpOnly). Klien wajib `credentials: 'include'` saat fetch dari SPA.

---

## 1. Health Check

### `GET /api/health`

**Response 200:**
```json
{
  "status": "ok",
  "app": "revocreative-cms"
}
```

---

## 2. Auth

### `POST /api/auth/login`

**Request body:**
```json
{
  "email": "admin@revocreative.local",
  "password": "admin12345"
}
```

**Response 200** (set cookie `revocms_session`):
```json
{
  "user": {
    "id": 1,
    "email": "admin@revocreative.local",
    "name": "Administrator",
    "role": "admin"
  }
}
```

**Response 401:**
```json
{ "error": "email atau password salah" }
```

### `POST /api/auth/logout`

Clear session cookie.

**Response 200:**
```json
{ "message": "logout berhasil" }
```

### `GET /api/auth/me` 🔒

**Response 200:**
```json
{
  "user": { "id": 1, "email": "...", "name": "...", "role": "admin" }
}
```

**Response 401:** session tidak valid / expired.

---

## 3. Settings

### `GET /api/settings` (public)

Return semua setting sebagai map key→value, untuk konsumsi frontend publik.

**Response 200:**
```json
{
  "company_name": "Revo Creative",
  "company_legal_name": "PT Rajawali Cakra Digdaya",
  "company_address": "Jl. Raya Pd. Gede No.14A, ...",
  "company_phone": "+62 856-7990-037",
  "company_email": "rajawalicakradigdaya@gmail.com",
  "company_whatsapp": "628567990037",
  "social_instagram": "https://www.instagram.com/revocreative_id/",
  "social_linkedin": "https://www.linkedin.com/company/101702844",
  "footer_tagline": "Membangun digital experience yang berdampak."
}
```

### `GET /api/admin/settings` 🔒

Return list lengkap dengan description & timestamps.

**Response 200:**
```json
[
  {
    "id": 1,
    "key": "company_name",
    "value": "Revo Creative",
    "description": "Nama brand yang ditampilkan di header/footer",
    "updated_at": "2026-05-25T10:30:00+07:00"
  },
  ...
]
```

### `PUT /api/admin/settings/:key` 🔒

Update value setting. `:key` harus sudah ada (tidak bisa create lewat sini).

**Request body:**
```json
{ "value": "+62 812-3456-7890" }
```

**Response 200:** object setting yang sudah ter-update.

**Response 404:** key tidak ditemukan.

---

## 4. Logos

Logo punya 3 tipe: `light`, `dark`, `favicon`. Hanya 1 logo per tipe yang `is_active` pada satu waktu.

### `GET /api/logos/active` (public)

Return map type → logo aktif.

**Response 200:**
```json
{
  "light": {
    "id": 5,
    "type": "light",
    "file_path": "/uploads/logos/abc123.webp",
    "is_active": true,
    "created_at": "...",
    "updated_at": "..."
  },
  "dark": { ... }
}
```

### `GET /api/admin/logos` 🔒

List semua logo (termasuk yang tidak aktif & history lama).

**Query param opsional:** `?type=light`

### `POST /api/admin/logos` 🔒

Register logo baru. Path file harus sudah ada di server (upload duluan via `/api/admin/upload`).

**Request body:**
```json
{
  "type": "light",
  "file_path": "/uploads/logos/abc123.webp",
  "activate": true
}
```

- `type`: `light` | `dark` | `favicon`
- `activate`: kalau `true`, langsung jadi aktif & non-aktifkan yang lain dengan type sama

**Response 201:** object logo.

### `PUT /api/admin/logos/:id/activate` 🔒

Aktifkan logo ini, non-aktifkan yang lain dengan type sama.

**Response 200:** object logo yang sudah `is_active: true`.

### `DELETE /api/admin/logos/:id` 🔒

Soft delete (data tetap di DB dengan `deleted_at` terisi).

**Response 200:**
```json
{ "message": "logo dihapus" }
```

---

## 5. Portfolio

### `GET /api/portfolio` (public)

List portfolio yang published, urut `display_order` ASC.

**Query param opsional:** `?category=website` (filter)

**Response 200:**
```json
[
  {
    "id": 1,
    "title": "ISI — Ikatan Surveyor Indonesia",
    "category": "website",
    "category_label": "Web + CMS + Membership",
    "description": "Re-frame website company profile ISI ...",
    "thumbnail_path": "/uploads/portfolio/xyz.jpg",
    "tags": ["ISI", "Membership System", "Event Calendar", "CMS"],
    "display_order": 1,
    "is_published": true,
    "created_at": "...",
    "updated_at": "..."
  }
]
```

### `GET /api/admin/portfolio` 🔒

List semua portfolio (termasuk unpublished).

### `GET /api/admin/portfolio/:id` 🔒

Detail portfolio.

### `POST /api/admin/portfolio` 🔒

**Request body:**
```json
{
  "title": "ISI — Ikatan Surveyor Indonesia",
  "category": "website",
  "category_label": "Web + CMS + Membership",
  "description": "Re-frame website company profile ISI...",
  "thumbnail_path": "/uploads/portfolio/xyz.jpg",
  "tags": ["ISI", "Membership System", "CMS"],
  "display_order": 1,
  "is_published": true
}
```

**Fields:**
- `title` *(required)*
- `category` *(required)* — bebas, untuk filter (mis. `website`, `app`, `dashboard`)
- `category_label` — label panjang untuk display
- `description` — markdown/plain text
- `thumbnail_path` — path file (dari endpoint upload)
- `tags` — array string
- `display_order` — angka, semakin kecil semakin di atas
- `is_published` — default `true`

**Response 201:** object portfolio.

### `PUT /api/admin/portfolio/:id` 🔒

Update semua field. Body sama dengan POST.

### `DELETE /api/admin/portfolio/:id` 🔒

Soft delete.

---

## 6. Products

### `GET /api/products` (public)

List products yang published.

### `GET /api/admin/products` 🔒

List semua.

### `GET /api/admin/products/:id` 🔒

### `POST /api/admin/products` 🔒

**Request body:**
```json
{
  "title": "Company Profile + CMS",
  "slug": "company-profile-cms",
  "description": "Website company profile dengan CMS...",
  "thumbnail_path": "/uploads/products/xyz.png",
  "features": ["Custom Domain", "CMS Editor", "SEO Ready", "Mobile Responsive"],
  "price": "Mulai Rp 3.5jt",
  "display_order": 1,
  "is_published": true
}
```

- `slug` *(required, unique)* — untuk URL `/products/<slug>`
- `price` *string* — supaya bisa "Mulai Rp 3jt" atau "Hubungi kami"

### `PUT /api/admin/products/:id` 🔒

### `DELETE /api/admin/products/:id` 🔒

---

## 7. File Upload

### `POST /api/admin/upload` 🔒 (multipart/form-data)

**Form fields:**
- `file` *(required)* — file binary
- `subdir` *(optional)* — folder pengelompokan, mis. `logos`, `portfolio`, `products`. Hanya `[a-z0-9_-]+`.

**Whitelist ekstensi:** `.png`, `.jpg`, `.jpeg`, `.webp`, `.svg`, `.gif`, `.ico`

**Max size:** dari env `UPLOAD_MAX_BYTES` (default 5MB).

**Validasi:** ekstensi + magic bytes (cek isi file beneran image).

**Response 201:**
```json
{
  "file_path": "/uploads/logos/abc12345-....webp",
  "original_name": "revo-logo-blue.webp",
  "size": 24576,
  "mime": "image/webp"
}
```

`file_path` ini disimpan ke DB sebagai `thumbnail_path` / logo `file_path` saat POST/PUT entity terkait.

**Response 400:**
- `file tidak ditemukan` — form tidak ada field `file`
- `file terlalu besar (max N bytes)`
- `tipe file tidak diizinkan. Hanya: png, jpg, jpeg, webp, svg, gif, ico`
- `isi file (...) tidak cocok dengan ekstensi (...)` — magic byte check gagal

### `GET /uploads/<file>` (public)

Akses langsung file yang sudah di-upload. Static serving.

Contoh: `http://localhost:8080/uploads/logos/abc123.webp`

---

## Workflow Lengkap — Frontend Upload + Set Logo Light

1. **Login** → `POST /api/auth/login` (cookie tersimpan)
2. **Upload file** → `POST /api/admin/upload` multipart, dapat `file_path: "/uploads/logos/xxx.webp"`
3. **Register logo** → `POST /api/admin/logos` dengan body:
   ```json
   { "type": "light", "file_path": "/uploads/logos/xxx.webp", "activate": true }
   ```
4. **Frontend publik** fetch → `GET /api/logos/active` → dapat URL logo aktif terbaru.

---

🔒 = perlu session cookie (login terlebih dahulu)
