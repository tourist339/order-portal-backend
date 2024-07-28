package association

import "time"

type RolePropertyAssociation struct {
	ID          string    `db:"id"`
	RoleID      string    `db:"role_id" db_opts:"not_null"`
	PropertyID  string    `db:"property_id" db_opts:"not_null"`
	RoleType    string    `db:"role_type" db_opts:"not_null"`
	RoleSubType string    `db:"worker_subtype"`
	CreatedAt   time.Time `db:"created_at" db_opts:"not_null"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at"`
}

const TABLE_NAME = "user_property_association"
const ROLE_TYPE_TENANT = "tenant"
const ROLE_TYPE_WORKER = "worker"
const ROLE_TYPE_MANAGER = "manager"
const ROLE_TYPE_OWNER = "owner"

var AcceptableWorkerRoleTypes = []string{"cleaner", "electrician", "plumber", "painter", "carpenter", "gardener", "handyman"}

var AcceptableRoleTypes = []string{ROLE_TYPE_TENANT, ROLE_TYPE_WORKER, ROLE_TYPE_MANAGER, ROLE_TYPE_OWNER}

func init() {
	for _, workerRoleType := range AcceptableWorkerRoleTypes {
		AcceptableRoleTypes = append(AcceptableRoleTypes, "worker:"+workerRoleType)
	}
}
