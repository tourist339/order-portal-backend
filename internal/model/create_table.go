package model

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

const PRIMARY_KEY_OPT = "primary_key"
const NOT_NULL_OPT = "not_null"

// TODO: implement options which will include constraints like UNIQUE and stuff like indexes
type Option func()
type Table struct {
	Name   string
	Fields []TemplateField
}
type TemplateField struct {
	Name string
	FieldArgs
}
type Fields map[string]FieldArgs
type FieldArgs struct {
	Type    string
	Options string
}

const sqlTemplate = `
CREATE TABLE public.{{ .Name }} (
	{{ range $index, $field := .Fields }}
	{{ $field.Name }} {{ $field.Type }}{{ $field.Options }},
	{{ end }}
    PRIMARY KEY (id)
);`

func (s *Service) CreateTable(tableName string, u any, opts ...Option) error {
	table := Table{
		Name: tableName,
	}
	fields, err := parseTableFields(u)
	if err != nil {
		return err
	}
	templateFields := make([]TemplateField, len(fields))
	i := 0
	for k, v := range fields {
		templateFields[i] = TemplateField{
			Name:      k,
			FieldArgs: v,
		}
		i += 1
	}
	table.Fields = templateFields
	b := &bytes.Buffer{}
	tmpl := template.Must(template.New("sqlTable").Funcs(template.FuncMap{
		"lastField": func(index int, fields []TemplateField) bool {
			return index == len(fields)-1
		},
	}).Parse(sqlTemplate))
	err = tmpl.Execute(b, table)
	if err != nil {
		return err
	}
	fmt.Printf("Create Query for table %s: %s\n", tableName, b.String())

	_, err = s.db.Exec(b.String())
	return err
}

func parseTableFields(u any) (Fields, error) {
	fields := make(map[string]FieldArgs)
	if reflect.TypeOf(u).Kind() != reflect.Ptr || reflect.TypeOf(u).Elem().Kind() != reflect.Struct {
		return nil, errors.New("Table struct is not a pointer to a struct")
	}
	typ := reflect.TypeOf(u).Elem()
	for i := 0; i < typ.NumField(); i++ {
		s := typ.Field(i)
		name := s.Tag.Get("db")
		if name == "" {
			return nil, errors.New(fmt.Sprintf("Table struct key %s does not have a db tag", s.Name))
		}
		opts := s.Tag.Get("db_opts")
		optionVal := ""
		if opts != "" {
			optMap := make(map[string]struct{})
			for _, k := range strings.Split(opts, ",") {
				optMap[k] = struct{}{}
			}
			if _, ok := optMap[NOT_NULL_OPT]; ok {
				optionVal = " NOT NULL"
			}
		}
		fieldType, err := goTypeToPostgresType(s.Type)
		if err != nil {
			return nil, err
		}
		fields[name] = FieldArgs{
			Type:    fieldType,
			Options: optionVal,
		}
	}
	if _, ok := fields["id"]; !ok {
		return nil, errors.New("Table struct must have a field 'id'")
	}
	return fields, nil
}

func goTypeToPostgresType(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.String:
		return "text", nil
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "integer", nil
	case reflect.Slice:
		typ, err := goTypeToPostgresType(t.Elem())
		if err != nil {
			return "", err
		}
		return typ + "[]", nil
	}
	if t == reflect.TypeOf(time.Time{}) {
		return "time with time zone", nil
	}
	return "", errors.New("Unsupported type: " + t.String())
}
