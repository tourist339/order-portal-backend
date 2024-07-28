package user

import "time"

type User struct {
	ID        string    `db:"id"`
	FirstName string    `db:"first_name" db_opts:"not_null"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email" db_opts:"not_null"`
	Password  string    `db:"password"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at" db_opts:"not_null"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

const TABLE_NAME = "users"
