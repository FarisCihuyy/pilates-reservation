# DIRO – Pilates Reservation App

Aplikasi reservasi kelas Pilates berbasis web yang memungkinkan user untuk:

- Memilih tanggal
- Memilih timeslot
- Memilih court/studio
- Melakukan pembayaran (Dummy / Midtrans)
- Mendapatkan konfirmasi reservasi

---

## Tech Stack

### Frontend

- Next.js (App Router)
- TypeScript
- Tailwind CSS

### Backend

- Golang (Gin)
- GORM
- PostgreSQL / MySQL
- Midtrans (Snap) / Dummy Payment

---

## Struktur Repository

```
diro-pilates/
│
├── frontend/
└── backend/
```

---

## 🔧 Backend Setup (Golang)

### 1. Masuk ke folder backend

```bash
cd backend
```

### 2. Install dependency

```bash
go mod tidy
```

### 3. Environment Variable

Buat file `.env`:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=diro_db

MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
MIDTRANS_BASE_URL=https://app.sandbox.midtrans.com/snap/v1
```

> Jika credential Midtrans kosong, sistem otomatis menggunakan **Dummy Payment**.

### 4. Jalankan backend

```bash
go run main.go
```

Backend berjalan di:

```
http://localhost:8080
```

---

## Frontend Setup (Next.js)

### 1. Masuk ke folder frontend

```bash
cd frontend
```

### 2. Install dependency

```bash
npm install
```

### 3. Environment Variable

Buat file `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 4. Jalankan frontend

```bash
npm run dev
```

Frontend berjalan di:

```
http://localhost:3000
```

---

## Authentication Flow

1. User register / login
2. Backend mengembalikan JWT
3. Token digunakan untuk reservasi & pembayaran

---

## Payment Flow

### Dummy Payment

- Digunakan jika Midtrans belum dikonfigurasi
- Redirect ke halaman simulasi pembayaran

### Midtrans (Optional)

- Otomatis aktif jika credential tersedia
- Menggunakan Midtrans Snap

---

## Reservation Status

| Status    | Deskripsi                  |
| --------- | -------------------------- |
| pending   | Reservasi dibuat           |
| confirmed | Pembayaran sukses          |
| cancelled | Pembayaran gagal / expired |

---

## Author

Faris Sulianto
