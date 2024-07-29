package model

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseDataStruct(t *testing.T) {
	t.Run("check array fields", func(t *testing.T) {
		type c struct {
			Array []string  `db:"array"`
			ID    string    `db:"id"`
			Time  time.Time `db:"time"`
			K     []string  `db:"k"`
		}
		k := &c{
			Array: []string{"a", "b", "c"},
			ID:    "fsdfs",
		}
		a, _, err := parseDataStruct(k)
		assert.NoError(t, err)
		val, ok := a["array"]
		assert.True(t, ok)
		fmt.Println("val", val)
		assert.NoError(t, err)

		assert.Equal(t, []string{"a", "b", "c"}, a["array"])
	})
}

func TestGenerateNullStruct(t *testing.T) {
	type customs struct {
		A string
		B []string
	}
	type check struct {
		abc  []int
		name string
		b    time.Time
		t    customs
	}
	abc := check{}
	fmt.Printf("New %v \n", generateNullTypeStruct(abc))
}
func TestGenerateNormalStruct(t *testing.T) {
	type test struct {
		A sql.NullString
		B sql.NullInt32
		C sql.NullTime
		D pq.StringArray
	}
	type yo struct {
		A string
		B int
		C time.Time
		D []string
	}

	k := &test{
		A: sql.NullString{String: "hey"},
		B: sql.NullInt32{Int32: 123},
		C: sql.NullTime{Time: time.Now()},
		D: pq.StringArray{"a", "b", "c"},
	}
	d := &yo{}

	generateNormalStruct(k, d)
	fmt.Printf("New %v \n", k)
	fmt.Printf("New A: %s B: %d\n", d.A, d.B)
}

func TestGenerateSQLString(t *testing.T) {
	s := &SelectQuery{
		TableName: "test",
		Fields:    []string{"a", "b", "c"},
		Where: []Condition{
			{
				Clause: "id",
				Param:  "U123",
			},
			{
				Clause: "name",
				Param:  "kk",
			},
		},
	}
	out, err := generateSelectQuery(s)
	assert.NoError(t, err)
	fmt.Println(out)
}
