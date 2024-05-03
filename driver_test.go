package polypheny

import (
	testing "testing"
)

func TestParseDSN(t *testing.T) {
	// wrong dsn format
	dsn := "localhost:20590"
	addr, user, password, err := parseDSN(dsn)
	if err == nil {
		t.Errorf("Expecting ClientError but got: %s %s %s", addr, user, password)
	}
	dsn = "localhost:20590,pa"
	addr, user, password, err = parseDSN(dsn)
	if err == nil {
		t.Errorf("Expecting ClientError but got: %s %s %s", addr, user, password)
	}
	dsn = "localhost:20590,pa:"
	addr, user, password, err = parseDSN(dsn)
	if err != nil {
		t.Errorf("Error. Got: %s %s %s %s", addr, user, password, err.Error())
	}
	if addr != "localhost:20590" || user != "pa" || password != "" {
		t.Errorf("Error. Got: %s %s %s", addr, user, password)
	}
}

func TestOpen(t *testing.T) {
	d := PolyphenyDriver{}
	conn, err := d.Open("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	conn.Close()
	_, err = d.Open("localhost:20590,pa")
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestOpenConnector(t *testing.T) {
	d := PolyphenyDriver{}
	_, err := d.OpenConnector("localhost:20590,pa:")
	if err != nil {
		t.Error(err)
	}
	_, err = d.OpenConnector("localhost:20590,pa")
	if err == nil {
		t.Error("Expecting error")
	}
}
