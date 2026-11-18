package mutation

import "database/sql"

func NewRepo(db *sql.DB) Repository {
	return &repoModule{
		db: db,
	}
}
