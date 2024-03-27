package polypheny

import (
	bytes "bytes"
	prism "polypheny/protos"
	testing "testing"
)

func TestNewConnection(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	defer conn.conn.Close()
	if conn.isConnected != statusServerConnected {
		t.Fatalf("Failed to make a connection. The current status is %v", conn.isConnected)
	}
}

func TestSerialize(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	username := "pa"
	password := ""
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &username,
				Password:        &password,
			},
		},
	}
	result := conn.serialize(&request)
	expected := []byte{154, 1, 8, 8, 2, 42, 2, 112, 97, 50, 0}
	if !bytes.Equal(result, expected) {
		t.Fatalf("Error when serializing a request got %v,but expected %v", result, expected)
	}
}

func TestSend(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	username := "pa"
	password := ""
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &username,
				Password:        &password,
			},
		},
	}
	conn.send(conn.serialize(&request))
	conn.recv()
}

func TestRecv(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	username := "pa"
	password := ""
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &username,
				Password:        &password,
			},
		},
	}
	conn.send(conn.serialize(&request))
	result := conn.recv()
	expected := []byte{16, 1, 98, 4, 8, 1, 16, 2}
	if !bytes.Equal(result, expected) {
		t.Fatalf("Error when receiving a response got %v,but expected %v", result, expected)
	}
}

func TestClose(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	conn.close()
	if conn.isConnected != statusDisconnected {
		t.Fatalf("Failed to disconnect, the current client status is %v", conn.isConnected)
	}
}

func TestHelperSendAndRecv(t *testing.T) {
	conn := newConnection("localhost:20590", "pa")
	username := "pa"
	password := ""
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &username,
				Password:        &password,
			},
		},
	}
	response := conn.helperSendAndRecv(&request)
	if response.GetConnectionResponse().GetIsCompatible() != true {
		t.Fatalf("The connection response is not correct")
	}
}

func TestHandleConnectRequest(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	if client.isConnected != statusPolyphenyConnected {
		t.Fatalf("Failed to connect to Polypheny, the current client status is %v", client.isConnected)
	}
}

func TestHandleConnectionPropertiesUpdateRequest(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	isAutoCommit := true
	client.handleConnectionPropertiesUpdateRequest(&isAutoCommit, nil)
}

func TestHandleConnectionCheckRequest(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	client.handleConnectionCheckRequest()
}

func TestHandleDisconnectRequest(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	client.handleDisconnectRequest()
	if client.isConnected != statusDisconnected {
		t.Fatalf("Failed to disconnect, the current client status is %v", client.isConnected)
	}
}

func TestMakeProtoValue1(t *testing.T) {
	var result *prism.ProtoValue
	var value interface{}
	value = true
	result = makeProtoValue(value)
	if result.GetBoolean().GetBoolean() != true {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetBoolean().GetBoolean())
	}
	value = int32(1)
	result = makeProtoValue(value)
	if result.GetInteger().GetInteger() != value.(int32) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetInteger().GetInteger())
	}
	value = int64(100000000000)
	result = makeProtoValue(value)
	if result.GetLong().GetLong() != value.(int64) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetLong().GetLong())
	}
	value = "Hello, world!"
	result = makeProtoValue(value)
	if result.GetString_().GetString_() != value.(string) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetString_().GetString_())
	}
}

func TestConvertProtoValue(t *testing.T) {
	var protoValue *prism.ProtoValue
	var result interface{}
	var expected interface{}
	expected = true
	protoValue = makeProtoValue(expected)
	result = convertProtoValue(protoValue)
	if result.(bool) != expected {
		t.Fatalf("Failed to convert, expected %v, but got %v", expected, result)
	}
}

func TestHandleExecuteUnparameterizedStatementRequestSql(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	query := UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "drop table if exists mytable",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "create table mytable(id int not null, yac int, primary key(id))",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "insert into mytable values(1, 1)",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "insert into mytable values(2, 2)",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "select * from mytable",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleCommitRequest()
	result := client.handleExecuteUnparameterizedStatementRequest(query)
	t.Log(result)
}

func TestHandleExecuteUnparameterizedStatementRequestMongo(t *testing.T) {
	address := "localhost:20590"
	username := "pa"
	password := ""
	client := handleConnectRequest(address, username, password)
	query := UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "drop table if exists mytable",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "create table mytable(id int not null, yac int, primary key(id))",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "insert into mytable values(1, 1)",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	query = UnparameterizedStatementRequest{
		language:      "sql",
		statement:     "insert into mytable values(2, 2)",
		fetchSize:     nil,
		namespaceName: nil,
	}
	client.handleExecuteUnparameterizedStatementRequest(query)
	client.handleCommitRequest()
	query = UnparameterizedStatementRequest{
		language:      "mongo",
		statement:     "db.mytable.find()",
		fetchSize:     nil,
		namespaceName: nil,
	}
	result := client.handleExecuteUnparameterizedStatementRequest(query)
	t.Log(result)
}
