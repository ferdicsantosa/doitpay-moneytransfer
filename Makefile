.PHONY: up down reset build loadtest verify-db e2e logs

up:
	docker compose up -d --build

down:
	docker compose down

reset:
	docker compose down -v && docker compose up -d --build

build:
	go build ./...

loadtest:
	go run ./cmd/loadtest -rps 20 -duration 5s -amount 100000

verify-db:
	docker exec doitpay-postgres psql -U postgres -d doitpay_transfer \
		-c "SELECT id, name, balance FROM merchants;" \
		-c "SELECT status, COUNT(*), COALESCE(SUM(amount), 0) AS total_amount FROM transfers GROUP BY status ORDER BY status;" \
		-c "SELECT type, COUNT(*), COALESCE(SUM(amount), 0) AS total_amount FROM balance_mutations GROUP BY type ORDER BY type;"

e2e:
	@{ \
		echo "=== E2E TEST RUN ($$(date)) ==="; \
		echo ""; \
		echo "--- 1. Health Check ---"; \
		curl -s -w "\nHTTP %{http_code}\n" http://localhost:8080/health; \
		echo ""; \
		echo "--- 2. Validation Error (400) ---"; \
		curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8080/transfer -H "Content-Type: application/json" -d '{"bank_code":"BCA","account_number":"123","amount":0}'; \
		echo ""; \
		echo "--- 3. Insufficient Balance (422) ---"; \
		curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8080/transfer -H "Content-Type: application/json" -d '{"bank_code":"BCA","account_number":"123","amount":999999999}'; \
		echo ""; \
		echo "--- 4. Idempotency Check (2 calls with same key) ---"; \
		IDEM_KEY="idem-e2e-$$(date +%s)"; \
		curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8080/transfer -H "Content-Type: application/json" -H "Idempotency-Key: $$IDEM_KEY" -d '{"bank_code":"BCA","account_number":"123456","amount":50000}'; \
		curl -s -w "\nHTTP %{http_code}\n" -X POST http://localhost:8080/transfer -H "Content-Type: application/json" -H "Idempotency-Key: $$IDEM_KEY" -d '{"bank_code":"BCA","account_number":"123456","amount":50000}'; \
		echo ""; \
		echo "--- 5. Load Test (20 RPS x 5s = 100 requests) ---"; \
		go run ./cmd/loadtest -rps 20 -duration 5s -amount 100000; \
		echo ""; \
		echo "--- 6. Waiting 8s for NSQ retries ---"; \
		sleep 8; \
		echo ""; \
		echo "--- 7. PostgreSQL Consistency Check ---"; \
		docker exec doitpay-postgres psql -U postgres -d doitpay_transfer \
			-c "SELECT id, name, balance FROM merchants;" \
			-c "SELECT status, COUNT(*), COALESCE(SUM(amount), 0) AS total_amount FROM transfers GROUP BY status ORDER BY status;" \
			-c "SELECT type, COUNT(*), COALESCE(SUM(amount), 0) AS total_amount FROM balance_mutations GROUP BY type ORDER BY type;"; \
	} | tee e2e-test.log

logs:
	docker compose logs -f api worker
