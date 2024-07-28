package tenant

import "time"

type Tenant struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	UnitID    string    `db:"unit_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

const TABLE_NAME = "tenant"
