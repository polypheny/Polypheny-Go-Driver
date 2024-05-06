package polypheny

import (
	context "context"
	testing "testing"
)

func TestQueryMongoContext(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer conn.(*PolyphenyConn).Close()
	_, err = conn.(*PolyphenyConn).QueryMongoContext(context.Background(), "db.emps.find()")
	if err != nil {
		t.Error(err)
	}
}
