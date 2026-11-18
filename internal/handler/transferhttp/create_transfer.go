package transferhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/entity"
	merchantrepo "github.com/ferdicsantosa/doitpay-moneytransfer/internal/repo/merchant"
	transferuc "github.com/ferdicsantosa/doitpay-moneytransfer/internal/usecase/transfer"
	"github.com/google/uuid"
)

func (h *handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req entity.CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IdempotencyKey == "" {
		req.IdempotencyKey = r.Header.Get("Idempotency-Key")
	}

	if err := validateCreateRequest(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.usecase.CreateTransfer(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, merchantrepo.ErrInsufficientBalance):
			writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, merchantrepo.ErrMerchantNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, transferuc.ErrBrokerUnavailable):
			writeJSONError(w, http.StatusServiceUnavailable, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}

func validateCreateRequest(req *entity.CreateTransferRequest) error {
	req.BankCode = strings.TrimSpace(req.BankCode)
	req.AccountNumber = strings.TrimSpace(req.AccountNumber)
	req.IdempotencyKey = strings.TrimSpace(req.IdempotencyKey)

	if req.BankCode == "" || req.AccountNumber == "" || req.Amount <= 0 {
		return fmt.Errorf("bank_code, account_number, and positive amount are required")
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = uuid.NewString()
	}
	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
