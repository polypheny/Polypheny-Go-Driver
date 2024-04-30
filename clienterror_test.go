package polypheny

import (
	testing "testing"
)

func TestClientError(t *testing.T) {
	err := func() error {
		return &ClientError{
			message: "test",
		}
	}()
	if err.Error() != "test" {
		t.Errorf("The error message should be 'test', but got %s", err.Error())
	}
}
