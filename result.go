package polypheny

// PolyphenyResult implements driver.Result
// and represents the results of queries that don't return data
type PolyphenyResult struct {
	lastInsertId int64
	rowsAffected int64
}

// LastInsertId returns the ID of the last inserted row
// TODO: does Polypheny support this at present?
func (r *PolyphenyResult) LastInsertId() (int64, error) {
	return r.lastInsertId, &ClientError{
		message: "This is not supported yet",
	}
}

// RowsAffected returns the number of affected rows of a query
func (r *PolyphenyResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}
