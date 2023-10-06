package polypheny

import (
	testing "testing"
)

func TestFlow(t *testing.T) {
	conn := Connect("localhost:20590")
	conn.Execute("SELECT * FROM emps", "sql")
	result := conn.Fetch()
	t.Log(result)
	conn.Close()
}
