# API Documentation

Dokumentasi lengkap untuk Pilates Reservation API.

## Base URL

```
http://localhost:8080
```

## 📋 Endpoints

### 1. Authentication

#### Register User

Mendaftarkan user baru.

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "phone": "081234567890"
}
```

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "081234567890",
      "is_active": true,
      "created_at": "2026-01-23T10:00:00Z"
    }
  }
}
```

**Error Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "Email already exists"
}
```

---

#### Login

Login user yang sudah terdaftar.

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

---

### 2. Browse Available Slots

#### Get Available Dates

Mendapatkan daftar tanggal yang tersedia (30 hari ke depan).

```http
GET /api/v1/dates
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "dates": [
      "2026-01-23",
      "2026-01-24",
      "2026-01-25",
      ...
    ]
  }
}
```

---

#### Get Timeslots

Mendapatkan daftar timeslot. Jika date disertakan, akan menampilkan availability info.

```http
GET /api/v1/timeslots
GET /api/v1/timeslots?date=2026-01-25
```

**Query Parameters:**

- `date` (optional): Format YYYY-MM-DD

**Success Response (200 OK):**

Without date:

```json
{
  "success": true,
  "data": {
    "timeslots": [
      {
        "id": 1,
        "time": "08:00",
        "duration": 60,
        "is_active": true
      },
      ...
    ]
  }
}
```

With date:

```json
{
  "success": true,
  "data": {
    "timeslots": [
      {
        "id": 1,
        "time": "08:00",
        "duration": 60,
        "is_active": true,
        "available": true,
        "booked_count": 1,
        "available_courts": 2
      },
      ...
    ]
  }
}
```

---

#### Get Available Courts

Mendapatkan daftar court yang tersedia untuk tanggal dan timeslot tertentu.

```http
GET /api/v1/courts?date=2026-01-25&timeslot_id=1
```

**Query Parameters:**

- `date` (required): Format YYYY-MM-DD
- `timeslot_id` (required): ID timeslot

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "courts": [
      {
        "id": 1,
        "name": "Studio A",
        "capacity": 10,
        "description": "Reformer Pilates - Premium equipment",
        "is_active": true,
        "available": true
      },
      {
        "id": 2,
        "name": "Studio B",
        "capacity": 8,
        "description": "Mat Pilates - Classic exercises",
        "is_active": true,
        "available": false
      }
    ]
  }
}
```

---

### 3. Reservations (Protected)

#### Create Reservation

Membuat reservasi baru. Requires authentication.

```http
POST /api/v1/reservations
Authorization: Bearer <token>
Content-Type: application/json

{
  "court_id": 1,
  "timeslot_id": 1,
  "date": "2026-01-25",
  "notes": "First time trying Pilates"
}
```

**Success Response (201 Created):**

```json
{
  "success": true,
  "message": "Reservation created successfully. Please proceed to payment.",
  "data": {
    "reservation": {
      "id": 1,
      "user_id": 1,
      "court_id": 1,
      "timeslot_id": 1,
      "date": "2026-01-25T00:00:00Z",
      "status": "pending",
      "notes": "First time trying Pilates",
      "court": {
        "id": 1,
        "name": "Studio A",
        "capacity": 10
      },
      "timeslot": {
        "id": 1,
        "time": "08:00",
        "duration": 60
      },
      "created_at": "2026-01-23T10:30:00Z"
    }
  }
}
```

**Error Responses:**

```json
// 400 - Past date
{
  "success": false,
  "error": "Cannot book past dates"
}

// 409 - Already booked
{
  "success": false,
  "error": "Court is already booked for this timeslot"
}
```

---

#### Get User Reservations

Mendapatkan semua reservasi user yang sedang login.

```http
GET /api/v1/reservations
Authorization: Bearer <token>
```

**Query Parameters:**

