package polypheny

import (
	testing "testing"
	protos "polypheny.com/protos"
	context "context"
	time "time"
	grpc "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
)

func TestProtoConnect(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "", true)
	if client.isConnected != true {
		t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
	}
}

func TestUpdateConnection(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "", true)
        if client.isConnected != true {
                t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
        }
	client.handleUpdateConnectionProperties(true)
}

func TestCheckConnection(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "", true)
        if client.isConnected != true {
                t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
        }
        client.handleConnectionCheckRequest()
}

func TestProtoDuplicateConnection(t *testing.T) {
	username := "pa"
	password := ""
	is_auto_commit := true
	client := handleConnectRequest("localhost:20590", username, password, is_auto_commit)
	request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID, // use the same uuid as the previous client
		Username: &username,
		Password: &password,
        }
	_, err := client.client.Connect(client.ctx, &request)
	if err == nil {
		t.Fatal("Expected an error when sending a connect request with a duplicated UUID")
	}
}

func TestProtoDisconnect(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "")
	client.handleDisconnectRequest()
	if client.isConnected != false {
		t.Fatal("Protoclient.isConnected is not false after disconnecting, however no error was caught")
	}
	// now we should be able to connect to the server with the previous uuid
	username := "pa"
        password := ""
	request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID, // use the same uuid as the previous client
		Username: &username,
                Password: &password,
        }
	// we cannot use client.ctx anymore, since it has been canceled
	// same with the client.Connect
	// therefore, we need to rebuild these objects
	conn, err := grpc.Dial("localhost:20590", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
        if err != nil {
                t.Fatalf("did not connect: %v", err)
        }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := protos.NewProtoInterfaceClient(conn)
	_, err = c.Connect(ctx, &request)
	if err != nil {
		//t.Fatal("Failed to re-connect with the same uuid")
		t.Fatal(err)
	}
}

func TestGetClientInfoProperties(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetClientInfoProperties()
        t.Log(result)
}

func TestGetDBMSVersion(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "")
	dbms, version, major, minor := client.handleGetDBMSVersion()
	t.Log(dbms, version, major, minor)
}

func TestGetSupportedLanguage(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        langs := client.handleGetSupportedLanguage()
        t.Log(langs)
}

func TestGetDatabases(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetDatabases()
        t.Log(result)
}

func TestGetTableTypes(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetTableTypes()
        t.Log(result)
}

func TestGetTypes(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetTypes()
	var names []string
	for _, v := range result {
		names = append(names, v.GetTypeName())
	}
        t.Log(names)
}

func TestGetUserDefinedTypes(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetUserDefinedTypes()
        t.Log(result)
}

func TestGetClientInfoPropertyMetas(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetClientInfoPropertyMetas()
        t.Log(result)
}

func TestSearchProcedures(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleSearchProcedures("sql")
	t.Log(result)
}

func TestSearchNamespaces(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleSearchNamespaces("")
        t.Log(result)
}

func TestGetNamespace(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetNamespace("public")
        t.Log(result)
}

func TestSearchEntities(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleSearchEntities("public")
        t.Log(result)
}

func TestGetSqlStringFunctions(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetSqlStringFunctions()
        t.Log(result)
}

func TestGetSqlSystemFunctions(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetSqlSystemFunctions()
        t.Log(result)
}

func TestGetSqlTimeDateFunctions(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetSqlTimeDateFunctions()
        t.Log(result)
}

func TestGetSqlNumericFunctions(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetSqlNumericFunctions()
        t.Log(result)
}

func TestGetSqlKeywords(t *testing.T) {
        client := handleConnectRequest("localhost:20590", "pa", "")
        result := client.handleGetSqlKeywords()
        t.Log(result)
}

func TestMakeProtoValue(t *testing.T) {
	t.Log(makeProtoValue(int32(1)))
	t.Log(makeProtoValue(1.1))
	t.Log(makeProtoValue(true))
	t.Log(makeProtoValue("Hello, world"))
	pv := makeProtoValue(int32(1))
	t.Log(pv.GetInteger().GetInteger())
	pv = makeProtoValue(1.1)
        t.Log(pv.GetFloat().GetFloat())
	pv = makeProtoValue(true)
        t.Log(pv.GetBoolean().GetBoolean())
	pv = makeProtoValue("Hello, world")
        t.Log(pv.GetString_().GetString_())
}
