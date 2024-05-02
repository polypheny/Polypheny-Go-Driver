package polypheny

import (
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