- `status` (optional): Filter by status (pending, confirmed, cancelled, completed)
- `from_date` (optional): Filter from date (YYYY-MM-DD)
- `to_date` (optional): Filter to date (YYYY-MM-DD)

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "reservations": [
      {
        "id": 1,
        "court": { "name": "Studio A" },
        "timeslot": { "time": "08:00" },
        "date": "2026-01-25T00:00:00Z",
        "status": "pending",
        "payment": {
          "amount": 100000,
          "status": "pending"
        }
      }
    ]
  }
}
```

---

#### Get Single Reservation

Mendapatkan detail reservasi berdasarkan ID.

```http
GET /api/v1/reservations/:id
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "reservation": {
      "id": 1,
      "user_id": 1,
      "court": {
        "id": 1,
        "name": "Studio A",
        "description": "Reformer Pilates"
      },
      "timeslot": {
        "id": 1,
        "time": "08:00",
        "duration": 60
      },
      "date": "2026-01-25T00:00:00Z",
      "status": "pending",
      "payment": null
    }
  }
}
```

---

#### Cancel Reservation

Membatalkan reservasi.

```http
PUT /api/v1/reservations/:id/cancel
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Reservation cancelled successfully",
  "data": {
    "reservation": {
      "id": 1,
      "status": "cancelled"
    }
  }
}
```

---

### 4. Payments (Protected)

#### Create Payment

Membuat transaksi pembayaran untuk reservasi.

```http
POST /api/v1/payments/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "reservation_id": 1
}
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Payment created successfully",
  "data": {
    "payment": {
      "id": 1,
      "reservation_id": 1,
      "amount": 100000,
      "status": "pending",
      "transaction_id": "TRX-abc123-1234567890"
    },
    "payment_url": "https://app.sandbox.midtrans.com/snap/v3/redirection/...",
    "snap_token": "snap_token_here",
    "client_key": "SB-Mid-client-..."
  }
}
```

**Usage:**

1. Redirect user ke `payment_url`, atau
2. Gunakan Snap.js dengan `snap_token` untuk embedded payment

---

#### Payment Callback

Webhook dari Midtrans untuk update status pembayaran.

```http
POST /api/v1/payments/callback
Content-Type: application/json

{
  "order_id": "TRX-abc123-1234567890",
  "transaction_status": "settlement",
  "transaction_id": "midtrans_trx_id",
  "status_code": "200",
  "gross_amount": "100000.00"
}
```

**Transaction Status:**

- `capture`, `settlement` → Payment successful
- `pending` → Payment pending
- `deny`, `expire`, `cancel` → Payment failed

---

#### Get Payment

Mendapatkan detail payment berdasarkan ID.

```http
GET /api/v1/payments/:id
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "payment": {
      "id": 1,
      "reservation_id": 1,
      "amount": 100000,
      "status": "paid",
      "transaction_id": "TRX-abc123-1234567890",
      "paid_at": "2026-01-23T10:35:00Z",
      "reservation": {
        "id": 1,
        "status": "confirmed"
      }
    }
  }
}
```

---

### 5. User Profile (Protected)

#### Get Profile

Mendapatkan profil user yang sedang login.

```http
GET /api/v1/profile
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "081234567890",
      "is_active": true,
      "created_at": "2026-01-20T10:00:00Z"
    }
  }
}
```

---

#### Update Profile

Update profil user.

```http
PUT /api/v1/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Updated",
  "phone": "081234567899"
}
```

---

### 6. Admin Endpoints

#### Get All Courts

```http
GET /api/v1/admin/courts
```

#### Create Court

```http
POST /api/v1/admin/courts
Content-Type: application/json

{
  "name": "Studio D",
  "capacity": 15,
  "description": "Group Class Studio",
  "is_active": true
}
```

#### Update Court

```http
PUT /api/v1/admin/courts/:id
Content-Type: application/json

{
  "capacity": 20
}
```

#### Delete Court

```http
DELETE /api/v1/admin/courts/:id
```

#### Similar endpoints for Timeslots

```http
GET    /api/v1/admin/timeslots
POST   /api/v1/admin/timeslots
PUT    /api/v1/admin/timeslots/:id
DELETE /api/v1/admin/timeslots/:id
```

#### Get Statistics

```http
GET /api/v1/admin/stats
```

---

## 🔒 Error Codes

| Code | Meaning                                 |
| ---- | --------------------------------------- |
| 200  | Success                                 |
| 201  | Created                                 |
| 400  | Bad Request - Invalid input             |
| 401  | Unauthorized - Invalid or missing token |
| 403  | Forbidden - No permission               |
| 404  | Not Found                               |
| 409  | Conflict - Resource already exists      |
| 500  | Internal Server Error                   |

---

## 📝 Notes

1. **Token Expiration**: JWT tokens expire after 24 hours
2. **Date Format**: Always use `YYYY-MM-DD` for dates
3. **Time Format**: Timeslots use `HH:MM` format (24-hour)
4. **Amount**: All amounts in IDR (Indonesian Rupiah)
5. **CORS**: Configured for `localhost:3000` and `localhost:3001`

---

## 🧪 Testing with cURL

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@test.com","password":"test123","phone":"081234567890"}'

# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}' \
  | jq -r '.data.token')

# Create reservation
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"court_id":1,"timeslot_id":1,"date":"2026-01-25"}'
```
