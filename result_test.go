package polypheny

import (
	testing "testing"
)

func TestLastInsertId(t *testing.T) {
	result := PolyphenyResult{}
	_, err := result.LastInsertId()
	if err == nil {
		t.Error("Expecting clienterror")
	}
}

func TestRowsAffected(t *testing.T) {
	result := PolyphenyResult{
		rowsAffected: 0,
	}
	_, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
	}
}
