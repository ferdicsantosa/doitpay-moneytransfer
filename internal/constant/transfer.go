package constant

const (
	DefaultMerchantID int64  = 1
	DefaultPartnerID  string = "mock-partner-01"
)

type TransferStatus string

const (
	StatusPending    TransferStatus = "PENDING"
	StatusProcessing TransferStatus = "PROCESSING"
	StatusRetrying   TransferStatus = "RETRYING"
	StatusSuccess    TransferStatus = "SUCCESS"
	StatusFailed     TransferStatus = "FAILED"
)

func (s TransferStatus) IsTerminal() bool {
	return s == StatusSuccess || s == StatusFailed
}

type MutationType string

const (
	MutationDebit  MutationType = "DEBIT"
	MutationCredit MutationType = "CREDIT"
)
