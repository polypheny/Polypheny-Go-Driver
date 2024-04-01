package polypheny

import (
	"context"
	"database/sql"
	testing "testing"
)

type emps struct {
	empid      int32
	deptno     int32
	name       string
	salary     int32
	commission int32
}

type mongo struct {
	key   string
	value any
}

func TestSQLFlow(t *testing.T) {
	db, _ := sql.Open("polypheny", "localhost:20590,pa:")
	rows, _ := db.QueryContext(context.Background(), "sql:select * from emps")
	t.Log(rows.Columns())
	for rows.Next() {
		emp := new(emps)
		rows.Scan(&emp.empid, &emp.deptno, &emp.name, &emp.salary, &emp.commission)
		t.Log(emp)
	}
}

func TestMongoFlow(t *testing.T) {
	db, _ := sql.Open("polypheny", "localhost:20590,pa:")
	rows, _ := db.QueryContext(context.Background(), "mongo:db.emps.find()")
	t.Log(rows.Columns())
	for rows.Next() {
		mongo := new(mongo)
		rows.Scan(&mongo.key, &mongo.value)
		t.Log(mongo)
	}
}
