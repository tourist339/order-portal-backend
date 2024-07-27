package model

import (
	"bytes"
	"text/template"
)

type Table struct {
	Name   string
	Fields []Column
}

type Column struct {
	Name     string
	DataType string
	Options  string // e.g., NULL, AUTO_INCREMENT, etc.
}

const sqlTemplate = `
CREATE TABLE {{ .Name }} (
	{{ range $index, $field := .Fields }}
	{{ $field.Name }} {{ $field.DataType }} {{ $field.Options }}{{ if not (lastField $index $.Fields) }},{{ end }}
	{{ end }}
);`

func (s *Service) CreateTable(schema Table) error {
	b := &bytes.Buffer{}
	tmpl := template.Must(template.New("sqlTable").Funcs(template.FuncMap{
		"lastField": func(index int, fields []Column) bool {
			return index == len(fields)-1
		},
	}).Parse(sqlTemplate))
	err := tmpl.Execute(b, schema)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(b.String())
	return err
}
