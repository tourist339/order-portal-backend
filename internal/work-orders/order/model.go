package order

type Order struct {
	ID          string   `db:"id"`
	PropertyID  string   `db:"property_id" db_opts:"not_null"`
	UnitID      string   `db:"unit_id"`
	WorkerID    string   `db:"worker_id"`
	TenantID    string   `db:"tenant_id" db_opts:"not_null"`
	Description string   `db:"description" db_opts:"not_null"`
	Pictures    []string `db:"picture"`
	Video       string   `db:"video"`
	CreatedAt   string   `db:"created_at" db_opts:"not_null"`
	UpdatedAt   string   `db:"updated_at"`
	DeletedAt   string   `db:"deleted_at"`
}
