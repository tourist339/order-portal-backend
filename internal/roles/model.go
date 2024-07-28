package roles

import "time"

type Role struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id" db_opts:"not_null"`
	CreatedAt time.Time `db:"created_at" db_opts:"not_null"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

const TABLE_NAME = "role"
