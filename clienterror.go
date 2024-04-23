package polypheny

type ClientError struct {
	message string
}

func (e *ClientError) Error() string {
	return e.message
}
