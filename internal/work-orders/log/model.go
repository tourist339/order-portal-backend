package log

type Log struct {
	ID          string `db:"id"`
	OrderID     string `db:"order_id" db_opts:"not_null"`
	Status      string `db:"status" db_opts:"not_null"`
	Description string `db:"description"`
	WorkerID    string `db:"worker_id" db_opts:"not_null"`
	CreatedAt   string `db:"created_at" db_opts:"not_null"`
	UpdatedAt   string `db:"updated_at"`
	DeletedAt   string `db:"deleted_at"`
}
