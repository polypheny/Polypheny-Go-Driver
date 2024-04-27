package polypheny

// These flow tests will be removed when the context support is done.

import (
	//"context"
	"context"
	"database/sql"
	testing "testing"
	"time"
)

type emps struct {
	empid      int32
	deptno     int32
	name       string
	salary     int32
	commission int32
}

type mongo struct {
	empid      int32
	commission int32
	salary     int32
	deptno     int32
	name       string
}

func TestSQLFlow(t *testing.T) {
	db, err := sql.Open("polypheny", "localhost:20590,pa:")
	t.Log(err)
	err = db.Ping()
	t.Log(err)
	rows, err := db.Query("sql:select * from emps")
	t.Log(err)
	t.Log(rows.Columns())
	for rows.Next() {
		emp := new(emps)
		rows.Scan(&emp.empid, &emp.deptno, &emp.name, &emp.salary, &emp.commission)
		t.Log(emp)
	}
}

func TestPingContext(t *testing.T) {
	db, err := sql.Open("polypheny", "localhost:20590,pa:")
	t.Log(err)
	err = db.Ping()
	t.Log(err)
	err = db.PingContext(context.Background())
	t.Log(err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Microsecond))
	err = db.PingContext(ctx)
	t.Log(err)
	defer cancel()
}

func TestMongoFlow(t *testing.T) {
	db, _ := sql.Open("polypheny", "localhost:20590,pa:")
	rows, _ := db.Query("mongo:db.emps.find()")
	t.Log(rows.Columns())
	for rows.Next() {
		mongo := new(mongo)
		rows.Scan(&mongo.name, &mongo.deptno, &mongo.salary, &mongo.empid, &mongo.commission)
		t.Log(mongo)
	}
}

type mytable struct {
	id  int32
	yac string
}

func TestExecFlow(t *testing.T) {
	db, _ := sql.Open("polypheny", "localhost:20590,pa:")
	result, _ := db.Exec("sql:drop table if exists mytable")
	t.Log(result.RowsAffected())
	result, _ = db.Exec("sql:create table mytable(id int not null, yac varchar(10), primary key(id))")
	t.Log(result.RowsAffected())
	result, _ = db.Exec("sql:insert into mytable values(1, 'hello')")
	t.Log(result.RowsAffected())
	result, _ = db.Exec("sql:insert into mytable values(2, 'world')")
	t.Log(result.RowsAffected())
	result, _ = db.Exec("sql:update mytable set yac = 'polypheny' where id in (select id from mytable where id = 1 or id = 2)")
	t.Log(result.RowsAffected())
	stmt, err := db.Prepare("sql:select * from mytable where id = ?")
	t.Log(err)
	rows, err := stmt.Query(1)
	t.Log(err)
	//rows, _ := db.Query("sql:select * from mytable")
	t.Log(rows.Columns())
	for rows.Next() {
		temp := new(mytable)
		rows.Scan(&temp.id, &temp.yac)
		t.Log(temp)
	}
	rows, _ = db.Query("sql:select * from mytable")
	t.Log(rows.Columns())
	for rows.Next() {
		temp := new(mytable)
		rows.Scan(&temp.id, &temp.yac)
		t.Log(temp)
	}
}
