package protocols

import "database/sql"

type Postgres interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
