package polypheny

import (
	context "context"
	driver "database/sql/driver"
	atomic "sync/atomic"
	testing "testing"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"

	proto "google.golang.org/protobuf/proto"
)

func TestIsValid(t *testing.T) {
	conn := PolyphenyConn{
		address:     "",
		username:    "",
		netConn:     nil,
		isConnected: atomic.Int32{},
	}
	conn.isConnected.Store(statusDisconnected)
	if conn.IsValid() != false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusDisconnected)
	}
	conn.isConnected.Store(statusServerConnected)
	if conn.IsValid() != false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusServerConnected)
	}
	conn.isConnected.Store(statusPolyphenyConnected)
	if conn.IsValid() == false {
		t.Errorf("Failed to set connection status, the status is %d, but should be %d.", conn.isConnected.Load(), statusPolyphenyConnected)
	}
}

func TestPingInternal(t *testing.T) {
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
	errChan := make(chan error)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
	conn.(*PolyphenyConn).isConnected.Store(statusDisconnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != driver.ErrBadConn {
		t.Errorf("Expected to receive an ErrBadConn but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusServerConnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err.Error() != "Ping: invalid connection to Polypheny server" {
		t.Errorf("Expected to receive a ClientError but got %s", err.Error())
	}
	conn.(*PolyphenyConn).isConnected.Store(statusPolyphenyConnected)
	go conn.(*PolyphenyConn).pingInternal(errChan)
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
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
	err = conn.(*PolyphenyConn).Ping(context.Background())
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = conn.(*PolyphenyConn).Ping(ctx)
	if err.Error() != "Context cancelled or timeout" {
		t.Errorf("Expected to receive a ClientError error but got %s", err.Error())
	}
}

func TestPrepare(t *testing.T) {
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
	stmt, err := conn.(*PolyphenyConn).Prepare("sql:SELECT * FROM emps WHERE name = ?")
	if err != nil {
		t.Error(err)
	}
	if len(stmt.(*PreparedStatement).args) != 1 {
		t.Errorf("Expected 1 arg in the prepared statement, but got %d", len(stmt.(*PreparedStatement).args))
	}
	err = stmt.Close()
	if err != nil {
		t.Error(err)
	}
	stmt, err = conn.(*PolyphenyConn).Prepare("sql:SELECT * FROM emps WHERE name = ? AND salary = ?")
	if err != nil {
		t.Error(err)
	}
	if len(stmt.(*PreparedStatement).args) != 2 {
		t.Errorf("Expected 2 args in the prepared statement, but got %d", len(stmt.(*PreparedStatement).args))
	}
	err = stmt.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).close()
	if err != nil {
		t.Error(err)
	}
	if conn.(*PolyphenyConn).isConnected.Load() != statusDisconnected {
		t.Errorf("The connection status should be %d, but got %d", statusDisconnected, conn.(*PolyphenyConn).isConnected.Load())
	}
}

func TestBegin(t *testing.T) {
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
	tx, err := conn.(*PolyphenyConn).Begin()
	if err != nil {
		t.Error(err)
	}
	if tx.(*PolyphenyTranaction).conn != conn {
		t.Error("Error when trys to start a transaction")
	}
}

func TestBeginTx(t *testing.T) {
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
	tx, err := conn.(*PolyphenyConn).BeginTx(context.Background(), driver.TxOptions{})
	if err != nil {
		t.Error(err)
	}
	if tx.(*PolyphenyTranaction).conn != conn {
		t.Error("Error when trys to start a transaction")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = conn.(*PolyphenyConn).BeginTx(ctx, driver.TxOptions{})
	if err != ctx.Err() {
		t.Error("context didn't work")
	}
}

func TestExec(t *testing.T) {
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
	_, err = conn.(*PolyphenyConn).Exec("sql:DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).Exec("sql:CREATE TABLE mytable(id int not null, primary key(id))", nil)
	if err != nil {
		t.Error(err)
	}
	result, err := conn.(*PolyphenyConn).Exec("sql:insert into mytable values(1)", nil)
	if err != nil {
		t.Error(err)
	}
	if result.(*PolyphenyResult).rowsAffected != 1 {
		t.Errorf("The number of affected rows should be 1 but got %d", result.(*PolyphenyResult).rowsAffected)
	}
}

func TestExecInternal(t *testing.T) {
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
	errChan := make(chan error)
	resultChan := make(chan *PolyphenyResult)
	var result *PolyphenyResult
	go conn.(*PolyphenyConn).execContextInternal("DROP TABLE IF EXISTS mytable", resultChan, errChan)
	result = <-resultChan
	err = <-errChan
	if err.Error() != "A query should have the following format: QueryLanguage:Query" || result != nil {
		t.Error("Expecting a ClientError")
	}
	go conn.(*PolyphenyConn).execContextInternal("sql:DROP TABLE IF EXISTS mytable", resultChan, errChan)
	result = <-resultChan
	err = <-errChan
	if err != nil {
		t.Error(err, result == nil)
	}
	go conn.(*PolyphenyConn).execContextInternal("sql:CREATE TABLE mytable(id int not null, primary key(id))", resultChan, errChan)
	result = <-resultChan
	err = <-errChan
	if err != nil {
		t.Error(err, result == nil)
	}
	go conn.(*PolyphenyConn).execContextInternal("sql:insert into mytable values(1)", resultChan, errChan)
	result = <-resultChan
	err = <-errChan
	if err != nil {
		t.Error(err, result == nil)
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if rowsaffected > 1 {
		t.Errorf("The number of affected rows should be 1 but got %d", rowsaffected)
	}
}

func TestExecContext(t *testing.T) {
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
	_, err = conn.(*PolyphenyConn).ExecContext(context.Background(), "sql:DROP TABLE IF EXISTS mytable", nil)
	if err != nil {
		t.Error(err)
	}
	// these queries are intentionally written in a wrong format
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = conn.(*PolyphenyConn).ExecContext(ctx, "DROP TABLE IF EXISTS mytable", nil)
	if err != ctx.Err() {
		t.Error(err)
	}
}

func TestQuery(t *testing.T) {
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
	rows, err := conn.(*PolyphenyConn).Query("sql:SELECT name FROM emps WHERE name = 'Bill'", nil)
	if err != nil {
		t.Error(err)
	}
	if len(rows.(*PolyphenyRows).Columns()) != 0 && rows.(*PolyphenyRows).Columns()[0] != "name" {
		t.Error("Error in Query")
	}
}

func TestQueryInternal(t *testing.T) {
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
	errChan := make(chan error)
	rowsChan := make(chan *PolyphenyRows)
	var result *PolyphenyRows
	go conn.(*PolyphenyConn).queryContextInternal("sql:SELECT name FROM emps WHERE name = 'Bill'", rowsChan, errChan)
	result = <-rowsChan
	err = <-errChan
	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("Error in queryContextInternal, result should not be nil")
	} else if len(result.Columns()) != 0 && result.Columns()[0] != "name" {
		t.Error("Error in queryContextInternal")
	}
}

func TestQueryContext(t *testing.T) {
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
	_, err = conn.(*PolyphenyConn).QueryContext(context.Background(), "sql:SELECT name FROM emps WHERE name = 'Bill'", nil)
	if err != nil {
		t.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = conn.(*PolyphenyConn).ExecContext(ctx, "sql:SELECT name FROM emps WHERE name = 'Bill'", nil)
	if err != ctx.Err() {
		t.Error(err)
	}
}

func TestSend(t *testing.T) {
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
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  "sql",
				Statement:     "DROP TABLE IF EXISTS mytable;CREATE TABLE mytable(id int not null, primary key(id));insert into mytable values(1);insert into mytable values(2);SELECT id FROM mytable where id = 2",
				FetchSize:     nil,
				NamespaceName: nil,
			},
		},
	}
	serializd, err := proto.Marshal(&request)
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).send(serializd, 1)
	if err.Error() != "PolyphenyConn.send(): the size of the serialized message is too large to be put in a byte array of size lengthSize" {
		t.Error(err)
	}
	request = prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	serializd, err = proto.Marshal(&request)
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).send(serializd, 8)
	if err != nil {
		t.Error(err)
	}
	conn.(*PolyphenyConn).recv(8)
}

func TestRecv(t *testing.T) {
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
	_, err = conn.(*PolyphenyConn).recv(16)
	if err == nil {
		t.Error("Expecting an err when trying to recv with lengthSize > 8")
	}
	request := prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	serializd, err := proto.Marshal(&request)
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).send(serializd, 8)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.(*PolyphenyConn).recv(8)
	if err != nil {
		t.Error(err)
	}
}

func TestNetConnClose(t *testing.T) {
	connector := Connector{
		address:  "localhost:20590",
		username: "pa",
		password: "",
	}
	conn, err := connector.Connect(context.Background())
	if err != nil {
		t.Error(err)
	}
	err = conn.(*PolyphenyConn).close()
	if err != nil {
		t.Error(err)
	}
}

func TestHelperSendAndRecv(t *testing.T) {
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
	request := prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	_, err = conn.(*PolyphenyConn).helperSendAndRecv(&request)
	if err != nil {
		t.Error(err)
	}
}
