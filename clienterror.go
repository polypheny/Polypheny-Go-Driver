package polypheny

// ClientError reports a driver or server error.
type ClientError struct {
	message      string
	sqlState     string
	errorCode    int32
	hasErrorCode bool
}

func (e *ClientError) Error() string {
	return e.message
}

// SQLState returns the server's SQLSTATE, or an empty string if none was supplied.
func (e *ClientError) SQLState() string {
	return e.sqlState
}

// ErrorCode returns the server's numeric error code and whether it was supplied.
// A supplied code may be zero.
func (e *ClientError) ErrorCode() (int32, bool) {
	return e.errorCode, e.hasErrorCode
}
