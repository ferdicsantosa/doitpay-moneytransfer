package partner

import (
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/helper"
)

func NewRepo() Repository {
	return &repoModule{
		cb: helper.NewCircuitBreaker(3, 5*time.Second),
	}
}
