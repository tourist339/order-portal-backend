package unit

import "time"

type Unit struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	PropertyID string    `db:"property_id"`
	TenantIDs  []string  `db:"tenant_ids"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	DeletedAt  time.Time `db:"deleted_at"`
}

const TABLE_NAME = "unit"
