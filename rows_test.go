package polypheny

import (
	"database/sql/driver"
	io "io"
	"sync/atomic"
	testing "testing"
)

func TestColumns(t *testing.T) {
	columns := []string{"id"}
	result := [][]any{{int32(1)}, {int32(2)}, {int32(3)}}
	rows := PolyphenyRows{
		columns:   columns,
		result:    result,
		readIndex: atomic.Int32{},
	}
	rows.readIndex.Store(0)
	for i, v := range columns {
		if v != rows.Columns()[i] {
			t.Errorf("Column mismatches: expecting %s but got %s", v, rows.Columns()[i])
		}
	}
}

func TestRowsClose(t *testing.T) {
	columns := []string{"id"}
	result := [][]any{{int32(1)}, {int32(2)}, {int32(3)}}
	rows := PolyphenyRows{
		columns:   columns,
		result:    result,
		readIndex: atomic.Int32{},
	}
	rows.readIndex.Store(0)
	rows.Close()
	if rows.readIndex.Load() != -1 {
		t.Error("Error closing result")
	}
}

func TestNext(t *testing.T) {
	columns := []string{"id"}
	result := [][]any{{int32(1)}, {int32(2)}, {int32(3)}}
	rows := PolyphenyRows{
		columns:   columns,
		result:    result,
		readIndex: atomic.Int32{},
	}
	rows.readIndex.Store(0)
	dest := make([]driver.Value, 1)
	err := rows.Next(dest)
	if err != nil {
		t.Error(err)
	}
	if len(dest) != 1 || dest[0] != int32(1) {
		t.Error("Error Next")
	}
	err = rows.Next(dest)
	if err != nil {
		t.Error(err)
	}
	if len(dest) != 1 || dest[0] != int32(2) {
		t.Error("Error Next")
	}
	err = rows.Next(dest)
	if err != nil {
		t.Error(err)
	}
	if len(dest) != 1 || dest[0] != int32(3) {
		t.Error("Error Next")
	}
	err = rows.Next(dest)
	if err != io.EOF {
		t.Error(err)
	}
	rows.Close()
	err = rows.Next(dest)
	if err.Error() != "The Rows iterator has been closed" {
		t.Error("Expecting an error indicting the rows has been closed")
	}
}
