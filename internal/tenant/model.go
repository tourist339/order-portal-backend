package tenant

import "time"

type Tenant struct {
	ID          string    `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	UnitID      string    `db:"unit_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at"`
}

const TABLE_NAME = "tenant"
