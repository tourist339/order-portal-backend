package model

import (
	"context"
	"fmt"
	"strings"
	"text/template"
)

type InsertField struct {
	Field string
	Value string
}

func (s *Service) Insert(ctx context.Context, tableName string, entity any) (string, error) {
	tx, err := s.GetTransaction(ctx)
	if err != nil {
		return "", err
	}
	data, _, err := parseDataStruct(entity)
	if err != nil {
		return "", err
	}
	placeholders := make([]string, 0)
	params := make([]interface{}, len(data))
	fields := make([]string, len(data))

	i := 0
	for key, val := range data {
		fields[i] = key
		params[i] = val
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		i += 1
	}

	// Create the SQL template
	tpl := `INSERT INTO {{.TableName}} ({{.Fields}}) VALUES ({{.Placeholders}}) RETURNING id`
	tmpl, err := template.New("insert").Parse(tpl)
	if err != nil {
		return "", err
	}

	// Execute the template with the data
	sqlString := new(strings.Builder)
	err = tmpl.Execute(sqlString, map[string]interface{}{
		"TableName":    tableName,
		"Fields":       strings.Join(fields, ", "),
		"Placeholders": strings.Join(placeholders, ", "),
	})
	if err != nil {
		return "", err
	}
	println(sqlString.String())
	row := tx.QueryRowContext(ctx, sqlString.String(), params...)
	if row.Err() != nil {
		return "", row.Err()
	}
	insertID := ""
	err = row.Scan(&insertID)
	// Execute the SQL statement
	return insertID, err
}

const UPDATE_TEMPLATE = `UPDATE {{.TableName}} SET{{.Fields}} WHERE id = {{.IDParamNumber}}`

func (s *Service) Update(ctx context.Context, tableName string, data any) error {
	tx, err := s.GetTransaction(ctx)
	if err != nil {
		return err
	}
	d, id, err := parseDataStruct(data)
	if err != nil {
		return err
	}
	fields := ""
	params := make([]interface{}, len(d))
	i := 0
	for field, param := range d {
		if field == "id" {
			continue
		}
		fields += fmt.Sprintf(" %s = $%d,", field, i+1)
		params[i] = param
		i += 1
	}
	params[i] = id
	fields = strings.TrimRight(fields, ",")
	tmpl, err := template.New("update").Parse(UPDATE_TEMPLATE)
	if err != nil {
		return err
	}
	sqlString := new(strings.Builder)
	err = tmpl.Execute(sqlString, map[string]interface{}{
		"TableName":     tableName,
		"Fields":        fields,
		"IDParamNumber": fmt.Sprintf("$%d", i+1),
	})
	if err != nil {
		return err
	}
	fmt.Println("UPDATE SQL", sqlString.String())
	_, err = tx.ExecContext(ctx, sqlString.String(), params...)
	return err
}
