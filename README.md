# Money Transfer (Disbursement) Service

A concurrency-safe, asynchronous money transfer (disbursement) backend service built in **Go** using Clean Architecture, **PostgreSQL**, **NSQ**, and **Redis**.

For detailed architecture diagrams, flows, and schema design, see [TechnicalDesign.md](./TechnicalDesign.md).

---

## Architecture Overview

- **API Service (`cmd/api`)**: Exposes `POST /transfer` and `GET /health`.
  - Validates the request payload at the HTTP handler layer.
  - Performs a lock-free idempotency pre-check (Redis cache with graceful fallback to PostgreSQL) so duplicate retries do not contend on database row locks.
  - Deducts merchant balance atomically inside a PostgreSQL transaction (`SELECT ... FOR UPDATE` + `CHECK (balance >= 0)`), inserts the `transfers` record (`PENDING`), and logs a `DEBIT` entry in `balance_mutations`.
  - Publishes `transfer_id` to NSQ (`disbursement.process`), compensating with an automatic refund if publishing fails.
- **Worker Service (`cmd/worker`)**: Consumes messages from NSQ concurrently.
  - Acquires a Redis distributed lock (`SET NX EX` + Lua ownership release) per `transfer_id` to prevent concurrent duplicate disbursements across workers.
  - Simulates an external payment partner (`100–800ms` delay, random success / retryable error / non-retryable error).
  - Protects partner calls with a **Circuit Breaker** (`CLOSED`/`OPEN`/`HALF-OPEN`) and retries transient errors via NSQ **exponential backoff with jitter**.
  - Updates the transfer to `SUCCESS` or retries via NSQ requeue on transient errors.
  - On permanent failure or exhausted retries, atomically marks the transfer `FAILED`, refunds the merchant balance, and records a `CREDIT` entry in `balance_mutations`.
  - Invalidates the Redis idempotency cache on terminal completion so subsequent retries receive the latest status.

---

## Prerequisites

- **Docker** & **Docker Compose**
- **Go 1.26+** (only required if running the load test or services locally outside Docker)

---

## Quick Start (Docker Compose)

Start the entire stack (`postgres`, `nsqd`, `redis`, `api`, and `worker`):

```bash
make up
# or: docker compose up -d --build
```

Verify the API is healthy:

```bash
curl -s http://localhost:8080/health
# {"status":"ok"}
```

The database is automatically initialized via `migrations/001_init.sql` with:
- Tables: `merchants`, `transfers`, `balance_mutations`
- Seeded merchant: `id = 1`, `name = 'Merchant Utama'`, `balance = 100000000` (100,000,000)

To reset the database and start fresh, view logs, or stop all containers:

```bash
make reset   # reset volumes & rebuild
make logs    # tail API & worker logs
make down    # stop containers
```

---

## API Usage

### `POST /transfer`

```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: req-001" \
  -d '{
    "bank_code": "BCA",
    "account_number": "1234567890",
    "amount": 100000
  }'
```

Response (`202 Accepted`):

```json
{
  "transfer_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING",
  "amount": 100000
}
```

---

## Load Test (`>= 10 RPS`)

Run the included concurrent load tester (defaults to `20 RPS` for `5s`, i.e. 100 concurrent requests), or run the full E2E suite (logged to `e2e-test.log`):

```bash
make loadtest
# or full E2E suite (health, 400, 422, idempotency, 20 RPS burst, DB check):
make e2e
```

Verify data consistency in PostgreSQL after the worker finishes processing:

```bash
make verify-db
```

Consistency invariant:
$$\text{merchants.balance} = 100{,}000{,}000 - \text{SUM}(\text{DEBIT}) + \text{SUM}(\text{CREDIT})$$

---

## Assumptions, Corners Cut & Future Improvements

### Assumptions
- Single hardcoded merchant (`merchant_id = 1`) with an initial balance of `100,000,000`.
- Currency amounts are represented as `BIGINT` (integer units) to avoid floating-point precision issues.
- If `idempotency_key` is omitted in both the request body and `Idempotency-Key` header, the API generates a UUID so burst tests without explicit keys work out of the box.

### Corners Cut (for Take-Home Scope)
- **No Authentication / Request Signing**: Merchant auth (JWT/API key) and partner webhook signatures are omitted.
- **Combined Mock Partner**: Mock disbursement and dummy webhook/notification live in a single `internal/repo/partner` package.
- **Limited Automated Tests**: Focused on clean architecture, concurrency safety, and end-to-end load testing rather than full unit test coverage.
- **Single-Node NSQ Broker**: Chose **NSQ** over Kafka for fast local setup, low resource footprint for reviewers, and native per-message delayed requeue (`RequeueWithoutBackoff`) for exponential backoff + jitter.
- **Best-Effort Publish Compensation**: If NSQ publishing fails after DB commit, the API compensates via `refundTransfer` (logging `CRITICAL` if compensation also fails).

### Production Improvements
- **Transactional Outbox Pattern & Reconciliation Sweeper**: Insert outbox events atomically within the balance-deduction transaction and relay them via a background publisher to eliminate stranded balances if the broker is down.
- **Dead Letter Queue (DLQ) & Worker Timeout**: Wrap each worker message in `context.WithTimeout` and route messages that exhaust `max_attempts` due to infra errors to a DLQ + alerting.
- **HTTP Hardening**: Add `http.MaxBytesReader` on request bodies, rate limiting, and merchant authentication/request signing.
- **Attempt History & Observability**: Add Prometheus/OpenTelemetry metrics, structured JSON logging, and an append-only `transfer_attempts` table (or log pipeline) for partner dispute auditing.
- **High-Throughput Scaling & Broker Migration**: `SELECT ... FOR UPDATE` on `merchants` comfortably handles hundreds of TPS per merchant. For extreme hot-row contention (thousands of RPS on a single merchant), adopt **Balance Sharding** or a **Partitioned Single-Writer** architecture. For multi-broker replication (`acks=all`), strict per-merchant ordering, and event replayability for financial auditing, migrate from single-node NSQ to **Kafka** (with retry topics + DLQ).
- **Automated Test Coverage**: Add unit tests (mocking repo/helper interfaces) and integration tests covering race conditions, idempotency, retries, and refunds.
