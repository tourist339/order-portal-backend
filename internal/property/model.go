package property

import (
	"backend/internal/model"
	"time"
)

var Table model.Table = model.Table{
	Name: "property",
	Fields: []model.Column{
		{
			Name:     "id",
			DataType: "VARCHAR(255)",
			Options:  "PRIMARY KEY",
		},
		{
			Name:     "company_id",
			DataType: "VARCHAR(255)",
			Options:  "NOT NULL UNIQUE",
		},
		{
			Name:     "address",
			DataType: "VARCHAR(255)",
			Options:  "NOT NULL",
		},
		{
			Name:     "unit_ids",
			DataType: "VARCHAR(255)[]",
		},
		{
			Name:     "team_ids",
			DataType: "VARCHAR(255)[]",
		},
	},
}

type Property struct {
	ID        string    `db:"id"`
	CompanyID string    `db:"company_id"`
	Address   string    `db:"address"`
	UnitIDs   []string  `db:"unit_ids"`
	TeamIDs   []string  `db:"team_ids"`
	Owner     string    `db:"owner"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
