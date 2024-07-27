package model

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"reflect"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

type Model interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	Write
	Query
}

func parseDataStruct(s interface{}) (map[string]interface{}, string, error) {
	ptyp := reflect.TypeOf(s)
	if ptyp.Kind() != reflect.Ptr {
		return nil, "", fmt.Errorf("%s is not a pointer", ptyp)
	}
	typ := ptyp.Elem()
	if typ.Kind() != reflect.Struct {
		return nil, "", fmt.Errorf("%s is not a struct", typ)
	}
	m := make(map[string]interface{})
	values := reflect.ValueOf(s).Elem()
	dataID := ""
	ok := false
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		if key := fld.Tag.Get("db"); key != "" {
			if values.Field(i).CanInterface() && !values.Field(i).IsZero() {
				val := values.Field(i).Interface()
				if typ.Field(i).Type.Kind() == reflect.Slice {
					switch typ.Field(i).Type.Elem().Kind() {
					case reflect.String:
						val = pq.StringArray(val.([]string))
					case reflect.Int:
						val = pq.Int32Array(val.([]int32))
					default:
						return nil, "", fmt.Errorf("Kind %s not supported", typ.Field(i).Type.Elem().Kind())
					}
				}
				m[key] = val
				if key == "id" {
					dataID, ok = val.(string)
					if !ok {
						return nil, "", fmt.Errorf("id is not a string")
					}
				}
			} else {
				if typ.Field(i).Type.Kind() == reflect.Slice {
					fmt.Println("Arr")
					switch typ.Field(i).Type.Elem().Kind() {
					case reflect.String:
						m[key] = pq.StringArray{}
					default:
						return nil, "", fmt.Errorf("Kind %s not supported", typ.Field(i).Type.Elem().Kind())
					}
				}
			}
		}
	}
	if dataID == "" {
		return nil, "", fmt.Errorf("id is empty")
	}
	return m, dataID, nil
}
