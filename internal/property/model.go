package property

import (
	"time"
)

type Property struct {
	ID        string    `db:"id" db_opts:"primary_key"`
	CompanyID string    `db:"company_id" db_opts:"not_null"`
	Address   string    `db:"address"`
	UnitIDs   []string  `db:"unit_ids"`
	TeamIDs   []string  `db:"team_ids"`
	Owner     string    `db:"owner"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
