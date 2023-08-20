package dbun

import "database/sql"

type SQLDBOpener interface {
	OpenDB(name string) (*sql.DB, string, error)
}
