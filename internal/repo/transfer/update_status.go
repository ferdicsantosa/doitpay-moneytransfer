package transfer

import (
	"context"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/constant"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func (r *repoModule) UpdateStatus(ctx context.Context, transferID string, status constant.TransferStatus, retryCount int, failureReason *string) error {
	const query = `
		UPDATE transfers
		SET status = $1,
		    retry_count = $2,
		    failure_reason = $3,
		    updated_at = NOW()
		WHERE id = $4
		  AND status NOT IN ('SUCCESS', 'FAILED')
	`

	exec := helper.GetExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, query, status, retryCount, failureReason, transferID)
	return err
}
