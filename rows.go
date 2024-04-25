package polypheny

import (
	driver "database/sql/driver"
	io "io"
)

// PolyphenyRows implements driver.Rows
// and represents the results of queries that return data
type PolyphenyRows struct {
	columns   []string
	result    [][]any
	readIndex int
}

// Columns returns the names of columns of a query
func (r *PolyphenyRows) Columns() []string {
	return r.columns
}

// Close will close the Rows iterator
func (r *PolyphenyRows) Close() error {
	r.readIndex = -1
	return nil
}

// Next iterates over the Rows
func (r *PolyphenyRows) Next(dest []driver.Value) error {
	if r.readIndex == -1 {
		return &ClientError{
			message: "The Rows iterator has been closed",
		}
	}
	if r.readIndex >= len(r.result) {
		return io.EOF
	}
	for i := range dest {
		dest[i] = r.result[r.readIndex][i]
	}
	r.readIndex++
	return nil
}
