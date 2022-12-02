package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testTimeout = time.Second * 10
)

func TestSelect(t *testing.T) {
	argStr := "hello"
	argInt := 123

	rows, err := Select[struct {
		ArgStr string
		ArgInt int
	}](testTimeout, `select $1::text arg_str, $2::int4 arg_int`, argStr, argInt)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, argStr, rows[0].ArgStr)
	assert.Equal(t, argInt, rows[0].ArgInt)
}

func TestJsonQueryRow(t *testing.T) {
	type body struct {
		TestStr   string         `json:"testStr"`
		TestInt   int            `json:"testInt"`
		TestArray []string       `json:"testArray"`
		TestMap   map[string]int `json:"testMap"`
	}

	srcBody := &body{
		TestStr:   "test str",
		TestInt:   123,
		TestArray: []string{"test", "array"},
		TestMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	}

	destBody, err := QueryRow[body](testTimeout, `select $1::jsonb`, srcBody)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, destBody, srcBody)
}
