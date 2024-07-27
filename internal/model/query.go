package model

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"reflect"
	"text/template"
	"time"
)

var ErrNotFound = errors.New("Not Found")

type SelectQuery struct {
	TableName string
	Fields    []string
	Where     []Condition // Dynamic WHERE clause conditions
}

type Condition struct {
	Clause string // Raw SQL query string for WHERE clause with placeholders
	Param  string // Parameters for dynamic WHERE conditions
}

const selectTemplate = `
SELECT {{ range $index, $field := .Fields }}
	{{ $field }}{{ if not (lastField $index $.Fields) }},{{ end }}
{{ end }}
FROM {{ .TableName }}
{{ if .Where }}
WHERE {{ range $index, $condition := .Where }}
	{{ $condition.Clause }}{{ if not (lastCondition $index $.Where) }} AND {{ end }}
{{ end }}
;
{{ else }}
;
{{ end }}
`

func (s *Service) GetByID(ctx context.Context, id, tableName string, fields []string, u any) error {
	return s.Get(ctx, &SelectQuery{
		TableName: tableName,
		Fields:    fields,
		Where:     []Condition{{Clause: "id=$1", Param: id}},
	}, u)
}

func (s *Service) Get(ctx context.Context, q *SelectQuery, u any) error {
	rows, err := s.getRows(ctx, q)
	if err != nil {
		return err
	}
	if rows.Next() == false {
		return ErrNotFound
	}
	nullStruct := generateNullTypeStruct(u)
	k := reflect.New(nullStruct).Interface()
	err = rows.StructScan(k)
	generateNormalStruct(k, u)
	if err != nil {
		fmt.Println(fmt.Errorf("Error Scanning rows into struct %s", err.Error()))
	}
	defer rows.Close()
	return err
}

func (s *Service) getRows(ctx context.Context, q *SelectQuery) (*sqlx.Rows, error) {
	tx, err := s.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}
	query := &bytes.Buffer{}
	tmpl := template.Must(template.New("selectQuery").Funcs(template.FuncMap{
		"lastField": func(index int, fields []string) bool {
			return index == len(fields)-1
		},
		"lastCondition": func(index int, conditions []Condition) bool {
			return index == len(conditions)-1
		},
	}).Parse(selectTemplate))

	err = tmpl.Execute(query, q)
	if err != nil {
		return nil, err
	}
	var params []interface{}
	for _, param := range q.Where {
		params = append(params, param.Param)
	}
	fmt.Println(query.String())
	rows, err := tx.QueryxContext(ctx, query.String(), params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func generateNullTypeStruct(u any) reflect.Type {
	typ := reflect.TypeOf(u)
	if typ.Kind() != reflect.Pointer {
		panic("Not a Pointer")
	}
	elem := typ.Elem()
	if elem.Kind() != reflect.Struct {
		panic("Not a Struct")
	}
	var _genInternal func(t any) reflect.Type
	_genInternal = func(t any) reflect.Type {
		tp, ok := t.(reflect.Type)
		if !ok {
			panic("Not a type")
		}
		fields := make([]reflect.StructField, 0, tp.NumField())
		for i := 0; i < tp.NumField(); i++ {
			temp := tp.Field(i)
			if temp.Type == reflect.TypeOf(time.Time{}) {
				temp.Type = reflect.TypeOf(sql.NullTime{})
			} else if temp.Type.Kind() == reflect.Struct {
				temp.Type = reflect.TypeOf(_genInternal(temp.Type))
			} else {
				temp.Type = convertToNullType(temp.Type)
			}
			fields = append(fields, temp)
		}
		return reflect.StructOf(fields)
	}
	return _genInternal(elem)
}

func generateNormalStruct(in, dest any) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr || reflect.TypeOf(in).Elem().Kind() != reflect.Struct {
		panic("Input is not a pointer to a struct")
	}
	if reflect.TypeOf(dest).Kind() != reflect.Ptr || reflect.TypeOf(dest).Elem().Kind() != reflect.Struct {
		panic("Destination is not a pointer to a struct")
	}
	inVal := reflect.ValueOf(in).Elem()
	inTyp := reflect.TypeOf(in).Elem()
	destVal := reflect.ValueOf(dest).Elem()
	for i := 0; i < inVal.NumField(); i++ {
		inFieldVal := inVal.Field(i)
		inFieldTyp := inTyp.Field(i)
		destField := destVal.FieldByName(inFieldTyp.Name)
		if !destField.IsValid() {
			panic(fmt.Sprintf("%s inp field not found in dest struct", inFieldTyp.Name))
		}
		destField.Set(convertNullToNormalValue(inFieldVal.Interface()))

	}

}
func convertNullToNormalValue(inp any) reflect.Value {
	switch inp.(type) {
	case sql.NullString:
		return reflect.ValueOf(inp.(sql.NullString).String)
	case sql.NullInt64:
		return reflect.ValueOf(int(inp.(sql.NullInt64).Int64))
	case sql.NullInt32:
		return reflect.ValueOf(int(inp.(sql.NullInt32).Int32))
	case sql.NullTime:
		return reflect.ValueOf(inp.(sql.NullTime).Time)
	case sql.NullBool:
		return reflect.ValueOf(inp.(sql.NullBool).Bool)
	default:
		return reflect.ValueOf(inp)
	}
}
func generateTypeStructFromNull(u any) reflect.Type {
	typ := reflect.TypeOf(u)
	if typ.Kind() != reflect.Pointer {
		panic("Not a Pointer")
	}
	elem := typ.Elem()
	if elem.Kind() != reflect.Struct {
		panic("Not a Struct")
	}
	var _genInternal func(t any) reflect.Type
	_genInternal = func(t any) reflect.Type {
		tp, ok := t.(reflect.Type)
		if !ok {
			panic("Not a type")
		}
		fields := make([]reflect.StructField, 0, tp.NumField())
		for i := 0; i < tp.NumField(); i++ {
			temp := tp.Field(i)
			temp.Type = convertNullTypeToType(temp.Type)
			fields = append(fields, temp)
		}
		return reflect.StructOf(fields)
	}
	return _genInternal(elem)
}

func convertNullTypeToType(t reflect.Type) reflect.Type {
	switch t {
	case reflect.TypeOf(sql.NullString{}):
		return reflect.TypeOf("")
	case reflect.TypeOf(sql.NullInt32{}):
		return reflect.TypeOf(0)
	case reflect.TypeOf(sql.NullInt64{}):
		return reflect.TypeOf(0)
	case reflect.TypeOf(sql.NullTime{}):
		return reflect.TypeOf(time.Time{})
	case reflect.TypeOf(pq.StringArray{}):
		return reflect.SliceOf(reflect.TypeOf(""))
	case reflect.TypeOf(pq.Int32Array{}):
		return reflect.SliceOf(reflect.TypeOf(0))
	case reflect.TypeOf(pq.Int64Array{}):
		return reflect.SliceOf(reflect.TypeOf(0))
	default:
		panic(fmt.Sprintf("Unsupported type %s", t.Kind().String()))
	}
}

func convertToNullType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Slice {
		switch t.Elem().Kind() {
		case reflect.Int:
			return reflect.TypeOf(pq.Int32Array{})
		case reflect.String:
			return reflect.TypeOf(pq.StringArray{})
		}
	}

	switch t.Kind() {
	case reflect.Int:
		return reflect.TypeOf(sql.NullInt32{})
	case reflect.String:
		return reflect.TypeOf(sql.NullString{})
	}

	panic("Unsupported type: " + t.Kind().String())
}
