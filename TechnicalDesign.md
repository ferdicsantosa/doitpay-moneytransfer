# Technical Design — Money Transfer (Disbursement) Service

## 1. Overview & Goals
Service ini adalah backend disbursement/money transfer yang menerima permintaan transfer dari merchant melalui HTTP API, memotong saldo merchant secara **concurrency-safe (ACID)** di awal, lalu memproses pengiriman dana ke payment provider secara **asynchronous** menggunakan message broker.

### Core Guarantees
- **No Race Condition / No Negative Balance**: Pemotongan saldo dilakukan secara atomik di dalam database transaction.
- **Safe Rollback**: Jika proses transfer gagal secara permanen (atau gagal publish ke broker), saldo merchant dikembalikan (*refund*) secara idempotent (maksimal 1x).
- **At-Least-Once Processing Safe**: Consumer aman terhadap duplikasi pesan maupun retry dari message broker.

---

## 2. Tech Stack
| Komponen | Pilihan | Alasan |
| :--- | :--- | :--- |
| **Language** | Golang | Performa tinggi, concurrency-friendly, idiomatic untuk microservice & NSQ. |
| **Database** | PostgreSQL | Mendukung ACID transaction dan row-level locking (`SELECT ... FOR UPDATE`) untuk menjamin konsistensi saldo. |
| **Message Broker** | NSQ (`nsqd` + `nsqlookupd`) | Ringan, native Go (`go-nsq`), mendukung requeue dengan backoff dan `MaxAttempts` untuk retry. |
| **Cache / Redis** | *Tidak digunakan* | Lihat [Section 6 (Keputusan Arsitektur: Redis)](#6-keputusan-arsitektur-apakah-perlu-redis). |

---

## 3. Architecture & Flows

### A. Main Flow (Happy Path)
```
Merchant/Client
   │
   │ POST /transfer (Idempotency-Key, amount, bank_code, account_number)
   ▼
[API Service (Go)]
   ├─ 1. Validasi payload (amount > 0, field wajib lengkap)
   ├─ 2. BEGIN DB Transaction
   │     ├─ Cek idempotency_key (jika sudah ada, return transaksi eksisting)
   │     ├─ Lock row merchant (`SELECT balance FROM merchants WHERE id = $1 FOR UPDATE`)
   │     ├─ Validasi `balance >= amount` (jika kurang -> ROLLBACK & return 422/400)
   │     ├─ Deduct balance (`UPDATE merchants SET balance = balance - $1 ...`)
   │     ├─ Insert ledger mutation (`DEBIT`)
   │     └─ Insert `transfers` (status = `PENDING`)
   ├─ 3. COMMIT DB Transaction
   ├─ 4. Publish message (`transfer_id`) ke NSQ topic `disbursement.process`
   └─ 5. Return `202 Accepted` (transfer_id, status: `PENDING`)
          │
          ▼
   [NSQ Broker]
          │
          ▼
[Background Worker / Consumer (Go)]
   ├─ 1. Ambil `transfers` dari DB; cek apakah status masih `PENDING` / `RETRYING`
   │     └─ Jika sudah terminal (`SUCCESS` / `FAILED`), ACK message & skip (idempotent)
   ├─ 2. Update status -> `PROCESSING`
   ├─ 3. Simulasi call ke Mock Payment Provider (random delay + random response)
   └─ 4. Jika Provider SUCCESS:
         ├─ Update `transfers` status -> `SUCCESS` (atomic check status != terminal)
         └─ Trigger dummy notification & dummy merchant webhook (async/best-effort)
```

### B. Negative & Failure Flows
1. **Validasi Input Gagal / Saldo Tidak Cukup**
   - Return `400 Bad Request` / `422 Unprocessable Entity`. Saldo tidak berubah.
2. **Gagal Publish ke NSQ (setelah DB Commit)**
   - **Mitigasi**: Agar dana tidak menggantung (*stranded balance*), API service langsung melakukan kompensasi: update status transaksi -> `FAILED` (`reason: broker_publish_failed`), refund saldo merchant, catat ledger `CREDIT`, lalu return `503 Service Unavailable`.
   - *(Production note: Dapat ditingkatkan menggunakan Transactional Outbox Pattern).*
3. **Consumer Menerima Pesan dengan Status Terminal (`SUCCESS` / `FAILED`)**
   - Jangan return error ke NSQ (karena akan memicu requeue berulang). Log warning dan **ACK (`return nil`)** agar pesan tidak berputar.
4. **Error dari Payment Provider — Retryable (misal timeout, 5xx, network error)**
   - Jika `attempt < max_retry`: update status -> `RETRYING`, increment `retry_count`, lalu `Requeue` message ke NSQ dengan backoff delay.
5. **Error dari Payment Provider — Non-Retryable (misal invalid account) atau Retry Habis (`attempt >= max_retry`)**
   - Jalankan **DB Transaction** untuk rollback saldo:
     - Lock row `transfers` (`FOR UPDATE`), pastikan status belum `FAILED`/`SUCCESS` (mencegah *double refund*).
     - Update `transfers` status -> `FAILED`, simpan `failure_reason`.
     - Kembalikan saldo: `UPDATE merchants SET balance = balance + $1`.
     - Insert ledger mutation (`CREDIT` / `REFUND`).
   - Trigger dummy notification & dummy merchant webhook dengan status `FAILED`.
   - ACK message ke NSQ.

---

## 4. Transaction State Machine

Menggunakan string status yang eksplisit agar mudah dibaca dan di-debug (dibandingkan angka yang mirip HTTP status code):

| Status | Deskripsi | Terminal? | Transisi Berikutnya |
| :--- | :--- | :---: | :--- |
| `PENDING` | Request valid, saldo sudah dipotong, menunggu diproses worker. | Tidak | `PROCESSING`, `FAILED` |
| `PROCESSING` | Worker sedang memanggil payment provider. | Tidak | `SUCCESS`, `RETRYING`, `FAILED` |
| `RETRYING` | Provider gagal (retryable), message di-requeue ke NSQ. | Tidak | `PROCESSING`, `FAILED` |
| `SUCCESS` | Transfer berhasil diproses oleh provider. | **Ya** | - |
| `FAILED` | Transfer gagal permanen / retry habis; saldo merchant sudah dikembalikan. | **Ya** | - |

> **Catatan**: Validasi request yang gagal di API layer (payload tidak valid / saldo tidak cukup) langsung dikembalikan sebagai HTTP error response tanpa perlu masuk ke tabel transaksi.

---

## 5. Database Schema (PostgreSQL)

> **Catatan Tipe Data**: Nominal uang menggunakan `BIGINT` (satuan terkecil, misal Rupiah penuh atau sen) — **tidak menggunakan `FLOAT`/`DOUBLE`** untuk menghindari masalah presisi.

### Table: `merchants`
| Column | Type | Constraint | Keterangan |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT` | `PRIMARY KEY` | Hardcoded `1` untuk MVP |
| `name` | `VARCHAR(100)` | `NOT NULL` | Nama merchant |
| `balance` | `BIGINT` | `NOT NULL, CHECK (balance >= 0)` | Guard di level DB agar saldo tidak pernah negatif |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | |

### Table: `transfers`
| Column | Type | Constraint | Keterangan |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | `PRIMARY KEY` | ID transaksi internal |
| `merchant_id` | `BIGINT` | `NOT NULL REFERENCES merchants(id)` | |
| `idempotency_key` | `VARCHAR(100)` | `NOT NULL` | Mencegah double deduction pada request duplikat |
| `partner_id` | `VARCHAR(50)` | `NOT NULL` | ID mock payment provider |
| `bank_code` | `VARCHAR(20)` | `NOT NULL` | Kode bank tujuan |
| `account_number` | `VARCHAR(50)` | `NOT NULL` | Nomor rekening tujuan |
| `amount` | `BIGINT` | `NOT NULL, CHECK (amount > 0)` | Nominal transfer |
| `status` | `VARCHAR(20)` | `NOT NULL` | `PENDING`, `PROCESSING`, `RETRYING`, `SUCCESS`, `FAILED` |
| `retry_count` | `INT` | `NOT NULL DEFAULT 0` | Jumlah percobaan retry |
| `failure_reason` | `TEXT` | `NULL` | Alasan gagal (jika `FAILED`) |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | |

*Index*: `UNIQUE INDEX idx_transfers_merchant_idempotency ON transfers (merchant_id, idempotency_key)`

### Table: `balance_mutations` (Audit Trail / Ledger)
| Column | Type | Constraint | Keterangan |
| :--- | :--- | :--- | :--- |
| `id` | `BIGINT` | `PRIMARY KEY GENERATED ALWAYS AS IDENTITY` | |
| `merchant_id` | `BIGINT` | `NOT NULL REFERENCES merchants(id)` | |
| `transfer_id` | `UUID` | `NOT NULL REFERENCES transfers(id)` | |
| `type` | `VARCHAR(20)` | `NOT NULL` | `DEBIT` (potong saldo) atau `CREDIT` (refund) |
| `amount` | `BIGINT` | `NOT NULL` | |
| `balance_before` | `BIGINT` | `NOT NULL` | |
| `balance_after` | `BIGINT` | `NOT NULL` | |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | |

---

## 6. Keputusan Arsitektur: Apakah Perlu Redis?

**Keputusan: Tidak diperlukan untuk scope sistem saat ini.**

Alasan:
1. **Source of Truth di PostgreSQL**: Concurrency control untuk pemotongan saldo sudah dijamin oleh ACID transaction dan row-level lock (`SELECT ... FOR UPDATE`) di PostgreSQL. Menaruh saldo di Redis justru berisiko menimbulkan *data inconsistency* antara cache dan database.
2. **Tidak Perlu Distributed Lock**: Selama transaksi database berjalan di satu cluster PostgreSQL, lock dari database sudah cukup menangani konkurensi tinggi (jauh di atas target minimal 10 RPS) tanpa menambah *failure point* baru dari Redis lock (seperti lock timeout/expiry).
3. **Idempotency**: Sudah bisa ditangani secara atomik dengan `UNIQUE INDEX (merchant_id, idempotency_key)` di PostgreSQL.

**Kapan Redis relevan di Production?** (Future Improvement):
- Rate limiting per merchant di API Gateway/middleware.
- Caching metadata merchant (nama, status aktif, callback URL) agar tidak selalu query ke DB.
- Fast-path idempotency check untuk request duplikat beruntun.

## 7. Cutting Corner
- all dummy partner is going to one repo, supposedly to have separates repo
- not implementing authorization, either to check the jwt, or also implement authorization to partner
- no full test coverage, no unit test