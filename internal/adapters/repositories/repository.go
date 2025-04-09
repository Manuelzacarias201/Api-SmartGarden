package repositories

import (
	"database/sql"
)

type MySQLRepository struct {
	db *sql.DB
}
