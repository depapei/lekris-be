# Lele Krispy POS Backend

Backend API untuk aplikasi Point of Sales **Lele Krispy**, dibangun dengan **Go (Golang)**, **GIN Framework**, dan **GORM** sebagai ORM. Menyediakan fitur CRUD untuk manajemen produk, supply, dan transaksi, serta sistem autentikasi berbasis JWT.

---

## 🚀 Tech Stack

| Komponen  | Teknologi            |
| --------- | -------------------- |
| Language  | Go 1.25.1            |
| Framework | GIN                  |
| ORM       | GORM + gorm-gen      |
| Database  | PostgreSQL           |
| Auth      | JWT (expired 1 hari) |

---

## 📋 Prasyarat

Pastikan lingkungan pengembangan Anda telah terinstal:

- [Go](https://go.dev/dl/) versi 1.25.1 atau lebih baru
- [PostgreSQL](https://www.postgresql.org/download/)
- Terminal / Command Prompt

---

## ⚙️ Instalasi & Konfigurasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd <nama-folder-proyek>
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Konfigurasi Environment Variables

Buat file `.env` di root proyek dan isi dengan konfigurasi berikut:

```env
SECRET_KEY="string_rahasia_anda"
DB_HOST="localhost"
DB_USER="postgres"
DB_PASSWORD="password_anda"
DB_PORT=5432
DB_DBNAME="lele_krispy_db"
PORT=":8080"
```

> ⚠️ **Penting**: Pastikan nilai `SECRET_KEY` diganti dengan string acak yang aman untuk produksi.

### 4. Generate Models (Opsional)

Jika Anda perlu memperbarui struktur model berdasarkan perubahan skema database, jalankan fungsi generator:

```bash
# Pastikan fungsi generateModel() dipanggil, atau jalankan via test/main sementara
go run main.go
```

> Model akan dihasilkan otomatis ke folder `./models` menggunakan `gorm-gen`.

### 5. Jalankan Aplikasi

```bash
go build -o main main.go
./main
```

Atau langsung jalankan tanpa build:

```bash
go run main.go
```

Aplikasi akan berjalan di `http://localhost:8080` (sesuai konfigurasi `PORT`).

---

## 🔐 Autentikasi

Sistem menggunakan **JWT (JSON Web Token)** dengan masa berlaku **1 hari**.

### Alur Login:

1. Client mengirim `POST /api/auth/login/` dengan kredensial admin.
2. Server memvalidasi dan mengembalikan `token` JWT.
3. Client menyertakan token di header `Authorization: Bearer <token>` untuk mengakses endpoint yang dilindungi.

---

## 📡 Daftar Endpoint API

### 🔓 Public Endpoint

| Method | Endpoint           | Deskripsi                               |
| ------ | ------------------ | --------------------------------------- |
| POST   | `/api/auth/login/` | Login admin untuk mendapatkan token JWT |

### 🔐 Protected Endpoints (Memerlukan Token JWT)

#### 📦 Products

| Method | Endpoint             | Deskripsi           |
| ------ | -------------------- | ------------------- |
| GET    | `/api/products/`     | Ambil daftar produk |
| GET    | `/api/products/:id/` | Ambil detail produk |
| POST   | `/api/products/`     | Tambah produk baru  |
| PUT    | `/api/products/:id/` | Update data produk  |
| DELETE | `/api/products/:id/` | Hapus produk        |

#### 📦 Supplies

| Method | Endpoint             | Deskripsi           |
| ------ | -------------------- | ------------------- |
| GET    | `/api/supplies/`     | Ambil daftar supply |
| GET    | `/api/supplies/:id/` | Ambil detail supply |
| POST   | `/api/supplies/`     | Tambah supply baru  |
| PUT    | `/api/supplies/:id/` | Update data supply  |
| DELETE | `/api/supplies/:id/` | Hapus supply        |

#### 💸 Transactions

| Method | Endpoint                               | Deskripsi              |
| ------ | -------------------------------------- | ---------------------- |
| GET    | `/api/transactions/`                   | Ambil daftar transaksi |
| GET    | `/api/transactions/:id/`               | Ambil detail transaksi |
| POST   | `/api/transactions/`                   | Buat transaksi baru    |
| PUT    | `/api/transactions/:id/`               | Update transaksi       |
| DELETE | `/api/transactions/:id/`               | Hapus transaksi        |
| GET    | `/api/transactions/payment-proof/:id/` | Ambil bukti pembayaran |

---

## 🗂️ Struktur Proyek

```
.
├── main.go              # Entry point, konfigurasi routes & middleware
├── controllers/         # Logic handler untuk setiap endpoint
├── models/              # Struct model (di-generate oleh gorm-gen)
├── middleware/          # Middleware (AuthMiddleware, dll)
├── uploads/             # Direktori untuk file upload (bukti pembayaran)
└── .env                 # Konfigurasi environment variables
```

> ⚙️ **Catatan**: Routing didefinisikan langsung di `main.go` menggunakan GIN router groups.

---

## 🗄️ Skema Database (Overview)

Proyek ini menggunakan relasi antar tabel sebagai berikut:

- `transactions` memiliki banyak `detail_transaction`
- `detail_transaction` memiliki satu `products`
- `transactions` memiliki relasi ke `user` untuk `created_by` dan `updated_by`

Relasi ini dikelola menggunakan fitur `FieldRelate` dari `gorm-gen`.

---

## 🛠️ Troubleshooting

| Masalah                          | Solusi                                                                   |
| -------------------------------- | ------------------------------------------------------------------------ |
| `connection refused` ke database | Pastikan PostgreSQL berjalan dan kredensial di `.env` sudah benar        |
| `401 Unauthorized`               | Pastikan token JWT valid dan disertakan di header `Authorization`        |
| Model tidak ter-generate         | Pastikan fungsi `generateModel()` dijalankan dan koneksi DB aktif        |
| Folder `uploads` tidak terbaca   | Pastikan middleware static file GIN sudah mengarah ke folder `./uploads` |

---

## 📝 Catatan Pengembangan

- Tidak ada unit test atau integration test saat ini.
- Migrasi skema tabel ditangani secara manual atau melalui `gorm-gen`.
- Untuk keamanan produksi, pertimbangkan menambahkan:
  - Rate limiting
  - Input validation yang lebih ketat
  - Logging & monitoring

---

> 🤝 **Dibuat oleh**: depapei / rangga  
> 📅 **Terakhir diperbarui**: 2026

_Jika ada pertanyaan atau kontribusi, silakan buka issue atau pull request._
