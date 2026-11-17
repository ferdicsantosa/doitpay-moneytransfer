package transfer

import (
	"database/sql"

	pkgnsq "github.com/ferdicsantosa/doitpay-moneytransfer/pkg/nsq"
)

func NewRepo(db *sql.DB, publisher pkgnsq.Publisher, topic string) Repository {
	return &repoModule{
		db:        db,
		publisher: publisher,
		topic:     topic,
	}
}
